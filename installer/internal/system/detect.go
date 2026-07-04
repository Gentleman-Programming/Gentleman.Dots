package system

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// atomicDistros contains os-release ID values that are known atomic/immutable.
// This is the authoritative source for detection.
var atomicDistroIDs = map[string]bool{
	"silverblue":      true,
	"fedora-silverblue": true,
	"kinoite":         true,
	"fedora-kinoite":  true,
	"bazzite":         true,
	"bluefin":         true,
	"aurora":          true,
	"sericea":         true,
	"wayblue":         true,
	"fedora-iot":      true,
}

type OSType int

const (
	OSMac OSType = iota
	OSLinux
	OSArch
	OSDebian  // Debian-based (Debian, Ubuntu, etc.)
	OSFedora  // Fedora/RHEL-based (Fedora, CentOS, RHEL, etc.)
	OSTermux  // Termux on Android
	OSUnknown
)

type SystemInfo struct {
	OS         OSType
	OSName     string
	IsWSL      bool
	IsARM      bool
	IsTermux   bool
	IsAtomic   bool   // Atomic/immutable distro (Silverblue, Bazzite, etc.)
	HomeDir    string
	HasBrew    bool
	HasPkg     bool // Termux package manager
	HasXcode   bool
	HasFlatpak bool
	UserShell  string
	Prefix     string // Termux $PREFIX or empty for other systems
}

func Detect() *SystemInfo {
	info := &SystemInfo{
		OS:      OSUnknown,
		OSName:  "Unknown",
		HomeDir: os.Getenv("HOME"),
		IsARM:   runtime.GOARCH == "arm64" || runtime.GOARCH == "arm",
		Prefix:  os.Getenv("PREFIX"),
	}

	// Check for Termux FIRST (it runs on Linux but is special)
	if isTermux() {
		info.OS = OSTermux
		info.OSName = "Termux"
		info.IsTermux = true
		info.HasPkg = checkPkg()
		info.HasBrew = false // Termux doesn't use Homebrew
		info.UserShell = detectCurrentShell()
		return info
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
		} else if isFedora() {
			info.OS = OSFedora
			info.OSName = "Fedora/RHEL"
		} else if isDebian() {
			info.OS = OSDebian
			info.OSName = "Debian/Ubuntu"
		}
	}

	info.IsAtomic = checkAtomic()
	if info.IsAtomic {
		info.OSName += " (Atomic)"
	}
	info.HasBrew = checkBrew()
	info.HasFlatpak = checkFlatpak()
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

func isFedora() bool {
	// Check for Fedora specifically
	if _, err := os.Stat("/etc/fedora-release"); err == nil {
		return true
	}
	// Check for RHEL/CentOS (also use dnf)
	if _, err := os.Stat("/etc/redhat-release"); err == nil {
		return true
	}
	return false
}

// isTermux detects if we're running in Termux on Android
func isTermux() bool {
	// Check TERMUX_VERSION environment variable
	if os.Getenv("TERMUX_VERSION") != "" {
		return true
	}
	// Check PREFIX contains termux path
	prefix := os.Getenv("PREFIX")
	if strings.Contains(prefix, "com.termux") {
		return true
	}
	// Check for Termux-specific paths
	if _, err := os.Stat("/data/data/com.termux"); err == nil {
		return true
	}
	return false
}

// checkPkg checks if Termux pkg command is available
func checkPkg() bool {
	_, err := exec.LookPath("pkg")
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

// checkAtomic detects if we're running on an atomic/immutable Linux distro.
// Priority: /run/ostree-booted (most reliable), rpm-ostree binary, os-release ID.
func checkAtomic() bool {
	if runtime.GOOS != "linux" {
		return false
	}

	// Most reliable indicator: /run/ostree-booted exists on ostree-based systems
	if _, err := os.Stat("/run/ostree-booted"); err == nil {
		return true
	}

	// Check for rpm-ostree binary (present on Fedora Atomic derivatives)
	if _, err := exec.LookPath("rpm-ostree"); err == nil {
		return true
	}

	// Check os-release for known atomic distro IDs
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return false
	}

	content := strings.ToLower(string(data))
	for id := range atomicDistroIDs {
		if strings.Contains(content, "id="+id) || strings.Contains(content, "id=\""+id+"\"") || strings.Contains(content, "variant_id="+id) || strings.Contains(content, "variant_id=\""+id+"\"") {
			return true
		}
	}

	return false
}

// checkFlatpak detects if the flatpak command is available
func checkFlatpak() bool {
	_, err := exec.LookPath("flatpak")
	return err == nil
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
