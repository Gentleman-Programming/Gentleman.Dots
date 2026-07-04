package system

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// shellPath caches the detected shell path
var (
	shellPath     string
	shellPathOnce sync.Once
)

// GetShell returns the path to an available shell (bash, sh, or zsh)
// Caches the result for subsequent calls
func GetShell() string {
	shellPathOnce.Do(func() {
		// In Termux, prefer sh to avoid fork/exec issues with bash
		// Go has known issues with fork/exec on Android
		if isTermux() {
			if path, err := exec.LookPath("sh"); err == nil {
				shellPath = path
				return
			}
		}

		// Try bash first (most compatible with our commands)
		if path, err := exec.LookPath("bash"); err == nil {
			shellPath = path
			return
		}
		// Fall back to sh (available on almost all Unix systems including Termux)
		if path, err := exec.LookPath("sh"); err == nil {
			shellPath = path
			return
		}
		// Last resort: zsh
		if path, err := exec.LookPath("zsh"); err == nil {
			shellPath = path
			return
		}
		// Default to sh if nothing found (let it fail with a clear error)
		shellPath = "sh"
	})
	return shellPath
}

// ExecError provides detailed error information for command execution
type ExecError struct {
	Command  string
	ExitCode int
	Stderr   string
	Stdout   string
	Wrapped  error
}

func (e *ExecError) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("command failed: %s\n", e.Command))
	sb.WriteString(fmt.Sprintf("exit code: %d\n", e.ExitCode))
	if e.Stderr != "" {
		sb.WriteString(fmt.Sprintf("stderr: %s\n", strings.TrimSpace(e.Stderr)))
	}
	if e.Stdout != "" && e.Stderr == "" {
		sb.WriteString(fmt.Sprintf("output: %s\n", strings.TrimSpace(e.Stdout)))
	}
	if e.Wrapped != nil {
		sb.WriteString(fmt.Sprintf("cause: %v", e.Wrapped))
	}
	return sb.String()
}

func (e *ExecError) Unwrap() error {
	return e.Wrapped
}

type ExecResult struct {
	Output   string
	Stderr   string
	Error    error
	ExitCode int
	Duration time.Duration
	Command  string
}

type ExecOptions struct {
	ShowOutput bool
	WorkDir    string
	Env        []string
	Timeout    time.Duration
}

// parseCommand splits a command string into executable and arguments
// This is a simple parser that handles basic quoting
func parseCommand(command string) (string, []string) {
	var args []string
	var current strings.Builder
	inQuote := false
	quoteChar := rune(0)

	for _, r := range command {
		if (r == '"' || r == '\'') && !inQuote {
			inQuote = true
			quoteChar = r
		} else if r == quoteChar && inQuote {
			inQuote = false
			quoteChar = 0
		} else if r == ' ' && !inQuote {
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		} else {
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}

	if len(args) == 0 {
		return "", nil
	}

	// In Termux, resolve executable path using $PREFIX/bin
	executable := args[0]
	if isTermux() {
		prefix := os.Getenv("PREFIX")
		if prefix != "" {
			fullPath := prefix + "/bin/" + executable
			if _, err := os.Stat(fullPath); err == nil {
				executable = fullPath
			}
		}
	}

	return executable, args[1:]
}

// Run executes a command and returns the result with detailed error information
func Run(command string, opts *ExecOptions) *ExecResult {
	if opts == nil {
		opts = &ExecOptions{}
	}

	start := time.Now()
	result := &ExecResult{
		Command: command,
	}

	ctx := context.Background()
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	// In Termux, execute commands directly without shell wrapper
	// Go has issues with fork/exec through shell on Android
	var cmd *exec.Cmd
	if isTermux() {
		executable, args := parseCommand(command)
		cmd = exec.CommandContext(ctx, executable, args...)
	} else {
		// Use available shell to run the command (bash, sh, or zsh)
		cmd = exec.CommandContext(ctx, GetShell(), "-c", command)
	}

	if opts.WorkDir != "" {
		cmd.Dir = opts.WorkDir
	}

	// Set environment
	cmd.Env = os.Environ()
	if len(opts.Env) > 0 {
		cmd.Env = append(cmd.Env, opts.Env...)
	}

	var stdout, stderr strings.Builder

	if opts.ShowOutput {
		// Stream output to stdout and capture it
		stdoutPipe, _ := cmd.StdoutPipe()
		stderrPipe, _ := cmd.StderrPipe()

		if err := cmd.Start(); err != nil {
			result.Error = &ExecError{
				Command: command,
				Wrapped: err,
			}
			result.Duration = time.Since(start)
			return result
		}

		go streamOutput(stdoutPipe, &stdout)
		go streamOutput(stderrPipe, &stderr)

		err := cmd.Wait()
		if err != nil {
			result.ExitCode = cmd.ProcessState.ExitCode()
			result.Error = &ExecError{
				Command:  command,
				ExitCode: result.ExitCode,
				Stdout:   stdout.String(),
				Stderr:   stderr.String(),
				Wrapped:  err,
			}
		}
	} else {
		// Capture stdout and stderr separately for better error messages
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			exitCode := 1
			if cmd.ProcessState != nil {
				exitCode = cmd.ProcessState.ExitCode()
			}
			result.ExitCode = exitCode
			result.Error = &ExecError{
				Command:  command,
				ExitCode: exitCode,
				Stdout:   stdout.String(),
				Stderr:   stderr.String(),
				Wrapped:  err,
			}
		}
	}

	result.Output = stdout.String()
	result.Stderr = stderr.String()
	result.Duration = time.Since(start)

	if cmd.ProcessState != nil {
		result.ExitCode = cmd.ProcessState.ExitCode()
	}

	return result
}

