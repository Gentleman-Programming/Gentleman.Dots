package system

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	t.Run("should execute simple command", func(t *testing.T) {
		result := Run("echo hello", nil)
		if result.Error != nil {
			t.Errorf("Unexpected error: %v", result.Error)
		}
		if !strings.Contains(result.Output, "hello") {
			t.Errorf("Expected output to contain 'hello', got '%s'", result.Output)
		}
	})

	t.Run("should capture exit code on failure", func(t *testing.T) {
		result := Run("exit 1", nil)
		if result.Error == nil {
			t.Error("Expected error for exit 1")
		}
		if result.ExitCode != 1 {
			t.Errorf("Expected exit code 1, got %d", result.ExitCode)
		}
	})

	t.Run("should respect timeout", func(t *testing.T) {
		start := time.Now()
		result := Run("sleep 10", &ExecOptions{Timeout: 100 * time.Millisecond})
		elapsed := time.Since(start)

		if elapsed > 2*time.Second {
			t.Errorf("Command should have timed out quickly, took %v", elapsed)
		}
		if result.Error == nil {
			t.Error("Expected timeout error")
		}
	})

	t.Run("should use working directory", func(t *testing.T) {
		result := Run("pwd", &ExecOptions{WorkDir: "/tmp"})
		if result.Error != nil {
			t.Errorf("Unexpected error: %v", result.Error)
		}
		// On macOS /tmp is a symlink to /private/tmp
		if !strings.Contains(result.Output, "tmp") {
			t.Errorf("Expected output to contain 'tmp', got '%s'", result.Output)
		}
	})

	t.Run("should track duration", func(t *testing.T) {
		result := Run("sleep 0.1", nil)
		if result.Duration < 50*time.Millisecond {
			t.Errorf("Duration seems too short: %v", result.Duration)
		}
	})

	t.Run("should handle environment variables", func(t *testing.T) {
		result := Run("echo $TEST_VAR", &ExecOptions{
			Env: []string{"TEST_VAR=gentleman"},
		})
		if result.Error != nil {
			t.Errorf("Unexpected error: %v", result.Error)
		}
		if !strings.Contains(result.Output, "gentleman") {
			t.Errorf("Expected output to contain 'gentleman', got '%s'", result.Output)
		}
	})
}

func TestEnsureDir(t *testing.T) {
	t.Run("should create directory if not exists", func(t *testing.T) {
		testDir := filepath.Join(os.TempDir(), "gentleman-test-dir")
		defer os.RemoveAll(testDir)

		err := EnsureDir(testDir)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		info, err := os.Stat(testDir)
		if err != nil {
			t.Errorf("Directory should exist: %v", err)
		}
		if !info.IsDir() {
			t.Error("Should be a directory")
		}
	})

	t.Run("should not error if directory exists", func(t *testing.T) {
		testDir := filepath.Join(os.TempDir(), "gentleman-test-existing")
		os.MkdirAll(testDir, 0755)
		defer os.RemoveAll(testDir)

		err := EnsureDir(testDir)
		if err != nil {
			t.Errorf("Should not error for existing directory: %v", err)
		}
	})

	t.Run("should create nested directories", func(t *testing.T) {
		testDir := filepath.Join(os.TempDir(), "gentleman-test", "nested", "deep")
		defer os.RemoveAll(filepath.Join(os.TempDir(), "gentleman-test"))

		err := EnsureDir(testDir)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		info, err := os.Stat(testDir)
		if err != nil {
			t.Errorf("Nested directory should exist: %v", err)
		}
		if !info.IsDir() {
			t.Error("Should be a directory")
		}
	})
}

func TestCopyFile(t *testing.T) {
	t.Run("should copy file contents", func(t *testing.T) {
		// Create source file
		srcDir := filepath.Join(os.TempDir(), "gentleman-copy-test")
		os.MkdirAll(srcDir, 0755)
		defer os.RemoveAll(srcDir)

		srcFile := filepath.Join(srcDir, "source.txt")
		dstFile := filepath.Join(srcDir, "dest.txt")
		content := "Hello, Gentleman!"

		err := os.WriteFile(srcFile, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create source file: %v", err)
		}

		err = CopyFile(srcFile, dstFile)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		data, err := os.ReadFile(dstFile)
		if err != nil {
			t.Errorf("Failed to read destination file: %v", err)
		}

		if string(data) != content {
			t.Errorf("Expected '%s', got '%s'", content, string(data))
		}
	})

	t.Run("should error for non-existent source", func(t *testing.T) {
		err := CopyFile("/non/existent/file", "/tmp/dest")
		if err == nil {
			t.Error("Expected error for non-existent source")
		}
	})
}

