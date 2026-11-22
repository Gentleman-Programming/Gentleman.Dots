package system

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestTapFormulaSync(t *testing.T) {
	formulaPath := filepath.Join("..", "..", "..", "homebrew-tap", "Formula", "gentleman-dots.rb")
	mainGoPath := filepath.Join("..", "..", "cmd", "gentleman-installer", "main.go")

	formulaBytes, err := os.ReadFile(formulaPath)
	if err != nil {
		t.Fatalf("read formula %s: %v", formulaPath, err)
	}
	formula := string(formulaBytes)

	mainBytes, err := os.ReadFile(mainGoPath)
	if err != nil {
		t.Fatalf("read main.go %s: %v", mainGoPath, err)
	}
	mainSrc := string(mainBytes)

	// Extract version "..." from formula
	versionRe := regexp.MustCompile(`(?m)^\s*version\s+"([^"]+)"`)
	versionMatch := versionRe.FindStringSubmatch(formula)
	if len(versionMatch) < 2 {
		t.Fatalf("could not find version declaration in %s", formulaPath)
	}
	formulaVersion := versionMatch[1]

	const expectedVersion = "2.12.2"
	semverRe := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	if !semverRe.MatchString(formulaVersion) {
		t.Errorf("formula version %q is not valid semver", formulaVersion)
	}
	if formulaVersion != expectedVersion {
		t.Errorf("formula version = %q, want %q", formulaVersion, expectedVersion)
	}

	// Extract var Version = "..." from main.go; release builds inject Version via ldflags,
	// so "dev" is expected in source, while any explicit version must match the formula.
	mainVersionRe := regexp.MustCompile(`(?m)^\s*var\s+Version\s*=\s*"([^"]+)"`)
	mainMatch := mainVersionRe.FindStringSubmatch(mainSrc)
	if len(mainMatch) < 2 {
		t.Fatalf("could not find var Version in %s", mainGoPath)
	}
	mainVersion := mainMatch[1]
	if mainVersion != "dev" && mainVersion != formulaVersion {
		t.Errorf("main.go Version = %q does not match formula version %q", mainVersion, formulaVersion)
	}

	// Validate URL and 64-character lowercase hex sha256 for all 4 platforms
	platforms := []string{
		"darwin-arm64",
		"darwin-amd64",
		"linux-arm64",
		"linux-amd64",
	}

	sha256Re := regexp.MustCompile(`^[0-9a-f]{64}$`)

	for _, platform := range platforms {
		t.Run(platform, func(t *testing.T) {
			blockRe := regexp.MustCompile(
				`url\s+"(https://github\.com/Gentleman-Programming/Gentleman\.Dots/releases/download/v#\{version\}/gentleman-installer-` +
					regexp.QuoteMeta(platform) +
					`)"\s*\n\s*sha256\s+"([^"]+)"`,
			)
			m := blockRe.FindStringSubmatch(formula)
			if len(m) < 3 {
				t.Fatalf("missing url+sha256 block for platform %s in %s", platform, formulaPath)
			}

			url := m[1]
			checksum := m[2]

			if !strings.Contains(url, `v#{version}/gentleman-installer-`+platform) {
				t.Errorf("unexpected URL %q for platform %s", url, platform)
			}
			if !sha256Re.MatchString(checksum) {
				t.Errorf("sha256 %q for %s is not a valid 64-char lowercase hex string", checksum, platform)
			}
			if checksum == strings.Repeat("0", 64) {
				t.Errorf("sha256 for %s is an all-zero placeholder", platform)
			}
		})
	}
}
