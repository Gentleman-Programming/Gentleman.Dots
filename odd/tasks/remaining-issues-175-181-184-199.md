# Fix remaining issues (#184, #199, #181, #175)

Objetivo: resolver los 4 últimos issues abiertos del repositorio: soporte robusto de paquetes y COPR en Fedora 44 (#184), orden de carga de `zsh-autocomplete` frente a `oh-my-zsh.sh` para evitar filas en blanco (#199), fuga de secuencias OSC en WezTerm + WSL2 + Zellij + Nushell (#181), y validación/sincronización de la fórmula de Homebrew tap con la versión del instalador (#175).
Branch: `fix/arch-cachyos-installer` (#184, #199) y `fix/nushell-wsl-osc-and-tap` (#181, #175) · Runner: `go test`

## Specs

### S1 — Resolución de paquetes por disponibilidad en Fedora + fallback de Carapace + COPR de Ghostty (#184)
En `installer/internal/tui/native_packages.go` e `installer/internal/tui/installer.go`:
- Agregar `dnfPackageAvailable`, `resolveFedoraPackages(packages string, available packageAvailabilityFunc)` e `installFedoraPackages(m *Model, stepID, packages, brewPackages string, onLog func(string)) *system.ExecResult` siguiendo exactamente el patrón de `installArchPackages` e `installDebianPackages` (y respetando `isDnfNoOp` cuando `dnf` retorna exit 1 con `"Nothing to do"`).
- En `installPlatformPackages`, la rama `m.SystemInfo.OS == system.OSFedora && packages.Fedora != ""` debe usar `installFedoraPackages(m, stepID, packages.Fedora, packages.Brew, onLog)`.
- En `stepInstallShell` (`installer.go`), si el sistema es `OSFedora` y `carapace` no quedó instalado tras `installPlatformPackages`, ejecutar `ensureCarapaceOnArch(m, stepID)` (renombrado o extendido para cubrir Arch y Fedora con el binario estático oficial SHA-256 verificado).
- En `stepInstallTerminal` (`installer.go`) para Ghostty en Fedora, habilitar el COPR mantenido `scottames/ghostty` (con fallback a `pgdev/ghostty` si el primero falla) antes de `dnf install -y ghostty`.

### S2 — Evitar las 4 filas en blanco de `zsh-autocomplete` 26.08.04 (#199)
En `GentlemanZsh/.zshrc`:
- `source $ZSH/oh-my-zsh.sh` debe ejecutarse **antes** de cargar `zsh-autocomplete.plugin.zsh`, `zsh-autosuggestions.zsh` y `zsh-syntax-highlighting.zsh`, para que el `compinit` de Oh My Zsh no pise los hooks `zle-line-init`/`_main_complete` ya registrados por `zsh-autocomplete` 26.08.04 (causa de las 4 filas en blanco sobre el menú con prompts multilínea como p10k).
- En `zsh 5.9` (`[[ $ZSH_VERSION == 5.9* ]]`), pre-registrar como widgets ZLE vacíos `recent-paths` y `menu-search` si no existen antes de cargar `zsh-syntax-highlighting`, evitando el warning `zsh-syntax-highlighting: unhandled ZLE widget 'recent-paths'`.
- Mantener los tests de `installer/internal/system/zsh_template_test.go` en verde y agregar un test que verifique que `source $ZSH/oh-my-zsh.sh` precede a `zsh-autocomplete.plugin.zsh` y que `zsh-syntax-highlighting.zsh` se carga después de `zsh-autocomplete.plugin.zsh`.

### S3 — Evitar fuga de secuencias OSC (`rgb:...`) en WezTerm + WSL2 + Zellij + Nushell (#181)
En `GentlemanNushell/config.nu` e `installer/internal/system/exec.go` (`PatchNushellForWM`):
- Cambiar `osc633: true` a `osc633: ("TERM_PROGRAM" in $env and $env.TERM_PROGRAM == "vscode")` (o `false` por defecto) en `GentlemanNushell/config.nu` para no emitir escapes propietarios de VS Code dentro de WezTerm/Zellij.
- En `start_multiplexer` (`GentlemanNushell/config.nu` y `PatchNushellForWM` en `installer/internal/system/exec.go`), usar `exec $MULTIPLEXER` en lugar de `run-external $MULTIPLEXER`, para que el proceso inicial de Nushell sea reemplazado por Zellij/Tmux/Herdr antes de inicializar Reedline y disputar las respuestas OSC 4/10/11 en stdin bajo WSL2/ConPTY.
- Actualizar los tests de `PatchNushellForWM` en `installer/internal/system/patch_test.go` e `installer/internal/tui/installation_steps_test.go`.

### S4 — Verificación de sincronía de `homebrew-tap/Formula/gentleman-dots.rb` (#175)
- Agregar un test en `installer/internal/system/tap_formula_test.go` que verifique que la versión declarada en `homebrew-tap/Formula/gentleman-dots.rb` coincide con `Version` del instalador (`2.12.2`) y que los 4 hashes `sha256` (darwin-arm64, darwin-amd64, linux-arm64, linux-amd64) tienen formato hexadecimal válido de 64 caracteres y URLs que apuntan a `v#{version}`.

## Tasks

| ID | Spec | Ruta | Estado | Evidencia |
|----|------|------|--------|-----------|
| T1 | S1, S2 | `gentle-ai-worker` | hecho | Fedora package resolution + Ghostty COPR + `.zshrc` load order (#184, #199) · commit `7f5d8ae` |
| T2 | S3, S4 | `gentle-ai-worker` | hecho | Nushell `exec` & OSC fix + Homebrew tap formula sync test (#181, #175) · commit `fdc73ad` |

## Log

### L1 — Pedido original (verbatim)
> "Solucionemoslo todos"

### L2 — Verificación (2026-11-21)
- T1 (`7f5d8ae`): `go test ./internal/tui/ -run "TestInstallPlatformPackages|TestFilterAvailablePackages|TestInstallArchPackages|TestInstallDebianPackages|TestInstallFedora"` OK · `go test ./internal/system/ -run "TestZshTemplate|TestPatchZsh"` OK · `zsh -n GentlemanZsh/.zshrc` en contenedor `gentleman-test-alpine` (Zsh 5.9) OK.
- T2 (`fdc73ad`): `go test -v ./internal/system/ -run "TestPatchNushell|TestTapFormula"` OK · `go test -v ./internal/tui/ -run "TestPatchNushell|TestStepInstallShellNushell|TestAllShellAndWMCombinations"` OK · suite `./internal/system/` completa OK (`5.430s`).
- Nota: un comando combinado (`go build` + `go vet` + `./internal/system/` + `./internal/tui/` completo + `docker run`) excedió el timeout de 30s durante `./internal/tui/` tras completar `./internal/system/`.

