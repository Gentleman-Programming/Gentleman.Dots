package tui

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
)

type packageCommandCall struct {
	runner  string
	command string
}

func withPackageCommandMocks(t *testing.T, sudoErr error) *[]packageCommandCall {
	t.Helper()

	originalPkg := runPkgInstallWithLogs
	originalSudo := runSudoWithLogs
	originalBrew := runBrewWithLogs

	calls := []packageCommandCall{}

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

	t.Cleanup(func() {
		runPkgInstallWithLogs = originalPkg
		runSudoWithLogs = originalSudo
		runBrewWithLogs = originalBrew
	})

	return &calls
}

func TestInstallPlatformPackagesFedoraFallsBackToBrewWhenNativeFails(t *testing.T) {
	calls := withPackageCommandMocks(t, errors.New("dnf failed"))

	m := &Model{SystemInfo: &system.SystemInfo{OS: system.OSFedora, HasBrew: true}}
	result := installPlatformPackages(m, "shell", platformPackages{
		Brew:   "fish carapace zoxide atuin starship",
		Fedora: "fish carapace zoxide atuin starship",
	}, nil)

	if result.Error != nil {
		t.Fatalf("expected brew fallback to succeed, got error: %v", result.Error)
	}

	expected := []packageCommandCall{
		{runner: "sudo", command: "dnf install -y fish carapace zoxide atuin starship"},
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
		Brew: "fish carapace zoxide atuin starship",
		Arch: "fish carapace zoxide atuin starship",
	}, nil)

	if result.Error != nil {
		t.Fatalf("expected brew fallback to succeed, got error: %v", result.Error)
	}

	expected := []packageCommandCall{
		{runner: "sudo", command: "pacman -S --needed --noconfirm fish carapace zoxide atuin starship"},
		{runner: "brew", command: "install fish carapace zoxide atuin starship"},
	}
	if !reflect.DeepEqual(*calls, expected) {
		t.Fatalf("calls = %#v, want %#v", *calls, expected)
	}
}

func TestIsDnfNoOp(t *testing.T) {
	tests := []struct {
		name   string
		result *system.ExecResult
		want   bool
	}{
		{
			name:   "nil result",
			result: nil,
			want:   false,
		},
		{
			name: "success exit 0",
			result: &system.ExecResult{
				ExitCode: 0,
				Output:   "Complete!",
			},
			want: false,
		},
		{
			name: "exit 1 with stderr Nothing to do (dnf4)",
			result: &system.ExecResult{
				ExitCode: 1,
				Stderr:   "Module or Group 'development-tools' is not available.\nError: Nothing to do.",
			},
			want: true,
		},
		{
			name: "exit 1 with stdout/output Nothing to do",
			result: &system.ExecResult{
				ExitCode: 1,
				Output:   "Nothing to do.",
			},
			want: true,
		},
		{
			name: "exit 1 real failure",
			result: &system.ExecResult{
				ExitCode: 1,
				Stderr:   "Error: Unable to find a match: invalidpackage",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isDnfNoOp(tt.result); got != tt.want {
				t.Errorf("isDnfNoOp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstallPlatformPackagesFedoraTreatsDnfNoOpAsSuccess(t *testing.T) {
	originalSudo := runSudoWithLogs
	originalBrew := runBrewWithLogs

	sudoCalls := 0
	brewCalls := 0

	runSudoWithLogs = func(command string, opts *system.ExecOptions, onLog system.LogCallback) *system.ExecResult {
		sudoCalls++
		return &system.ExecResult{
			Command:  command,
			ExitCode: 1,
			Error:    errors.New("exit status 1"),
			Stderr:   "Error: Nothing to do.",
		}
	}
	runBrewWithLogs = func(args string, opts *system.ExecOptions, onLog system.LogCallback) *system.ExecResult {
		brewCalls++
		return &system.ExecResult{Command: args}
	}

	t.Cleanup(func() {
		runSudoWithLogs = originalSudo
		runBrewWithLogs = originalBrew
	})

	m := &Model{SystemInfo: &system.SystemInfo{OS: system.OSFedora, HasBrew: true}}
	result := installPlatformPackages(m, "shell", platformPackages{
		Brew:   "fish carapace zoxide atuin starship",
		Fedora: "fish carapace zoxide atuin starship",
	}, nil)

	if result.Error != nil {
		t.Fatalf("expected dnf no-op to be treated as success, got error: %v", result.Error)
	}
	if sudoCalls != 1 {
		t.Errorf("sudoCalls = %d, want 1", sudoCalls)
	}
	if brewCalls != 0 {
		t.Errorf("brewCalls = %d, want 0 (should not fall back to brew when dnf says nothing to do)", brewCalls)
	}
}
