package tui

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// aurOnlyArchPackages are package names that do not exist in the official Arch
// repositories (core/extra) and are only available through the AUR. Because
// pacman transactions are atomic, a single unresolvable name aborts the whole
// installation, so none of them may appear in an Arch package list.
var aurOnlyArchPackages = []string{
	"carapace",
	"zsh-theme-powerlevel10k",
}

var archFieldPattern = regexp.MustCompile(`(?m)^\s*Arch:\s*"([^"]*)",`)

// TestArchPackageListsHaveNoAUROnlyNames inspects the real Arch package list
// literals in installer.go. The other platform_packages tests build their
// fixtures inline, so they cannot catch a regression in the actual source
// strings; this test reads the source directly and is the regression guard for
// the arch-package-install spec.
func TestArchPackageListsHaveNoAUROnlyNames(t *testing.T) {
	source, err := os.ReadFile("installer.go")
	if err != nil {
		t.Fatalf("read installer.go: %v", err)
	}

	matches := archFieldPattern.FindAllStringSubmatch(string(source), -1)
	if len(matches) == 0 {
		t.Fatal("found no Arch package list literals in installer.go")
	}

	for _, match := range matches {
		list := match[1]
		for _, pkg := range strings.Fields(list) {
			for _, aurOnly := range aurOnlyArchPackages {
				if pkg == aurOnly {
					t.Errorf("Arch package list %q contains AUR-only package %q; pacman cannot resolve it and the atomic transaction aborts", list, aurOnly)
				}
			}
		}
	}
}

// TestArchPackageListsMatchExpected pins the exact expected content of every
// Arch package list so an unintended addition or removal fails loudly.
func TestArchPackageListsMatchExpected(t *testing.T) {
	// Every Arch package list in installer.go, in source order. All names were
	// verified against the official Arch package database (core/extra).
	expected := []string{
		"fish zoxide atuin starship",
		"zsh zoxide atuin zsh-autosuggestions zsh-syntax-highlighting zsh-autocomplete",
		"nushell zoxide atuin jq bash starship",
		"tmux",
		"zellij",
		"nodejs npm",
		"neovim git gcc fzf fd ripgrep coreutils bat curl lazygit tree-sitter",
	}

	source, err := os.ReadFile("installer.go")
	if err != nil {
		t.Fatalf("read installer.go: %v", err)
	}

	matches := archFieldPattern.FindAllStringSubmatch(string(source), -1)
	got := make([]string, 0, len(matches))
	for _, match := range matches {
		got = append(got, match[1])
	}

	if len(got) != len(expected) {
		t.Fatalf("expected %d Arch package lists, got %d: %q", len(expected), len(got), got)
	}
	for i := range expected {
		if got[i] != expected[i] {
			t.Errorf("Arch package list %d:\n got  %q\n want %q", i, got[i], expected[i])
		}
	}
}
