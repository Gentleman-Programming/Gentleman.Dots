package tui

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
)

// carapaceReleaseVersion is the upstream static release used when the host's
// repositories do not provide carapace.
const carapaceReleaseVersion = "v1.8.2"

// carapaceChecksums maps GOARCH to the SHA256 of the matching
// carapace-bin_<version>_linux_<arch>.tar.gz release asset.
var carapaceChecksums = map[string]string{
	"amd64": "1b3dc4cab51361d3443400aefa4314ad8d44d8414bf4b1d18f41a6fe71b1b0e4",
	"arm64": "1e44575128e35554eec39e8b8d04aad0177ef74696ccd5be48f85cd2672a7126",
}

// pacmanPackageAvailable reports whether a package name resolves in the host's
// configured repositories. Arch and CachyOS expose different package sets for
// the same name (zsh-theme-powerlevel10k is AUR-only on Arch but packaged in
// CachyOS), so availability is queried instead of assumed.
var pacmanPackageAvailable = func(name string) bool {
	return system.Run("pacman -Si "+name, nil).Error == nil
}

// aptPackageAvailable reports whether a package name resolves in the host's
// configured Debian/Ubuntu repositories. `apt-cache show` exits non-zero when
// no package of that name is known.
var aptPackageAvailable = func(name string) bool {
	return system.Run("apt-cache show "+name, nil).Error == nil
}

// dnfPackageAvailable reports whether a package name resolves in the host's
// configured Fedora repositories. `dnf info` exits non-zero when no package
// of that name is known.
var dnfPackageAvailable = func(name string) bool {
	return system.Run("dnf info "+name, nil).Error == nil
}

// filterAvailablePackages splits pkgs into the names the host can install and
// the names it cannot, preserving the original order.
func filterAvailablePackages(pkgs []string, available func(string) bool) (ok []string, missing []string) {
	for _, name := range pkgs {
		if available(name) {
			ok = append(ok, name)
			continue
		}
		missing = append(missing, name)
	}
	return ok, missing
}

// installAvailablePackages installs only the names the host can actually
// resolve. A hardcoded distro list rots as soon as one package moves to the
// AUR, to a PPA, or out of the repositories entirely, and an atomic package
// transaction then installs nothing at all.
func installAvailablePackages(packages string, available func(string) bool, nativeCommand func([]string) string, brewPackages string, hasBrew bool, onLog func(string)) *system.ExecResult {
	requested := strings.Fields(packages)
	if len(requested) == 0 {
		return &system.ExecResult{Error: fmt.Errorf("no packages requested")}
	}

	ok, missing := filterAvailablePackages(requested, available)
	for _, name := range missing {
		if onLog != nil {
			onLog(fmt.Sprintf("Warning: %s is not available in the configured repositories; skipping it", name))
		}
	}
	if len(ok) == 0 {
		return &system.ExecResult{Error: fmt.Errorf("none of the requested packages are available: %s", packages)}
	}

	result := runSudoWithLogs(nativeCommand(ok), nil, onLog)
	if result.Error == nil || !hasBrew || brewPackages == "" {
		return result
	}

	return runBrewWithLogs("install "+brewPackages, nil, onLog)
}

// installArchPackages installs an Arch package list without letting a name that
// pacman cannot resolve abort the whole transaction.
func installArchPackages(archPackages, brewPackages string, hasBrew bool, onLog func(string)) *system.ExecResult {
	return installAvailablePackages(archPackages, pacmanPackageAvailable,
		func(pkgs []string) string {
			return "pacman -S --needed --noconfirm " + strings.Join(pkgs, " ")
		},
		brewPackages, hasBrew, onLog)
}

// installDebianPackages installs a Debian/Ubuntu package list without letting a
// name that apt cannot resolve abort the whole transaction. This is the same
// defect class as the Arch lists: `starship` is not packaged by Debian or
// Ubuntu, so every shell profile failed at this step.
func installDebianPackages(debianPackages, brewPackages string, hasBrew bool, onLog func(string)) *system.ExecResult {
	return installAvailablePackages(debianPackages, aptPackageAvailable,
		func(pkgs []string) string {
			return "apt-get install -y " + strings.Join(pkgs, " ")
		},
		brewPackages, hasBrew, onLog)
}

// installFedoraPackages installs a Fedora package list without letting a name
// that dnf cannot resolve (such as carapace) abort the whole transaction.
func installFedoraPackages(fedoraPackages, brewPackages string, hasBrew bool, onLog func(string)) *system.ExecResult {
	return installAvailablePackages(fedoraPackages, dnfPackageAvailable,
		func(pkgs []string) string {
			return "dnf install -y " + strings.Join(pkgs, " ")
		},
		brewPackages, hasBrew, onLog)
}

// installCarapaceBinary installs carapace from the official static release when
// the host repositories do not provide it. It mirrors installHerdrBinary:
// download, verify the pinned SHA256, then place the binary on PATH.
func installCarapaceBinary(stepID string) error {
	if system.CommandExists("carapace") {
		SendLog(stepID, "carapace already installed")
		return nil
	}

	if system.IsDryRun() {
		SendLog(stepID, fmt.Sprintf("DRY RUN: would install carapace %s from the official release", carapaceReleaseVersion))
		return nil
	}

	expectedSHA256, ok := carapaceChecksums[runtime.GOARCH]
	if !ok {
		return fmt.Errorf("unsupported carapace architecture: %s", runtime.GOARCH)
	}

	homeDir := os.Getenv("HOME")
	binDir := filepath.Join(homeDir, ".local", "bin")
	if err := system.EnsureDir(binDir); err != nil {
		return err
	}

	tmpDir, err := os.MkdirTemp("", "carapace-install-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	version := strings.TrimPrefix(carapaceReleaseVersion, "v")
	url := fmt.Sprintf(
		"https://github.com/carapace-sh/carapace-bin/releases/download/%s/carapace-bin_%s_linux_%s.tar.gz",
		carapaceReleaseVersion, version, runtime.GOARCH,
	)
	archive := filepath.Join(tmpDir, "carapace-bin.tar.gz")

	SendLog(stepID, "Downloading carapace release...")
	result := system.RunWithLogs(fmt.Sprintf("curl -fsSL %q -o %q", url, archive), nil, func(line string) {
		SendLog(stepID, line)
	})
	if result.Error != nil {
		return result.Error
	}

	data, err := os.ReadFile(archive)
	if err != nil {
		return err
	}
	actualSHA256 := sha256.Sum256(data)
	if hex.EncodeToString(actualSHA256[:]) != expectedSHA256 {
		return fmt.Errorf("carapace checksum mismatch for %s", url)
	}

	SendLog(stepID, "Extracting carapace...")
	result = system.RunWithLogs(fmt.Sprintf("tar -xzf %q -C %q carapace", archive, tmpDir), nil, func(line string) {
		SendLog(stepID, line)
	})
	if result.Error != nil {
		return result.Error
	}

	dest := filepath.Join(binDir, "carapace")
	if err := system.CopyFile(filepath.Join(tmpDir, "carapace"), dest); err != nil {
		return err
	}

	return os.Chmod(dest, 0o755)
}
