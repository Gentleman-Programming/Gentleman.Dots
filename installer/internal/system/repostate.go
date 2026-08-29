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

	if strings.TrimSpace(stdout) == "" {
		return RepoCleanCheckout
	}
	return RepoDirtyCheckout
}
