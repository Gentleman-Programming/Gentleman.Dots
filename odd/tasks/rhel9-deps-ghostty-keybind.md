# Fix Ghostty Shift+Enter keybind (#217) and dnf4 deps no-op exit (#219)

Objetivo:
1. Eliminar el keybind obsoleto `shift+enter=text:\x1b\r` en Ghostty que rompe saltos de línea y pegado multilínea en Claude Code / Herdr (#217).
2. Tolerar el código de salida 1 de dnf4 (`Error: Nothing to do.`) cuando todas las dependencias ya están instaladas en RHEL 9 / AlmaLinux / CentOS 9 (#219).

Branch: `fix/rhel9-deps-ghostty` · Base: `origin/main` · Runner: `go test`

## Specs

### S1 — Remover keybind `shift+enter` en `GentlemanGhostty/config` (#217)
Ghostty 1.3+ y Herdr negocian el protocolo de teclado de Kitty y envían eventos nativos de `Shift+Enter`. El keybind de texto `text:\x1b\r` reemplaza ese evento con `ESC+CR`, haciendo que Claude Code envíe el mensaje en vez de insertar una nueva línea.
Aceptación:
- `GentlemanGhostty/config` no contiene `shift+enter=text:\x1b\r`.

### S2 — Tolerar dnf4 no-op exit 1 con "Nothing to do" (#219)
En distros con dnf4 (RHEL 9 / AlmaLinux / Rocky / CentOS 9), si todas las dependencias solicitadas ya están satisfechas, `dnf install` sale con exit 1 y `Error: Nothing to do.`. El instalador debe reconocer esto como éxito y no abortar.
Aceptación:
- `isDnfNoOp` retorna `true` cuando `ExitCode == 1` y stderr/stdout contiene `"Nothing to do"`.
- `stepInstallDeps` no falla en Fedora/RHEL cuando dnf retorna exit 1 con `"Nothing to do"`.
- `runNativeWithBrewFallback` reconoce `isDnfNoOp` y no falla ni cae innecesariamente a Homebrew.

## Tasks

| ID | Spec | Ruta | Estado | Evidencia |
|----|------|------|--------|-----------|
| T1 | S1 | inline | hecho | `GentlemanGhostty/config` (eliminado `shift+enter=text:\x1b\r`) |
| T2 | S2 | inline | hecho | `installer.go` (`isDnfNoOp`) + tests en `platform_packages_test.go` |
| T3 | S1, S2 | inline | hecho | `go test ./internal/tui/` pasando, `go build ./...` OK, `go vet ./...` limpio |

## Log

### L1 — Pedido original (verbatim)
> "liquidemos #217, 219"

### L2 — Verificación (2026-11-21)
- `GentlemanGhostty/config`: ya no contiene `shift+enter=text:\x1b\r`. Ghostty y Herdr usan el protocolo Kitty para enviar `Shift+Enter` nativo sin romper saltos de línea ni pegado en Claude Code (#217).
- `isDnfNoOp`: testeado con suite table-driven (exit 0, exit 1 con "Nothing to do" en stderr y stdout, exit 1 con error real).
- `stepInstallDeps` y `runNativeWithBrewFallback`: toleran dnf4 exit 1 sin abortar y sin caer innecesariamente a brew (#219).
- `go build ./...` OK, `go vet ./...` limpio, suite de tests pasando sin regresiones.

