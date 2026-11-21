package tui

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
)

// setPacmanUnavailable makes the named packages unresolvable for Arch,
// emulating a repository set that does not provide them.
func setPacmanUnavailable(t *testing.T, names ...string) {
	t.Helper()

	blocked := make(map[string]bool, len(names))
	for _, name := range names {
		blocked[name] = true
	}

	original := pacmanPackageAvailable
	pacmanPackageAvailable = func(name string) bool { return !blocked[name] }
	t.Cleanup(func() { pacmanPackageAvailable = original })
}

// setAptUnavailable makes the named packages unresolvable for Debian/Ubuntu.
func setAptUnavailable(t *testing.T, names ...string) {
	t.Helper()

	blocked := make(map[string]bool, len(names))
	for _, name := range names {
		blocked[name] = true
	}

	original := aptPackageAvailable
	aptPackageAvailable = func(name string) bool { return !blocked[name] }
	t.Cleanup(func() { aptPackageAvailable = original })
}

func TestFilterAvailablePackagesKeepsOrder(t *testing.T) {
	ok, missing := filterAvailablePackages(
		[]string{"zsh", "carapace", "zoxide"},
		func(name string) bool { return name != "carapace" },
	)

	if !reflect.DeepEqual(ok, []string{"zsh", "zoxide"}) {
		t.Fatalf("available = %#v, want [zsh zoxide]", ok)
	}
	if !reflect.DeepEqual(missing, []string{"carapace"}) {
		t.Fatalf("missing = %#v, want [carapace]", missing)
	}
}

func TestInstallArchPackagesSkipsUnavailableNameWithoutAborting(t *testing.T) {
	calls := withPackageCommandMocks(t, nil)
	setPacmanUnavailable(t, "carapace")

	var warnings []string
	m := &Model{SystemInfo: &system.SystemInfo{OS: system.OSArch}}
	result := installPlatformPackages(m, "shell", platformPackages{
		Brew: "fish carapace zoxide atuin starship",
		Arch: "fish carapace zoxide atuin starship",
	}, func(line string) { warnings = append(warnings, line) })

	if result.Error != nil {
		t.Fatalf("expected install to succeed without carapace, got: %v", result.Error)
	}

	expected := []packageCommandCall{
		{runner: "sudo", command: "pacman -S --needed --noconfirm fish zoxide atuin starship"},
	}
	if !reflect.DeepEqual(*calls, expected) {
		t.Fatalf("calls = %#v, want %#v", *calls, expected)
	}

	if len(warnings) != 1 || !strings.Contains(warnings[0], "carapace") {
		t.Fatalf("warnings = %#v, want one warning naming carapace", warnings)
	}
}

// CachyOS packages zsh-theme-powerlevel10k while vanilla Arch does not, so the
// same shell profile must keep the theme where it exists and drop it where it
// does not, without any distro-specific branch.
func TestInstallArchPackagesKeepsPowerlevel10kWhenPackaged(t *testing.T) {
	calls := withPackageCommandMocks(t, nil)
	setPacmanUnavailable(t, "carapace")

	m := &Model{SystemInfo: &system.SystemInfo{OS: system.OSArch}}
	result := installPlatformPackages(m, "shell", platformPackages{
		Brew: "zsh carapace zoxide zsh-theme-powerlevel10k",
		Arch: "zsh carapace zoxide zsh-theme-powerlevel10k",
	}, nil)

	if result.Error != nil {
		t.Fatalf("expected install to succeed, got: %v", result.Error)
	}

	expected := []packageCommandCall{
		{runner: "sudo", command: "pacman -S --needed --noconfirm zsh zoxide zsh-theme-powerlevel10k"},
	}
	if !reflect.DeepEqual(*calls, expected) {
		t.Fatalf("calls = %#v, want %#v", *calls, expected)
	}
}

func TestInstallArchPackagesDropsPowerlevel10kOnVanillaArch(t *testing.T) {
	calls := withPackageCommandMocks(t, nil)
	setPacmanUnavailable(t, "carapace", "zsh-theme-powerlevel10k")

	m := &Model{SystemInfo: &system.SystemInfo{OS: system.OSArch}}
	result := installPlatformPackages(m, "shell", platformPackages{
		Brew: "zsh carapace zoxide zsh-theme-powerlevel10k",
		Arch: "zsh carapace zoxide zsh-theme-powerlevel10k",
	}, nil)

	if result.Error != nil {
		t.Fatalf("expected install to succeed, got: %v", result.Error)
	}

	expected := []packageCommandCall{
		{runner: "sudo", command: "pacman -S --needed --noconfirm zsh zoxide"},
	}
	if !reflect.DeepEqual(*calls, expected) {
		t.Fatalf("calls = %#v, want %#v", *calls, expected)
	}
}

func TestInstallArchPackagesFailsWhenNothingIsAvailable(t *testing.T) {
	calls := withPackageCommandMocks(t, nil)
	setPacmanUnavailable(t, "fish", "zoxide")

	m := &Model{SystemInfo: &system.SystemInfo{OS: system.OSArch}}
	result := installPlatformPackages(m, "shell", platformPackages{
		Arch: "fish zoxide",
	}, nil)

	if result.Error == nil {
		t.Fatal("expected an error when no requested package is available")
	}
	if len(*calls) != 0 {
		t.Fatalf("expected no package manager invocation, got %#v", *calls)
	}
}

