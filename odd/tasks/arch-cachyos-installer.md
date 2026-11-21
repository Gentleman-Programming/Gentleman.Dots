# Arch/CachyOS installer hardening

Objetivo: que el installer de Gentleman.Dots complete en Arch y CachyOS sin abortar por paquetes AUR-only, sin perder funcionalidad (carapace, powerlevel10k) y sin destruir configs del usuario.
Branch: `fix/arch-cachyos-installer` (fork MrDylanDev/Gentleman.Dots) · Base: `59aee78` + merge de PR #197 (`395d88e`) · Runner: `go test` (Go 1.27.4)

## Specs

### S1 — Un nombre irresoluble no debe abortar la instalación en Arch
El paso de instalación de paquetes en Arch NO debe dejar que un nombre que `pacman` no resuelve aborte la transacción completa. Debe: filtrar por disponibilidad real en los repos configurados (`pacman -Si <name>`), instalar el conjunto resoluble en una sola transacción, emitir un warning visible por cada nombre omitido, y fallar solo si el binario del shell elegido no quedó disponible.
Aceptación: en un host Arch sin repos extra, `fish`, `zsh`, `nushell` + `zoxide atuin starship` instalan y el paso termina OK; `carapace` ausente no aborta.

### S2 — carapace debe seguir disponible en Arch/CachyOS
`carapace` no existe en los repos oficiales de Arch ni en CachyOS (verificado: `pacman -Ss '^carapace'` vacío; AUR: `carapace-bin` 1.8.2). Cuando no esté en repos, el installer debe instalarlo desde el release estático oficial verificando SHA256, reutilizando el patrón ya existente de `installHerdrBinary`.
Aceptación: con carapace ausente del repo, el binario queda en `~/.local/bin/carapace` ejecutable; checksum incorrecto ⇒ error y binario borrado.

### S3 — CachyOS no debe perder powerlevel10k
`zsh-theme-powerlevel10k` NO está en repos oficiales de Arch, pero SÍ está en el repo `cachyos` (verificado: `cachyos/zsh-theme-powerlevel10k 1:1.20.18-1`). El filtro de S1 debe instalarlo donde exista y omitirlo con warning donde no, sin hardcodear distro.
Aceptación: en CachyOS se instala; en Arch vanilla se omite con warning y el paso no falla.

### S4 — `~/.config/fish` no debe destruirse
Antes de copiar la configuración de Fish, si `~/.config/fish` existe debe respaldarse; el `config.fish` del usuario no debe perderse sin backup ni aviso.
Aceptación: con `~/.config/fish/config.fish` propio, tras el paso existe un backup que contiene ese contenido original y el usuario ve el aviso.

## Tasks

| ID | Spec | Ruta | Estado | Evidencia |
|----|------|------|--------|-----------|
| T1 | S1, S3 | inline | hecho | `arch_packages.go` (filtro por disponibilidad) · commit `80c1b0d` |
| T2 | S2 | inline | hecho | `installCarapaceBinary` v1.8.2 + SHA256 verificado · `80c1b0d` |
| T3 | S1-S3 wiring en `stepInstallShell` | inline | hecho | `installer.go` caso Arch + hook carapace · `80c1b0d` |
| T4 | S4 | inline | hecho | `config_preserve.go` · `80c1b0d` |
| T5 | suite + build + diff diferencial | inline | hecho | ver L5 |

## Log

### L1 — Pedido original (verbatim)
> "Hagamos un Fork y solucionemoslo los de arch y cachyos"

### L2 — Baseline observado (2026-11-21)
- `go test ./internal/system/` → ok.
- `go test ./internal/tui/` → FAIL pre-existente: `TestWelcomeScreenGolden`, `TestOSSelectGolden` (golden de entorno). Se registran como fallos conocidos, no causados por este cambio.
- Tras mergear PR #197: mismos 2 fallos, `go build ./...` OK ⇒ sin regresión.

### L3 — Verificación en la máquina del usuario (CachyOS, x86_64)
- `pacman -Ss '^carapace'` → vacío ⇒ AUR-only. AUR: `carapace-bin` 1.8.2, `carapace` 1.8.2.
- `pacman -Ss powerlevel10k` → `cachyos/zsh-theme-powerlevel10k 1:1.20.18-1` ⇒ presente en CachyOS.
- API Arch (`q=powerlevel`) → vacío ⇒ ausente en Arch oficial.
- Conclusión: el mismo listado de paquetes se comporta distinto en Arch y CachyOS ⇒ la lista hardcodeada es la causa raíz; el filtro por disponibilidad (S1) resuelve ambos casos.

### L4 — Decisiones
- Se reutiliza PR #197 (merge `395d88e`) para #193/#195/dry-run en vez de reescribirlo; autoría preservada.
- Se descarta el enfoque de PR #202 (borrar `zsh-theme-powerlevel10k`) porque rompe el tema en CachyOS.
- Se adopta el enfoque de PR #185 (carapace desde release estático) por ser el único que conserva la funcionalidad.
- Riesgo: ítem 1 (datos: configs del usuario) + ítem 5 (delivery/installer).

### L5 — Verificación (2026-11-21)
- `go vet ./internal/tui/ ./internal/system/` → limpio.
- `go build ./...` → OK.
- `go test ./internal/...` → system OK, trainer OK; tui solo falla en `TestWelcomeScreenGolden` y `TestOSSelectGolden`, idénticos al baseline ⇒ sin regresión.
- 9 tests nuevos en verde: filtro de disponibilidad, skip de carapace, p10k presente en CachyOS, p10k ausente en Arch, nada disponible ⇒ error, fallback a brew tras filtrar, y 3 de preservación de config (backup, skip con backup global, skip en dry-run).
- Sonda real en CachyOS con `pacman -Si` (cero mutaciones): `carapace` OMITIDO; `zsh-theme-powerlevel10k` OK; fish/zsh/nushell/zoxide/atuin/starship/zsh-autocomplete/zsh-autosuggestions/zsh-syntax-highlighting OK.
- Commit del work unit: `80c1b0d` (427 inserciones, 1 borrado). `odd/` queda sin trackear a propósito para no ensuciar un eventual PR upstream.
- Pendiente (decisión humana): push al fork y PR upstream.
