package system

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPatchZshForWM(t *testing.T) {
	// Sample .zshrc content similar to the real one
	zshrcContent := `# Enable Powerlevel10k instant prompt.
if [[ -r "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh" ]]; then
  source "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh"
fi

export ZSH="$HOME/.oh-my-zsh"

WM_VAR="/$TMUX"
# change with ZELLIJ
WM_CMD="tmux"
# change with zellij

function start_if_needed() {
    if [[ $- == *i* ]] && [[ -z "${WM_VAR#/}" ]] && [[ -t 1 ]]; then
        exec $WM_CMD
    fi
}

alias fzfbat='fzf --preview="bat --theme=gruvbox-dark --color=always {}"'

eval "$(fzf --zsh)"
eval "$(zoxide init zsh)"

start_if_needed
`

	tests := []struct {
		name           string
		wm             string
		installNvim    bool
		wantContain    []string
		wantNotContain []string
	}{
		{
			name:        "WM none should remove all WM-related lines",
			wm:          "none",
			installNvim: true,
			wantContain: []string{
				"eval \"$(zoxide init zsh)\"",
				"alias fzfbat",
			},
			wantNotContain: []string{
				"WM_VAR",
				"WM_CMD",
				"start_if_needed",
				"function start_if_needed",
				"exec $WM_CMD",
			},
		},
		{
			name:        "WM none without nvim should wrap fzf",
			wm:          "none",
			installNvim: false,
			wantContain: []string{
				"if command -v fzf",
				"eval \"$(fzf --zsh)\"",
				"fi",
			},
			wantNotContain: []string{
				"WM_VAR",
				"start_if_needed",
			},
		},
		{
			name:        "WM zellij should replace tmux with zellij",
			wm:          "zellij",
			installNvim: true,
			wantContain: []string{
				"WM_VAR=\"$ZELLIJ\"",
				"WM_CMD=\"zellij\"",
				"start_if_needed",
			},
			wantNotContain: []string{
				"WM_VAR=\"/$TMUX\"",
				"WM_CMD=\"tmux\"",
				"# change with ZELLIJ",
			},
		},
		{
			name:        "WM tmux should keep original content",
			wm:          "tmux",
			installNvim: true,
			wantContain: []string{
				"WM_VAR=\"/$TMUX\"",
				"WM_CMD=\"tmux\"",
				"start_if_needed",
			},
			wantNotContain: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			tmpDir := t.TempDir()
			zshrcPath := filepath.Join(tmpDir, ".zshrc")

			if err := os.WriteFile(zshrcPath, []byte(zshrcContent), 0644); err != nil {
				t.Fatalf("Failed to write temp file: %v", err)
			}

			// Run the patch
			if err := PatchZshForWM(zshrcPath, tt.wm, tt.installNvim); err != nil {
				t.Fatalf("PatchZshForWM failed: %v", err)
			}

			// Read result
			result, err := os.ReadFile(zshrcPath)
			if err != nil {
				t.Fatalf("Failed to read patched file: %v", err)
			}
			content := string(result)

			// Check expected content
			for _, want := range tt.wantContain {
				if !strings.Contains(content, want) {
					t.Errorf("Expected content to contain %q, but it didn't.\nContent:\n%s", want, content)
				}
			}

			// Check content that should not be present
			for _, notWant := range tt.wantNotContain {
				if strings.Contains(content, notWant) {
					t.Errorf("Expected content NOT to contain %q, but it did.\nContent:\n%s", notWant, content)
				}
			}
		})
	}
}

