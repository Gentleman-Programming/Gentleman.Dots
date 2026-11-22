# Port fish-theme and herdr-theme selectors to main (#218)

Objetivo: portar a `main` los selectores de temas `fish-theme` y `herdr-theme` con los perfiles `gentleman` y `gentleman-cute`, para que los usuarios sin Nix puedan alternar entre el tema oscuro neutral y el perfil graphite-rose.
Branch: `feat/theme-selectors-port` · Base: `origin/main` · Runner: `fish` + `bash` + `go test`

## Specs

### S1 — Port de selector y temas para Herdr (`herdr-theme`)
- `herdr/config-base.toml` con configuración base de atajos compartidos.
- `herdr/profiles/gentleman.toml` y `herdr/profiles/gentleman-cute.toml`.
- Script ejecutable `herdr/herdr-theme <gentleman|gentleman-cute>` que combina atómicamente la base y el perfil hacia `~/.config/herdr/config.toml` y ejecuta `herdr server reload-config`.
- En `installer.go`, el paso de Herdr copia la base, perfiles y el script ejecutable a `~/.local/bin/herdr-theme`.

### S2 — Port de selector y temas para Fish, Starship y Atuin (`fish-theme`)
- `GentlemanFish/fish/functions/fish-theme.fish`: función Fish que cambia temas de Fish, perfil de Starship y tema de Atuin.
- `GentlemanFish/fish/themes/{gentleman,gentleman-cute}.fish`: perfiles de colores de sintaxis, fzf y `LS_COLORS`.
- `starship/{gentleman,gentleman-cute}.toml`: perfiles estáticos de Starship para ambos esquemas de color.
- `atuin/themes/{gentleman,gentleman-cute}.toml`: temas para Atuin.
- En `installer.go`, el paso de Fish copia perfiles de Starship a `~/.config/starship/` y temas de Atuin a `~/.config/atuin/themes/`.
- Tests de regresión pasando: `ls-theme.sh`, `fish-theme-ls.sh`, `fish-theme-atuin.sh`.

## Tasks

| ID | Spec | Ruta | Estado | Evidencia |
|----|------|------|--------|-----------|
| T1 | S1 | inline | hecho | `herdr/config-base.toml`, `herdr/profiles/`, `herdr/herdr-theme` + prueba funcional |
| T2 | S2 | inline | hecho | `fish-theme.fish`, `themes/*.fish`, `starship/`, `atuin/` + tests en `GentlemanFish/tests/` |
| T3 | S1, S2 | inline | hecho | wiring en `installer.go` + `go build ./...` + `go vet ./...` + `go test` pasando |

## Log

### L1 — Pedido original (verbatim)
> "vamos con #218"
