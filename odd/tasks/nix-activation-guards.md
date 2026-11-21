# Harden Nix activation hooks (#207, #208)

Objetivo: evitar que los hooks de activación de Home Manager fallen o aborten el switch por PATHs incompletos o errores transitorios de red.
Branch: `fix/nix-activation-guards` · Base: `origin/nix-migration` · Runner: `bash -n`

## Specs

### S1 — Exportar PATH de Homebrew en `herdr.nix` (#207)
`home.activation.installHerdr` corre en un subshell con PATH mínimo que no incluye Homebrew (`/opt/homebrew/bin`, `/usr/local/bin`, `/home/linuxbrew/.linuxbrew/bin`). Debe anteponer las rutas estándar de Homebrew antes de verificar `command -v herdr` y `command -v brew`.
Aceptación:
- En macOS (`Darwin`), antepone `/opt/homebrew/bin:/usr/local/bin`.
- En Linux, antepone `/home/linuxbrew/.linuxbrew/bin`.
- Si `herdr` o `brew` existen en esas rutas, se detectan correctamente y no se emite el warning falso `"Homebrew not found"`.

### S2 — Proteger `installOpenCode` contra fallas de red y remover `gh copilot` deprecado (#208)
`home.activation.installOpenCode` corre bajo `set -e`. Si la API de GitHub rate-limita la IP o no hay red durante `curl -fsSL https://opencode.ai/install | bash`, el switch completo de Home Manager aborta con exit 1.
Aceptación:
- La invocación del instalador de OpenCode debe estar protegida para emitir un aviso en caso de error sin abortar la activación (`|| echo ...`).
- Se elimina el bloque de `gh extension install github/gh-copilot`, el cual es código muerto/deprecado y fue archivado upstream.

## Tasks

| ID | Spec | Ruta | Estado | Evidencia |
|----|------|------|--------|-----------|
| T1 | S1 | inline | hecho | `herdr.nix` PATH export (#207) |
| T2 | S2 | inline | hecho | `opencode.nix` guard de red + limpieza copilot (#208) |
| T3 | S1, S2 | inline | hecho | validación de sintaxis bash (`bash -n`) limpia |

## Log

### L1 — Pedido original (verbatim)
> "Vamos a solucionarlo #207, #208" -> "1"
Contexto: rama dedicada `fix/nix-activation-guards` basada en `origin/nix-migration`.

### L2 — Verificación (2026-11-21)
- `bash -n` ejecutado sobre ambos bloques de activación: sin errores de sintaxis.
- En `herdr.nix`: antepone `/opt/homebrew/bin:/usr/local/bin` en Darwin y `/home/linuxbrew/.linuxbrew/bin` en Linux, resolviendo el falso warning de "Homebrew not found" y permitiendo autoinstalar Herdr.
- En `opencode.nix`: protege el pipe `curl ... | bash` con warning en stderr/stdout sin propagar exit 1 a Home Manager; elimina bloque muerto de `github/gh-copilot`.

