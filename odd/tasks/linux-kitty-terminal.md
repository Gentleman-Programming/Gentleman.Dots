# Add Kitty to Linux TUI options and installer backend (#190)

Objetivo: permitir que los usuarios de Linux seleccionen Kitty desde la TUI del instalador, e instalar el binario automáticamente a través de los paquetes de la plataforma (Arch, Fedora, Debian/Ubuntu) si falta.
Branch: `fix/linux-kitty-installer` · Base: `origin/main` · Runner: `go test`

## Specs

### S1 — Opciones de terminal en Linux TUI (#190)
`model.go` en `ScreenTerminalSelect` excluía a Kitty cuando `Choices.OS == "linux"`, a pesar de que la CLI `--terminal=kitty` y la documentación la soportan. Debe ofrecer `Kitty` en Linux al igual que en macOS.
Aceptación:
- `GetCurrentOptions()` incluye `"Kitty"` en la posición 2 en Linux.
- Actualización de tests en `model_test.go`, `comprehensive_test.go`, `integration_test.go` e `installation_steps_test.go`.

### S2 — Instalación de Kitty en el backend de Linux
En `installer.go` caso `"kitty"`, si Kitty no está instalado y no es macOS, el instalador no ejecutaba ningún comando y logueaba erróneamente `"Kitty already installed"`. Debe instalar `kitty` mediante `installPlatformPackages` (Arch, Fedora, Debian, Brew).
Aceptación:
- `stepInstallTerminal` invoca `installPlatformPackages` con `kitty` en Linux cuando no está instalado.
- Test unitario en `platform_packages_test.go`.

## Tasks

| ID | Spec | Ruta | Estado | Evidencia |
|----|------|------|--------|-----------|
| T1 | S1 | inline | hecho | `model.go` + tests en `model_test.go`, `comprehensive_test.go`, `integration_test.go`, `installation_steps_test.go` |
| T2 | S2 | inline | hecho | `installer.go` + `platform_packages_test.go` |
| T3 | S1, S2 | inline | hecho | suite pasando, `go build ./...` OK, `go vet ./...` limpio |

## Log

### L1 — Pedido original (verbatim)
> "#190, #190, #218"