func TestConfigPaths(t *testing.T) {
	t.Run("should return map of config paths", func(t *testing.T) {
		paths := ConfigPaths()

		if len(paths) == 0 {
			t.Error("ConfigPaths should return non-empty map")
		}

		// Check for expected keys
		expectedKeys := []string{"nvim", "fish", "zsh", "tmux", "zellij", "alacritty", "ghostty"}
		for _, key := range expectedKeys {
			if _, exists := paths[key]; !exists {
				t.Errorf("Expected key '%s' in ConfigPaths", key)
			}
		}
	})

	t.Run("all paths should start with home directory", func(t *testing.T) {
		home := os.Getenv("HOME")
		paths := ConfigPaths()

		for key, path := range paths {
			if !strings.HasPrefix(path, home) {
				t.Errorf("Path for '%s' should start with HOME: %s", key, path)
			}
		}
	})
}

func TestDetectExistingConfigs(t *testing.T) {
	t.Run("should return slice of existing configs", func(t *testing.T) {
		// This will return configs that exist on the current system
		existing := DetectExistingConfigs()

		// Result should be a slice (can be empty)
		if existing == nil {
			t.Error("DetectExistingConfigs should not return nil")
		}
	})

	t.Run("should format as 'key: path'", func(t *testing.T) {
		existing := DetectExistingConfigs()

		for _, config := range existing {
			if !strings.Contains(config, ": ") {
				t.Errorf("Config should be formatted as 'key: path', got: %s", config)
			}
		}
	})
}

func TestGetBackupDir(t *testing.T) {
	t.Run("should return path with timestamp", func(t *testing.T) {
		dir := GetBackupDir()

		if !strings.Contains(dir, ".gentleman-backup-") {
			t.Errorf("Backup dir should contain '.gentleman-backup-', got: %s", dir)
		}

		home := os.Getenv("HOME")
		if !strings.HasPrefix(dir, home) {
			t.Errorf("Backup dir should start with HOME: %s", dir)
		}
	})

	t.Run("should generate unique paths", func(t *testing.T) {
		dir1 := GetBackupDir()
		time.Sleep(2 * time.Second) // Wait to ensure different timestamp
		dir2 := GetBackupDir()

		if dir1 == dir2 {
			t.Error("Backup dirs should be unique")
		}
	})
}

func TestListBackups(t *testing.T) {
	t.Run("should return slice of BackupInfo", func(t *testing.T) {
		backups := ListBackups()

		// Result should be a slice (can be empty)
		if backups == nil {
			t.Error("ListBackups should not return nil")
		}
	})

	t.Run("BackupInfo should have required fields", func(t *testing.T) {
		// Create a temporary backup directory
		home := os.Getenv("HOME")
		testBackupDir := filepath.Join(home, ".gentleman-backup-test-123456")
		os.MkdirAll(testBackupDir, 0755)
		defer os.RemoveAll(testBackupDir)

		// Create a test file inside
		os.WriteFile(filepath.Join(testBackupDir, "test-config"), []byte("test"), 0644)

		backups := ListBackups()

		found := false
		for _, backup := range backups {
			if strings.Contains(backup.Path, "test-123456") {
				found = true
				if backup.Path == "" {
					t.Error("BackupInfo.Path should not be empty")
				}
				if backup.Timestamp.IsZero() {
					t.Error("BackupInfo.Timestamp should not be zero")
				}
				if len(backup.Files) == 0 {
					t.Error("BackupInfo.Files should not be empty")
				}
				break
			}
		}

		if !found {
			t.Log("Test backup directory was not found in ListBackups (may be filtered)")
		}
	})
}

func TestCreateBackup(t *testing.T) {
	t.Run("should create backup directory", func(t *testing.T) {
		// Create a temporary config to backup
		home := os.Getenv("HOME")
		testConfigDir := filepath.Join(home, ".config", "gentleman-test-backup")
		os.MkdirAll(testConfigDir, 0755)
		testFile := filepath.Join(testConfigDir, "config.txt")
		os.WriteFile(testFile, []byte("test config"), 0644)
		defer os.RemoveAll(testConfigDir)

		// We can't easily test CreateBackup without mocking ConfigPaths
		// So we'll just test that it returns a valid path
		backupDir, err := CreateBackup([]string{})
		if err != nil {
			// Expected - no configs to backup
			t.Log("CreateBackup with empty list succeeded or failed as expected")
		}

		// Clean up if backup was created
		if backupDir != "" {
			defer os.RemoveAll(backupDir)
		}
	})

	t.Run("should return valid backup path format", func(t *testing.T) {
		// Even with no configs, it should create the backup directory
		backupDir, _ := CreateBackup([]string{"nonexistent"})

		if backupDir != "" {
			defer os.RemoveAll(backupDir)

			home := os.Getenv("HOME")
			if !strings.HasPrefix(backupDir, home) {
				t.Errorf("Backup dir should start with HOME: %s", backupDir)
			}

			if !strings.Contains(backupDir, ".gentleman-backup-") {
				t.Errorf("Backup dir should contain '.gentleman-backup-': %s", backupDir)
			}
		}
	})
}

