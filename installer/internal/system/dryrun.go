package system

import (
	"os"
)

// IsDryRun reports whether the installer was started with --dry-run, which
// cmd/gentleman-installer/main.go signals by setting GENTLEMAN_DRY_RUN=1.
// It is a direct env lookup (no caching) so tests can set and unset the
// variable with t.Setenv.
//
// It reads ONLY GENTLEMAN_DRY_RUN. --test sets GENTLEMAN_TEST_MODE instead
// (see setupTestMode) and intentionally keeps performing package-manager,
// sudo and network operations, so it must never be treated as a dry run.
func IsDryRun() bool {
	return os.Getenv("GENTLEMAN_DRY_RUN") == "1"
}

// dryRunResult is the synthetic success returned by the dry-run gate. It
// keeps every existing caller on its success path without executing
// anything: upstream #190 requires a dry run to perform zero mutations
// while still letting the flow run to completion.
func dryRunResult(command string) *ExecResult {
	return &ExecResult{Command: command, ExitCode: 0, Error: nil}
}