func streamOutput(r io.Reader, w *strings.Builder) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		w.WriteString(line + "\n")
	}
}

// RunSudo runs a command with sudo
func RunSudo(command string, opts *ExecOptions) *ExecResult {
	return Run("sudo "+command, opts)
}

// RunBrew runs a brew command
func RunBrew(args string, opts *ExecOptions) *ExecResult {
	brewPath := GetBrewPrefix() + "/bin/brew"
	return Run(brewPath+" "+args, opts)
}

// RunRpmOstree runs an rpm-ostree command (for atomic distros)
func RunRpmOstree(args string, opts *ExecOptions) *ExecResult {
	return Run("rpm-ostree "+args, opts)
}

// RunFlatpak runs a flatpak command
func RunFlatpak(args string, opts *ExecOptions) *ExecResult {
	return Run("flatpak "+args, opts)
}

// RunPkg runs a Termux pkg command (install packages)
func RunPkg(args string, opts *ExecOptions) *ExecResult {
	return Run("pkg "+args, opts)
}

// RunPkgWithLogs runs a Termux pkg command with log streaming
func RunPkgWithLogs(args string, opts *ExecOptions, logFunc func(string)) *ExecResult {
	return RunWithLogs("pkg "+args, opts, logFunc)
}

// RunPkgInstall runs pkg install with -y flag for non-interactive installs
func RunPkgInstall(packages string, opts *ExecOptions, logFunc func(string)) *ExecResult {
	return RunWithLogs("pkg install -y "+packages, opts, logFunc)
}

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("copy file %s: source is a directory", src)
	}

	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}

// CopyDir recursively copies a directory using native Go (shell-independent)
func CopyDir(src, dst string) error {
	// Clean paths - remove trailing /* or /. if present
	src = strings.TrimSuffix(strings.TrimSuffix(src, "/*"), "/.")
	walkRoot, err := filepath.EvalSymlinks(src)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		}
		walkRoot = src
	}

	rootInfo, err := os.Stat(walkRoot)
	if err != nil {
		return err
	}
	if !rootInfo.IsDir() {
		return fmt.Errorf("copy dir %s: source is not a directory", src)
	}
	if err := os.MkdirAll(dst, rootInfo.Mode()); err != nil {
		return err
	}

	return filepath.Walk(walkRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		resolvedInfo, err := os.Stat(path)
		if err != nil {
			return err
		}

		// Calculate destination path
		relPath, err := filepath.Rel(walkRoot, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if resolvedInfo.IsDir() {
			return os.MkdirAll(dstPath, resolvedInfo.Mode())
		}

		// Copy file
		return CopyFile(path, dstPath)
	})
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// BackupInfo contains information about a backup
type BackupInfo struct {
	Path      string
	Timestamp time.Time
	Files     []string
}