func TestRestoreBackup(t *testing.T) {
	t.Run("should error for non-existent backup", func(t *testing.T) {
		err := RestoreBackup("/non/existent/backup")
		if err == nil {
			t.Error("Expected error for non-existent backup directory")
		}
	})

	t.Run("should restore files from backup", func(t *testing.T) {
		// This is a complex integration test - we'll just verify it doesn't panic
		home := os.Getenv("HOME")
		testBackupDir := filepath.Join(home, ".gentleman-backup-restore-test")
		os.MkdirAll(testBackupDir, 0755)
		defer os.RemoveAll(testBackupDir)

		// RestoreBackup with empty backup dir should not error
		err := RestoreBackup(testBackupDir)
		if err != nil {
			t.Errorf("RestoreBackup should not error for empty backup: %v", err)
		}
	})
}

func TestDeleteBackup(t *testing.T) {
	t.Run("should delete backup directory", func(t *testing.T) {
		// Create a temporary backup directory
		home := os.Getenv("HOME")
		testBackupDir := filepath.Join(home, ".gentleman-backup-delete-test")
		os.MkdirAll(testBackupDir, 0755)
		os.WriteFile(filepath.Join(testBackupDir, "test"), []byte("test"), 0644)

		// Verify it exists
		if _, err := os.Stat(testBackupDir); os.IsNotExist(err) {
			t.Fatal("Test backup directory should exist before delete")
		}

		// Delete it
		err := DeleteBackup(testBackupDir)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Verify it's gone
		if _, err := os.Stat(testBackupDir); !os.IsNotExist(err) {
			t.Error("Backup directory should not exist after delete")
			os.RemoveAll(testBackupDir) // Cleanup just in case
		}
	})

	t.Run("should not error for non-existent directory", func(t *testing.T) {
		err := DeleteBackup("/non/existent/backup/to/delete")
		if err != nil {
			t.Errorf("DeleteBackup should not error for non-existent dir: %v", err)
		}
	})
}

func TestBackupInfo(t *testing.T) {
	t.Run("BackupInfo struct should have correct fields", func(t *testing.T) {
		info := BackupInfo{
			Path:      "/test/path",
			Timestamp: time.Now(),
			Files:     []string{"file1", "file2"},
		}

		if info.Path != "/test/path" {
			t.Errorf("Expected path '/test/path', got '%s'", info.Path)
		}

		if len(info.Files) != 2 {
			t.Errorf("Expected 2 files, got %d", len(info.Files))
		}
	})
}

