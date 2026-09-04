package system

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RepoState classifies the on-disk state of a directory that the installer
// is about to destructively remove, so the delete guard can fail closed
// instead of running an unconditional rm -rf.
type RepoState int

const (
	// RepoAbsent means the directory does not exist on disk.
	RepoAbsent RepoState = iota
	// RepoCleanCheckout means the directory is a git repository with no
	// uncommitted changes and no untracked files.
	RepoCleanCheckout
	// RepoDirtyCheckout means the directory is a git repository with
	// uncommitted changes and/or untracked files.
	RepoDirtyCheckout
	// RepoNotAGitRepo means the directory exists but is not a git
	// repository (no .git).
	RepoNotAGitRepo
	// RepoUnknown means the git status could not be determined for any
	// reason (git binary missing, command error, unparseable output).
	RepoUnknown
	// RepoUnpublishedWork means the working tree is clean but the
	// repository still holds work that exists nowhere else: commits that
	// are on no remote, a branch whose publication cannot be proven, or a
	// stash. A clean working tree is not proof that nothing would be lost.
	RepoUnpublishedWork
)

// runGitStatus runs `git -C <dir> status --porcelain` and returns its raw
// stdout. It is a package-level var-function seam so tests can override and
// restore it, matching the pattern used elsewhere in this codebase
// (runPkgInstallWithLogs/runSudoWithLogs/runBrewWithLogs).
//
// The directory is passed as a discrete argument to git -C, never
// interpolated into a shell command string.
var runGitStatus = func(dir string) (stdout string, err error) {
	cmd := exec.Command("git", "-C", dir, "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// runGitLines runs an arbitrary read-only git command in dir and returns its
// raw stdout. Like runGitStatus it is a var-function seam for tests, and the
// directory is passed as a discrete argument, never interpolated into a
// shell string.
var runGitLines = func(dir string, args ...string) (stdout string, err error) {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// hasUnpublishedWork reports whether a clean checkout still holds work that
// deleting it would destroy. It fails closed: anything it cannot prove
// published counts as unpublished.
func hasUnpublishedWork(absPath string) (bool, error) {
	// A stash is invisible to `git status --porcelain`.
	stashes, err := runGitLines(absPath, "stash", "list")
	if err != nil {
		return true, err
	}
	if strings.TrimSpace(stashes) != "" {
		return true, nil
	}

	// HEAD must be contained by at least one remote-tracking branch. This
	// also covers a branch with no upstream and a repository with no
	// remotes at all: both yield no containing remote branch, so both are
	// treated as unpublished.
	remotes, err := runGitLines(absPath, "branch", "--remotes", "--contains", "HEAD")
	if err != nil {
		return true, err
	}
	return strings.TrimSpace(remotes) == "", nil
}

// InspectRepoDir classifies the git status of path, resolving it to an
// absolute path first. The defect this guards against was an rm -rf on a
// relative path, so the absolutization MUST happen before any git
// invocation or deletion decision.
func InspectRepoDir(path string) RepoState {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return RepoUnknown
	}

	info, statErr := os.Stat(absPath)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			return RepoAbsent
		}
		return RepoUnknown
	}
	if !info.IsDir() {
		return RepoUnknown
	}

	gitDir := filepath.Join(absPath, ".git")
	if _, err := os.Stat(gitDir); err != nil {
		if os.IsNotExist(err) {
			return RepoNotAGitRepo
		}
		return RepoUnknown
	}

	stdout, err := runGitStatus(absPath)
	if err != nil {
		return RepoUnknown
	}

	if strings.TrimSpace(stdout) != "" {
		return RepoDirtyCheckout
	}

	// The working tree is clean, which is NOT the same as "nothing would be
	// lost": committed-but-unpushed work and stashes survive a clean status.
	unpublished, err := hasUnpublishedWork(absPath)
	if err != nil {
		return RepoUnknown
	}
	if unpublished {
		return RepoUnpublishedWork
	}
	return RepoCleanCheckout
}
