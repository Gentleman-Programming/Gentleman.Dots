package system

import "fmt"

// Notify is the default reporter for messages that have no dedicated log
// callback (e.g. CopyFile/CopyDir skip warnings). It is a package-level
// var-function so callers (tests, the TUI) can override and restore it,
// matching the existing var-function seam pattern used by
// runPkgInstallWithLogs/runSudoWithLogs/runBrewWithLogs in
// internal/tui/installer.go.
var Notify = func(msg string) { fmt.Println(msg) }

// report routes msg to logFunc when provided, falling back to Notify
// otherwise. This lets call sites that already stream to a log callback
// (e.g. RunWithLogs) keep messages on that channel, while call sites with
// no callback (e.g. CopyDir) still surface messages somewhere visible.
func report(logFunc func(string), msg string) {
	if logFunc != nil {
		logFunc(msg)
	} else {
		Notify(msg)
	}
}
