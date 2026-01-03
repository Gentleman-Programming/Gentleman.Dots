package system

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type OSType int

const (
	OSMac OSType = iota
	OSLinux
	OSArch
	OSDebian // Debian-based (Debian, Ubuntu, etc.)
	OSUnknown
)

type SystemInfo struct {
	OS        OSType
	OSName    string
	IsWSL     bool
	IsARM     bool
	HomeDir   string
	HasBrew   bool
	HasXcode  bool
	UserShell string
}

func Detect() *SystemInfo {
	info := &SystemInfo{
		OS:      OSUnknown,
		OSName:  "Unknown",
		HomeDir: os.Getenv("HOME"),
		IsARM:   runtime.GOARCH == "arm64" || runtime.GOARCH == "arm",
	}

	switch runtime.GOOS {
	case "darwin":
		info.OS = OSMac
		info.OSName = "macOS"
		info.HasXcode = checkXcode()
	case "linux":
		info.OS = OSLinux
		info.OSName = "Linux"
		info.IsWSL = checkWSL()

		if isArchLinux() {
			info.OS = OSArch
			info.OSName = "Arch Linux"
		} else if isDebian() {
			info.OS = OSDebian
			info.OSName = "Debian/Ubuntu"
		}
	}

	info.HasBrew = checkBrew()
	info.UserShell = detectCurrentShell()

	return info
}

func checkWSL() bool {
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	content := strings.ToLower(string(data))
	return strings.Contains(content, "microsoft") || strings.Contains(content, "wsl")
}

func isArchLinux() bool {
	_, err := os.Stat("/etc/arch-release")
	return err == nil
}

func isDebian() bool {
	_, err := os.Stat("/etc/debian_version")
	return err == nil
}

func checkBrew() bool {
	_, err := exec.LookPath("brew")
	return err == nil
}

func checkXcode() bool {
	cmd := exec.Command("xcode-select", "-p")
	return cmd.Run() == nil
}

func detectCurrentShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "unknown"
	}
	parts := strings.Split(shell, "/")
	return parts[len(parts)-1]
}

// CommandExists checks if a command is available in PATH
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// GetBrewPrefix returns the homebrew prefix path
func GetBrewPrefix() string {
	if runtime.GOOS == "darwin" {
		// Apple Silicon (arm64) uses /opt/homebrew
		// Intel (amd64) uses /usr/local
		if runtime.GOARCH == "arm64" {
			return "/opt/homebrew"
		}
		return "/usr/local"
	}
	return "/home/linuxbrew/.linuxbrew"
}
