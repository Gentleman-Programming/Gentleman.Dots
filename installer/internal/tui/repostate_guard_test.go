package tui

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// initGitFixture runs `git init` plus committer identity config in dir so
// fixtures can commit without relying on any global/user git config.
func initGitFixture(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll(%s): %v", dir, err)
	}
	runGitCmd(t, dir, "init")
	runGitCmd(t, dir, "config", "user.email", "test@example.com")
	runGitCmd(t, dir, "config", "user.name", "Test User")
}

func runGitCmd(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}

// cleanGitFixture builds a checkout that is genuinely safe to delete: the
// working tree is clean AND every commit is published to a remote. A clean
// working tree alone is not enough - a repository with no remote holds
// commits that exist nowhere else, which the guard must refuse to remove.
func cleanGitFixture(t *testing.T, dir string) {
	t.Helper()
	origin := dir + ".origin.git"
	runGitCmd(t, filepath.Dir(dir), "init", "--bare", "-b", "main", origin)

	initGitFixture(t, dir)
	filePath := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(filePath, []byte("hello"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	runGitCmd(t, dir, "add", "file.txt")
	runGitCmd(t, dir, "commit", "-m", "initial commit")
	runGitCmd(t, dir, "remote", "add", "origin", origin)
	runGitCmd(t, dir, "push", "-u", "origin", "HEAD:main")
}

func dirtyGitFixture(t *testing.T, dir string) {
	t.Helper()
	cleanGitFixture(t, dir)
	untracked := filepath.Join(dir, "untracked.txt")
	if err := os.WriteFile(untracked, []byte("new"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}

// chdirTo changes the process working directory to dir and restores the
// original on cleanup. stepCloneRepo/stepCleanup operate on the relative
// "Gentleman.Dots" path, matching production behavior, so tests must
// exercise them from a controlled cwd.
func chdirTo(t *testing.T, dir string) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })
}

// --- removeRepoDirIfSafe unit-level coverage (proceed cases) ---

func TestRemoveRepoDirIfSafe_AbsentIsNoOp(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "does-not-exist")

	err := removeRepoDirIfSafe("clone", dir)

	if err != nil {
		t.Fatalf("removeRepoDirIfSafe(absent) = %v, want nil", err)
	}
}

func TestRemoveRepoDirIfSafe_DeletesCleanCheckout(t *testing.T) {
	dir := t.TempDir()
	cleanGitFixture(t, dir)

	err := removeRepoDirIfSafe("clone", dir)

	if err != nil {
		t.Fatalf("removeRepoDirIfSafe(clean) = %v, want nil", err)
	}
	if _, statErr := os.Stat(dir); !os.IsNotExist(statErr) {
		t.Fatalf("expected %s to be removed, stat err = %v", dir, statErr)
	}
}

func TestRemoveRepoDirIfSafe_RefusesDirtyCheckout(t *testing.T) {
	dir := t.TempDir()
	dirtyGitFixture(t, dir)

	err := removeRepoDirIfSafe("clone", dir)

	if err == nil {
		t.Fatalf("removeRepoDirIfSafe(dirty) = nil, want refusal error")
	}
	if _, statErr := os.Stat(dir); statErr != nil {
		t.Fatalf("expected %s to remain, stat err = %v", dir, statErr)
	}
}

func TestRemoveRepoDirIfSafe_RefusesNotAGitRepo(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "file.txt"), []byte("hi"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	err := removeRepoDirIfSafe("clone", dir)

	if err == nil {
		t.Fatalf("removeRepoDirIfSafe(not-a-repo) = nil, want refusal error")
	}
	if _, statErr := os.Stat(dir); statErr != nil {
		t.Fatalf("expected %s to remain, stat err = %v", dir, statErr)
	}
}

func TestRemoveRepoDirIfSafe_RefusalMessageNamesAbsolutePath(t *testing.T) {
	base := t.TempDir()
	relDir := "some-repo"
	dir := filepath.Join(base, relDir)
	dirtyGitFixture(t, dir)

	chdirTo(t, base)

	err := removeRepoDirIfSafe("clone", relDir)

	if err == nil {
		t.Fatalf("removeRepoDirIfSafe(dirty, relative) = nil, want refusal error")
	}
	wantAbs, absErr := filepath.Abs(dir)
	if absErr != nil {
		t.Fatalf("filepath.Abs: %v", absErr)
	}
	if !strings.Contains(err.Error(), wantAbs) {
		t.Fatalf("refusal error %q does not name absolute path %q", err.Error(), wantAbs)
	}
}

