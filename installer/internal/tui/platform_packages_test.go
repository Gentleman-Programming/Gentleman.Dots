package tui

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
)

type packageCommandCall struct {
	runner  string
	command string
}

type packageCheckResult int

const (
	checkExists packageCheckResult = iota
	checkMissing
	checkInfraFail
	checkInfraExit
)

func withPackageCommandMocks(t *testing.T, sudoErr error) *[]packageCommandCall {
	t.Helper()
	return withPackageCommandMocksAndCheck(t, sudoErr, nil)
}

func withPackageCommandMocksAndCheck(t *testing.T, sudoErr error, check func(string) packageCheckResult) *[]packageCommandCall {
	t.Helper()

	originalPkg := runPkgInstallWithLogs
	originalSudo := runSudoWithLogs
	originalBrew := runBrewWithLogs
	originalCheck := runCheckCommand
	originalGitHub := installGitHubPackageFn

	calls := []packageCommandCall{}

	if check == nil {
		check = func(string) packageCheckResult { return checkExists }
	}

	runPkgInstallWithLogs = func(packages string, opts *system.ExecOptions, onLog func(string)) *system.ExecResult {
		calls = append(calls, packageCommandCall{runner: "pkg", command: packages})
		return &system.ExecResult{Command: packages}
	}
	runSudoWithLogs = func(command string, opts *system.ExecOptions, onLog system.LogCallback) *system.ExecResult {
		calls = append(calls, packageCommandCall{runner: "sudo", command: command})
		return &system.ExecResult{Command: command, Error: sudoErr}
	}
	runBrewWithLogs = func(args string, opts *system.ExecOptions, onLog system.LogCallback) *system.ExecResult {
		calls = append(calls, packageCommandCall{runner: "brew", command: args})
		return &system.ExecResult{Command: args}
	}
	runCheckCommand = func(command string, opts *system.ExecOptions) *system.ExecResult {
		calls = append(calls, packageCommandCall{runner: "check", command: command})
		pkg := checkCommandPackage(command)
		switch check(pkg) {
		case checkMissing:
			// Model the real system.Run contract: a non-zero exit always
			// carries a non-nil Error.
			if strings.HasPrefix(command, "LC_ALL=C pacman ") || strings.HasPrefix(command, "pacman ") {
				return &system.ExecResult{Command: command, ExitCode: 1, Error: errors.New("exit status 1"), Stderr: "error: target not found: " + pkg}
			}
			if strings.HasPrefix(command, "LC_ALL=C dnf ") || strings.HasPrefix(command, "dnf ") {
				return &system.ExecResult{Command: command, ExitCode: 1, Error: errors.New("exit status 1"), Stderr: "No match for argument: " + pkg}
			}
			return &system.ExecResult{Command: command, ExitCode: 100, Error: errors.New("exit status 100")}
		case checkInfraFail:
			// Start failure: Error is set but the process never ran, so
			// ExitCode stays 0.
			return &system.ExecResult{Command: command, Error: errors.New("could not query package database")}
		case checkInfraExit:
			// The manager ran but failed without the "not found" marker
			// (locked database, missing metadata, no network).
			return &system.ExecResult{Command: command, ExitCode: 1, Error: errors.New("exit status 1"), Stderr: "error: failed to synchronize all databases"}
		default:
			return &system.ExecResult{Command: command, ExitCode: 0, Output: "Name : " + pkg}
		}
	}
	installGitHubPackageFn = func(m *Model, stepID, pkg string, onLog func(string)) error {
		calls = append(calls, packageCommandCall{runner: "github", command: pkg})
		return nil
	}

	t.Cleanup(func() {
		runPkgInstallWithLogs = originalPkg
		runSudoWithLogs = originalSudo
		runBrewWithLogs = originalBrew
		runCheckCommand = originalCheck
		installGitHubPackageFn = originalGitHub
	})

	return &calls
}

func checkCommandPackage(command string) string {
	fields := strings.Fields(command)
	return fields[len(fields)-1]
}

