# Autostart del WM con powerlevel10k (#210)

Objetivo: que tmux/zellij/herdr arranquen solos al abrir la terminal con Zsh + powerlevel10k.
Branch: `fix/arch-cachyos-installer` · Base: `8ee434e` · Runner: `go test`

## Specs

### S1 — El lanzador del WM debe ejecutarse antes del instant prompt de p10k
El bloque de instant prompt de powerlevel10k redirige stdout a un archivo temporal mientras corre el resto de `.zshrc` (`internal/p10k.zsh:6501`: `exec ... 0<&$fd_null 1>&$fd_out`). Después de ese punto `[[ -t 1 ]]` es falso, así que el guard del lanzador nunca se cumple y el WM nunca arranca. El lanzador debe estar **antes** del bloque de instant prompt.
Aceptación: la línea de `start_if_needed` precede a la línea del `p10k-instant-prompt`; y `PatchZshForWM` preserva ese orden para tmux, zellij y herdr.

### S2 — Elegir "none" sigue eliminando el lanzador
El patching debe seguir borrando el bloque completo cuando no se elige WM, ahora que está al principio del archivo.
Aceptación: tras `PatchZshForWM(..., "none", ...)` no queda `start_if_needed`, `WM_CMD` ni `WM_VAR`.

## Tasks

| ID | Spec | Ruta | Estado | Evidencia |
|----|------|------|--------|-----------|
| T1 | S1 | inline | hecho | `GentlemanZsh/.zshrc` (bloque movido al inicio) |
| T2 | S1, S2 | inline | hecho | `installer/internal/system/zsh_template_test.go` (4 tests) |

## Log

### L1 — Pedido original (verbatim)
> "Sigamos"

Contexto: continuación de "¿Seguimos con alguno de esos 6 ... o publicamos lo que tenemos?" — el usuario eligió seguir; se tomó #210/#211 por ser un bug real con PR ya escrito.

### L2 — Verificación de la premisa
- Fuente de p10k (`internal/p10k.zsh`, `_p9k_dump_instant_prompt`, línea 6501): el bloque generado hace `exec {__p9k_fd_0}<&0 {__p9k_fd_1}>&1 {__p9k_fd_2}>&2 0<&$fd_null 1>&$fd_out`. Es el shell principal con fd 0 en `/dev/null` y fd 1/2 en un archivo temporal, no un subshell. La conclusión del issue es correcta aunque la palabra "subshell" sea imprecisa.
- El guard del instant prompt exige `-t 0 && -t 1 && -t 2`; después de la redirección los tres son falsos, por eso el lanzador no dispara.
- `PatchZshForWM` es position-independent (escanea todas las líneas buscando `WM_VAR=`, `WM_CMD=`, `function start_if_needed` y la llamada), así que mover el bloque no rompe el patching.

### L3 — Verificación (2026-11-21)
- `go build ./...` OK · `go vet ./internal/system/` limpio.
- 4 tests nuevos en verde: orden launcher/instant-prompt en el template real, una sola declaración/llamada, orden preservado tras `PatchZshForWM` para tmux/zellij/herdr, y `none` elimina todo el bloque.
- Sonda empírica en pty con bash (semántica de fd idéntica): con el lanzador antes de `exec 1>archivo`, "WM STARTED" aparece en la terminal; con el lanzador después, no aparece en ningún lado. Control positivo y negativo.
- Suite `./internal/...`: solo fallan `TestWelcomeScreenGolden` y `TestOSSelectGolden` (pre-existentes).
- zsh no está instalado en esta máquina, por eso la sonda usó bash; la propiedad `-t 1` sobre fd es POSIX.

### L4 — Fuera de alcance
`~/.zshrc` se sobrescribe sin backup (mismo patrón que fish, que ya corregimos). Eso es #206 y queda como paso siguiente.
