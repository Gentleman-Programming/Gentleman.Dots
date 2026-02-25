---
name: gentleman-installer
description: >
  Installation step patterns for Gentleman.Dots TUI installer.
  Trigger: When editing installer.go, adding installation steps, or modifying the installation flow.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when:
- Adding new installation steps
- Modifying existing tool installations
- Working on backup/restore functionality
- Implementing non-interactive mode support
- Adding new OS/platform support

---

## Critical Patterns

### Pattern 1: InstallStep Structure

All steps follow this structure in `model.go`:

```go
type InstallStep struct {
    ID          string      // Unique identifier: "terminal", "shell", etc.
    Name        string      // Display name: "Install Fish"
    Description string      // Short description
    Status      StepStatus  // Pending, Running, Done, Failed, Skipped
    Progress    float64     // 0.0 - 1.0
    Error       error       // Error if failed
    Interactive bool        // Needs terminal control (sudo, chsh)
}
```

### Pattern 2: Step Registration in SetupInstallSteps

Steps MUST be registered in `SetupInstallSteps()` in `model.go`:

```go
func (m *Model) SetupInstallSteps() {
    m.Steps = []InstallStep{}

    // Conditional step based on user choice
    if m.Choices.SomeChoice {
        m.Steps = append(m.Steps, InstallStep{
            ID:          "newstep",
            Name:        "Install Something",
            Description: "Description here",
            Status:      StatusPending,
            Interactive: false, // true if needs sudo/password
        })
    }
}
```

### Pattern 3: Step Execution in executeStep

All step logic goes in `installer.go`:

```go
func executeStep(stepID string, m *Model) error {
    switch stepID {
    case "newstep":
        return stepNewStep(m)
    // ... other cases
    default:
        return fmt.Errorf("unknown step: %s", stepID)
    }
}

func stepNewStep(m *Model) error {
    stepID := "newstep"

    SendLog(stepID, "Starting installation...")

    // Check if already installed
    if system.CommandExists("newtool") {
        SendLog(stepID, "Already installed, skipping...")
        return nil
    }

    // Install based on OS
    var result *system.ExecResult
    if m.SystemInfo.IsTermux {
        result = system.RunPkgInstall("newtool", nil, func(line string) {
            SendLog(stepID, line)
        })
    } else {
        result = system.RunBrewWithLogs("install newtool", nil, func(line string) {
            SendLog(stepID, line)
        })
    }

    if result.Error != nil {
        return wrapStepError("newstep", "Install NewTool",
            "Failed to install NewTool",
            result.Error)
    }

    SendLog(stepID, "✓ NewTool installed")
    return nil
}
```

### Pattern 4: Interactive Steps (sudo/password required)

Mark step as Interactive and use `runInteractiveStep`:

```go
// In SetupInstallSteps:
m.Steps = append(m.Steps, InstallStep{
    ID:          "interactive_step",
    Name:        "Configure System",
    Description: "Requires password",
    Status:      StatusPending,
    Interactive: true,  // KEY: marks as interactive
})

// In runNextStep (update.go):
if step.Interactive {
    return runInteractiveStep(step.ID, &m)
}
```

---

## Decision Tree

```
Adding new tool installation?
├── Add step to SetupInstallSteps() with conditions
├── Add case in executeStep() switch
├── Create step{Name}() function in installer.go
├── Handle all OS variants (Mac, Linux, Arch, Debian, Termux)
├── Use SendLog() for progress updates
└── Return wrapStepError() on failure

Step needs password/sudo?
├── Set Interactive: true in InstallStep
├── Use system.RunSudo() or system.RunSudoWithLogs()
└── Use tea.ExecProcess for full terminal control

Step should be conditional?
├── Check m.Choices.{option} before appending
├── Check m.SystemInfo for OS-specific logic
└── Use StatusSkipped if conditions not met
```

---

## Code Examples

### Example 1: OS-Specific Installation

```go
func stepInstallTool(m *Model) error {
    stepID := "tool"

    if !system.CommandExists("tool") {
        SendLog(stepID, "Installing tool...")

        var result *system.ExecResult
        switch {
        case m.SystemInfo.IsTermux:
            result = system.RunPkgInstall("tool", nil, logFunc(stepID))
        case m.SystemInfo.OS == system.OSArch:
            result = system.RunSudoWithLogs("pacman -S --noconfirm tool", nil, logFunc(stepID))
        case m.SystemInfo.OS == system.OSMac:
            result = system.RunBrewWithLogs("install tool", nil, logFunc(stepID))
        default: // Debian/Ubuntu
            result = system.RunBrewWithLogs("install tool", nil, logFunc(stepID))
        }

        if result.Error != nil {
            return wrapStepError("tool", "Install Tool",
                "Failed to install tool",
                result.Error)
        }
    }

    // Copy configuration
    SendLog(stepID, "Copying configuration...")
    homeDir := os.Getenv("HOME")
    if err := system.CopyDir(filepath.Join("Gentleman.Dots", "ToolConfig/*"),
        filepath.Join(homeDir, ".config/tool/")); err != nil {
        return wrapStepError("tool", "Install Tool",
            "Failed to copy configuration",
            err)
    }

    SendLog(stepID, "✓ Tool configured")
    return nil
}

func logFunc(stepID string) func(string) {
    return func(line string) {
        SendLog(stepID, line)
    }
}
```

### Example 2: Error Wrapping Pattern

```go
func wrapStepError(stepID, stepName, description string, cause error) error {
    return &StepError{
        StepID:      stepID,
        StepName:    stepName,
        Description: description,
        Cause:       cause,
    }
}

// Usage:
if result.Error != nil {
    return wrapStepError("terminal", "Install Alacritty",
        "Failed to install Alacritty. Check your internet connection.",
        result.Error)
}
```

### Example 3: Config Patching

```go
// Patch config based on user choices
func stepInstallShell(m *Model) error {
    // ... install shell ...

    // Patch config for window manager choice
    configPath := filepath.Join(homeDir, ".config/fish/config.fish")
    if err := system.PatchFishForWM(configPath, m.Choices.WindowMgr, m.Choices.InstallNvim); err != nil {
        return wrapStepError("shell", "Install Fish",
            "Failed to configure window manager in shell",
            err)
    }

    return nil
}
```

---

## Logging Pattern

Always use `SendLog` for step progress:

```go
SendLog(stepID, "Starting...")           // Start
SendLog(stepID, "Downloading...")        // Progress
SendLog(stepID, "  → file.txt")          // Sub-item
SendLog(stepID, "✓ Step completed")      // Success
```

---

## Commands

```bash
cd installer && go build ./cmd/gentleman-installer           # Build
./gentleman-installer --help                                  # Show help
./gentleman-installer --non-interactive --shell=fish         # Non-interactive
GENTLEMAN_VERBOSE=1 ./gentleman-installer --non-interactive  # Verbose logs
```

---

## Resources

- **Steps**: See `installer/internal/tui/installer.go` for step implementations
- **Model**: See `installer/internal/tui/model.go` for SetupInstallSteps
- **System**: See `installer/internal/system/exec.go` for command execution
- **Non-interactive**: See `installer/internal/tui/non_interactive.go` for CLI mode