func TestPatchFishForWM(t *testing.T) {
	fishContent := `if status is-interactive
    # Commands to run in interactive sessions can go here
end

if test (uname) = Darwin
    if test -f /opt/homebrew/bin/brew
        set BREW_BIN /opt/homebrew/bin/brew
    end
end

if not set -q TMUX
    tmux
end

#if not set -q ZELLIJ 
#  zellij
#end

starship init fish | source
zoxide init fish | source
fzf --fish | source

alias fzfbat='fzf --preview="bat --theme=gruvbox-dark --color=always {}"'
`

	tests := []struct {
		name           string
		wm             string
		installNvim    bool
		wantContain    []string
		wantNotContain []string
	}{
		{
			name:        "WM none should remove tmux block",
			wm:          "none",
			installNvim: true,
			wantContain: []string{
				"starship init fish | source",
				"zoxide init fish | source",
			},
			wantNotContain: []string{
				"if not set -q TMUX",
				"tmux",
				"#if not set -q ZELLIJ",
			},
		},
		{
			name:        "WM none without nvim should wrap fzf",
			wm:          "none",
			installNvim: false,
			wantContain: []string{
				"if command -v fzf",
				"fzf --fish | source",
				"end",
			},
			wantNotContain: []string{
				"if not set -q TMUX",
			},
		},
		{
			name:        "WM zellij should uncomment zellij and remove tmux",
			wm:          "zellij",
			installNvim: true,
			wantContain: []string{
				"if not set -q ZELLIJ",
				"zellij",
			},
			wantNotContain: []string{
				"if not set -q TMUX",
				"#if not set -q ZELLIJ",
			},
		},
		{
			name:        "WM tmux should keep original",
			wm:          "tmux",
			installNvim: true,
			wantContain: []string{
				"if not set -q TMUX",
				"tmux",
			},
			wantNotContain: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			fishPath := filepath.Join(tmpDir, "config.fish")

			if err := os.WriteFile(fishPath, []byte(fishContent), 0644); err != nil {
				t.Fatalf("Failed to write temp file: %v", err)
			}

			if err := PatchFishForWM(fishPath, tt.wm, tt.installNvim); err != nil {
				t.Fatalf("PatchFishForWM failed: %v", err)
			}

			result, err := os.ReadFile(fishPath)
			if err != nil {
				t.Fatalf("Failed to read patched file: %v", err)
			}
			content := string(result)

			for _, want := range tt.wantContain {
				if !strings.Contains(content, want) {
					t.Errorf("Expected content to contain %q, but it didn't.\nContent:\n%s", want, content)
				}
			}

			for _, notWant := range tt.wantNotContain {
				if strings.Contains(content, notWant) {
					t.Errorf("Expected content NOT to contain %q, but it did.\nContent:\n%s", notWant, content)
				}
			}
		})
	}
}

func TestPatchNushellForWM(t *testing.T) {
	nushellContent := `$env.config = {
    show_banner: false
}

def fzfbat [] {
  fzf --preview "bat --theme=gruvbox-dark --color=always {}" 
}

source ~/.zoxide.nu
source ~/.cache/carapace/init.nu

let MULTIPLEXER = "tmux" 
let MULTIPLEXER_ENV_PREFIX = "TMUX"

def start_multiplexer [] {
  if $MULTIPLEXER_ENV_PREFIX not-in ($env | columns) {
    run-external $MULTIPLEXER
  }
}

start_multiplexer
`

	tests := []struct {
		name           string
		wm             string
		wantContain    []string
		wantNotContain []string
	}{
		{
			name: "WM none should remove multiplexer code",
			wm:   "none",
			wantContain: []string{
				"$env.config",
				"def fzfbat",
				"source ~/.zoxide.nu",
			},
			wantNotContain: []string{
				"let MULTIPLEXER",
				"def start_multiplexer",
				"start_multiplexer",
				"run-external",
			},
		},
		{
			name: "WM zellij should replace tmux with zellij",
			wm:   "zellij",
			wantContain: []string{
				"let MULTIPLEXER = \"zellij\"",
				"let MULTIPLEXER_ENV_PREFIX = \"ZELLIJ\"",
				"start_multiplexer",
			},
			wantNotContain: []string{
				"let MULTIPLEXER = \"tmux\"",
				"let MULTIPLEXER_ENV_PREFIX = \"TMUX\"",
			},
		},
		{
			name: "WM tmux should keep original",
			wm:   "tmux",
			wantContain: []string{
				"let MULTIPLEXER = \"tmux\"",
				"let MULTIPLEXER_ENV_PREFIX = \"TMUX\"",
				"start_multiplexer",
			},
			wantNotContain: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			nuPath := filepath.Join(tmpDir, "config.nu")

			if err := os.WriteFile(nuPath, []byte(nushellContent), 0644); err != nil {
				t.Fatalf("Failed to write temp file: %v", err)
			}

			if err := PatchNushellForWM(nuPath, tt.wm); err != nil {
				t.Fatalf("PatchNushellForWM failed: %v", err)
			}

			result, err := os.ReadFile(nuPath)
			if err != nil {
				t.Fatalf("Failed to read patched file: %v", err)
			}
			content := string(result)

			for _, want := range tt.wantContain {
				if !strings.Contains(content, want) {
					t.Errorf("Expected content to contain %q, but it didn't.\nContent:\n%s", want, content)
				}
			}

			for _, notWant := range tt.wantNotContain {
				if strings.Contains(content, notWant) {
					t.Errorf("Expected content NOT to contain %q, but it did.\nContent:\n%s", notWant, content)
				}
			}
		})
	}
}

func TestPatchZshForWM_FileNotFound(t *testing.T) {
	err := PatchZshForWM("/nonexistent/path/.zshrc", "none", true)
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestPatchFishForWM_FileNotFound(t *testing.T) {
	err := PatchFishForWM("/nonexistent/path/config.fish", "none", true)
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestPatchNushellForWM_FileNotFound(t *testing.T) {
	err := PatchNushellForWM("/nonexistent/path/config.nu", "none")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}