func TestInstallPlatformPackagesFedoraFallsBackToBrewWhenNativeFails(t *testing.T) {
	calls := withPackageCommandMocks(t, errors.New("dnf failed"))

	m := &Model{SystemInfo: &system.SystemInfo{OS: system.OSFedora, HasBrew: true}}
	result := installPlatformPackages(m, "shell", platformPackages{
		Brew:   "fish carapace zoxide atuin starship",
		Fedora: "fish zoxide atuin starship",
		GitHub: "carapace",
	}, nil)

	if result.Error != nil {
		t.Fatalf("expected brew fallback to succeed, got error: %v", result.Error)
	}

	expected := []packageCommandCall{
		{runner: "check", command: "LC_ALL=C dnf info fish"},
		{runner: "check", command: "LC_ALL=C dnf info zoxide"},
		{runner: "check", command: "LC_ALL=C dnf info atuin"},
		{runner: "check", command: "LC_ALL=C dnf info starship"},
		{runner: "check", command: "LC_ALL=C dnf info carapace"},
		{runner: "sudo", command: "dnf install -y fish zoxide atuin starship carapace"},
		{runner: "brew", command: "install fish carapace zoxide atuin starship"},
	}
	if !reflect.DeepEqual(*calls, expected) {
		t.Fatalf("calls = %#v, want %#v", *calls, expected)
	}
}

func TestInstallPlatformPackagesDebianWithBrewUsesBrewDirectly(t *testing.T) {
	calls := withPackageCommandMocks(t, nil)

	m := &Model{SystemInfo: &system.SystemInfo{OS: system.OSDebian, HasBrew: true}}
	result := installPlatformPackages(m, "shell", platformPackages{
		Brew:   "fish carapace zoxide atuin starship",
		Debian: "fish zoxide starship",
	}, nil)

	if result.Error != nil {
		t.Fatalf("expected brew install to succeed, got error: %v", result.Error)
	}

	expected := []packageCommandCall{
		{runner: "brew", command: "install fish carapace zoxide atuin starship"},
	}
	if !reflect.DeepEqual(*calls, expected) {
		t.Fatalf("calls = %#v, want %#v", *calls, expected)
	}
}

func TestInstallPlatformPackagesDebianWithoutBrewUsesApt(t *testing.T) {
	calls := withPackageCommandMocks(t, nil)

	m := &Model{SystemInfo: &system.SystemInfo{OS: system.OSDebian, HasBrew: false}}
	result := installPlatformPackages(m, "shell", platformPackages{
		Brew:   "fish carapace zoxide atuin starship",
		Debian: "fish zoxide starship",
	}, nil)

	if result.Error != nil {
		t.Fatalf("expected apt install to succeed, got error: %v", result.Error)
	}

	expected := []packageCommandCall{
		{runner: "check", command: "apt-cache show fish"},
		{runner: "check", command: "apt-cache show zoxide"},
		{runner: "check", command: "apt-cache show starship"},
		{runner: "sudo", command: "apt-get install -y fish zoxide starship"},
	}
	if !reflect.DeepEqual(*calls, expected) {
		t.Fatalf("calls = %#v, want %#v", *calls, expected)
	}
}

func TestInstallPlatformPackagesArchFallsBackToBrewWhenNativeFails(t *testing.T) {
	calls := withPackageCommandMocks(t, errors.New("pacman failed"))

	m := &Model{SystemInfo: &system.SystemInfo{OS: system.OSArch, HasBrew: true}}
	result := installPlatformPackages(m, "shell", platformPackages{
		Brew:   "fish carapace zoxide atuin starship",
		Arch:   "fish zoxide atuin starship",
		GitHub: "carapace",
	}, nil)

	if result.Error != nil {
		t.Fatalf("expected brew fallback to succeed, got error: %v", result.Error)
	}

	expected := []packageCommandCall{
		{runner: "check", command: "LC_ALL=C pacman -Si fish"},
		{runner: "check", command: "LC_ALL=C pacman -Si zoxide"},
		{runner: "check", command: "LC_ALL=C pacman -Si atuin"},
		{runner: "check", command: "LC_ALL=C pacman -Si starship"},
		{runner: "check", command: "LC_ALL=C pacman -Si carapace"},
		{runner: "sudo", command: "pacman -S --needed --noconfirm fish zoxide atuin starship carapace"},
		{runner: "brew", command: "install fish carapace zoxide atuin starship"},
	}
	if !reflect.DeepEqual(*calls, expected) {
		t.Fatalf("calls = %#v, want %#v", *calls, expected)
	}
}

