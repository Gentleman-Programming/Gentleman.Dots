# Preserve an existing Oh My Zsh checkout (#187)

Objetivo: el installer no debe escribir dentro de `~/.oh-my-zsh` cuando ya existe, porque es un checkout git que se autoactualiza.
Branch: `fix/arch-cachyos-installer` · Base: `ad70637` · Runner: `go test`

## Specs

### S1 — Un `~/.oh-my-zsh` existente no debe tocarse
`system.CopyDir` escribía el snapshot vendorizado (`GentlemanZsh/.oh-my-zsh`, 1019 archivos, ~11 MB) sobre `~/.oh-my-zsh` sin comparar ni preguntar. Sobre un checkout real eso marca los archivos versionados como modificados, y `omz update` (`git pull --rebase` con `rebase.autoStash true`) choca al hacer pop del autostash. El installer debe copiar el snapshot **solo cuando no existe** un `~/.oh-my-zsh`.
Aceptación: con un checkout existente, su contenido queda byte a byte igual y no se agregan archivos del snapshot; sin omz, el snapshot se copia.

### S2 — Un checkout enlazado simbólicamente cuenta como presente
`PathExists` resuelve symlinks, así que un `~/.oh-my-zsh` que apunta a un dotfiles repo cuenta como instalado y no se pisa el destino real.
Aceptación: con symlink a un directorio real, el contenido del directorio real queda intacto.

## Tasks

| ID | Spec | Ruta | Estado | Evidencia |
|----|------|------|--------|-----------|
| T1 | S1, S2 | inline | hecho | `PathExists` en `system/exec.go` + `installOhMyZshIfMissing` en `config_preserve.go` |
| T2 | S1, S2 | inline | hecho | 3 tests en `config_preserve_test.go` |

## Log

### L1 — Pedido original (verbatim)
> "seguimos"

Contexto: continuación de "¿seguimos con #187 (`~/.oh-my-zsh`) antes de publicar nada?" — el usuario eligió seguir con #187.

### L2 — Hallazgos
- El snapshot vendorizado **no contiene código propio**: `themes/` son temas upstream de omz y `custom/` son los `example.*` de omz. `PathExists` no existía.
- `.zshrc` (línea 117) hace `source $ZSH/oh-my-zsh.sh`, así que omz debe existir, pero no debe sobrescribirse.
- PR #188 propone borrar el snapshot completo (106.280 líneas) e instalar omz con el script oficial solo si falta. Resuelve también el problema de frescura, pero agrega dependencia de red y un diff de 1025 archivos.
- Decisión: cerrar el invariante reportado en #187 (nunca escribir sobre un checkout ajeno) con un diff mínimo y sin dependencia de red. El borrado del snapshot queda como decisión aparte, porque cambia el comportamiento de instalación desde cero (ver L3).

### L3 — Limitación conocida y pendiente
Instalación desde cero sigue copiando el snapshot, que **no es un repo git**, así que `omz update` no funciona ahí tampoco. Eso lo resuelve el enfoque de #188 (instalar omz oficialmente). No se incluyó acá para no meter una dependencia de red ni un borrado masivo en el mismo cambio.

### L4 — Verificación (2026-11-21)
- `go build ./...` OK · `go vet` limpio · mis archivos pasan `gofmt -l` (vacío).
- `go test ./internal/...` solo falla en `TestWelcomeScreenGolden` y `TestOSSelectGolden` (pre-existentes).
- 3 tests nuevos en verde: copia cuando falta, checkout existente intacto sin archivos agregados, symlink tratado como presente.
- `installer/internal/system/detect.go` aparece en `gofmt -l`, pero no tiene diff mío: ya venía así upstream. No se toca para no ensuciar el cambio.
