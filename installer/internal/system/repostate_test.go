package system

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// initGitRepo runs `git init` plus committer identity config in dir so
// fixtures can commit without relying on any global/user git config.
func initGitRepo(t *testing.T, dir string) {
	t.Helper()
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "Test User")
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}

func TestInspectRepoDir_Absent(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "does-not-exist")

	got := InspectRepoDir(dir)

	if got != RepoAbsent {
		t.Fatalf("InspectRepoDir(%q) = %v, want RepoAbsent", dir, got)
	}
}

func TestInspectRepoDir_CleanCheckout(t *testing.T) {
	dir := t.TempDir()
	initGitRepo(t, dir)

	filePath := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(filePath, []byte("hello"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	runGit(t, dir, "add", "file.txt")
	runGit(t, dir, "commit", "-m", "initial commit")

	got := InspectRepoDir(dir)

	if got != RepoCleanCheckout {
		t.Fatalf("InspectRepoDir(%q) = %v, want RepoCleanCheckout", dir, got)
	}
}

func TestInspectRepoDir_DirtyCheckout_Unstaged(t *testing.T) {
	dir := t.TempDir()
	initGitRepo(t, dir)

	filePath := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(filePath, []byte("hello"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	runGit(t, dir, "add", "file.txt")
	runGit(t, dir, "commit", "-m", "initial commit")

	// Modify the tracked file without staging - unstaged change.
	if err := os.WriteFile(filePath, []byte("hello world"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	got := InspectRepoDir(dir)

	if got != RepoDirtyCheckout {
		t.Fatalf("InspectRepoDir(%q) = %v, want RepoDirtyCheckout (unstaged)", dir, got)
	}
}

func TestInspectRepoDir_DirtyCheckout_Staged(t *testing.T) {
	dir := t.TempDir()
	initGitRepo(t, dir)

	filePath := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(filePath, []byte("hello"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	runGit(t, dir, "add", "file.txt")
	runGit(t, dir, "commit", "-m", "initial commit")

	// Stage a change but do not commit it.
	if err := os.WriteFile(filePath, []byte("hello world"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	runGit(t, dir, "add", "file.txt")

	got := InspectRepoDir(dir)

	if got != RepoDirtyCheckout {
		t.Fatalf("InspectRepoDir(%q) = %v, want RepoDirtyCheckout (staged)", dir, got)
	}
}

func TestInspectRepoDir_DirtyCheckout_Untracked(t *testing.T) {
	dir := t.TempDir()
	initGitRepo(t, dir)

	filePath := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(filePath, []byte("hello"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	runGit(t, dir, "add", "file.txt")
	runGit(t, dir, "commit", "-m", "initial commit")

	// Add a new untracked file, never staged or committed.
	untrackedPath := filepath.Join(dir, "untracked.txt")
	if err := os.WriteFile(untrackedPath, []byte("new"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	got := InspectRepoDir(dir)

	if got != RepoDirtyCheckout {
		t.Fatalf("InspectRepoDir(%q) = %v, want RepoDirtyCheckout (untracked)", dir, got)
	}
}

func TestInspectRepoDir_NotAGitRepo(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "file.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	got := InspectRepoDir(dir)

	if got != RepoNotAGitRepo {
		t.Fatalf("InspectRepoDir(%q) = %v, want RepoNotAGitRepo", dir, got)
	}
}

func TestInspectRepoDir_Unknown_GitStatusError(t *testing.T) {
	dir := t.TempDir()
	// Directory exists and looks like a git repo (has a .git dir) so
	// InspectRepoDir reaches the runGitStatus call, whose override then
	// fails to simulate git being unavailable or erroring.
	if err := os.Mkdir(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}

	origRunGitStatus := runGitStatus
	t.Cleanup(func() { runGitStatus = origRunGitStatus })
	runGitStatus = func(d string) (string, error) {
		return "", errors.New("git: command not found")
	}

	got := InspectRepoDir(dir)

	if got != RepoUnknown {
		t.Fatalf("InspectRepoDir(%q) = %v, want RepoUnknown", dir, got)
	}
}

func TestInspectRepoDir_AbsolutizesRelativePath(t *testing.T) {
	base := t.TempDir()
	relName := "relative-repo"
	dir := filepath.Join(base, relName)
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}
	if err := os.Mkdir(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}

	origRunGitStatus := runGitStatus
	t.Cleanup(func() { runGitStatus = origRunGitStatus })

	var receivedDir string
	runGitStatus = func(d string) (string, error) {
		receivedDir = d
		return "", nil
	}

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	if err := os.Chdir(base); err != nil {
		t.Fatalf("Chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(origWd) })

	// Call with a relative path - runGitStatus must receive an absolute one.
	InspectRepoDir(relName)

	if !filepath.IsAbs(receivedDir) {
		t.Fatalf("runGitStatus received non-absolute path %q, want an absolute path", receivedDir)
	}

	wantAbs, err := filepath.Abs(filepath.Join(base, relName))
	if err != nil {
		t.Fatalf("filepath.Abs: %v", err)
	}
	if receivedDir != wantAbs {
		t.Fatalf("runGitStatus received %q, want %q", receivedDir, wantAbs)
	}
}
