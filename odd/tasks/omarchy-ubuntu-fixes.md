# Fix Omarchy Herdr autostart (#198) and Ubuntu 24.04 zshrc guards (#216)

Objetivo:
1. Detectar Omarchy OS para evitar generar el autostart incondicional de Herdr en los shells (Fish, Zsh, Nushell), respetando los atajos de teclado nativos de Omarchy (`Super+Return` para terminal limpia vs `Super+Ctrl+Return` para Herdr) (#198).
2. Proteger las integraciones de herramientas en `GentlemanZsh/.zshrc` (`carapace`, `fzf`, `zoxide`, `atuin`) con verificaciones de existencia para evitar errores al abrir el shell en Ubuntu 24.04 y entornos con paquetes parciales (#216).

Branch: `fix/omarchy-ubuntu-fixes` · Base: `origin/main` · Runner: `go test`

## Specs

### S1 — Detección de Omarchy OS y omisión de autostart de Herdr (#198)
- `detect.go`: detecta Omarchy OS revisando `/etc/omarchy-release` o `ID=omarchy` en `/etc/os-release`. Setea `IsOmarchy = true` en `SystemInfo`.
- `installer.go`: en `stepInstallShell` para Fish, Zsh y Nushell, si `IsOmarchy` es true y el WM seleccionado es `"herdr"`, se omite el bloque de autostart (`effectiveWM = "none"`).
Aceptación:
- En Omarchy, Herdr se instala y configura, pero no se inyecta el hook de autostart en el shell.

### S2 — Guards para herramientas externas en `GentlemanZsh/.zshrc` (#216)
- Carapace: se envuelve en `if command -v carapace >/dev/null 2>&1`.
- fzf: se verifica `command -v fzf`. Si `fzf --zsh` falla (versión < 0.48 en apt de Ubuntu 24.04), hace fallback a los scripts del paquete `/usr/share/doc/fzf/examples/`.
- zoxide y atuin: se guardan con `command -v ... >/dev/null 2>&1 && eval ...`.
Aceptación:
- Si alguna herramienta no está instalada o es una versión antigua de apt, el shell inicia limpio sin emitir mensajes de error `command not found`.

## Tasks

| ID | Spec | Ruta | Estado | Evidencia |
|----|------|------|--------|-----------|
| T1 | S1 | inline | hecho | `detect.go` + `detect_test.go` + `installer.go` |
| T2 | S2 | inline | hecho | `GentlemanZsh/.zshrc` |
| T3 | S1, S2 | inline | hecho | `go build ./...`, `go vet ./...`, tests pasando |

## Log

### L1 — Pedido original (verbatim)
> "#198, #216"