// ConfigPaths returns all config paths that Gentleman.Dots will modify
func ConfigPaths() map[string]string {
	home := os.Getenv("HOME")
	return map[string]string{
		"nvim":      home + "/.config/nvim",
		"fish":      home + "/.config/fish",
		"zsh":       home + "/.zshrc",
		"zsh_p10k":  home + "/.p10k.zsh",
		"oh-my-zsh": home + "/.oh-my-zsh",
		"nushell":   home + "/.config/nushell",
		"tmux":      home + "/.tmux.conf",
		"zellij":    home + "/.config/zellij",
		"herdr":     home + "/.config/herdr",
		"alacritty": home + "/.config/alacritty",
		"wezterm":   home + "/.wezterm.lua",
		"kitty":     home + "/.config/kitty",
		"ghostty":   home + "/.config/ghostty",
		"starship":  home + "/.config/starship.toml",
	}
}

// DetectExistingConfigs checks which config files/directories already exist
func DetectExistingConfigs() []string {
	existing := []string{}
	for name, path := range ConfigPaths() {
		if _, err := os.Stat(path); err == nil {
			existing = append(existing, name+": "+path)
		}
	}
	return existing
}

// GetBackupDir returns the backup directory path with timestamp
func GetBackupDir() string {
	home := os.Getenv("HOME")
	timestamp := time.Now().Format("2006-01-02-150405")
	return home + "/.gentleman-backup-" + timestamp
}

// ListBackups returns all existing backups
func ListBackups() []BackupInfo {
	home := os.Getenv("HOME")
	backups := []BackupInfo{}

	entries, err := os.ReadDir(home)
	if err != nil {
		return backups
	}

	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), ".gentleman-backup-") {
			backupPath := home + "/" + entry.Name()
			info, err := entry.Info()
			if err != nil {
				continue
			}

			// List files in backup
			files := []string{}
			subEntries, _ := os.ReadDir(backupPath)
			for _, sub := range subEntries {
				files = append(files, sub.Name())
			}

			backups = append(backups, BackupInfo{
				Path:      backupPath,
				Timestamp: info.ModTime(),
				Files:     files,
			})
		}
	}

	return backups
}

// CreateBackup creates a backup of existing configs
func CreateBackup(configs []string) (string, error) {
	backupDir := GetBackupDir()
	if err := EnsureDir(backupDir); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	configPaths := ConfigPaths()

	for _, configKey := range configs {
		// Extract key from "key: path" format if present
		key := configKey
		if idx := strings.Index(configKey, ":"); idx > 0 {
			key = configKey[:idx]
		}

		srcPath, exists := configPaths[key]
		if !exists {
			continue
		}

		// Check if source exists
		info, err := os.Stat(srcPath)
		if err != nil {
			continue // File doesn't exist, skip
		}

		// Determine destination path
		dstPath := backupDir + "/" + key

		if info.IsDir() {
			// Copy directory
			if err := CopyDir(srcPath, dstPath); err != nil {
				return backupDir, fmt.Errorf("failed to backup %s: %w", key, err)
			}
		} else {
			// Copy file
			if err := CopyFile(srcPath, dstPath); err != nil {
				return backupDir, fmt.Errorf("failed to backup %s: %w", key, err)
			}
		}
	}

	return backupDir, nil
}

// RestoreBackup restores configs from a backup directory
func RestoreBackup(backupDir string) error {
	configPaths := ConfigPaths()

	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return fmt.Errorf("failed to read backup directory: %w", err)
	}

	for _, entry := range entries {
		key := entry.Name()
		dstPath, exists := configPaths[key]
		if !exists {
			continue
		}

		srcPath := backupDir + "/" + key

		// Remove current config
		os.RemoveAll(dstPath)

		srcInfo, err := os.Stat(srcPath)
		if err != nil {
			continue
		}

		if srcInfo.IsDir() {
			if err := CopyDir(srcPath, dstPath); err != nil {
				return fmt.Errorf("failed to restore %s: %w", key, err)
			}
		} else {
			if err := CopyFile(srcPath, dstPath); err != nil {
				return fmt.Errorf("failed to restore %s: %w", key, err)
			}
		}
	}

	return nil
}

