package tui

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
)

// stubGitOnPATH replaces git on PATH with a stub that intercepts only
// `git clone`: clone invocations are recorded in logPath and fail,
// while every other git subcommand (notably `status --porcelain`, which
// InspectRepoDir relies on) delegates to the real git binary.
// stepCloneRepo clones through GetShell() -c "git clone ...", so the
// shell (resolved independently) finds only the stub: any real clone
// attempt is recorded and fails.
//
// This keeps the tests hermetic (no network dependency) and lets them
// assert zero clone spawn under dry run: after the fix the stub log
// must be empty. Fixtures must be created BEFORE calling this helper so
// the real git still works for `git init`.
func stubGitOnPATH(t *testing.T, logPath string) {
	t.Helper()
	realGit, err := exec.LookPath("git")
	if err != nil {
		t.Skipf("git not available: %v", err)
	}
	binDir := filepath.Join(t.TempDir(), "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(binDir): %v", err)
	}
	stub := `#!/bin/sh
for a in "$@"; do
  if [ "$a" = "clone" ]; then
    echo "git $@" >> ` + logPath + `
    exit 1
  fi
done
exec ` + realGit + ` "$@"
`
	if err := os.WriteFile(filepath.Join(binDir, "git"), []byte(stub), 0o755); err != nil {
		t.Fatalf("WriteFile(git stub): %v", err)
	}
	// The shell itself must still be resolvable; only git is stubbed.
	if err := os.Symlink("/bin/sh", filepath.Join(binDir, "sh")); err != nil {
		t.Fatalf("Symlink(sh): %v", err)
	}
	t.Setenv("PATH", binDir)
}

// stubLogContains reports whether the stub git recorded any invocation.
func stubLogContains(t *testing.T, logPath string) string {
	t.Helper()
	data, err := os.ReadFile(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ""
		}
		t.Fatalf("ReadFile(stub log): %v", err)
	}
	return string(data)
}

// --- stepCloneRepo under dry run (tasks 4.7 / 4.8) ---

// Regression class from the governing rule "dry-run suppresses effects,
// never derivations": when the clone is suppressed the Gentleman.Dots
// directory legitimately does not exist, so the post-clone os.Stat
// verification in stepCloneRepo must be gated as well. Without that gate
// a dry run reports a false failure.
func TestStepCloneRepo_DryRunDoesNotReportFalseFailure(t *testing.T) {
	base := t.TempDir() // no Gentleman.Dots directory exists
	repoDir := filepath.Join(base, "Gentleman.Dots")
	chdirTo(t, base)

	stubLog := filepath.Join(base, "stub-git-invocations.log")
	stubGitOnPATH(t, stubLog)
	t.Setenv("GENTLEMAN_DRY_RUN", "1")

	err := stepCloneRepo(&Model{})

	if err != nil {
		t.Fatalf("stepCloneRepo() under dry-run = %v, want nil (clone suppressed, no false failure)", err)
	}
	if invoked := stubLogContains(t, stubLog); invoked != "" {
		t.Fatalf("stepCloneRepo() under dry-run attempted a real git invocation:\n%s", invoked)
	}
	if _, statErr := os.Stat(repoDir); !os.IsNotExist(statErr) {
		t.Fatalf("stepCloneRepo() under dry-run created Gentleman.Dots, stat err: %v", statErr)
	}
}

// When an existing clean checkout is present, stepCloneRepo's delete site
// must print the plan and keep the directory, then the (suppressed) clone
// succeeds and the step returns nil.
func TestStepCloneRepo_DryRunKeepsExistingCleanCheckout(t *testing.T) {
	base := t.TempDir()
	repoDir := filepath.Join(base, "Gentleman.Dots")
	cleanGitFixture(t, repoDir) // must run BEFORE stubGitOnPATH
	chdirTo(t, base)

	stubLog := filepath.Join(base, "stub-git-invocations.log")
	stubGitOnPATH(t, stubLog)
	t.Setenv("GENTLEMAN_DRY_RUN", "1")

	err := stepCloneRepo(&Model{})

	if err != nil {
		t.Fatalf("stepCloneRepo() under dry-run = %v, want nil", err)
	}
	if invoked := stubLogContains(t, stubLog); invoked != "" {
		t.Fatalf("stepCloneRepo() under dry-run attempted a real git invocation:\n%s", invoked)
	}
	// The user's checkout must survive byte-for-byte, not be replaced by
	// a real clone.
	data, readErr := os.ReadFile(filepath.Join(repoDir, "file.txt"))
	if readErr != nil {
		t.Fatalf("stepCloneRepo() under dry-run deleted the existing clean checkout: %v", readErr)
	}
	if string(data) != "hello" {
		t.Fatalf("existing checkout content changed: got %q, want %q", string(data), "hello")
	}
}

// --- removeRepoDirIfSafe under dry run (tasks 4.9 / 4.10) ---

// Unit 3 replaced the shell `rm -rf` with os.RemoveAll inside
// removeRepoDirIfSafe, so the Run/RunWithLogs gate does NOT cover that
// deletion. A dry run must print the planned removal instead of deleting
// the user's clean checkout.
func TestRemoveRepoDirIfSafe_DryRunDoesNotDeleteAndPrintsPlan(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "repo")
	cleanGitFixture(t, dir)

	t.Setenv("GENTLEMAN_DRY_RUN", "1")

	var notified []string
	originalNotify := system.Notify
	system.Notify = func(msg string) { notified = append(notified, msg) }
	t.Cleanup(func() { system.Notify = originalNotify })

	err := removeRepoDirIfSafe("clone", dir)

	if err != nil {
		t.Fatalf("removeRepoDirIfSafe() under dry-run = %v, want nil", err)
	}
	if _, statErr := os.Stat(dir); statErr != nil {
		t.Fatalf("removeRepoDirIfSafe() under dry-run deleted the clean checkout %s, stat err: %v", dir, statErr)
	}

	wantAbs, absErr := filepath.Abs(dir)
	if absErr != nil {
		t.Fatalf("filepath.Abs: %v", absErr)
	}
	joined := strings.Join(notified, "\n")
	if !strings.Contains(joined, "[dry-run] would remove:") || !strings.Contains(joined, wantAbs) {
		t.Errorf("removeRepoDirIfSafe() under dry-run did not print the planned removal naming %s, got: %v", wantAbs, notified)
	}
}

// --- stepCleanup under dry run (task 4.9) ---

func TestStepCleanup_DryRunDoesNotDelete(t *testing.T) {
	base := t.TempDir()
	repoDir := filepath.Join(base, "Gentleman.Dots")
	cleanGitFixture(t, repoDir)
	chdirTo(t, base)
	t.Setenv("GENTLEMAN_DRY_RUN", "1")

	err := stepCleanup(&Model{})

	if err != nil {
		t.Fatalf("stepCleanup() under dry-run = %v, want nil", err)
	}
	if _, statErr := os.Stat(repoDir); statErr != nil {
		t.Fatalf("stepCleanup() under dry-run deleted %s, stat err: %v", repoDir, statErr)
	}
}
