---
name: gentleman-system
description: >
  System detection and command execution patterns for Gentleman.Dots.
  Trigger: When editing files in installer/internal/system/, adding OS support, or modifying command execution.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when:
- Adding support for new operating systems
- Modifying OS detection logic
- Working with command execution (sudo, brew, pkg)
- Adding new system checks
- Implementing backup/restore functionality

---

## Critical Patterns

### Pattern 1: OSType Enum

All OS types are defined in `detect.go`:

```go
type OSType int

const (
    OSMac OSType = iota
    OSLinux
    OSArch
    OSDebian    // Debian-based (Debian, Ubuntu)
    OSTermux    // Termux on Android
    OSUnknown
)
```

### Pattern 2: SystemInfo Structure

Detection results in `SystemInfo` struct:

```go
type SystemInfo struct {
    OS        OSType
    OSName    string
    IsWSL     bool
    IsARM     bool
    IsTermux  bool
    HomeDir   string
    HasBrew   bool
    HasPkg    bool    // Termux package manager
    HasXcode  bool
    UserShell string
    Prefix    string  // Termux $PREFIX or empty
}
```

### Pattern 3: OS Detection Priority

Termux is checked FIRST (runs on Linux but is special):

```go
func Detect() *SystemInfo {
    info := &SystemInfo{...}

    // Check Termux FIRST
    if isTermux() {
        info.OS = OSTermux
        info.IsTermux = true
        info.HasPkg = checkPkg()
        return info
    }

    // Then check standard OS
    switch runtime.GOOS {
    case "darwin":
        info.OS = OSMac
    case "linux":
        if isArchLinux() {
            info.OS = OSArch
        } else if isDebian() {
            info.OS = OSDebian
        }
    }
    return info
}
```

### Pattern 4: Command Execution Functions

Use the right function for each context:

```go
// Basic command (no sudo, no logs)
system.Run("git clone ...", nil)

// With real-time logs
system.RunWithLogs("git clone ...", nil, func(line string) {
    SendLog(stepID, line)
})

// Homebrew commands
system.RunBrewWithLogs("install fish", nil, logFunc)

// Sudo commands (password prompt)
system.RunSudo("apt-get install -y git", nil)
system.RunSudoWithLogs("pacman -S git", nil, logFunc)

// Termux pkg commands (no sudo needed)
system.RunPkgInstall("fish git", nil, logFunc)
system.RunPkgWithLogs("update", nil, logFunc)
```

---

## Decision Tree

```
Adding new OS support?
├── Add OSType constant in detect.go
├── Add detection function (isNewOS())
├── Update Detect() with priority order
├── Update SystemInfo if new fields needed
└── Add OS case in installer.go steps

Running a command?
├── Needs sudo? → RunSudo() or RunSudoWithLogs()
├── Needs brew? → RunBrewWithLogs()
├── On Termux? → RunPkgInstall() or RunPkgWithLogs()
├── Needs logs? → RunWithLogs()
└── Simple exec? → Run()

Checking if tool exists?
├── Use CommandExists("toolname")
└── Returns bool
```

---

## Code Examples

### Example 1: Termux Detection

```go
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
```

### Example 2: Platform-Specific Execution

```go
func installTool(m *Model) error {
    var result *system.ExecResult

    switch {
    case m.SystemInfo.IsTermux:
        // Termux: use pkg (no sudo)
        result = system.RunPkgInstall("tool", nil, logFunc)

    case m.SystemInfo.OS == system.OSArch:
        // Arch: use pacman with sudo
        result = system.RunSudoWithLogs("pacman -S --noconfirm tool", nil, logFunc)

    case m.SystemInfo.OS == system.OSMac:
        // macOS: use Homebrew
        result = system.RunBrewWithLogs("install tool", nil, logFunc)

    case m.SystemInfo.OS == system.OSDebian:
        // Debian/Ubuntu: use Homebrew (installed by us)
        result = system.RunBrewWithLogs("install tool", nil, logFunc)

    default:
        return fmt.Errorf("unsupported OS: %v", m.SystemInfo.OS)
    }

    return result.Error
}
```

### Example 3: Homebrew Prefix Detection

```go
func GetBrewPrefix() string {
    if runtime.GOOS == "darwin" {
        // Apple Silicon uses /opt/homebrew
        // Intel uses /usr/local
        if runtime.GOARCH == "arm64" {
            return "/opt/homebrew"
        }
        return "/usr/local"
    }
    return "/home/linuxbrew/.linuxbrew"
}
```

### Example 4: File Operations

```go
// Ensure directory exists
if err := system.EnsureDir(filepath.Join(homeDir, ".config/tool")); err != nil {
    return err
}

// Copy single file
if err := system.CopyFile(src, dst); err != nil {
    return err
}

// Copy directory contents
if err := system.CopyDir("Gentleman.Dots/Config/*", destDir+"/"); err != nil {
    return err
}
```

---

## ExecResult Structure

```go
type ExecResult struct {
    Output   string  // stdout
    Stderr   string  // stderr
    ExitCode int     // exit code
    Error    error   // error if any
}

// Usage
result := system.Run("command", nil)
if result.Error != nil {
    // Handle error
}
if result.ExitCode != 0 {
    // Non-zero exit
}
```

---

## Commands

```bash
cd installer && go test ./internal/system/...   # Run system tests
cd installer && go test -run TestDetect         # Test OS detection
cd installer && go test -run TestExec           # Test command execution
```

---

## Resources

- **Detection**: See `installer/internal/system/detect.go` for OS detection
- **Execution**: See `installer/internal/system/exec.go` for command running
- **Backup**: See `installer/internal/system/backup.go` for backup/restore
- **Tests**: See `installer/internal/system/*_test.go` for patterns
