package tui

import (
	"os"
	"path/filepath"
	"testing"
)

// writeTestFile creates a file with its parent directories.
func writeTestFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestInstallOhMyZshCopiesBundledSnapshotWhenMissing(t *testing.T) {
	root := t.TempDir()
	repoDir := filepath.Join(root, "repo")
	homeDir := filepath.Join(root, "home")

	bundled := filepath.Join(repoDir, "GentlemanZsh", ".oh-my-zsh")
	writeTestFile(t, filepath.Join(bundled, "lib", "git.zsh"), "bundled\n")

	if err := installOhMyZshIfMissing("shell", repoDir, homeDir); err != nil {
		t.Fatalf("installOhMyZshIfMissing() error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(homeDir, ".oh-my-zsh", "lib", "git.zsh"))
	if err != nil {
		t.Fatalf("expected bundled snapshot to be copied: %v", err)
	}
	if string(data) != "bundled\n" {
		t.Fatalf("content = %q, want %q", string(data), "bundled\n")
	}
}

// The reported defect: writing the snapshot over an existing checkout marks
// tracked files as modified and breaks `omz update` on the autostash pop.
func TestInstallOhMyZshLeavesExistingCheckoutUntouched(t *testing.T) {
	root := t.TempDir()
	repoDir := filepath.Join(root, "repo")
	homeDir := filepath.Join(root, "home")

	bundled := filepath.Join(repoDir, "GentlemanZsh", ".oh-my-zsh")
	writeTestFile(t, filepath.Join(bundled, "lib", "git.zsh"), "bundled\n")
	writeTestFile(t, filepath.Join(bundled, "upstream-only.zsh"), "new upstream file\n")

	existing := filepath.Join(homeDir, ".oh-my-zsh")
	writeTestFile(t, filepath.Join(existing, "lib", "git.zsh"), "user checkout\n")
	writeTestFile(t, filepath.Join(existing, ".git", "HEAD"), "ref: refs/heads/master\n")

	if err := installOhMyZshIfMissing("shell", repoDir, homeDir); err != nil {
		t.Fatalf("installOhMyZshIfMissing() error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(existing, "lib", "git.zsh"))
	if err != nil {
		t.Fatalf("read existing checkout: %v", err)
	}
	if string(data) != "user checkout\n" {
		t.Fatalf("tracked file was overwritten: got %q, want %q", string(data), "user checkout\n")
	}

	if _, err := os.Stat(filepath.Join(existing, "upstream-only.zsh")); !os.IsNotExist(err) {
		t.Fatalf("expected no files from the bundled snapshot to be added, stat err = %v", err)
	}
}

func TestInstallOhMyZshTreatsSymlinkedCheckoutAsPresent(t *testing.T) {
	root := t.TempDir()
	repoDir := filepath.Join(root, "repo")
	homeDir := filepath.Join(root, "home")

	bundled := filepath.Join(repoDir, "GentlemanZsh", ".oh-my-zsh")
	writeTestFile(t, filepath.Join(bundled, "lib", "git.zsh"), "bundled\n")

	realDir := filepath.Join(root, "dotfiles", "oh-my-zsh")
	writeTestFile(t, filepath.Join(realDir, "lib", "git.zsh"), "symlinked checkout\n")

	if err := os.MkdirAll(homeDir, 0o755); err != nil {
		t.Fatalf("mkdir home: %v", err)
	}
	if err := os.Symlink(realDir, filepath.Join(homeDir, ".oh-my-zsh")); err != nil {
		t.Fatalf("symlink: %v", err)
	}

	if err := installOhMyZshIfMissing("shell", repoDir, homeDir); err != nil {
		t.Fatalf("installOhMyZshIfMissing() error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(realDir, "lib", "git.zsh"))
	if err != nil {
		t.Fatalf("read symlinked checkout: %v", err)
	}
	if string(data) != "symlinked checkout\n" {
		t.Fatalf("symlinked checkout was overwritten: got %q", string(data))
	}
}

func TestPreserveExistingConfigFileCopiesRegularFile(t *testing.T) {
	homeDir := t.TempDir()
	zshrc := filepath.Join(homeDir, ".zshrc")
	const original = "export MY_CUSTOM_VAR=gentleman\n"
	if err := os.WriteFile(zshrc, []byte(original), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	m := &Model{}
	if err := preserveExistingConfig(m, "shell", "Zsh (.zshrc)", zshrc); err != nil {
		t.Fatalf("preserveExistingConfig() error: %v", err)
	}

	backups, err := filepath.Glob(zshrc + ".backup-*")
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	if len(backups) != 1 {
		t.Fatalf("backups = %#v, want exactly one", backups)
	}

	data, err := os.ReadFile(backups[0])
	if err != nil {
		t.Fatalf("read backup: %v", err)
	}
	if string(data) != original {
		t.Fatalf("backup content = %q, want %q", string(data), original)
	}
}

func TestPreserveExistingConfigFileSkipsWhenGlobalBackupRan(t *testing.T) {
	homeDir := t.TempDir()
	zshrc := filepath.Join(homeDir, ".zshrc")
	if err := os.WriteFile(zshrc, []byte("custom\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	m := &Model{BackupDir: filepath.Join(t.TempDir(), "gentleman-backup")}
	if err := preserveExistingConfig(m, "shell", "Zsh (.zshrc)", zshrc); err != nil {
		t.Fatalf("preserveExistingConfig() error: %v", err)
	}

	backups, _ := filepath.Glob(zshrc + ".backup-*")
	if len(backups) != 0 {
		t.Fatalf("backups = %#v, want none when global backup ran", backups)
	}
}

func TestPreserveExistingConfigFileSkipsUnderDryRun(t *testing.T) {
	t.Setenv("GENTLEMAN_DRY_RUN", "1")

	homeDir := t.TempDir()
	zshrc := filepath.Join(homeDir, ".zshrc")
	if err := os.WriteFile(zshrc, []byte("custom\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	m := &Model{}
	if err := preserveExistingConfig(m, "shell", "Zsh (.zshrc)", zshrc); err != nil {
		t.Fatalf("preserveExistingConfig() error: %v", err)
	}

	backups, _ := filepath.Glob(zshrc + ".backup-*")
	if len(backups) != 0 {
		t.Fatalf("backups = %#v, want none under dry-run", backups)
	}
}

func TestPreserveExistingConfigPowerlevel10kFile(t *testing.T) {
	homeDir := t.TempDir()
	p10k := filepath.Join(homeDir, ".p10k.zsh")
	const customP10k = "# custom powerlevel10k theme settings\n"
	if err := os.WriteFile(p10k, []byte(customP10k), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	m := &Model{}
	if err := preserveExistingConfig(m, "shell", "Powerlevel10k (.p10k.zsh)", p10k); err != nil {
		t.Fatalf("preserveExistingConfig() error: %v", err)
	}

	backups, err := filepath.Glob(p10k + ".backup-*")
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	if len(backups) != 1 {
		t.Fatalf("backups = %#v, want exactly one", backups)
	}

	data, err := os.ReadFile(backups[0])
	if err != nil {
		t.Fatalf("read backup: %v", err)
	}
	if string(data) != customP10k {
		t.Fatalf("backup content = %q, want %q", string(data), customP10k)
	}
}