func TestInstallArchPackagesFallsBackToBrewAfterFiltering(t *testing.T) {
	calls := withPackageCommandMocks(t, errors.New("pacman failed"))
	setPacmanUnavailable(t, "carapace")

	m := &Model{SystemInfo: &system.SystemInfo{OS: system.OSArch, HasBrew: true}}
	result := installPlatformPackages(m, "shell", platformPackages{
		Brew: "fish carapace zoxide",
		Arch: "fish carapace zoxide",
	}, nil)

	if result.Error != nil {
		t.Fatalf("expected brew fallback to succeed, got: %v", result.Error)
	}

	expected := []packageCommandCall{
		{runner: "sudo", command: "pacman -S --needed --noconfirm fish zoxide"},
		{runner: "brew", command: "install fish carapace zoxide"},
	}
	if !reflect.DeepEqual(*calls, expected) {
		t.Fatalf("calls = %#v, want %#v", *calls, expected)
	}
}

func TestPreserveExistingConfigBacksUpWhenNoGlobalBackupRan(t *testing.T) {
	fishDir := filepath.Join(t.TempDir(), "fish")
	if err := os.MkdirAll(fishDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	const original = "set -gx MINE 1\n"
	if err := os.WriteFile(filepath.Join(fishDir, "config.fish"), []byte(original), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	m := &Model{}
	if err := preserveExistingConfig(m, "shell", "Fish", fishDir); err != nil {
		t.Fatalf("preserveExistingConfig() error: %v", err)
	}

	backups, err := filepath.Glob(fishDir + ".backup-*")
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	if len(backups) != 1 {
		t.Fatalf("backups = %#v, want exactly one", backups)
	}

	data, err := os.ReadFile(filepath.Join(backups[0], "config.fish"))
	if err != nil {
		t.Fatalf("read backup: %v", err)
	}
	if string(data) != original {
		t.Fatalf("backup content = %q, want %q", string(data), original)
	}
}

func TestPreserveExistingConfigSkipsWhenGlobalBackupRan(t *testing.T) {
	fishDir := filepath.Join(t.TempDir(), "fish")
	if err := os.MkdirAll(fishDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	m := &Model{BackupDir: filepath.Join(t.TempDir(), "gentleman-backup")}
	if err := preserveExistingConfig(m, "shell", "Fish", fishDir); err != nil {
		t.Fatalf("preserveExistingConfig() error: %v", err)
	}

	backups, _ := filepath.Glob(fishDir + ".backup-*")
	if len(backups) != 0 {
		t.Fatalf("backups = %#v, want none when the global backup covers the config", backups)
	}
}

func TestPreserveExistingConfigSkipsUnderDryRun(t *testing.T) {
	t.Setenv("GENTLEMAN_DRY_RUN", "1")

	fishDir := filepath.Join(t.TempDir(), "fish")
	if err := os.MkdirAll(fishDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	m := &Model{}
	if err := preserveExistingConfig(m, "shell", "Fish", fishDir); err != nil {
		t.Fatalf("preserveExistingConfig() error: %v", err)
	}

	backups, _ := filepath.Glob(fishDir + ".backup-*")
	if len(backups) != 0 {
		t.Fatalf("backups = %#v, want none under dry-run", backups)
	}
}

// Ubuntu/Debian do not package starship, so every shell profile failed at this
// step (#215). The install must proceed with the names apt can resolve.
func TestInstallDebianPackagesSkipsUnavailableStarship(t *testing.T) {
	calls := withPackageCommandMocks(t, nil)
	setAptUnavailable(t, "starship")

	var warnings []string
	m := &Model{SystemInfo: &system.SystemInfo{OS: system.OSDebian}}
	result := installPlatformPackages(m, "shell", platformPackages{
		Brew:   "fish starship zoxide",
		Debian: "fish zoxide starship",
	}, func(line string) { warnings = append(warnings, line) })

	if result.Error != nil {
		t.Fatalf("expected install to succeed without starship, got: %v", result.Error)
	}

	expected := []packageCommandCall{
		{runner: "sudo", command: "apt-get install -y fish zoxide"},
	}
	if !reflect.DeepEqual(*calls, expected) {
		t.Fatalf("calls = %#v, want %#v", *calls, expected)
	}

	if len(warnings) != 1 || !strings.Contains(warnings[0], "starship") {
		t.Fatalf("warnings = %#v, want one warning naming starship", warnings)
	}
}

func TestInstallDebianPackagesFailsWhenNothingIsAvailable(t *testing.T) {
	calls := withPackageCommandMocks(t, nil)
	setAptUnavailable(t, "fish", "zoxide", "starship")

	m := &Model{SystemInfo: &system.SystemInfo{OS: system.OSDebian}}
	result := installPlatformPackages(m, "shell", platformPackages{
		Debian: "fish zoxide starship",
	}, nil)

	if result.Error == nil {
		t.Fatal("expected an error when no requested package is available")
	}
	if len(*calls) != 0 {
		t.Fatalf("expected no package manager invocation, got %#v", *calls)
	}
}