// DeleteBackup removes a backup directory
func DeleteBackup(backupDir string) error {
	return os.RemoveAll(backupDir)
}

// LogCallback is a function that receives log lines during command execution
type LogCallback func(line string)

// RunWithLogs executes a command and streams output to a callback function
// This allows the TUI to display real-time installation progress
func RunWithLogs(command string, opts *ExecOptions, onLog LogCallback) *ExecResult {
	if opts == nil {
		opts = &ExecOptions{}
	}

	start := time.Now()
	result := &ExecResult{
		Command: command,
	}

	ctx := context.Background()
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	// In Termux, execute commands directly without shell wrapper
	var cmd *exec.Cmd
	if isTermux() {
		executable, args := parseCommand(command)
		cmd = exec.CommandContext(ctx, executable, args...)
	} else {
		cmd = exec.CommandContext(ctx, GetShell(), "-c", command)
	}

	if opts.WorkDir != "" {
		cmd.Dir = opts.WorkDir
	}

	cmd.Env = os.Environ()
	if len(opts.Env) > 0 {
		cmd.Env = append(cmd.Env, opts.Env...)
	}

	var stdout, stderr strings.Builder

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		result.Error = &ExecError{Command: command, Wrapped: err}
		result.Duration = time.Since(start)
		return result
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		result.Error = &ExecError{Command: command, Wrapped: err}
		result.Duration = time.Since(start)
		return result
	}

	if err := cmd.Start(); err != nil {
		result.Error = &ExecError{Command: command, Wrapped: err}
		result.Duration = time.Since(start)
		return result
	}

	// Stream stdout with callback
	done := make(chan struct{})
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			stdout.WriteString(line + "\n")
			if onLog != nil {
				onLog(line)
			}
		}
		done <- struct{}{}
	}()

	// Stream stderr with callback
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			stderr.WriteString(line + "\n")
			if onLog != nil {
				onLog(line)
			}
		}
		done <- struct{}{}
	}()

	// Wait for both streams to complete
	<-done
	<-done

	err = cmd.Wait()
	if err != nil {
		exitCode := 1
		if cmd.ProcessState != nil {
			exitCode = cmd.ProcessState.ExitCode()
		}
		result.ExitCode = exitCode
		result.Error = &ExecError{
			Command:  command,
			ExitCode: exitCode,
			Stdout:   stdout.String(),
			Stderr:   stderr.String(),
			Wrapped:  err,
		}
	}

	result.Output = stdout.String()
	result.Stderr = stderr.String()
	result.Duration = time.Since(start)

	if cmd.ProcessState != nil {
		result.ExitCode = cmd.ProcessState.ExitCode()
	}

	return result
}

// RunBrewWithLogs runs a brew command with log streaming
func RunBrewWithLogs(args string, opts *ExecOptions, onLog LogCallback) *ExecResult {
	brewPath := GetBrewPrefix() + "/bin/brew"
	return RunWithLogs(brewPath+" "+args, opts, onLog)
}

// RunSudoWithLogs runs a sudo command with log streaming
func RunSudoWithLogs(command string, opts *ExecOptions, onLog LogCallback) *ExecResult {
	return RunWithLogs("sudo "+command, opts, onLog)
}

// RunBrewWithLogs runs a brew command with log streaming (already defined above, keeping the first one)

// RunRpmOstreeWithLogs runs an rpm-ostree command with log streaming
func RunRpmOstreeWithLogs(args string, opts *ExecOptions, onLog LogCallback) *ExecResult {
	return RunWithLogs("rpm-ostree "+args, opts, onLog)
}

// RunFlatpakWithLogs runs a flatpak command with log streaming
func RunFlatpakWithLogs(args string, opts *ExecOptions, onLog LogCallback) *ExecResult {
	return RunWithLogs("flatpak "+args, opts, onLog)
}

