package system

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

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

	// Use bash -c to run the command
	cmd := exec.CommandContext(ctx, "bash", "-c", command)

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

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}

// CopyDir recursively copies a directory
func CopyDir(src, dst string) error {
	result := Run(fmt.Sprintf("cp -r %s %s", src, dst), nil)
	return result.Error
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

	cmd := exec.CommandContext(ctx, "bash", "-c", command)

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
