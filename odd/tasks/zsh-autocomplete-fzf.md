# Integrate zsh-autocomplete menu navigation with fzf (#205)

Objetivo: permitir que la navegación del menú de `zsh-autocomplete` (Tab / Shift+Tab) conviva de forma determinista con el trigger `**<Tab>` de `fzf` sin pisarse.
Branch: `fix/arch-cachyos-installer` · Base: `8d97ca9` · Runner: `go test`

## Specs

### S1 — Widget despachador `gentleman-fzf-or-menu`
Tras `eval "$(fzf --zsh)"`, `GentlemanZsh/.zshrc` debe definir un widget ZLE `gentleman-fzf-or-menu` que:
- Si el buffer termina con el trigger de completado de fzf (`${FZF_COMPLETION_TRIGGER-'**'}`), invoque `zle fzf-completion`.
- De lo contrario, invoque `zle menu-select` para navegar el menú de `zsh-autocomplete`.
Aceptación:
- `bindkey '^I'` se enlaza a `gentleman-fzf-or-menu`.
- En el keymap `menuselect`, `^I` (Tab) avanza con `menu-complete`.
- Si `${terminfo[kcbt]}` está definido (Shift+Tab), se enlaza a `menu-select` y en `menuselect` a `reverse-menu-complete`.

### S2 — Ubicación y precedencia en `.zshrc`
El bloque despachador debe instalarse inmediatamente después de `eval "$(fzf --zsh)"` para capturar el widget y no ser anulado por integraciones previas.
Aceptación:
- `gentleman-fzf-or-menu` se define y enlaza después de `eval "$(fzf --zsh)"`.

## Tasks

| ID | Spec | Ruta | Estado | Evidencia |
|----|------|------|--------|-----------|
| T1 | S1, S2 | inline | hecho | test `TestZshTemplateIntegratesAutocompleteWithFzf` en `zsh_template_test.go` (`44b10b5`) |
| T2 | S1, S2 | inline | hecho | widget `gentleman-fzf-or-menu` y bindings en `GentlemanZsh/.zshrc` (`44b10b5`) |
| T3 | S1, S2 | inline | hecho | `go build ./...` OK, `go vet ./...` limpio, suite pasando (`44b10b5`) |

## Log

### L1 — Pedido original (verbatim)
> "Dejemoslo asi vamos con #205"

### L2 — Verificación (2026-11-21)
- `go build ./...` OK.
- `go vet ./...` limpio.
- `gofmt -l` limpio.
- `TestZshTemplateIntegratesAutocompleteWithFzf`: valida snippets del widget `gentleman-fzf-or-menu`, `zle fzf-completion`, `zle menu-select`, bindings de Tab y Shift+Tab en `menuselect`, y precedencia posterior a `eval "$(fzf --zsh)"`.
- Suite general: `installer/internal/system/` (100% pasando con `TMPDIR=/tmp`) e `installer/internal/tui/` pasando (excepto golden pre-existentes).