func TestInstallPlatformPackagesArchMissingPackageInstallsFromGitHub(t *testing.T) {
	calls := withPackageCommandMocksAndCheck(t, nil, func(pkg string) packageCheckResult {
		if pkg == "carapace" {
			return checkMissing
		}
		return checkExists
	})

	m := &Model{SystemInfo: &system.SystemInfo{OS: system.OSArch, HasBrew: true}}
	result := installPlatformPackages(m, "shell", platformPackages{
		Brew:   "fish carapace zoxide atuin starship",
		Arch:   "fish zoxide atuin starship",
		GitHub: "carapace",
	}, nil)

	if result.Error != nil {
		t.Fatalf("expected success, got error: %v", result.Error)
	}

	// Availability checks run over the combined native+GitHub list (carapace
	// is probed even though pacman is not expected to have it), the native
	// command receives the filtered list, and the missing package is handed
	// to the GitHub installer.
	expected := []packageCommandCall{
		{runner: "check", command: "LC_ALL=C pacman -Si fish"},
		{runner: "check", command: "LC_ALL=C pacman -Si zoxide"},
		{runner: "check", command: "LC_ALL=C pacman -Si atuin"},
		{runner: "check", command: "LC_ALL=C pacman -Si starship"},
		{runner: "check", command: "LC_ALL=C pacman -Si carapace"},
		{runner: "sudo", command: "pacman -S --needed --noconfirm fish zoxide atuin starship"},
		{runner: "github", command: "carapace"},
	}
	if !reflect.DeepEqual(*calls, expected) {
		t.Fatalf("calls = %#v, want %#v", *calls, expected)
	}
}

func TestInstallPlatformPackagesArchMissingPackageWithoutGitHubFallbackFails(t *testing.T) {
	calls := withPackageCommandMocksAndCheck(t, nil, func(pkg string) packageCheckResult {
		if pkg == "atuin" {
			return checkMissing
		}
		return checkExists
	})

	m := &Model{SystemInfo: &system.SystemInfo{OS: system.OSArch, HasBrew: true}}
	result := installPlatformPackages(m, "shell", platformPackages{
		Brew:   "fish carapace zoxide atuin starship",
		Arch:   "fish zoxide atuin starship",
		GitHub: "carapace",
	}, nil)

	if result.Error == nil {
		t.Fatal("expected error for package without GitHub fallback")
	}
	if !strings.Contains(result.Error.Error(), "atuin") {
		t.Errorf("error should mention the package name, got: %v", result.Error)
	}

	for _, call := range *calls {
		if call.runner == "sudo" || call.runner == "github" {
			t.Errorf("unexpected call after missing package error: %#v", call)
		}
	}
}

func TestInstallPlatformPackagesCheckInfrastructureFailureDegradesToNative(t *testing.T) {
	calls := withPackageCommandMocksAndCheck(t, nil, func(pkg string) packageCheckResult {
		if pkg == "carapace" {
			return checkInfraFail
		}
		return checkExists
	})

	m := &Model{SystemInfo: &system.SystemInfo{OS: system.OSArch, HasBrew: true}}
	result := installPlatformPackages(m, "shell", platformPackages{
		Brew:   "fish carapace zoxide atuin starship",
		Arch:   "fish zoxide atuin starship",
		GitHub: "carapace",
	}, nil)

	if result.Error != nil {
		t.Fatalf("expected degraded native install to succeed, got error: %v", result.Error)
	}

	// When the query itself fails (the manager process never started), the
	// full native list goes to the manager and no GitHub installer runs.
	expected := []packageCommandCall{
		{runner: "check", command: "LC_ALL=C pacman -Si fish"},
		{runner: "check", command: "LC_ALL=C pacman -Si zoxide"},
		{runner: "check", command: "LC_ALL=C pacman -Si atuin"},
		{runner: "check", command: "LC_ALL=C pacman -Si starship"},
		{runner: "check", command: "LC_ALL=C pacman -Si carapace"},
		{runner: "sudo", command: "pacman -S --needed --noconfirm fish zoxide atuin starship"},
	}
	if !reflect.DeepEqual(*calls, expected) {
		t.Fatalf("calls = %#v, want %#v", *calls, expected)
	}
}

