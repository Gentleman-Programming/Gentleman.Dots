package system

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// zshTemplatePath resolves GentlemanZsh/.zshrc from this package directory.
func zshTemplatePath(t *testing.T) string {
	t.Helper()
	return filepath.Join("..", "..", "..", "GentlemanZsh", ".zshrc")
}

func readZshTemplate(t *testing.T) string {
	t.Helper()

	content, err := os.ReadFile(zshTemplatePath(t))
	if err != nil {
		t.Fatalf("read zsh template: %v", err)
	}
	return string(content)
}

// launcherPositions returns the 0-based line of the start_if_needed invocation
// and of the Powerlevel10k instant prompt guard, or -1 when absent.
func launcherPositions(lines []string) (launcher, instantPrompt int) {
	launcher, instantPrompt = -1, -1

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if launcher < 0 && trimmed == "start_if_needed" {
			launcher = i
		}
		if instantPrompt < 0 && strings.Contains(trimmed, "p10k-instant-prompt-") {
			instantPrompt = i
		}
	}
	return launcher, instantPrompt
}

// The Powerlevel10k instant prompt redirects stdout to a file while the rest
// of .zshrc runs, so `[[ -t 1 ]]` is false after it and the launcher's exec
// never fires. The launcher therefore has to come first (#210).
func TestZshTemplateStartsWMAboveInstantPrompt(t *testing.T) {
	launcher, instantPrompt := launcherPositions(strings.Split(readZshTemplate(t), "\n"))

	if launcher < 0 {
		t.Fatal("zsh template does not invoke start_if_needed")
	}
	if instantPrompt < 0 {
		t.Skip("zsh template no longer uses the Powerlevel10k instant prompt")
	}
	if launcher > instantPrompt {
		t.Errorf("start_if_needed (line %d) comes after the instant prompt (line %d): window manager autostart is dead",
			launcher+1, instantPrompt+1)
	}
}

// A duplicated launcher would run once before the instant prompt and once
// after it, where it can never exec. Exactly one definition and one call.
func TestZshTemplateDeclaresLauncherOnce(t *testing.T) {
	content := readZshTemplate(t)

	if got := strings.Count(content, "function start_if_needed()"); got != 1 {
		t.Errorf("function start_if_needed declarations = %d, want 1", got)
	}
	if got := strings.Count(content, "\nstart_if_needed\n"); got != 1 {
		t.Errorf("start_if_needed invocations = %d, want 1", got)
	}
}

// The installer patches .zshrc after copying it, so the patch must keep the
// launcher above the instant prompt for every supported window manager.
func TestPatchZshForWMKeepsLauncherAboveInstantPrompt(t *testing.T) {
	tests := []struct {
		name string
		wm   string
		want string
	}{
		{name: "tmux", wm: "tmux", want: `WM_CMD="tmux"`},
		{name: "zellij", wm: "zellij", want: `WM_CMD="zellij"`},
		{name: "herdr", wm: "herdr", want: `WM_CMD="herdr"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), ".zshrc")
			if err := os.WriteFile(path, []byte(readZshTemplate(t)), 0o644); err != nil {
				t.Fatalf("write temp zshrc: %v", err)
			}

			if err := PatchZshForWM(path, tt.wm, true); err != nil {
				t.Fatalf("PatchZshForWM(%q): %v", tt.wm, err)
			}

			patched, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read patched zshrc: %v", err)
			}

			content := string(patched)
			launcher, instantPrompt := launcherPositions(strings.Split(content, "\n"))
			if launcher < 0 {
				t.Fatalf("patched zshrc lost the start_if_needed invocation")
			}
			if launcher > instantPrompt {
				t.Errorf("patched zshrc moved start_if_needed (line %d) below the instant prompt (line %d)",
					launcher+1, instantPrompt+1)
			}
			if !strings.Contains(content, tt.want) {
				t.Errorf("patched zshrc does not contain %q", tt.want)
			}
		})
	}
}

// Selecting no window manager must remove the launcher entirely, including
// the block now placed above the instant prompt.
func TestPatchZshForWMNoneRemovesLauncher(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".zshrc")
	if err := os.WriteFile(path, []byte(readZshTemplate(t)), 0o644); err != nil {
		t.Fatalf("write temp zshrc: %v", err)
	}

	if err := PatchZshForWM(path, "none", true); err != nil {
		t.Fatalf("PatchZshForWM(none): %v", err)
	}

	patched, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read patched zshrc: %v", err)
	}

	content := string(patched)
	for _, unwanted := range []string{"start_if_needed", "WM_CMD", "WM_VAR"} {
		if strings.Contains(content, unwanted) {
			t.Errorf("patched zshrc still contains %q after choosing no window manager", unwanted)
		}
	}
}

// User-maintained customizations must be loaded from ~/.zshrc.d/*.zsh in lexical
// order so they survive installer updates (#206).
func TestZshTemplateLoadsUserExtensionsDirectory(t *testing.T) {
	content := readZshTemplate(t)

	requiredSnippets := []string{
		`"$HOME"/.zshrc.d/*.zsh(N)`,
		`[[ -r "$config_file" ]] && source "$config_file"`,
		`unset config_file`,
	}

	for _, snippet := range requiredSnippets {
		if !strings.Contains(content, snippet) {
			t.Errorf("zsh template missing required extension loader snippet: %q", snippet)
		}
	}

	// Extensions must load after core integrations and p10k.
	p10kPos := strings.Index(content, "source ~/.p10k.zsh")
	extPos := strings.Index(content, `"$HOME"/.zshrc.d/*.zsh(N)`)
	if p10kPos >= 0 && extPos >= 0 && extPos < p10kPos {
		t.Errorf("user extension loader (pos %d) must appear after p10k customization (pos %d)", extPos, p10kPos)
	}
}
