# Guía de Instalación: Neovim para Registros de Enfermería (N4N)

## 📋 Introducción

Esta guía te ayudará a instalar y configurar **GentlemanNvim** como una herramienta poderosa para crear registros de enfermería, aprovechando la velocidad de edición de Neovim combinada con plugins especializados para redacción de texto médico.

### ¿Por qué Neovim para Enfermería?

- **⚡ Velocidad**: Edición ultra-rápida sin tocar el mouse
- **📝 Plantillas**: Snippets personalizables para notas médicas recurrentes
- **🔍 Búsqueda potente**: Encuentra información instantáneamente en tus registros
- **💾 Sincronización**: Integración con Obsidian para tus notas clínicas
- **🤖 IA integrada**: Claude Code para asistencia en redacción
- **🌐 Portabilidad**: Especialmente en Windows, lleva tu entorno en un USB

---

## 🪟 Instalación para Windows (Modo Portátil)

Esta configuración **NO requiere privilegios de administrador** y puede ejecutarse desde una carpeta o USB.

### Paso 1: Preparar la Estructura de Carpetas

Crea una carpeta base donde instalarás todo (ejemplo: `C:\PortableDevTools` o `E:\PortableDevTools`):

```powershell
# En PowerShell o CMD
mkdir C:\PortableDevTools
cd C:\PortableDevTools
```

### Paso 2: Descargar e Instalar Git Portable

