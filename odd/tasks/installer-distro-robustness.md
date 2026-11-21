# Installer distro robustness (#194, #215)

Objetivo: que el CI no mienta sobre el resultado de los tests, y que el paso de shell funcione en Debian/Ubuntu igual que en Arch.
Branch: `fix/arch-cachyos-installer` · Base: `80c1b0d` · Runner: `go test` + `bash -n`

## Specs

### S1 — Un test que falla dentro de un pipe no puede reportarse como PASS (#194)
`installer/e2e/docker-test.sh` corre bajo `set -e` sin `pipefail`, así que `run_image | tee` devuelve el status de `tee` (siempre 0) y un test fallido se cuenta como aprobado. Debe habilitar `pipefail`, lo que exige bash (POSIX sh no lo soporta). El cambio de intérprete no debe alterar la salida de color: los `echo "$COLOR..."` dependían del `echo` de dash, que interpreta escapes.
Aceptación: `(exit 1) | tee /dev/null` aborta con status 1; `echo "\033[32mok"` sigue emitiendo el byte ESC (033).

### S2 — El paso de shell no debe fallar en Debian/Ubuntu por un paquete inexistente (#215)
Las listas `Debian:` incluyen `starship`, que no está empaquetado en Debian/Ubuntu, así que `apt-get install -y` falla y **los tres perfiles de shell (fish, zsh, nushell) abortan siempre**. El instalador debe resolver cada nombre antes de instalar e instalar solo los resolubles, con warning visible por cada omitido.
Aceptación: con `starship` irresoluble, `apt-get install -y fish zoxide` se ejecuta, hay un warning nombrando `starship` y el paso termina OK; si ningún paquete resuelve, el paso falla sin invocar apt.

## Tasks

| ID | Spec | Ruta | Estado | Evidencia |
|----|------|------|--------|-----------|
| T1 | S1 | inline | hecho | `docker-test.sh` · ver L2 |
| T2 | S2 | inline | hecho | `native_packages.go` + `installer.go` · ver L2 |
| T3 | S1+S2 suite | inline | hecho | ver L2 |

## Log

### L1 — Pedido original (verbatim)
> "Arranquemos"

Contexto: continuación de "¿Arrancamos por #194 + #215 (mismo patrón, impacto inmediato)?" — el usuario eligió ese orden.

### L2 — Verificación (2026-11-21)
- `bash -n installer/e2e/docker-test.sh` → OK.
- Pipe fallido detectado: `bash -c 'set -eo pipefail; (exit 1) | tee /dev/null'` → exit 1 (antes: 0 ⇒ PASS falso).
- `xpg_echo` preserva color: `echo "\033[32mok"` emite `033 [ 3 2 m` (byte ESC 033) ⇒ sin regresión visual al pasar de dash a bash.
- `go build ./...` OK · `go vet` limpio · `go test ./internal/...` solo falla en `TestWelcomeScreenGolden` y `TestOSSelectGolden` (pre-existentes, idénticos al baseline).
- 2 tests nuevos en verde: skip de `starship` en Debian con warning, y error cuando nada resuelve sin invocar apt.
- Auditoría de pipes bajo `pipefail`: el único que necesitaba `|| true` es la extracción de `[FAIL]` (grep sale 1 sin coincidencias y habría abortado el script).
- No se agregó `set -u` a propósito: hay variables opcionales (`platform_flag`, `tag_suffix`) que quedarían sin definir.