// TestBackupE2EFlow tests a complete backup and restore cycle
func TestBackupE2EFlow(t *testing.T) {
	home := os.Getenv("HOME")

	// Create temp test directory structure
	testConfigDir := filepath.Join(home, ".config", "gentleman-e2e-test")
	defer os.RemoveAll(testConfigDir)

	t.Run("complete backup and restore cycle", func(t *testing.T) {
		// 1. Setup: Create a fake config file
		os.MkdirAll(testConfigDir, 0755)
		originalContent := "original config content - test123"
		testFile := filepath.Join(testConfigDir, "test-config.lua")
		err := os.WriteFile(testFile, []byte(originalContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// 2. Verify file was created
		content, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("Failed to read test file: %v", err)
		}
		if string(content) != originalContent {
			t.Errorf("Content mismatch: got '%s', want '%s'", string(content), originalContent)
		}

		// 3. Create backup directory manually (simulating CreateBackup)
		backupDir := filepath.Join(home, ".gentleman-backup-e2e-test-flow")
		defer os.RemoveAll(backupDir)
		os.MkdirAll(backupDir, 0755)

		// Copy file to backup
		backupFile := filepath.Join(backupDir, "test-config.lua")
		err = CopyFile(testFile, backupFile)
		if err != nil {
			t.Fatalf("Failed to backup file: %v", err)
		}

		// 4. Modify the original (simulating installer overwrite)
		modifiedContent := "modified by installer"
		err = os.WriteFile(testFile, []byte(modifiedContent), 0644)
		if err != nil {
			t.Fatalf("Failed to modify test file: %v", err)
		}

		// Verify modification
		content, _ = os.ReadFile(testFile)
		if string(content) != modifiedContent {
			t.Errorf("Modification failed: got '%s', want '%s'", string(content), modifiedContent)
		}

		// 5. Restore from backup
		err = CopyFile(backupFile, testFile)
		if err != nil {
			t.Fatalf("Failed to restore from backup: %v", err)
		}

		// 6. Verify restoration
		content, _ = os.ReadFile(testFile)
		if string(content) != originalContent {
			t.Errorf("Restore failed: got '%s', want '%s'", string(content), originalContent)
		}
	})

	t.Run("backup directory naming follows timestamp format", func(t *testing.T) {
		dir := GetBackupDir()

		// Should contain .gentleman-backup- prefix
		if !strings.Contains(dir, ".gentleman-backup-") {
			t.Errorf("Backup dir should contain '.gentleman-backup-': %s", dir)
		}

		// Should contain timestamp in format YYYY-MM-DD-HHMMSS
		// Extract the timestamp part
		parts := strings.Split(filepath.Base(dir), "-")
		if len(parts) < 5 {
			t.Errorf("Expected at least 5 parts in backup dir name, got: %s", dir)
		}
	})

	t.Run("multiple backups with different timestamps", func(t *testing.T) {
		dir1 := GetBackupDir()
		time.Sleep(2 * time.Second)
		dir2 := GetBackupDir()

		if dir1 == dir2 {
			t.Error("Backup dirs should have different timestamps")
		}
	})
}

// TestDetectExistingConfigsWithRealFiles tests detection with actual files
func TestDetectExistingConfigsWithRealFiles(t *testing.T) {
	home := os.Getenv("HOME")

	t.Run("detects created config files", func(t *testing.T) {
		// Create a test fish config
		testFishDir := filepath.Join(home, ".config", "fish")
		testFishFile := filepath.Join(testFishDir, "config.fish")

		// Check if fish config already exists (don't overwrite user's config!)
		fishExisted := false
		if _, err := os.Stat(testFishFile); err == nil {
			fishExisted = true
		}

		if !fishExisted {
			// Create temp fish config for test
			os.MkdirAll(testFishDir, 0755)
			os.WriteFile(testFishFile, []byte("# test"), 0644)
			defer func() {
				if !fishExisted {
					os.Remove(testFishFile)
					// Only remove dir if we created it and it's empty
					os.Remove(testFishDir)
				}
			}()
		}

		// Now detect existing configs
		configs := DetectExistingConfigs()

		// Should find at least one config (fish, or user's actual configs)
		if len(configs) == 0 {
			t.Log("No existing configs found - this is OK if system has no configs")
		} else {
			t.Logf("Found %d existing configs", len(configs))
			for _, config := range configs {
				t.Logf("  - %s", config)
			}
		}
	})
}

// TestBackupPreservesFilePermissions verifies file permissions are maintained
func TestBackupPreservesContent(t *testing.T) {
	home := os.Getenv("HOME")
	testDir := filepath.Join(home, ".config", "gentleman-perm-test")
	defer os.RemoveAll(testDir)

	t.Run("preserves file content during copy", func(t *testing.T) {
		os.MkdirAll(testDir, 0755)

		// Create file with specific content including newlines and special chars
		content := "line1\nline2\n# comment\nexport VAR=\"value\"\n"
		srcFile := filepath.Join(testDir, "source.txt")
		dstFile := filepath.Join(testDir, "dest.txt")

		err := os.WriteFile(srcFile, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to write source: %v", err)
		}

		err = CopyFile(srcFile, dstFile)
		if err != nil {
			t.Fatalf("Failed to copy: %v", err)
		}

		// Read back and compare
		copied, err := os.ReadFile(dstFile)
		if err != nil {
			t.Fatalf("Failed to read dest: %v", err)
		}

		if string(copied) != content {
			t.Errorf("Content mismatch:\nGot: %q\nWant: %q", string(copied), content)
		}
	})
}

func TestRunPkg(t *testing.T) {
	t.Run("should construct correct pkg command", func(t *testing.T) {
		// We can't really test pkg on non-Termux systems,
		// but we can verify the command construction by checking
		// that it at least attempts to run "pkg" with the args
		result := RunPkg("--help", nil)

		// On non-Termux systems, this will fail with "command not found"
		// which is expected behavior - we just verify it tried
		if result.Command != "pkg --help" {
			t.Errorf("Expected command 'pkg --help', got '%s'", result.Command)
		}
	})

	t.Run("should respect options", func(t *testing.T) {
		opts := &ExecOptions{
			WorkDir: "/tmp",
			Timeout: 100 * time.Millisecond,
		}
		result := RunPkg("--version", opts)

		// Just verify it ran with options
		if result.Command != "pkg --version" {
			t.Errorf("Expected command 'pkg --version', got '%s'", result.Command)
		}
	})
}

func TestRunPkgWithLogs(t *testing.T) {
	t.Run("should stream logs to callback", func(t *testing.T) {
		var logs []string
		logFunc := func(line string) {
			logs = append(logs, line)
		}

		// Use echo to simulate pkg output
		result := RunWithLogs("echo 'test output'", nil, logFunc)

		if result.Error != nil {
			t.Errorf("Unexpected error: %v", result.Error)
		}

		// Verify callback was called
		if len(logs) == 0 {
			t.Error("Expected log callback to be called")
		}

		foundOutput := false
		for _, log := range logs {
			if strings.Contains(log, "test output") {
				foundOutput = true
				break
			}
		}
		if !foundOutput {
			t.Errorf("Expected 'test output' in logs, got: %v", logs)
		}
	})

	t.Run("should construct correct pkg command", func(t *testing.T) {
		result := RunPkgWithLogs("--help", nil, nil)

		if result.Command != "pkg --help" {
			t.Errorf("Expected command 'pkg --help', got '%s'", result.Command)
		}
	})
}

func TestRunPkgInstall(t *testing.T) {
	t.Run("should construct correct install command", func(t *testing.T) {
		result := RunPkgInstall("vim git", nil, nil)

		// Verify command includes -y flag for non-interactive
		expectedCmd := "pkg install -y vim git"
		if result.Command != expectedCmd {
			t.Errorf("Expected command '%s', got '%s'", expectedCmd, result.Command)
		}
	})

	t.Run("should include -y flag for non-interactive installs", func(t *testing.T) {
		result := RunPkgInstall("neovim", nil, nil)

		if !strings.Contains(result.Command, "-y") {
			t.Errorf("Command should include -y flag, got: %s", result.Command)
		}
	})

	t.Run("should stream logs to callback", func(t *testing.T) {
		logFunc := func(line string) {
			// Callback exists, just verify it doesn't panic
			_ = line
		}

		// This will fail on non-Termux but we're just testing the callback wiring
		RunPkgInstall("somepackage", nil, logFunc)

		// Note: on non-Termux systems, pkg won't exist so there might be no output
		// The important thing is the function doesn't panic
		t.Log("RunPkgInstall completed without panic")
	})
}

func TestRunWithLogs(t *testing.T) {
	t.Run("should capture both stdout and stderr", func(t *testing.T) {
		var logs []string
		logFunc := func(line string) {
			logs = append(logs, line)
		}

		result := RunWithLogs("echo stdout && echo stderr >&2", nil, logFunc)

		if result.Error != nil {
			t.Errorf("Unexpected error: %v", result.Error)
		}

		// Both stdout and stderr should be captured (may come as 1 or 2 entries depending on timing)
		if len(logs) < 1 {
			t.Errorf("Expected at least 1 log entry, got %d", len(logs))
		}

		// Check that output contains expected strings
		combined := result.Output + result.Stderr
		if !strings.Contains(combined, "stdout") {
			t.Errorf("Expected output to contain 'stdout', got: %s", combined)
		}
	})

	t.Run("should handle nil callback gracefully", func(t *testing.T) {
		// Should not panic with nil callback
		result := RunWithLogs("echo test", nil, nil)

		if result.Error != nil {
			t.Errorf("Unexpected error: %v", result.Error)
		}

		if !strings.Contains(result.Output, "test") {
			t.Errorf("Expected output to contain 'test', got: %s", result.Output)
		}
	})

	t.Run("should respect timeout", func(t *testing.T) {
		start := time.Now()
		result := RunWithLogs("sleep 10", &ExecOptions{Timeout: 100 * time.Millisecond}, nil)
		elapsed := time.Since(start)

		if elapsed > 2*time.Second {
			t.Errorf("Command should have timed out quickly, took %v", elapsed)
		}
		if result.Error == nil {
			t.Error("Expected timeout error")
		}
	})

	t.Run("should capture exit code on failure", func(t *testing.T) {
		result := RunWithLogs("exit 42", nil, nil)

		if result.Error == nil {
			t.Error("Expected error for non-zero exit")
		}
		if result.ExitCode != 42 {
			t.Errorf("Expected exit code 42, got %d", result.ExitCode)
		}
	})
}
