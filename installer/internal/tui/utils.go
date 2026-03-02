package tui

import (
	"os/exec"
	"runtime"
)

// openURL opens a URL or file path in the system's default application
func openURL(url string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Start()
	case "linux":
		return exec.Command("xdg-open", url).Start()
	}
	return nil
}
