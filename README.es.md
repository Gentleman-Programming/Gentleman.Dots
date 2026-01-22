# Gentleman.Dots

> ℹ️ **Actualización (enero 2026)**: OpenCode ahora soporta suscripciones **Claude Max/Pro** mediante el plugin `opencode-anthropic-auth` (incluido en esta configuración).
> Tanto **Claude Code** como **OpenCode** funcionan con tu suscripción de Claude.
> *Nota: este workaround es estable por ahora, pero Anthropic podría bloquearlo en el futuro.*

📄 Leer en: [English](README.md) | **Español**

## Tabla de Contenidos

* [¿Qué es esto?](#qué-es-esto)
* [Inicio rápido](#inicio-rápido)
* [Plataformas soportadas](#plataformas-soportadas)
* [🎮 Entrenador de Maestría en Vim](#-entrenador-de-maestría-en-vim)
* [Documentación](#documentación)
* [Resumen de herramientas](#resumen-de-herramientas)
* [Bleeding Edge](#bleeding-edge)
* [Estructura del proyecto](#estructura-del-proyecto)
* [Soporte](#soporte)

---

## Vista previa

### Instalador TUI

<img width="1424" height="1536" alt="Instalador TUI" src="https://github.com/user-attachments/assets/1db56d3b-a8c0-4885-82aa-c5ec04af4ac0" />

### Showcase

<img width="3840" height="2160" alt="Showcase del entorno de desarrollo" src="https://github.com/user-attachments/assets/fff14c05-9676-4e04-b05e-dab5e3cf300a" />

---

## ¿Qué es esto?

Una configuración completa de entorno de desarrollo que incluye:

* **Neovim** con LSP, autocompletado y asistentes de IA (Claude Code, Gemini, OpenCode)
* **Shells**: Fish, Zsh, Nushell
* **Multiplexores de terminal**: Tmux, Zellij
* **Emuladores de terminal**: Alacritty, WezTerm, Kitty, Ghostty

---

## Inicio rápido

### Opción 1: Homebrew (Recomendado)

```bash
brew install Gentleman-Programming/tap/gentleman-dots
gentleman-dots
```

### Opción 2: Descarga directa

```bash
# macOS Apple Silicon
curl -fsSL https://github.com/Gentleman-Programming/Gentleman.Dots/releases/latest/download/gentleman-installer-darwin-arm64 -o gentleman.dots

# macOS Intel
curl -fsSL https://github.com/Gentleman-Programming/Gentleman.Dots/releases/latest/download/gentleman-installer-darwin-amd64 -o gentleman.dots

# Linux x86_64
curl -fsSL https://github.com/Gentleman-Programming/Gentleman.Dots/releases/latest/download/gentleman-installer-linux-amd64 -o gentleman.dots

# Linux ARM64 (Raspberry Pi, etc.)
curl -fsSL https://github.com/Gentleman-Programming/Gentleman.Dots/releases/latest/download/gentleman-installer-linux-arm64 -o gentleman.dots

# Luego ejecutar
chmod +x gentleman.dots
./gentleman.dots
```

### Opción 3: Termux (Android)

Termux requiere compilar el instalador localmente (la cross-compilación de Go a Android tiene limitaciones).

```bash
# 1. Instalar dependencias
pkg update && pkg upgrade
pkg install git golang

# 2. Clonar el repositorio
git clone https://github.com/Gentleman-Programming/Gentleman.Dots.git
cd Gentleman.Dots/installer

# 3. Compilar y ejecutar
go build -o ~/gentleman-installer ./cmd/gentleman-installer
cd ~
./gentleman-installer
```

| Soporte en Termux                 | Estado                                               |
| --------------------------------- | ---------------------------------------------------- |
| Shells (Fish, Zsh, Nushell)       | ✅ Disponible                                         |
| Multiplexores (Tmux, Zellij)      | ✅ Disponible                                         |
| Neovim con configuración completa | ✅ Disponible                                         |
| Nerd Fonts                        | ✅ Instaladas automáticamente en `~/.termux/font.ttf` |
| Emuladores de terminal            | ❌ No aplica                                          |
| Homebrew                          | ❌ Usa `pkg`                                          |

> **Tip:** Después de la instalación, reiniciá Termux para aplicar la fuente y luego ejecutá `tmux` o `zellij` para iniciar el entorno configurado.

El instalador TUI te guía para seleccionar tus herramientas preferidas y maneja toda la configuración automáticamente.

> **Usuarios de Windows:** primero debés configurar WSL. Ver la [Guía de instalación manual](docs/manual-installation.md#windows-wsl).

---

## Plataformas soportadas

| Plataforma            | Arquitectura          | Método de instalación       | Gestor de paquetes |
| --------------------- | --------------------- | --------------------------- | ------------------ |
| macOS                 | Apple Silicon (ARM64) | Homebrew, descarga directa  | Homebrew           |
| macOS                 | Intel (x86_64)        | Homebrew, descarga directa  | Homebrew           |
| Linux (Ubuntu/Debian) | x86_64, ARM64         | Homebrew, descarga directa  | Homebrew           |
| Linux (Fedora/RHEL)   | x86_64, ARM64         | Descarga directa            | dnf                |
| Linux (Arch)          | x86_64                | Homebrew, descarga directa  | Homebrew           |
| Windows               | WSL                   | Descarga directa (ver docs) | Homebrew           |
| Android               | Termux (ARM64)        | Compilación local           | pkg                |

---

## 🎮 Entrenador de Maestría en Vim

¡Aprendé Vim de forma divertida! El instalador incluye un entrenador interactivo estilo RPG con:

| Módulo                   | Teclas cubiertas                         |
| ------------------------ | ---------------------------------------- |
| 🔤 Movimiento horizontal | `w`, `e`, `b`, `f`, `t`, `0`, `$`, `^`   |
| ↕️ Movimiento vertical   | `j`, `k`, `G`, `gg`, `{`, `}`            |
| 📦 Objetos de texto      | `iw`, `aw`, `i"`, `a(`, `it`, `at`       |
| ✂️ Cambiar y repetir     | `d`, `c`, `dd`, `cc`, `D`, `C`, `x`      |
| 🔄 Sustitución           | `r`, `R`, `s`, `S`, `~`, `gu`, `gU`, `J` |
| 🎬 Macros y registros    | `qa`, `@a`, `@@`, `"ay`, `"+p`           |
| 🔍 Regex / Búsqueda      | `/`, `?`, `n`, `N`, `*`, `#`, `\v`       |

Cada módulo incluye 15 lecciones progresivas, modo práctica con selección inteligente de ejercicios, jefes finales y seguimiento de XP.

Podés iniciarlo desde el menú principal: **Vim Mastery Trainer**

---

## Documentación

| Documento                                                     | Descripción                                          |
| ------------------------------------------------------------- | ---------------------------------------------------- |
| [Guía del instalador TUI](docs/tui-installer.md)              | Funciones interactivas, navegación, backup y restore |
| [Instalación manual](docs/manual-installation.md)             | Configuración paso a paso para todas las plataformas |
| [Keymaps de Neovim](docs/neovim-keymaps.md)                   | Referencia completa de atajos                        |
| [Configuración de IA](docs/ai-configuration.md)               | Claude Code, OpenCode, Copilot y más                 |
| [Especificación del entrenador Vim](docs/vim-trainer-spec.md) | Detalles técnicos del entrenador                     |
| [Testing con Docker](docs/docker-testing.md)                  | Tests E2E con contenedores                           |
| [Contribuir](docs/contributing.md)                            | Setup de desarrollo, sistema de skills y releases    |

---

## Resumen de herramientas

### Emuladores de terminal

| Herramienta   | Descripción                                  |
| ------------- | -------------------------------------------- |
| **Ghostty**   | Acelerado por GPU, nativo y ultra rápido     |
| **Kitty**     | Rico en funcionalidades, renderizado por GPU |
| **WezTerm**   | Configurable con Lua, multiplataforma        |
| **Alacritty** | Minimalista, escrito en Rust                 |

### Shells

| Herramienta | Descripción                                |
| ----------- | ------------------------------------------ |
| **Nushell** | Datos estructurados y pipelines modernos   |
| **Fish**    | Amigable y con excelentes defaults         |
| **Zsh**     | Altamente personalizable, compatible POSIX |

### Multiplexores

| Herramienta | Descripción                           |
| ----------- | ------------------------------------- |
| **Tmux**    | Probado en batalla, ampliamente usado |
| **Zellij**  | Moderno, plugins WebAssembly          |

### Editor

| Herramienta | Descripción                             |
| ----------- | --------------------------------------- |
| **Neovim**  | Config LazyVim con LSP, completado e IA |

### Prompts

| Herramienta  | Descripción                            |
| ------------ | -------------------------------------- |
| **Starship** | Prompt multi-shell con integración Git |

---

## Bleeding Edge

¿Querés las últimas funcionalidades experimentales de mi workflow diario (solo macOS)?

Mirá la rama [`nix-migration`](https://github.com/Gentleman-Programming/Gentleman.Dots/tree/nix-migration).

Contiene configuraciones de vanguardia que luego pasan a `main` cuando se estabilizan.

---

## Estructura del proyecto

```
Gentleman.Dots/
├── installer/               # Instalador TUI en Go
│   ├── cmd/                 # Punto de entrada
│   ├── internal/            # TUI, sistema y entrenador
│   └── e2e/                 # Tests E2E con Docker
├── docs/                    # Documentación
├── skills/                  # Skills de agentes IA
│
├── GentlemanNvim/           # Configuración Neovim
├── GentlemanClaude/         # Config Claude Code + skills
├── GentlemanOpenCode/       # Config OpenCode
│
├── GentlemanFish/
├── GentlemanZsh/
├── GentlemanNushell/
├── GentlemanTmux/
├── GentlemanZellij/
│
├── GentlemanGhostty/
├── GentlemanKitty/
├── alacritty.toml
├── .wezterm.lua
│
└── starship.toml
```

---

## Soporte

* **Issues**: GitHub Issues
* **Discord**: Gentleman Programming Community
* **YouTube**: @GentlemanProgramming
* **Twitch**: GentlemanProgramming

---

## Licencia

Licencia MIT — libre de usar, modificar y compartir.

**¡Feliz coding!** 🎩