func TestInstallPlatformPackagesArchInfraExitWithoutPatternDegradesToNative(t *testing.T) {
	calls := withPackageCommandMocksAndCheck(t, nil, func(pkg string) packageCheckResult {
		if pkg == "carapace" {
			return checkInfraExit
		}
		return checkExists
	})

	m := &Model{SystemInfo: &system.SystemInfo{OS: system.OSArch, HasBrew: true}}
	result := installPlatformPackages(m, "shell", platformPackages{
		Brew:   "fish carapace zoxide atuin starship",
		Arch:   "fish zoxide atuin starship",
		GitHub: "carapace",
	}, nil)

	if result.Error != nil {
		t.Fatalf("expected degraded native install to succeed, got error: %v", result.Error)
	}

	// A non-zero exit WITHOUT the "target not found" marker is an
	// infrastructure failure, not a missing package: degrade to the full
	// native list and never run the GitHub installer.
	expected := []packageCommandCall{
		{runner: "check", command: "LC_ALL=C pacman -Si fish"},
		{runner: "check", command: "LC_ALL=C pacman -Si zoxide"},
		{runner: "check", command: "LC_ALL=C pacman -Si atuin"},
		{runner: "check", command: "LC_ALL=C pacman -Si starship"},
		{runner: "check", command: "LC_ALL=C pacman -Si carapace"},
		{runner: "sudo", command: "pacman -S --needed --noconfirm fish zoxide atuin starship"},
	}
	if !reflect.DeepEqual(*calls, expected) {
		t.Fatalf("calls = %#v, want %#v", *calls, expected)
	}
}

func TestInstallPlatformPackagesCachyosPowerlevel10kNativeCarapaceGitHub(t *testing.T) {
	calls := withPackageCommandMocksAndCheck(t, nil, func(pkg string) packageCheckResult {
		if pkg == "carapace" {
			return checkMissing
		}
		return checkExists
	})

	m := &Model{SystemInfo: &system.SystemInfo{OS: system.OSArch, HasBrew: true}}
	result := installPlatformPackages(m, "shell", platformPackages{
		Brew:   "zsh carapace zoxide atuin zsh-autosuggestions zsh-syntax-highlighting zsh-autocomplete powerlevel10k",
		Arch:   "zsh zoxide atuin zsh-autosuggestions zsh-syntax-highlighting zsh-autocomplete",
		GitHub: "carapace zsh-theme-powerlevel10k",
	}, nil)

	if result.Error != nil {
		t.Fatalf("expected success, got error: %v", result.Error)
	}

	// CachyOS-style repos provide zsh-theme-powerlevel10k natively, so it
	// lands in `existing` and is installed via pacman; carapace is missing
	// and is routed to its GitHub installer.
	expected := []packageCommandCall{
		{runner: "check", command: "LC_ALL=C pacman -Si zsh"},
		{runner: "check", command: "LC_ALL=C pacman -Si zoxide"},
		{runner: "check", command: "LC_ALL=C pacman -Si atuin"},
		{runner: "check", command: "LC_ALL=C pacman -Si zsh-autosuggestions"},
		{runner: "check", command: "LC_ALL=C pacman -Si zsh-syntax-highlighting"},
		{runner: "check", command: "LC_ALL=C pacman -Si zsh-autocomplete"},
		{runner: "check", command: "LC_ALL=C pacman -Si carapace"},
		{runner: "check", command: "LC_ALL=C pacman -Si zsh-theme-powerlevel10k"},
		{runner: "sudo", command: "pacman -S --needed --noconfirm zsh zoxide atuin zsh-autosuggestions zsh-syntax-highlighting zsh-autocomplete zsh-theme-powerlevel10k"},
		{runner: "github", command: "carapace"},
	}
	if !reflect.DeepEqual(*calls, expected) {
		t.Fatalf("calls = %#v, want %#v", *calls, expected)
	}
}

