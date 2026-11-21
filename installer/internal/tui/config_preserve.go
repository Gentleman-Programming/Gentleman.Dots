package tui

import (
	"fmt"
	"os"
	"time"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
)

// preserveExistingConfig keeps a copy of an existing configuration path before
// the installer overwrites it. When the run already produced a global backup
// (m.BackupDir), that backup is the copy and nothing extra is created.
//
// This is what keeps a declined backup from turning into data loss: choosing
// "no backup" must not authorize the installer to destroy a user's config.
func preserveExistingConfig(m *Model, stepID, label, path string) error {
	if system.IsDryRun() {
		return nil
	}

	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return nil
	}

	if m.BackupDir != "" {
		return nil
	}

	backupPath := fmt.Sprintf("%s.backup-%s", path, time.Now().Format("20060102-150405"))
	if err := system.CopyDir(path, backupPath); err != nil {
		return err
	}

	SendLog(stepID, fmt.Sprintf("Existing %s configuration preserved at %s", label, backupPath))
	return nil
}