// PatchZshForWM modifies .zshrc based on window manager choice.
func PatchZshForWM(zshrcPath string, wm string, installNvim bool) error {
	content, err := os.ReadFile(zshrcPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	inStartIfNeeded := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.Contains(line, `eval "$(fzf --zsh)"`) {
			if !installNvim {
				newLines = append(newLines, `if command -v fzf &> /dev/null; then`)
				newLines = append(newLines, line)
				newLines = append(newLines, `fi`)
			} else {
				newLines = append(newLines, line)
			}
			continue
		}

		if strings.HasPrefix(trimmed, "function start_if_needed") {
			inStartIfNeeded = !strings.Contains(trimmed, "}")
			if wm != "none" {
				newLines = append(newLines, line)
			}
			continue
		}
		if inStartIfNeeded {
			if wm != "none" {
				if strings.Contains(trimmed, "[[ $- == *i* ]]") && strings.Contains(trimmed, "WM_VAR") {
					newLines = append(newLines, `    if [[ $- == *i* ]] && command -v "$WM_CMD" >/dev/null 2>&1 && [[ -z "${WM_VAR#/}" ]] && [[ -z "$TMUX" ]] && [[ -z "$ZELLIJ" ]] && [[ -z "$HERDR_ENV" ]] && [[ -t 1 ]]; then`)
				} else {
					newLines = append(newLines, line)
				}
			}
			if trimmed == "}" {
				inStartIfNeeded = false
			}
			continue
		}
		if trimmed == `WM_VAR="/$TMUX"` || trimmed == `WM_VAR="$ZELLIJ"` || trimmed == `WM_VAR="$HERDR_ENV"` {
			switch wm {
			case "tmux":
				newLines = append(newLines, `WM_VAR="/$TMUX"`)
			case "zellij":
				newLines = append(newLines, `WM_VAR="$ZELLIJ"`)
			case "herdr":
				newLines = append(newLines, `WM_VAR="$HERDR_ENV"`)
			}
			continue
		}
		if trimmed == `WM_CMD="tmux"` || trimmed == `WM_CMD="zellij"` || trimmed == `WM_CMD="herdr"` {
			switch wm {
			case "tmux":
				newLines = append(newLines, `WM_CMD="tmux"`)
			case "zellij":
				newLines = append(newLines, `WM_CMD="zellij"`)
			case "herdr":
				newLines = append(newLines, `WM_CMD="herdr"`)
			}
			continue
		}
		if trimmed == "# change with ZELLIJ" || trimmed == "# change with zellij" || trimmed == "# change with HERDR" || trimmed == "# change with herdr" {
			continue
		}
		if trimmed == "start_if_needed" && wm == "none" {
			continue
		}

		newLines = append(newLines, line)
	}

	return os.WriteFile(zshrcPath, []byte(strings.Join(newLines, "\n")), 0644)
}

// PatchFishForWM modifies config.fish based on window manager choice.
func PatchFishForWM(configPath string, wm string, installNvim bool) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	insertedMultiplexerBlock := false
	inMultiplexerBlock := false
	multiplexerBlockDepth := 0
	inCommentedMultiplexerBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.Contains(line, "fzf --fish | source") {
			if !installNvim {
				newLines = append(newLines, `if command -v fzf &> /dev/null`)
				newLines = append(newLines, line)
				newLines = append(newLines, `end`)
			} else {
				newLines = append(newLines, line)
			}
			continue
		}

		if inMultiplexerBlock {
			if strings.HasPrefix(trimmed, "if ") {
				multiplexerBlockDepth++
			}
			if trimmed == "end" {
				multiplexerBlockDepth--
				if multiplexerBlockDepth <= 0 {
					inMultiplexerBlock = false
				}
			}
			continue
		}

		if inCommentedMultiplexerBlock {
			if trimmed == "#end" {
				inCommentedMultiplexerBlock = false
			}
			continue
		}

		startsMultiplexerBlock := trimmed == "# Start tmux/zellij (tmux first, zellij as fallback)" ||
			trimmed == "# Start selected terminal multiplexer" ||
			trimmed == "if not set -q TMUX" ||
			strings.HasPrefix(trimmed, "if status is-interactive; and") && (strings.Contains(trimmed, "TMUX") || strings.Contains(trimmed, "ZELLIJ") || strings.Contains(trimmed, "HERDR_ENV"))

		if startsMultiplexerBlock {
			if !insertedMultiplexerBlock && wm != "none" {
				newLines = append(newLines, fishMultiplexerBlock(wm)...)
				insertedMultiplexerBlock = true
			}
			if strings.HasPrefix(trimmed, "if ") {
				inMultiplexerBlock = true
				multiplexerBlockDepth = 1
			}
			continue
		}

		if strings.HasPrefix(trimmed, "#if not set -q ZELLIJ") || strings.HasPrefix(trimmed, "#if not set -q HERDR_ENV") {
			inCommentedMultiplexerBlock = true
			continue
		}

		newLines = append(newLines, line)
	}

	if !insertedMultiplexerBlock && wm != "none" {
		newLines = append(newLines, fishMultiplexerBlock(wm)...)
	}

	return os.WriteFile(configPath, []byte(strings.Join(newLines, "\n")), 0644)
}