**🔗 Link de descarga**: [Git for Windows Portable](https://github.com/git-for-windows/git/releases/latest)

1. Descarga el archivo: `PortableGit-X.XX.X-64-bit.7z.exe` (ejemplo: `PortableGit-2.43.0-64-bit.7z.exe`)
2. Ejecuta el archivo en la carpeta `C:\PortableDevTools\Git`
3. Espera a que extraiga todos los archivos

**Verificación**:
```powershell
C:\PortableDevTools\Git\bin\git.exe --version
# Debería mostrar: git version 2.43.0 (o superior)
```

### Paso 3: Descargar e Instalar Node.js Portable

**🔗 Link de descarga**: [Node.js Downloads](https://nodejs.org/en/download/prebuilt-binaries)

1. Ve a la página de descargas de Node.js
2. Descarga la versión **LTS** (Long Term Support): `node-v22.x.x-win-x64.zip`
3. Extrae el archivo ZIP en `C:\PortableDevTools\NodeJS`

**Estructura resultante**:
```
C:\PortableDevTools\NodeJS\
├── node.exe
├── npm
├── npm.cmd
└── node_modules\
```

**Verificación**:
```powershell
C:\PortableDevTools\NodeJS\node.exe --version
# Debería mostrar: v22.x.x

C:\PortableDevTools\NodeJS\npm.cmd --version
# Debería mostrar: 10.x.x
```

### Paso 4: Descargar e Instalar Neovim Portable

**🔗 Link de descarga**: [Neovim Releases](https://github.com/neovim/neovim/releases/latest)

1. Descarga el archivo: `nvim-win64.zip`
2. Extrae el contenido en `C:\PortableDevTools\Neovim`

**Estructura resultante**:
```
C:\PortableDevTools\Neovim\
├── bin\
│   └── nvim.exe
├── lib\
└── share\
```

**Verificación**:
```powershell
C:\PortableDevTools\Neovim\bin\nvim.exe --version
# Debería mostrar: NVIM v0.10.x
```

### Paso 5: Descargar Fuente Nerd Font

**🔗 Link de descarga**: [Iosevka Term Nerd Font](https://github.com/ryanoasis/nerd-fonts/releases/latest)

1. Descarga el archivo: `IosevkaTerm.zip`
2. Extrae el ZIP
3. Selecciona todos los archivos `.ttf`
4. Click derecho → **Instalar** (si tienes permisos) o **Instalar para el usuario actual**

**Fuentes alternativas** (si prefieres otra):
- [JetBrainsMono Nerd Font](https://github.com/ryanoasis/nerd-fonts/releases/latest) - `JetBrainsMono.zip`
- [FiraCode Nerd Font](https://github.com/ryanoasis/nerd-fonts/releases/latest) - `FiraCode.zip`

### Paso 6: Descargar Ripgrep (para búsquedas)

**🔗 Link de descarga**: [Ripgrep Releases](https://github.com/BurntSushi/ripgrep/releases/latest)

1. Descarga el archivo: `ripgrep-XX.X.X-x86_64-pc-windows-msvc.zip`
2. Extrae el contenido en `C:\PortableDevTools\ripgrep`

### Paso 7: Descargar fd (para búsqueda de archivos)

**🔗 Link de descarga**: [fd Releases](https://github.com/sharkdp/fd/releases/latest)

1. Descarga el archivo: `fd-vX.X.X-x86_64-pc-windows-msvc.zip`
2. Extrae el contenido en `C:\PortableDevTools\fd`

### Paso 8: Configurar Variables de Entorno Portables

Crea un archivo `launch-nvim.bat` en `C:\PortableDevTools\`:

```batch
@echo off
REM === Configuración de Rutas Portables ===
set PORTABLE_ROOT=%~dp0

REM Git
set PATH=%PORTABLE_ROOT%Git\bin;%PATH%

REM Node.js y npm
set PATH=%PORTABLE_ROOT%NodeJS;%PATH%

REM Neovim
set PATH=%PORTABLE_ROOT%Neovim\bin;%PATH%

REM ripgrep
set PATH=%PORTABLE_ROOT%ripgrep;%PATH%

REM fd
set PATH=%PORTABLE_ROOT%fd;%PATH%

REM Configuración de Neovim
set XDG_CONFIG_HOME=%PORTABLE_ROOT%config
set XDG_DATA_HOME=%PORTABLE_ROOT%data
set XDG_STATE_HOME=%PORTABLE_ROOT%state
set XDG_CACHE_HOME=%PORTABLE_ROOT%cache

REM Crear carpetas si no existen
if not exist "%XDG_CONFIG_HOME%" mkdir "%XDG_CONFIG_HOME%"
if not exist "%XDG_DATA_HOME%" mkdir "%XDG_DATA_HOME%"
if not exist "%XDG_STATE_HOME%" mkdir "%XDG_STATE_HOME%"
if not exist "%XDG_CACHE_HOME%" mkdir "%XDG_CACHE_HOME%"

REM Lanzar PowerShell con el entorno configurado
echo ======================================
echo  Entorno Portable Neovim Cargado
echo ======================================
echo.
echo Git:     %PORTABLE_ROOT%Git\bin
echo Node:    %PORTABLE_ROOT%NodeJS
echo Neovim:  %PORTABLE_ROOT%Neovim\bin
echo Config:  %XDG_CONFIG_HOME%
echo.
echo Ejecuta 'nvim' para iniciar Neovim
echo ======================================
echo.

powershell.exe -NoExit
```

**Para usar**: Haz doble click en `launch-nvim.bat` cada vez que quieras usar Neovim.

### Paso 9: Clonar la Configuración GentlemanNvim

Abre el archivo `launch-nvim.bat` que acabas de crear y ejecuta:

```powershell
# Navega a la carpeta de configuración
cd $env:XDG_CONFIG_HOME

# Clona el repositorio
git clone https://github.com/fegome90-cmd/n4n.dots.git temp-clone

# Copia solo la carpeta de Neovim
cp -r temp-clone/GentlemanNvim/nvim nvim

# Limpia el temporal
rm -r -Force temp-clone

# Verifica la instalación
ls nvim
```

### Paso 10: Primera Ejecución y Instalación de Plugins

```powershell
# Inicia Neovim por primera vez
nvim

# Neovim automáticamente:
# 1. Instalará Lazy.nvim (el gestor de plugins)
# 2. Descargará e instalará todos los plugins
# 3. Configurará LSP servers vía Mason
#
# Esto puede tomar 3-5 minutos la primera vez
# Verás una ventana con barras de progreso
```

**Si algo sale mal**, presiona `:` y escribe:
```vim
:Lazy sync
```

### Paso 11: Instalar Herramientas Adicionales con Mason

Dentro de Neovim, presiona `:` y ejecuta:

```vim
:Mason
```

Instala las siguientes herramientas (usa `/` para buscar, `Enter` para instalar):

**Para texto y markdown**:
- `marksman` (LSP para Markdown)
- `markdown-toc` (Tabla de contenidos)
- `markdownlint` (Linter para Markdown)
- `prettier` (Formateador)

**Diccionarios y corrección**:
- `vale` (Linter de prosa para escritura médica)
- `ltex-ls` (LanguageTool para corrección gramatical)

---

## 🍎 Instalación para macOS

La instalación en Mac es más sencilla gracias a Homebrew.

### Paso 1: Instalar Homebrew (si no lo tienes)

Abre Terminal y ejecuta:

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

Sigue las instrucciones en pantalla. Al finalizar, agrega Homebrew al PATH:

```bash
# Para Apple Silicon (M1/M2/M3)
echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zprofile
eval "$(/opt/homebrew/bin/brew shellenv)"

# Para Intel Macs
echo 'eval "$(/usr/local/bin/brew shellenv)"' >> ~/.zprofile
eval "$(/usr/local/bin/brew shellenv)"
```

### Paso 2: Instalar Todas las Dependencias

```bash
# Neovim
brew install neovim

# Git (ya viene en macOS, pero actualiza a la última versión)
brew install git

# Node.js (versión LTS)
brew install node@22

# Herramientas de búsqueda
brew install ripgrep fd fzf

# Utilidades adicionales
brew install tree-sitter lazygit
```

### Paso 3: Instalar Fuente Nerd Font

```bash
# Instalar Iosevka Term Nerd Font
brew tap homebrew/cask-fonts
brew install font-iosevka-term-nerd-font

# Alternativas
# brew install font-jetbrains-mono-nerd-font
# brew install font-fira-code-nerd-font
```

### Paso 4: Clonar la Configuración

```bash
# Backup de configuración existente (si la tienes)
[ -d ~/.config/nvim ] && mv ~/.config/nvim ~/.config/nvim.backup.$(date +%Y%m%d)

# Crear carpeta de configuración
mkdir -p ~/.config

# Clonar el repositorio
git clone https://github.com/fegome90-cmd/n4n.dots.git ~/temp-n4n-dots

# Copiar solo la configuración de Neovim
cp -r ~/temp-n4n-dots/GentlemanNvim/nvim ~/.config/nvim

# Limpiar temporal
rm -rf ~/temp-n4n-dots
```

### Paso 5: Primera Ejecución

```bash
# Abre Neovim
nvim

# Lazy.nvim se instalará automáticamente junto con todos los plugins
# Espera 3-5 minutos mientras se descargan e instalan
```

### Paso 6: Instalar Herramientas con Mason

Dentro de Neovim, presiona `:` y ejecuta:

```vim
:Mason
```

Instala las siguientes herramientas (usa `/` para buscar, `i` para instalar):

**Para texto y markdown**:
- `marksman`
- `markdown-toc`
- `markdownlint`
- `prettier`

**Para escritura médica**:
- `vale`
- `ltex-ls`

---

## 📝 Configuración de Plugins para Redacción de Texto Médico

Esta configuración ya incluye varios plugins útiles para escritura. Aquí te explico cómo aprovecharlos para registros de enfermería.

### 1. Obsidian.nvim - Tu Sistema de Notas Clínicas

**Ya está instalado y configurado**. Para usarlo:

```vim
" Buscar o crear nota
<leader>of   " Find note (buscar nota existente)
<leader>on   " New note (crear nueva nota)
<leader>oq   " Quick search (búsqueda rápida en todas las notas)
<leader>ot   " Insert template (insertar plantilla)

" Ejemplo de uso:
" Presiona: <leader> on
" Escribe: Registro-Turno-Noche-2025-01-24
" Presiona Enter
```

**Configurar tu carpeta de notas**:

Edita el archivo `~/.config/nvim/lua/plugins/obsidian.lua`:

```lua
workspaces = {
  {
    name = "registros-enfermeria",
    path = "~/Documents/RegistrosEnfermeria", -- Cambia esta ruta
  },
},
```

### 2. Snippets Personalizados para Enfermería

Crea el archivo `~/.config/nvim/snippets/markdown.json`:

```json
{
  "Registro de Enfermería": {
    "prefix": "regenferm",
    "body": [
      "# Registro de Enfermería - ${1:Fecha}",
      "",
      "## Datos del Paciente",
      "- **Nombre**: ${2:Nombre completo}",
      "- **Edad**: ${3:Edad}",
      "- **Cama**: ${4:Número de cama}",
      "- **Diagnóstico**: ${5:Diagnóstico principal}",
      "",
      "## Signos Vitales",
      "- **Presión Arterial**: ${6:120/80} mmHg",
      "- **Frecuencia Cardíaca**: ${7:72} lpm",
      "- **Temperatura**: ${8:36.5}°C",
      "- **Saturación O2**: ${9:98}%",
      "- **Frecuencia Respiratoria**: ${10:16} rpm",
      "",
      "## Evaluación",
      "${11:Estado general del paciente...}",
      "",
      "## Intervenciones",
      "- ${12:Primera intervención}",
      "- ${13:Segunda intervención}",
      "",
      "## Observaciones",
      "${14:Observaciones adicionales...}",
      "",
      "---",
      "*Registrado por*: ${15:Tu nombre} | *Turno*: ${16:Mañana/Tarde/Noche}"
    ],
    "description": "Plantilla completa de registro de enfermería"
  },
  "Signos Vitales Rápido": {
    "prefix": "sv",
    "body": [
      "**Signos Vitales** (${1:HH:MM})",
      "- PA: ${2:120/80} mmHg | FC: ${3:72} lpm | T: ${4:36.5}°C | SatO2: ${5:98}% | FR: ${6:16} rpm"
    ],
    "description": "Entrada rápida de signos vitales"
  },
  "Administración de Medicamento": {
    "prefix": "med",
    "body": [
      "### Medicamento Administrado",
      "- **Hora**: ${1:HH:MM}",
      "- **Medicamento**: ${2:Nombre del medicamento}",
      "- **Dosis**: ${3:Dosis}",
      "- **Vía**: ${4:Oral/IV/IM/SC}",
      "- **Observaciones**: ${5:Sin reacciones adversas}"
    ],
    "description": "Registro de administración de medicamento"
  },
  "Nota de Evolución": {
    "prefix": "evol",
    "body": [
      "## Nota de Evolución - ${1:Fecha} ${2:Hora}",
      "",
      "**Subjetivo**: ${3:Lo que el paciente refiere...}",
      "",
      "**Objetivo**: ${4:Hallazgos objetivos...}",
      "",
      "**Análisis**: ${5:Interpretación...}",
      "",
      "**Plan**: ${6:Plan de cuidados...}"
    ],
    "description": "Nota SOAP para evolución del paciente"
  },
  "Incidente/Evento Adverso": {
    "prefix": "incidente",
    "body": [
      "# ⚠️ REPORTE DE INCIDENTE",
      "",
      "**Fecha y Hora**: ${1:YYYY-MM-DD HH:MM}",
      "**Paciente**: ${2:Nombre/Identificación}",
      "**Tipo de Incidente**: ${3:Caída/Medicación/Otro}",
      "",
      "## Descripción del Incidente",
      "${4:Descripción detallada de lo ocurrido...}",
      "",
      "## Acciones Inmediatas",
      "- ${5:Primera acción tomada}",
      "- ${6:Segunda acción tomada}",
      "",
      "## Notificaciones",
      "- Médico de guardia: ${7:Sí/No}",
      "- Supervisor: ${8:Sí/No}",
      "",
      "## Estado Actual del Paciente",
      "${9:Estado post-incidente...}",
      "",
      "**Reportado por**: ${10:Tu nombre}"
    ],
    "description": "Reporte de incidente o evento adverso"
  }
}
```

**Para usar los snippets**:
1. En modo INSERT, escribe el prefijo (ejemplo: `regenferm`)
2. Verás un popup de autocompletado
3. Presiona `Tab` para seleccionar
4. Navega entre campos con `Tab` y `Shift+Tab`

### 3. Spell Checking en Español

Agrega esto a tu `~/.config/nvim/lua/config/options.lua`:

```lua
-- Corrección ortográfica en español e inglés
vim.opt.spelllang = { "es", "en" }
vim.opt.spell = true

-- Palabras médicas personalizadas
vim.opt.spellfile = vim.fn.stdpath("config") .. "/spell/medical.utf-8.add"
```

**Para agregar términos médicos al diccionario**:
```vim
" Coloca el cursor sobre una palabra subrayada
zg    " Agregar palabra al diccionario personal
zug   " Quitar palabra del diccionario
]s    " Ir a siguiente error ortográfico
[s    " Ir a error anterior
z=    " Ver sugerencias de corrección
```

### 4. Vale - Linter de Prosa Médica

Crea un archivo `.vale.ini` en tu carpeta de registros:

```ini
StylesPath = .vale/styles

MinAlertLevel = suggestion

[*.md]
BasedOnStyles = write-good, proselint

# Ignorar ciertos términos médicos
Vale.Terms = NO
```

### 5. Plantillas con Fechas Automáticas

La configuración de Obsidian ya incluye plantillas. Crea archivos en `~/Documents/RegistrosEnfermeria/templates/`:

**`turno-diario.md`**:
```markdown
---
fecha: {{date}}
hora_inicio: {{time}}
turno: [Mañana/Tarde/Noche]
---

# Registro de Turno - {{date}}

## Pacientes Asignados
1. [ ] Paciente 1 - Cama X
2. [ ] Paciente 2 - Cama Y
3. [ ] Paciente 3 - Cama Z

## Pendientes del Turno
- [ ] Ronda inicial de signos vitales
- [ ] Administración de medicamentos 08:00
- [ ] Administración de medicamentos 12:00
- [ ] Ronda de evaluación
- [ ] Documentación completada

## Observaciones Generales
{{cursor}}

---
*Turno iniciado*: {{time}}
```

### 6. Atajos de Teclado Útiles para Escritura

```vim
" Ya configurados en GentlemanNvim:

" Navegación rápida
<leader>ff   " Buscar archivo
<leader>fg   " Buscar en contenido (grep)
<leader>fb   " Buscar en buffers abiertos

" Obsidian (tus notas)
<leader>of   " Buscar nota
<leader>on   " Nueva nota
<leader>oq   " Búsqueda rápida
<leader>ot   " Insertar plantilla

" Edición
gcc          " Comentar/descomentar línea
gbc          " Comentar/descomentar bloque (visual mode)
<leader>w    " Guardar archivo

" Claude Code (IA para ayuda)
<leader>ac   " Toggle Claude panel
<leader>af   " Focus en Claude
<leader>aa   " Enviar selección a Claude (visual mode)

" Markdown
gx           " Abrir link bajo cursor
[[           " Crear/seguir link interno
```

---

## 🏥 Flujo de Trabajo Recomendado para Registros de Enfermería

### Inicio de Turno

1. Abre Neovim: `nvim` (Windows) o desde Terminal (Mac)
2. Crea nota del turno: `<leader>on` → escribe `Turno-2025-01-24-Noche`
3. Inserta plantilla: `<leader>ot` → selecciona `turno-diario.md`
4. Completa la información básica

### Durante el Turno

**Para cada paciente**:
```vim
" Crear nota rápida
<leader>on  " → Paciente-101-Juan-Perez

" Dentro de la nota, usar snippets:
regenferm   " → Plantilla completa de registro
sv          " → Signos vitales rápidos
med         " → Administración de medicamento
evol        " → Nota de evolución
```

**Búsqueda rápida**:
```vim
" Buscar en todas las notas
<leader>oq  " → Escribe: "diabetes" o "Juan Perez"

" Buscar en carpeta actual
<leader>fg  " → Escribe término de búsqueda
```

### Fin de Turno

1. Revisa pendientes: `<leader>of` → abre nota del turno
2. Marca todas las tareas completadas (`[x]`)
3. Agrega observaciones finales
4. Guarda todo: `<leader>w`

### Casos Especiales

**Incidente**:
```vim
<leader>on       " → Incidente-2025-01-24-Caida-Cama15
incidente<Tab>   " → Usa el snippet de incidente
" Completa todos los campos con Tab
```

**Consulta con IA (Claude)**:
```vim
" Selecciona texto en visual mode (V o v)
<leader>aa  " → Envía a Claude para análisis

" Ejemplo: "Claude, revisa esta nota y sugiere mejoras"
```

---

## 🎯 Plugins Adicionales Recomendados

Si quieres expandir aún más las capacidades, considera agregar:

### 1. **vim-table-mode** - Tablas en Markdown

Edita `~/.config/nvim/lua/plugins/table-mode.lua`:

```lua
return {
  "dhruvasagar/vim-table-mode",
  ft = "markdown",
  config = function()
    vim.g.table_mode_corner = "|"
    vim.g.table_mode_corner_corner = "|"
    vim.g.table_mode_header_fillchar = "-"
  end,
  keys = {
    { "<leader>tm", "<cmd>TableModeToggle<cr>", desc = "Toggle Table Mode" },
  },
}
```

**Uso**:
```markdown
| Hora | PA | FC | T |
|-|-|-|-|
| 08:00 | 120/80 | 72 | 36.5 |
| 12:00 | 118/78 | 70 | 36.6 |
```

### 2. **headlines.nvim** - Mejor visualización de Markdown

Edita `~/.config/nvim/lua/plugins/headlines.lua`:

```lua
return {
  "lukas-reineke/headlines.nvim",
  dependencies = "nvim-treesitter/nvim-treesitter",
  ft = "markdown",
  config = function()
    require("headlines").setup({
      markdown = {
        headline_highlights = { "Headline1", "Headline2", "Headline3" },
        fat_headlines = true,
        fat_headline_upper_string = "▃",
        fat_headline_lower_string = "🬂",
      },
    })
  end,
}
```

### 3. **clipboard-image.nvim** - Pegar imágenes desde portapapeles

Útil para agregar fotos de heridas, equipos, etc.

```lua
return {
  "ekickx/clipboard-image.nvim",
  ft = "markdown",
  config = function()
    require("clipboard-image").setup({
      default = {
        img_dir = "images",
        img_name = function()
          return os.date("%Y-%m-%d-%H-%M-%S")
        end,
      },
    })
  end,
  keys = {
    { "<leader>pi", "<cmd>PasteImg<cr>", desc = "Paste Image from Clipboard" },
  },
}
```

---

## 🔧 Solución de Problemas Comunes

### Windows Portátil

**Problema**: "nvim no es reconocido como comando"
**Solución**: Asegúrate de ejecutar siempre `launch-nvim.bat` antes de usar nvim

**Problema**: "Error al instalar plugins"
**Solución**:
```vim
:Lazy clear
:Lazy sync
```

**Problema**: "Node.js no encontrado"
**Solución**: Verifica que `launch-nvim.bat` tenga las rutas correctas

### Mac

**Problema**: "neovim: command not found"
**Solución**:
```bash
brew install neovim
# O actualiza el PATH
echo 'export PATH="/opt/homebrew/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

**Problema**: Fuente no se muestra correctamente
**Solución**: Configura la fuente en tu terminal:
- Terminal.app: Preferencias → Perfiles → Texto → Fuente → "IosevkaTerm Nerd Font"
- iTerm2: Preferences → Profiles → Text → Font → "IosevkaTerm Nerd Font"

### Ambos Sistemas

**Problema**: Lazy.nvim no instala plugins
**Solución**:
```vim
:checkhealth lazy
:Lazy restore
```

**Problema**: LSP no funciona
**Solución**:
```vim
:checkhealth lsp
:Mason
" Reinstala el LSP que falla
```

---

## 📚 Recursos y Documentación

### Aprender Neovim
- **Tutor integrado**: Ejecuta `nvim +Tutor` en la terminal
- **Cheatsheet**: `:help quickref`
- [Neovim Documentation](https://neovim.io/doc/)

### Plugins Clave
- [Lazy.nvim](https://github.com/folke/lazy.nvim) - Gestor de plugins
- [Obsidian.nvim](https://github.com/epwalsh/obsidian.nvim) - Sistema de notas
- [Mason.nvim](https://github.com/williamboman/mason.nvim) - LSP/herramientas
- [Claude Code](https://github.com/anthropics/claude-code) - IA integrada

### Markdown
- [Markdown Guide](https://www.markdownguide.org/)
- [Obsidian Help](https://help.obsidian.md/)

---

## 🚀 Próximos Pasos

1. **Practica los atajos de teclado**: El poder de Neovim está en no usar el mouse
2. **Personaliza los snippets**: Adapta las plantillas a tu flujo de trabajo específico
3. **Integra con Obsidian**: Sincroniza tus notas entre dispositivos
4. **Explora Claude Code**: Usa IA para mejorar tus registros (`<leader>ac`)
5. **Crea tu diccionario médico**: Agrega términos con `zg` mientras escribes

---

## 💡 Tips Finales

- **Modo INSERT** (`i`): Para escribir
- **Modo NORMAL** (`Esc`): Para navegar y comandos
- **Modo VISUAL** (`v` o `V`): Para seleccionar texto
- **Guardar**: `<leader>w` o `:w`
- **Salir**: `:q` (o `:q!` para salir sin guardar)
- **Ayuda**: `:help <tema>` (ejemplo: `:help markdown`)

**Recuerda**: La curva de aprendizaje inicial vale la pena. En 2-3 semanas serás mucho más rápido que con cualquier editor tradicional.

---

**Documentación creada para el proyecto N4N (Neovim for Nurses)**
*Última actualización: 2025-01-24*
