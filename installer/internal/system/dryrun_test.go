package system

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// captureNotify overrides Notify for the duration of a test and returns the
// captured messages, restoring the original on cleanup.
func captureNotify(t *testing.T) *[]string {
	t.Helper()
	var messages []string
	original := Notify
	Notify = func(msg string) { messages = append(messages, msg) }
	t.Cleanup(func() { Notify = original })
	return &messages
}

// enableDryRun activates the dry-run env var for the duration of a test.
func enableDryRun(t *testing.T) {
	t.Helper()
	t.Setenv("GENTLEMAN_DRY_RUN", "1")
}

// assertDryRunResult checks the full dry-run contract for one executed
// entry point: synthetic success, no side effect, and the planned
// operation printed on the given channel.
func assertDryRunResult(t *testing.T, name, wantCmd, sentinel string, result *ExecResult, channel []string) {
	t.Helper()
	if result.Error != nil {
		t.Fatalf("%s() under dry-run returned error: %v, want synthetic success", name, result.Error)
	}
	if result.ExitCode != 0 {
		t.Errorf("%s() under dry-run ExitCode = %d, want 0", name, result.ExitCode)
	}
	if result.Command != wantCmd {
		t.Errorf("%s() under dry-run Command = %q, want %q", name, result.Command, wantCmd)
	}
	if _, err := os.Stat(sentinel); !os.IsNotExist(err) {
		t.Fatalf("%s() under dry-run executed the command (sentinel %s exists)", name, sentinel)
	}
	want := "[dry-run] would run: " + wantCmd
	found := false
	for _, msg := range channel {
		if msg == want {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("%s() under dry-run did not print the plan %q, got: %v", name, want, channel)
	}
}

// --- IsDryRun env semantics (tasks 4.11 / 4.12) ---

func TestIsDryRun_DefaultsToFalse(t *testing.T) {
	t.Setenv("GENTLEMAN_DRY_RUN", "")
	if IsDryRun() {
		t.Fatal("IsDryRun() = true, want false when GENTLEMAN_DRY_RUN is absent")
	}
}

func TestIsDryRun_TrueOnlyForGentlemanDryRunOne(t *testing.T) {
	// --test sets GENTLEMAN_TEST_MODE (see setupTestMode in
	// cmd/gentleman-installer/main.go) and MUST NOT be mistaken for a dry
	// run: test mode still performs package-manager, sudo and network
	// operations.
	t.Setenv("GENTLEMAN_DRY_RUN", "")
	t.Setenv("GENTLEMAN_TEST_MODE", "1")
	if IsDryRun() {
		t.Fatal("IsDryRun() = true with only GENTLEMAN_TEST_MODE set; --test must not be a dry run")
	}

	t.Setenv("GENTLEMAN_DRY_RUN", "1")
	if !IsDryRun() {
		t.Fatal("IsDryRun() = false, want true when GENTLEMAN_DRY_RUN=1 (even in test mode)")
	}
}

func TestTestModeStillExecutes(t *testing.T) {
	// Behavioral half of 4.11: with --test env active but no dry-run env,
	// Run must really execute (this is what makes --test NOT a substitute
	// for --dry-run).
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("GENTLEMAN_DRY_RUN", "")
	t.Setenv("GENTLEMAN_TEST_MODE", "1")

	sentinel := filepath.Join(home, "test-mode-sentinel")
	result := Run("touch "+sentinel, nil)

	if result.Error != nil {
		t.Fatalf("Run() in test mode returned error: %v", result.Error)
	}
	if _, err := os.Stat(sentinel); err != nil {
		t.Fatalf("test mode must still execute: sentinel %s was not created: %v", sentinel, err)
	}
}

// --- Central gate on the two base exec functions (task 4.1) ---

func TestRun_DryRunExecutesNothing(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	enableDryRun(t)
	notified := captureNotify(t)

	sentinel := filepath.Join(home, "dry-run-sentinel")
	result := Run("touch "+sentinel, nil)
	assertDryRunResult(t, "Run", "touch "+sentinel, sentinel, result, *notified)
}

func TestRunWithLogs_DryRunExecutesNothing(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	enableDryRun(t)

	sentinel := filepath.Join(home, "dry-run-sentinel-logs")
	var logs []string
	logFunc := func(line string) { logs = append(logs, line) }

	result := RunWithLogs("touch "+sentinel, nil, logFunc)
	assertDryRunResult(t, "RunWithLogs", "touch "+sentinel, sentinel, result, logs)
}

func TestRunWithLogs_DryRunFallsBackToNotifyWhenNoCallback(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	enableDryRun(t)
	notified := captureNotify(t)

	sentinel := filepath.Join(home, "dry-run-sentinel-nocb")
	result := RunWithLogs("touch "+sentinel, nil, nil)
	assertDryRunResult(t, "RunWithLogs(nil cb)", "touch "+sentinel, sentinel, result, *notified)
}

func TestRun_DryRunGateSitsBeforeTermuxDirectExec(t *testing.T) {
	// Run's Termux branch calls parseCommand before spawning; the gate
	// must sit before it so a dry run never parses or spawns. On this
	// (non-Termux) host the proof is behavioral: a command whose
	// executable cannot exist still returns synthetic success.
	home := t.TempDir()
	t.Setenv("HOME", home)
	enableDryRun(t)
	captureNotify(t)

	sentinel := filepath.Join(home, "parse-gate-sentinel")
	result := Run("definitely-not-a-real-binary-xyz "+sentinel, nil)

	if result.Error != nil {
		t.Fatalf("Run() under dry-run must not even parse the command; got error: %v", result.Error)
	}
	if _, err := os.Stat(sentinel); !os.IsNotExist(err) {
		t.Fatalf("Run() under dry-run executed a nonexistent binary (sentinel %s exists)", sentinel)
	}
}

// --- All 7 wrappers by delegation (task 4.2) ---

// Each wrapper subtest runs a command whose execution would create a
// sentinel file. Under dry run the sentinel must not exist, the result
// must be a synthetic success, and the full composed command (including
// the "sudo "/"pkg "/"brew " prefix the wrapper adds) must be printed.
func TestRunWrappers_DryRunExecuteNothing(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	enableDryRun(t)
	notified := captureNotify(t)

	brewPath := GetBrewPrefix() + "/bin/brew"

	t.Run("RunSudo", func(t *testing.T) {
		sentinel := filepath.Join(home, "sentinel-runsudo")
		result := RunSudo("touch "+sentinel, nil)
		assertDryRunResult(t, "RunSudo", "sudo touch "+sentinel, sentinel, result, *notified)
	})
	t.Run("RunBrew", func(t *testing.T) {
		sentinel := filepath.Join(home, "sentinel-runbrew")
		result := RunBrew("touch "+sentinel, nil)
		assertDryRunResult(t, "RunBrew", brewPath+" touch "+sentinel, sentinel, result, *notified)
	})
	t.Run("RunPkg", func(t *testing.T) {
		sentinel := filepath.Join(home, "sentinel-runpkg")
		result := RunPkg("touch "+sentinel, nil)
		assertDryRunResult(t, "RunPkg", "pkg touch "+sentinel, sentinel, result, *notified)
	})
	t.Run("RunPkgWithLogs", func(t *testing.T) {
		sentinel := filepath.Join(home, "sentinel-runpkgwithlogs")
		var logs []string
		result := RunPkgWithLogs("touch "+sentinel, nil, func(line string) { logs = append(logs, line) })
		assertDryRunResult(t, "RunPkgWithLogs", "pkg touch "+sentinel, sentinel, result, logs)
	})
	t.Run("RunPkgInstall", func(t *testing.T) {
		sentinel := filepath.Join(home, "sentinel-runpkginstall")
		var logs []string
		result := RunPkgInstall(sentinel, nil, func(line string) { logs = append(logs, line) })
		assertDryRunResult(t, "RunPkgInstall", "pkg install -y "+sentinel, sentinel, result, logs)
	})
	t.Run("RunBrewWithLogs", func(t *testing.T) {
		sentinel := filepath.Join(home, "sentinel-runbrewwithlogs")
		var logs []string
		result := RunBrewWithLogs("touch "+sentinel, nil, func(line string) { logs = append(logs, line) })
		assertDryRunResult(t, "RunBrewWithLogs", brewPath+" touch "+sentinel, sentinel, result, logs)
	})
	t.Run("RunSudoWithLogs", func(t *testing.T) {
		sentinel := filepath.Join(home, "sentinel-runsudowithlogs")
		var logs []string
		result := RunSudoWithLogs("touch "+sentinel, nil, func(line string) { logs = append(logs, line) })
		assertDryRunResult(t, "RunSudoWithLogs", "sudo touch "+sentinel, sentinel, result, logs)
	})
}

// Control: without the dry-run env the same Run call really executes.
// This proves the sentinel assertions above are meaningful.
func TestRun_WithoutDryRunReallyExecutes(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("GENTLEMAN_DRY_RUN", "")

	sentinel := filepath.Join(home, "control-sentinel")
	result := Run("touch "+sentinel, nil)

	if result.Error != nil {
		t.Fatalf("Run() without dry-run returned error: %v", result.Error)
	}
	if _, err := os.Stat(sentinel); err != nil {
		t.Fatalf("control failed: sentinel %s was not created by a real execution: %v", sentinel, err)
	}
}

// --- Filesystem mutation paths (tasks 4.4 / 4.5) ---

func TestCreateBackup_DryRunCreatesNothingButPlans(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	enableDryRun(t)
	notified := captureNotify(t)

	// Real configs present so the plan has something to enumerate.
	if err := os.WriteFile(filepath.Join(home, ".zshrc"), []byte("export ZSH=1"), 0o644); err != nil {
		t.Fatalf("write .zshrc: %v", err)
	}
	nvimDir := filepath.Join(home, ".config", "nvim")
	if err := os.MkdirAll(nvimDir, 0o755); err != nil {
		t.Fatalf("mkdir nvim: %v", err)
	}
	if err := os.WriteFile(filepath.Join(nvimDir, "init.lua"), []byte("print('hi')"), 0o644); err != nil {
		t.Fatalf("write init.lua: %v", err)
	}

	backupDir, err := CreateBackup([]string{"zsh: " + filepath.Join(home, ".zshrc"), "nvim: " + nvimDir})
	if err != nil {
		t.Fatalf("CreateBackup() under dry-run returned error: %v, want nil", err)
	}
	if backupDir == "" {
		t.Fatal("CreateBackup() under dry-run returned empty backupDir, want the real GetBackupDir() path")
	}
	if !strings.HasPrefix(backupDir, home) || !strings.Contains(backupDir, ".gentleman-backup-") {
		t.Fatalf("CreateBackup() under dry-run backupDir = %q, want a real GetBackupDir() path under %s", backupDir, home)
	}
	if _, statErr := os.Stat(backupDir); !os.IsNotExist(statErr) {
		t.Fatalf("CreateBackup() under dry-run created %s on disk, stat err: %v", backupDir, statErr)
	}

	// The plan must name the paths that would be backed up.
	joined := strings.Join(*notified, "\n")
	for _, want := range []string{
		"[dry-run] would copy: " + filepath.Join(home, ".zshrc") + " -> " + backupDir + "/zsh",
		"[dry-run] would copy: " + nvimDir + " -> " + backupDir + "/nvim",
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("CreateBackup() plan missing %q, got: %s", want, joined)
		}
	}
}

