# Linux Installation Fixes

Investigation and fixes for two bugs affecting the installer on Linux and WSL.
Reported in issue #137. Fixed in PRs #148 and #149.

---

## How to Build and Run Locally

### Prerequisites

Install Go (required to compile from source):

```bash
curl -fsSL https://go.dev/dl/go1.23.4.linux-amd64.tar.gz | sudo tar -C /usr/local -xzf -
echo 'export PATH="/usr/local/go/bin:$PATH"' >> ~/.bashrc && source ~/.bashrc
```

### Compile

```bash
cd installer
go build -o ~/gentleman.dots ./cmd/gentleman-installer
```

### Run

```bash
~/gentleman.dots
```

### Test with Docker (no Go required)

```bash
docker build -f Dockerfile.test -t gentleman-test --no-cache .
docker run -it --rm gentleman-test
```

---

## Bug 1 — Homebrew fails with "Bash is required to interpret this script."

### Symptom

On Debian/Ubuntu/WSL, the Homebrew installation step fails immediately:

```
🍺 Installing Homebrew package manager...
   (You may be prompted for your password)

Bash is required to interpret this script.
```

### Root Cause

`getHomebrewScript()` in `installer/internal/tui/interactive.go` generated a script
with `#!/bin/sh` that called the Homebrew installer via:

```sh
/bin/sh -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

On Debian/Ubuntu, `/bin/sh` is symlinked to **dash**, not bash. Homebrew's `install.sh`
checks `$BASH_VERSION` at the very top and aborts if it is empty — which it always is
under dash.

### Why `stepInstallHomebrew` in `installer.go` was not the issue

The Homebrew step is registered with `Interactive: true` in `model.go`:

```go
m.Steps = append(m.Steps, InstallStep{
    ID:          "homebrew",
    Interactive: true,
    ...
})
```

This means `runNextStep()` always routes it through `runInteractiveStep()` →
`getHomebrewScript()` → `tea.ExecProcess`. The `stepInstallHomebrew` function in
`installer.go` is never called for this path.

### Fix — `installer/internal/tui/interactive.go`

```diff
- script := fmt.Sprintf(`#!/bin/sh
+ script := fmt.Sprintf(`#!/bin/bash
  set -e
  ...
- /bin/sh -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
+ NONINTERACTIVE=1 /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

Two changes:
1. Shebang changed to `#!/bin/bash`.
2. Installer call uses explicit `/bin/bash` and `NONINTERACTIVE=1` to skip
   post-install interactive prompts that previously caused a false `exit status 1`.

### PR

**#148** — `fix(homebrew): use bash and NONINTERACTIVE=1 in interactive install script`

---

## Bug 2 — Shell change fails with "chsh: PAM: Authentication failure"

### Symptom

After Homebrew and shell installation, the step that sets the default shell fails:

```
🔐 Changing default shell...
   (You may need to enter your password)

Password:
chsh: PAM: Authentication failure
```

### Root Cause

`getSetShellScript()` uses `chsh -s "$SHELL_PATH"` which requires PAM authentication.
In Docker containers and minimal Linux environments, PAM is not fully configured, so
this always fails — even when the user has full `sudo` access with `NOPASSWD`.

### Fix — `installer/internal/tui/interactive.go`

```diff
- chsh -s "$SHELL_PATH"
-
- echo ""
- echo "✅ Default shell changed to $SHELL_PATH"
+ if chsh -s "$SHELL_PATH" 2>/dev/null; then
+     echo ""
+     echo "✅ Default shell changed to $SHELL_PATH"
+ elif sudo usermod -s "$SHELL_PATH" "$(whoami)" 2>/dev/null; then
+     echo ""
+     echo "✅ Default shell changed to $SHELL_PATH (via usermod)"
+ else
+     echo ""
+     echo "⚠️  Could not change default shell automatically."
+     echo "   Run manually: chsh -s $SHELL_PATH"
+ fi
```

Strategy:
1. Try `chsh` first — works on real machines where PAM is properly configured.
2. Fall back to `sudo usermod` — works in Docker and environments with sudo access.
3. If both fail — show a clear manual instruction instead of crashing.

### PR

**#149** — `fix(setshell): fallback to sudo usermod when chsh fails PAM auth`

---

## Test Result

Full installation with both fixes applied on Ubuntu 22.04 (Docker):

```
✨ Installation Complete! ✨

  • OS: linux
  • Terminal: alacritty
  • Shell: fish
  • Window Manager: tmux
  • Font: Iosevka Term Nerd Font
  • Editor: Neovim with Gentleman config
```

---

## Files Modified

| File | Change |
|------|--------|
| `installer/internal/tui/interactive.go` | Homebrew script: `#!/bin/sh` → `#!/bin/bash`, `/bin/sh` → `NONINTERACTIVE=1 /bin/bash` |
| `installer/internal/tui/interactive.go` | Shell change: `chsh` with fallback to `sudo usermod` |