func fishMultiplexerBlock(wm string) []string {
	switch wm {
	case "zellij":
		return []string{
			"# Start selected terminal multiplexer",
			"if status is-interactive; and command -q zellij; and not set -q TMUX; and not set -q ZELLIJ; and not set -q HERDR_ENV",
			"    zellij attach -c main",
			"end",
		}
	case "herdr":
		return []string{
			"# Start selected terminal multiplexer",
			"if status is-interactive; and command -q herdr; and not set -q HERDR_ENV; and not set -q TMUX; and not set -q ZELLIJ",
			"    herdr; or echo \"⚠️  Herdr failed to start; continuing in Fish.\"",
			"end",
		}
	default:
		return []string{
			"# Start selected terminal multiplexer",
			"if status is-interactive; and command -q tmux; and not set -q TMUX; and not set -q ZELLIJ; and not set -q HERDR_ENV",
			"    tmux new-session -A -s main",
			"end",
		}
	}
}

// PatchNushellForWM modifies config.nu based on window manager choice.
func PatchNushellForWM(configPath string, wm string) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	inMultiplexerBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "let MULTIPLEXER =") {
			switch wm {
			case "tmux", "zellij", "herdr":
				newLines = append(newLines, fmt.Sprintf(`let MULTIPLEXER = "%s"`, wm))
			}
			continue
		}
		if strings.HasPrefix(trimmed, "let MULTIPLEXER_ENV_PREFIX =") {
			switch wm {
			case "tmux":
				newLines = append(newLines, `let MULTIPLEXER_ENV_PREFIX = "TMUX"`)
			case "zellij":
				newLines = append(newLines, `let MULTIPLEXER_ENV_PREFIX = "ZELLIJ"`)
			case "herdr":
				newLines = append(newLines, `let MULTIPLEXER_ENV_PREFIX = "HERDR_ENV"`)
			}
			continue
		}
		if strings.HasPrefix(trimmed, "def start_multiplexer") {
			inMultiplexerBlock = !strings.Contains(trimmed, "}")
			if wm != "none" {
				newLines = append(newLines, line)
			}
			continue
		}
		if inMultiplexerBlock {
			if wm != "none" {
				if strings.HasPrefix(trimmed, "if $MULTIPLEXER_ENV_PREFIX not-in") {
					newLines = append(newLines, `  let active_env = ($env | columns)`)
					newLines = append(newLines, `  if (which $MULTIPLEXER | is-not-empty) and $MULTIPLEXER_ENV_PREFIX not-in $active_env and "TMUX" not-in $active_env and "ZELLIJ" not-in $active_env and "HERDR_ENV" not-in $active_env {`)
				} else {
					newLines = append(newLines, line)
				}
			}
			if trimmed == "}" {
				inMultiplexerBlock = false
			}
			continue
		}
		if trimmed == "start_multiplexer" && wm == "none" {
			continue
		}

		newLines = append(newLines, line)
	}

	return os.WriteFile(configPath, []byte(strings.Join(newLines, "\n")), 0644)
}