func TestCopyFile_DryRunWritesNothingButPlans(t *testing.T) {
	src := filepath.Join(t.TempDir(), "source.txt")
	if err := os.WriteFile(src, []byte("hello"), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}
	dst := filepath.Join(t.TempDir(), "does-not-exist", "source.txt")
	enableDryRun(t)
	notified := captureNotify(t)

	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile() under dry-run returned error: %v, want nil", err)
	}
	if _, statErr := os.Stat(dst); !os.IsNotExist(statErr) {
		t.Fatalf("CopyFile() under dry-run wrote %s, stat err: %v", dst, statErr)
	}
	if _, statErr := os.Stat(filepath.Dir(dst)); !os.IsNotExist(statErr) {
		t.Fatalf("CopyFile() under dry-run created parent dir %s, stat err: %v", filepath.Dir(dst), statErr)
	}
	want := "[dry-run] would copy: " + src + " -> " + dst
	if !strings.Contains(strings.Join(*notified, "\n"), want) {
		t.Errorf("CopyFile() plan missing %q, got: %v", want, *notified)
	}
}

func TestCopyDir_DryRunWritesNothingButWalksAndPlans(t *testing.T) {
	src := filepath.Join(t.TempDir(), "src")
	if err := os.MkdirAll(filepath.Join(src, "nested"), 0o755); err != nil {
		t.Fatalf("mkdir src: %v", err)
	}
	if err := os.WriteFile(filepath.Join(src, "a.txt"), []byte("a"), 0o644); err != nil {
		t.Fatalf("write a.txt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(src, "nested", "b.txt"), []byte("b"), 0o644); err != nil {
		t.Fatalf("write b.txt: %v", err)
	}
	dst := filepath.Join(t.TempDir(), "dst")
	enableDryRun(t)
	notified := captureNotify(t)

	if err := CopyDir(src, dst); err != nil {
		t.Fatalf("CopyDir() under dry-run returned error: %v, want nil", err)
	}
	if _, statErr := os.Stat(dst); !os.IsNotExist(statErr) {
		t.Fatalf("CopyDir() under dry-run created %s on disk, stat err: %v", dst, statErr)
	}

	// CopyDir must still walk the source tree, so the plan is real:
	// one would-copy line per regular file, with the destination paths.
	joined := strings.Join(*notified, "\n")
	for _, want := range []string{
		"[dry-run] would copy: " + filepath.Join(src, "a.txt") + " -> " + filepath.Join(dst, "a.txt"),
		"[dry-run] would copy: " + filepath.Join(src, "nested", "b.txt") + " -> " + filepath.Join(dst, "nested", "b.txt"),
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("CopyDir() plan missing %q, got: %s", want, joined)
		}
	}
}