// --- stepCloneRepo wiring (task 3.5) ---

func TestStepCloneRepo_RefusesOnDirtyCheckout(t *testing.T) {
	base := t.TempDir()
	dirtyGitFixture(t, filepath.Join(base, "Gentleman.Dots"))
	chdirTo(t, base)

	err := stepCloneRepo(&Model{})

	if err == nil {
		t.Fatalf("stepCloneRepo() = nil, want refusal error on dirty checkout")
	}
	stepErr, ok := err.(*StepError)
	if !ok {
		t.Fatalf("stepCloneRepo() error type = %T, want *StepError", err)
	}
	if stepErr.StepID != "clone" {
		t.Fatalf("stepErr.StepID = %q, want %q", stepErr.StepID, "clone")
	}
	wantAbs, _ := filepath.Abs(filepath.Join(base, "Gentleman.Dots"))
	if !strings.Contains(err.Error(), wantAbs) {
		t.Fatalf("stepCloneRepo() error %q does not name absolute path %q", err.Error(), wantAbs)
	}
}

func TestStepCloneRepo_RefusesOnNotAGitRepo(t *testing.T) {
	base := t.TempDir()
	repoDir := filepath.Join(base, "Gentleman.Dots")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repoDir, "file.txt"), []byte("hi"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	chdirTo(t, base)

	err := stepCloneRepo(&Model{})

	if err == nil {
		t.Fatalf("stepCloneRepo() = nil, want refusal error on non-git directory")
	}
	if _, statErr := os.Stat(repoDir); statErr != nil {
		t.Fatalf("expected %s to remain undeleted, stat err = %v", repoDir, statErr)
	}
}

// --- stepCleanup wiring (task 3.6) ---

func TestStepCleanup_WarnsAndReturnsNilOnDirtyCheckout(t *testing.T) {
	base := t.TempDir()
	repoDir := filepath.Join(base, "Gentleman.Dots")
	dirtyGitFixture(t, repoDir)
	chdirTo(t, base)

	err := stepCleanup(&Model{})

	if err != nil {
		t.Fatalf("stepCleanup() = %v, want nil (non-critical refusal)", err)
	}
	if _, statErr := os.Stat(repoDir); statErr != nil {
		t.Fatalf("expected %s to remain (refused delete), stat err = %v", repoDir, statErr)
	}
}

func TestStepCleanup_DeletesOnCleanCheckout(t *testing.T) {
	base := t.TempDir()
	repoDir := filepath.Join(base, "Gentleman.Dots")
	cleanGitFixture(t, repoDir)
	chdirTo(t, base)

	err := stepCleanup(&Model{})

	if err != nil {
		t.Fatalf("stepCleanup() = %v, want nil", err)
	}
	if _, statErr := os.Stat(repoDir); !os.IsNotExist(statErr) {
		t.Fatalf("expected %s to be removed, stat err = %v", repoDir, statErr)
	}
}

func TestStepCleanup_NoOpOnAbsent(t *testing.T) {
	base := t.TempDir()
	chdirTo(t, base)

	err := stepCleanup(&Model{})

	if err != nil {
		t.Fatalf("stepCleanup() = %v, want nil on absent directory", err)
	}
}

// --- non-interactive deterministic abort (task 3.7) ---

func TestStepCloneRepo_NonInteractiveDirtyCheckoutAbortsWithoutPrompt(t *testing.T) {
	base := t.TempDir()
	dirtyGitFixture(t, filepath.Join(base, "Gentleman.Dots"))
	chdirTo(t, base)

	SetNonInteractiveMode(true)
	t.Cleanup(func() { SetNonInteractiveMode(false) })

	done := make(chan error, 1)
	go func() {
		done <- stepCloneRepo(&Model{})
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Fatalf("stepCloneRepo() = nil, want deterministic refusal error")
		}
	case <-time.After(5 * time.Second):
		// A generous but bounded window: any interactive prompt would
		// block indefinitely on stdin, so this only trips for a real hang.
		t.Fatal("stepCloneRepo() did not return - it appears to be blocking on a prompt")
	}
}
