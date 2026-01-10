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
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}

// CopyDir recursively copies a directory using native Go (shell-independent)
func CopyDir(src, dst string) error {
	// Clean paths - remove trailing /* or /. if present
	src = strings.TrimSuffix(strings.TrimSuffix(src, "/*"), "/.")

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
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

// PatchZshForWM modifies .zshrc based on window manager choice
// If wm is "none", removes WM-related lines
// If wm is "zellij", changes tmux references to zellij
// If wm is "tmux", leaves as-is (default)
func PatchZshForWM(zshrcPath string, wm string, installNvim bool) error {
	content, err := os.ReadFile(zshrcPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string

	// Lines to remove if WM is "none"
	wmLinesToRemove := map[string]bool{
		`WM_VAR="/$TMUX"`:      true,
		`WM_VAR="$ZELLIJ"`:     true,
		`WM_CMD="tmux"`:        true,
		`WM_CMD="zellij"`:      true,
		`# change with ZELLIJ`: true,
		`# change with zellij`: true,
		`start_if_needed`:      true,
	}

	inStartIfNeeded := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Handle fzf line - wrap with command check if nvim not installed
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

		// Handle WM based on choice
		if wm == "none" {
			// Skip WM-related lines entirely
			if wmLinesToRemove[trimmed] {
				continue
			}
			// Skip function start_if_needed block
			if strings.HasPrefix(trimmed, "function start_if_needed") {
				inStartIfNeeded = true
				continue
			}
			if inStartIfNeeded {
				if trimmed == "}" {
					inStartIfNeeded = false
				}
				continue
			}
			// Skip the call to start_if_needed
			if trimmed == "start_if_needed" {
				continue
			}
			newLines = append(newLines, line)
		} else if wm == "zellij" {
			// Replace tmux with zellij
			modified := line
			if strings.Contains(line, `WM_VAR="/$TMUX"`) {
				modified = `WM_VAR="$ZELLIJ"`
			} else if strings.Contains(line, `WM_CMD="tmux"`) {
				modified = `WM_CMD="zellij"`
			} else if strings.Contains(line, "# change with ZELLIJ") {
				continue // Remove this comment
			} else if strings.Contains(line, "# change with zellij") {
				continue // Remove this comment
			}
			newLines = append(newLines, modified)
		} else {
			// tmux is default, keep as-is
			newLines = append(newLines, line)
		}
	}

	return os.WriteFile(zshrcPath, []byte(strings.Join(newLines, "\n")), 0644)
}

// PatchFishForWM modifies config.fish based on window manager choice
func PatchFishForWM(configPath string, wm string, installNvim bool) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string

	inTmuxBlock := false
	inZellijBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Handle fzf line
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

		if wm == "none" {
			// Skip tmux block
			if trimmed == "if not set -q TMUX" {
				inTmuxBlock = true
				continue
			}
			if inTmuxBlock {
				if trimmed == "end" {
					inTmuxBlock = false
				}
				continue
			}
			// Skip zellij commented block
			if strings.HasPrefix(trimmed, "#if not set -q ZELLIJ") || strings.HasPrefix(trimmed, "#  zellij") || strings.HasPrefix(trimmed, "#end") {
				continue
			}
			newLines = append(newLines, line)
		} else if wm == "zellij" {
			// Comment out tmux, uncomment zellij
			if trimmed == "if not set -q TMUX" {
				inTmuxBlock = true
				continue
			}
			if inTmuxBlock {
				if trimmed == "end" {
					inTmuxBlock = false
				}
				continue
			}
			// Uncomment zellij block
			if trimmed == "#if not set -q ZELLIJ" {
				newLines = append(newLines, "if not set -q ZELLIJ")
				inZellijBlock = true
				continue
			}
			if inZellijBlock && strings.HasPrefix(trimmed, "#") {
				if trimmed == "#end" {
					newLines = append(newLines, "end")
					inZellijBlock = false
				} else {
					// Remove leading # and space
					newLines = append(newLines, strings.TrimPrefix(trimmed, "#"))
				}
				continue
			}
			newLines = append(newLines, line)
		} else {
			// tmux - keep as is
			newLines = append(newLines, line)
		}
	}

	return os.WriteFile(configPath, []byte(strings.Join(newLines, "\n")), 0644)
}

// PatchNushellForWM modifies config.nu based on window manager choice
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

		if wm == "none" {
			// Skip multiplexer-related lines
			if strings.HasPrefix(trimmed, "let MULTIPLEXER =") ||
				strings.HasPrefix(trimmed, "let MULTIPLEXER_ENV_PREFIX =") {
				continue
			}
			if strings.HasPrefix(trimmed, "def start_multiplexer") {
				inMultiplexerBlock = true
				continue
			}
			if inMultiplexerBlock {
				if trimmed == "}" {
					inMultiplexerBlock = false
				}
				continue
			}
			if trimmed == "start_multiplexer" {
				continue
			}
			newLines = append(newLines, line)
		} else if wm == "zellij" {
			// Replace tmux with zellij
			if strings.Contains(line, `let MULTIPLEXER = "tmux"`) {
				newLines = append(newLines, `let MULTIPLEXER = "zellij"`)
			} else if strings.Contains(line, `let MULTIPLEXER_ENV_PREFIX = "TMUX"`) {
				newLines = append(newLines, `let MULTIPLEXER_ENV_PREFIX = "ZELLIJ"`)
			} else {
				newLines = append(newLines, line)
			}
		} else {
			// tmux - keep as is
			newLines = append(newLines, line)
		}
	}

	return os.WriteFile(configPath, []byte(strings.Join(newLines, "\n")), 0644)
}
