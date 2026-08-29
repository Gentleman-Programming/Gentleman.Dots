package system

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// gitFixture builds a real git repository with one commit and a fake remote
// that the branch tracks, so the checkout starts fully published.
func gitFixture(t *testing.T) (repo string) {
	t.Helper()

	root := t.TempDir()
	origin := filepath.Join(root, "origin.git")
	repo = filepath.Join(root, "work")

	run := func(dir string, args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.com",
			"GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.com",
		)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}

	run(root, "init", "--bare", "-b", "main", origin)
	run(root, "clone", "--quiet", origin, repo)

	if err := os.WriteFile(filepath.Join(repo, "a.txt"), []byte("a\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	run(repo, "add", "a.txt")
	run(repo, "commit", "-qm", "first")
	run(repo, "push", "--quiet", "-u", "origin", "main")

	return repo
}

func gitRun(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.com",
		"GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.com",
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

// A checkout whose work is fully pushed is the only case that stays safe to
// delete.
func TestInspectRepoDir_PublishedCleanCheckoutIsDeletable(t *testing.T) {
	repo := gitFixture(t)

	if got := InspectRepoDir(repo); got != RepoCleanCheckout {
		t.Fatalf("InspectRepoDir = %v, want RepoCleanCheckout for a fully published clean checkout", got)
	}
}

// The defect that destroyed real work: committing everything made the
// checkout clean, so the guard classified it as deletable even though the
// commits existed nowhere else.
func TestInspectRepoDir_UnpushedCommitIsNotDeletable(t *testing.T) {
	repo := gitFixture(t)

	if err := os.WriteFile(filepath.Join(repo, "b.txt"), []byte("b\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, repo, "add", "b.txt")
	gitRun(t, repo, "commit", "-qm", "local only")

	got := InspectRepoDir(repo)
	if got == RepoCleanCheckout {
		t.Fatal("InspectRepoDir = RepoCleanCheckout for a checkout holding an unpushed commit; the guard would delete unrecoverable work")
	}
	if got != RepoUnpublishedWork {
		t.Fatalf("InspectRepoDir = %v, want RepoUnpublishedWork", got)
	}
}

// The guard must not over-block: a branch with no upstream whose HEAD is
// already contained by a remote branch loses no commits if removed. Only the
// branch name would go, and that is not work.
func TestInspectRepoDir_UnpushedBranchOfPublishedCommitIsDeletable(t *testing.T) {
	repo := gitFixture(t)
	gitRun(t, repo, "checkout", "-qb", "side-branch")

	if got := InspectRepoDir(repo); got != RepoCleanCheckout {
		t.Fatalf("InspectRepoDir = %v, want RepoCleanCheckout: HEAD is already on a remote", got)
	}
}

// A stash holds work that no commit and no status line reveals.
func TestInspectRepoDir_StashedWorkIsNotDeletable(t *testing.T) {
	repo := gitFixture(t)

	if err := os.WriteFile(filepath.Join(repo, "a.txt"), []byte("changed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, repo, "stash", "push", "-q", "-m", "wip")

	// The working tree is clean again at this point, which is exactly what
	// makes this case dangerous.
	got := InspectRepoDir(repo)
	if got == RepoCleanCheckout {
		t.Fatal("InspectRepoDir = RepoCleanCheckout for a checkout with a stash; the stashed work would be destroyed")
	}
	if got != RepoUnpublishedWork {
		t.Fatalf("InspectRepoDir = %v, want RepoUnpublishedWork", got)
	}
}

// An uncommitted change still classifies as dirty, not as unpublished work.
func TestInspectRepoDir_DirtyStillReportsDirty(t *testing.T) {
	repo := gitFixture(t)

	if err := os.WriteFile(filepath.Join(repo, "a.txt"), []byte("changed\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if got := InspectRepoDir(repo); got != RepoDirtyCheckout {
		t.Fatalf("InspectRepoDir = %v, want RepoDirtyCheckout", got)
	}
}
