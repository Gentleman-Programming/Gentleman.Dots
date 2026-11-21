# Preserve user Zsh customizations across updates (#206)

Objetivo: permitir que las personalizaciones locales de Zsh sobrevivan a las actualizaciones del instalador mediante `~/.zshrc.d/*.zsh`, y que `~/.zshrc` y `~/.p10k.zsh` existentes se respalden automáticamente antes de ser sobrescritos.
Branch: `fix/arch-cachyos-installer` · Base: `ec06887` · Runner: `go test`

## Specs

### S1 — Carga de extensiones en `~/.zshrc.d/*.zsh`
`GentlemanZsh/.zshrc` debe incluir un cargador para archivos de extensión mantenidos por el usuario en `~/.zshrc.d/*.zsh(N)` en orden léxico.
Aceptación:
- Si `~/.zshrc.d` no existe o está vacío, no se emite ningún error en zsh.
- Si existen archivos legibles `.zsh`, se cargan (`source`) en orden léxico.
- La variable de iteración `config_file` queda desasignada (`unset`) al finalizar.
- Se ubica al final de `.zshrc`, tras las integraciones centrales (fzf, zoxide, atuin, p10k).

### S2 — Preservación automática de archivos existentes en `preserveExistingConfig`
`preserveExistingConfig` soportaba únicamente directorios (`system.CopyDir`). Debe soportar también archivos regulares (`system.CopyFile`) para respaldar configs individuales como `~/.zshrc` y `~/.p10k.zsh` con timestamp `.backup-<timestamp>` cuando no haya un backup global activo.
Aceptación:
- Si el path es un archivo regular existente y `m.BackupDir == ""`, se copia con `system.CopyFile` a `<path>.backup-<timestamp>`.
- Si `m.BackupDir != ""` o estamos en `dry-run`, no se crea copia extra.

### S3 — Respaldo de `~/.zshrc` y `~/.p10k.zsh` en `stepInstallShell`
Antes de copiar `GentlemanZsh/.zshrc` y `GentlemanZsh/.p10k.zsh` a `homeDir`, el instalador debe llamar a `preserveExistingConfig` para cada archivo.
Aceptación: un `~/.zshrc` o `~/.p10k.zsh` previo no se destruye sin backup ni aviso.

## Tasks

| ID | Spec | Ruta | Estado | Evidencia |
|----|------|------|--------|-----------|
| T1 | S1 | inline | hecho | `GentlemanZsh/.zshrc` loader + test en `zsh_template_test.go` (`8d97ca9`) |
| T2 | S2 | inline | hecho | `config_preserve.go` soporte para archivos + 4 tests en `config_preserve_test.go` (`8d97ca9`) |
| T3 | S3 | inline | hecho | llamadas en `installer.go` (`stepInstallShell` case "zsh") (`8d97ca9`) |
| T4 | S1-S3 | inline | hecho | `go build ./...` OK, `go vet ./...` limpio, suite de tests passing (`8d97ca9`) |

## Log

### L1 — Pedido original (verbatim)
> "1" (en respuesta a "Atacar el issue #206: aplicar a ~/.zshrc el mismo patrón de backup/preservación que ya implementamos para fish")

### L2 — Hallazgos y diseño
- Issue #206 (ricardo-rod) reporta que `.zshrc` se sobrescribe en cada actualización, perdiendo customizaciones locales del usuario (como menús de autocompletado y aliases).
- Propuesta oficial de #206: proveer `~/.zshrc.d/*.zsh` como punto de extensión estable del usuario. Se usa el glob qualifier `(N)` de Zsh para que no falle si la carpeta no existe o está vacía, y se itera en orden léxico desasignando `config_file`.
- Se coloca al final de `.zshrc` después de p10k, zoxide, fzf y atuin, permitiendo sobreescribir o extender cualquier integración previa.
- Además, `preserveExistingConfig` se generalizó: antes solo hacía `CopyDir` si era directorio. Ahora detecta si es `info.Mode().IsRegular()` y usa `system.CopyFile`, respaldando `.zshrc` y `.p10k.zsh` antes de que el instalador los copie.

### L3 — Verificación (2026-11-21)
- `go build ./...` OK.
- `go vet ./...` limpio.
- `gofmt -l` limpio en todos los archivos modificados.
- Tests nuevos y modificados pasando:
  - `TestZshTemplateLoadsUserExtensionsDirectory`: verifica presencia de snippets del loader y orden respecto a p10k.
  - `TestPreserveExistingConfigFileCopiesRegularFile`: verifica copia timestamped de archivo regular y preservación exacta de contenido.
  - `TestPreserveExistingConfigFileSkipsWhenGlobalBackupRan`: verifica que no se duplica si ya hubo backup global.
  - `TestPreserveExistingConfigFileSkipsUnderDryRun`: verifica no-op en dry run.
  - `TestPreserveExistingConfigPowerlevel10kFile`: verifica preservación de `.p10k.zsh`.
- Suite general: `installer/internal/system/` (100% pasando con `TMPDIR=/tmp`) e `installer/internal/tui/` (todos pasando excepto golden pre-existentes).