func TestInstallPlatformPackagesDebianMissingPackageInstallsFromGitHub(t *testing.T) {
	calls := withPackageCommandMocksAndCheck(t, nil, func(pkg string) packageCheckResult {
		if pkg == "carapace" {
			return checkMissing
		}
		return checkExists
	})

	m := &Model{SystemInfo: &system.SystemInfo{OS: system.OSDebian, HasBrew: false}}
	result := installPlatformPackages(m, "shell", platformPackages{
		Debian: "fish zoxide starship",
		GitHub: "carapace",
	}, nil)

	if result.Error != nil {
		t.Fatalf("expected success, got error: %v", result.Error)
	}

	expected := []packageCommandCall{
		{runner: "check", command: "apt-cache show fish"},
		{runner: "check", command: "apt-cache show zoxide"},
		{runner: "check", command: "apt-cache show starship"},
		{runner: "check", command: "apt-cache show carapace"},
		{runner: "sudo", command: "apt-get install -y fish zoxide starship"},
		{runner: "github", command: "carapace"},
	}
	if !reflect.DeepEqual(*calls, expected) {
		t.Fatalf("calls = %#v, want %#v", *calls, expected)
	}
}

func TestInstallPowerlevel10k(t *testing.T) {
	if _, err := os.Stat("/usr/share/powerlevel10k/powerlevel10k.zsh-theme"); err == nil {
		t.Skip("Powerlevel10k already installed on this system")
	}

	originalSudo := runSudoWithLogs
	defer func() { runSudoWithLogs = originalSudo }()

	var gotCommand string
	runSudoWithLogs = func(command string, opts *system.ExecOptions, onLog system.LogCallback) *system.ExecResult {
		gotCommand = command
		return &system.ExecResult{Command: command}
	}

	m := &Model{SystemInfo: &system.SystemInfo{}}
	if err := installPowerlevel10k(m, "shell", nil); err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	want := "git clone --depth 1 https://github.com/romkatv/powerlevel10k /usr/share/powerlevel10k"
	if gotCommand != want {
		t.Errorf("command = %q, want %q", gotCommand, want)
	}
}

func TestExtractTarGzBinary(t *testing.T) {
	t.Run("extracts the binary entry", func(t *testing.T) {
		dir := t.TempDir()
		tarGzPath := filepath.Join(dir, "release.tar.gz")
		writeTestTarGz(t, tarGzPath)

		dest := filepath.Join(dir, "carapace")
		if err := extractTarGzBinary(tarGzPath, "carapace", dest); err != nil {
			t.Fatalf("extract failed: %v", err)
		}

		data, err := os.ReadFile(dest)
		if err != nil {
			t.Fatalf("failed to read extracted binary: %v", err)
		}
		if string(data) != "binary-content" {
			t.Errorf("extracted content = %q, want %q", string(data), "binary-content")
		}
	})

	t.Run("fails when the entry is missing", func(t *testing.T) {
		dir := t.TempDir()
		tarGzPath := filepath.Join(dir, "release.tar.gz")
		writeTestTarGz(t, tarGzPath)

		err := extractTarGzBinary(tarGzPath, "missing", filepath.Join(dir, "out"))
		if err == nil {
			t.Fatal("expected error for missing entry")
		}
	})
}

func writeTestTarGz(t *testing.T, path string) {
	t.Helper()

	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create archive: %v", err)
	}
	defer f.Close()

	gz := gzip.NewWriter(f)
	defer gz.Close()

	tw := tar.NewWriter(gz)
	defer tw.Close()

	entries := []struct {
		name    string
		content string
	}{
		{"LICENSE", "MIT\n"},
		{"carapace", "binary-content"},
		{"README.md", "carapace release\n"},
	}
	for _, entry := range entries {
		hdr := &tar.Header{
			Name: entry.name,
			Mode: 0755,
			Size: int64(len(entry.content)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatalf("failed to write header: %v", err)
		}
		if _, err := tw.Write([]byte(entry.content)); err != nil {
			t.Fatalf("failed to write entry: %v", err)
		}
	}
}
