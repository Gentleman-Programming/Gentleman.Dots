# Gentleman.Dots

## Descripción

Este repositorio contiene configuraciones personalizadas para el entorno de desarrollo en Neovim, incluyendo plugins específicos y keymaps para mejorar la productividad. Se hace uso de [LazyVim](https://github.com/LazyVim/LazyVim) como un conjunto preconfigurado de plugins y ajustes para facilitar el uso de Neovim.

## Carpeta `GentlemanNvim`

### Transpaso de configuraciónes

```bash
git clone https://github.com/Gentleman-Programming/Gentleman.Dots
cp -r Gentleman.Dots/GentlemanNvim/* ~/.config
```

Reinicia Neovim para aplicar los cambios.

### Carpeta `plugins`

#### Archivo `codeium.lua`

Este archivo configura el plugin [codeium.vim](https://github.com/Exafunction/codeium.vim), proporcionando atajos de teclado para aceptar, completar y limpiar sugerencias.

```lua
return {
  "Exafunction/codeium.vim",
  config = function()
    vim.keymap.set("i", "<C-g>", function()
      return vim.fn["codeium#Accept"]()
    end, { expr = true })

    vim.keymap.set("i", "<C-l>", function()
      return vim.fn["codeium#CycleCompletions"](1)
    end, { expr = true })

    vim.keymap.set("i", "<C-M>", function()
      return vim.fn["codeium#Complete"]()
    end, { expr = true })

    vim.keymap.set("i", "<C-x>", function()
      return vim.fn["codeium#Clear"]()
    end, { expr = true })
  end,
}
```

#### Archivo `colorscheme.lua`

Este archivo configura el esquema de colores utilizando el plugin [nvim](https://github.com/catppuccin/nvim). Se elige el tema "catppuccin" con opciones específicas.

```lua
return {
  {
    "catppuccin/nvim",
    name = "catppuccin",
    lazy = false,
    opts = {
      transparent_background = true,
      flavour = "mocha",
    },
    integrations = {
      cmp = true,
      gitsigns = true,
      nvimtree = true,
      treesitter = true,
      notify = false,
      mini = {
        enabled = true,
        indentscope_color = "",
      },
      -- Otras integraciones de plugins pueden encontrarse en el enlace proporcionado
    },
  },
  {
    "LazyVim/LazyVim",
    opts = {
      colorscheme = "catppuccin",
    },
  },
}
```

#### Archivo `editor.lua`

Este archivo configura varios plugins para mejorar la experiencia de edición, como resaltar patrones en archivos Markdown y herramientas de búsqueda avanzada con Telescope.

```lua
1. **mini.hipatterns**:
   - Plugin: `echasnovski/mini.hipatterns`
   - Evento: `BufReadPre`
   - Configuración:
     - Se configura un resaltador para los colores HSL. Los colores HSL en el código se resaltarán con un color de fondo que corresponda al color HSL.

2. **git.nvim**:
   - Plugin: `dinhhuy258/git.nvim`
   - Evento: `BufReadPre`
   - Configuración:
     - Se configuran los atajos de teclado para abrir una ventana de blame (`<Leader>gb`) y para abrir un archivo o carpeta en el repositorio de Git (`<Leader>go`).

3. **telescope.nvim**:
   - Plugin: `telescope.nvim`
   - Dependencias: `nvim-telescope/telescope-fzf-native.nvim` y `nvim-telescope/telescope-file-browser.nvim`
   - Configuración:
     - Se configuran una serie de atajos de teclado para varias funcionalidades, como buscar archivos, buscar una cadena en el directorio actual, listar búferes abiertos, listar etiquetas de ayuda y reanudar el selector de telescope anterior.
     - También se configura para abrir un navegador de archivos con el camino del búfer actual con `<Leader>sf`.
   - Configuración adicional:
     - Se configura para que los resultados se envuelvan, la estrategia de diseño sea horizontal, la posición del prompt sea en la parte superior y la estrategia de clasificación sea ascendente.
     - Se configura para que el selector de diagnósticos tenga el tema "ivy", el modo inicial sea "normal" y el corte de vista previa sea 9999.
     - Se configura para que el navegador de archivos tenga el tema "dropdown", secuestre netrw y se utilice en su lugar, y tenga sus propios mapeos.
```

#### Archivo `harpoon.lua`

Este archivo configura el plugin [harpoon](https://github.com/ThePrimeagen/harpoon) para facilitar la navegación entre archivos marcados.

```lua
return {
  "ThePrimeagen/harpoon",
  lazy = false,
  dependencies = {
    "nvim-lua/plenary.nvim",
  },
  config = true,
  keys = {
    { "<leader>hm", "<cmd>lua require('harpoon.mark').add_file()<cr>", desc = "Marcar archivo con harpoon" },
    { "<leader>hn", "<cmd>lua require('harpoon.ui').nav_next()<cr>", desc = "Ir al siguiente marcador de harpoon" },
    { "<leader>hp", "<cmd>lua require('harpoon.ui').nav_prev()<cr>", desc = "Ir al marcador de harpoon anterior" },
    { "<leader>ha", "<cmd>lua require('harpoon.ui').toggle_quick_menu()<cr>", desc = "Mostrar marcadores de harpoon" },
  },
}
```

#### Archivo `telescope.lua`

Este archivo configura el plugin [Telescope](https://github.com/nvim-telescope/telescope.nvim) para realizar búsquedas avanzadas en archivos y otros recursos.

```lua
return {
  "nvim-telescope/telescope.nvim",
  opts = {
    defaults = {
      layout_strategy = "vertical",
      layout_config = { preview_cutoff = 6 },
    },
  },
}
```

#### Archivo `ui.lua`

Este archivo configura varios plugins para mejorar la interfaz de usuario, incluyendo notificaciones, animaciones, líneas de buffers y líneas de estado.

```lua

1. **Noice.nvim**:
   - Plugin: `folke/noice.nvim`
   - Configuración:
     - Se añade una ruta a la configuración de noice para filtrar mensajes de notificación con el texto "No information available". Estos mensajes se omitirán.
     - Se establecen autocmds para detectar cuando la ventana de Neovim gana o pierde foco. Esto se utiliza para determinar si la interfaz está enfocada o no.
     - Se añade una ruta adicional para mostrar notificaciones en el sistema cuando Neovim pierde el foco.

2. **Nvim-notify**:
   - Plugin: `rcarriga/nvim-notify`
   - Configuración:
     - Se establece el color de fondo y el tiempo de espera para las notificaciones.

3. **Mini.animate**:
   - Plugin: `echasnovski/mini.animate`
   - Configuración:
     - Se deshabilita la animación de desplazamiento (`scroll`).

4. **Bufferline.nvim**:
   - Plugin: `akinsho/bufferline.nvim`
   - Configuración:
     - Se definen atajos de teclado para cambiar entre pestañas.
     - Se configuran opciones para mostrar iconos y pestañas de cierre.

5. **Lualine.nvim**:
   - Plugin: `nvim-lualine/lualine.nvim`
   - Configuración:
     - Se configura el tema "catppuccin" para la línea de estado (statusline).

6. **Incline.nvim**:
   - Plugin: `b0o/incline.nvim`
   - Configuración:
     - Se configuran colores y opciones visuales para el resaltado de nombres de archivo en la línea de estado.

7. **Zen-mode.nvim**:
   - Plugin: `folke/zen-mode.nvim`
   - Configuración:
     - Se configuran atajos de teclado para activar el "Modo Zen", que oculta elementos de la interfaz de usuario para centrarse en la edición de texto.

8. **Dashboard-nvim**:
   - Plugin: `nvimdev/dashboard-nvim`
   - Configuración:
     - Se establece un logo personalizado para el tablero de inicio de Neovim.

```

### Archivo `keymaps.lua`

Este archivo define algunas keymaps personalizadas para mejorar la navegación y manipulación del texto en modo insertar.

```lua
vim.keymap.set("i", "<C-d>", "<C-d>zz")
vim.keymap.set("i", "<C-u>", "<C-u>zz")
vim.keymap.set("i", "<C-b>", "<C-o>de")
```
## Carpeta `GentlemanTmux`

Contiene configuraciones para el entorno de tmux, para instalarlo e utilizarlo se debe realizar la siguiente serie de pasos:

### Instalación de Tmux

#### HomeBrew (recomendado)
```brew install tmux```

#### Ubuntu/Debian
```
sudo apt-get update
sudo apt-get install tmux
```
#### Fedora
```sudo dnf -y install tmux```


### Transpaso de configuraciones
```bash
git clone https://github.com/Gentleman-Programming/Gentleman.Dots
cp -r Gentleman.Dots/GentlemanTmux/* ~/
```

### Iniciar Tmux

#### Lo ponemos en marcha
```bash
tmux
```
#### Cargamos la configuración
```bash
tmux source-file ~/.tmux.conf
```

### Explicación de la configuración 


1. **Configuración del Shell Predeterminado:**
   ```bash
   set-option -g default-shell /usr/bin/fish
   ```
   Establece el shell predeterminado que Tmux utilizará como `/usr/bin/fish`.

2. **Configuración de Plugins:**
   ```bash
   set -g @plugin 'catppuccin/tmux'
   set -g @catppuccin_flavour 'macchiato'
   ```
   - Se utiliza el plugin 'catppuccin/tmux'.
   - Se configura el sabor (`flavour`) del plugin como 'macchiato'.

   ```bash
   set -g @plugin 'tmux-plugins/tpm'
   set -g @plugin 'tmux-plugins/tmux-sensible'
   set -g @plugin 'tmux-plugins/tmux-resurrect'
   ```
   - Otros plugins utilizados, como el Plugin Manager de Tmux (`tpm`) y plugins sensibles por defecto.

   ```bash
   set -g default-terminal "tmux-256color"
   set-option -ga terminal-overrides ",xterm*:Tc"
   ```
   - Configuración del tipo de terminal predeterminado y algunas configuraciones adicionales para la terminación.

3. **Configuración de la Apariencia de las Ventanas y Paneles:**
   ```bash
   set -g @catppuccin_window_left_separator "█"
   set -g @catppuccin_window_right_separator "█ "
   set -g @catppuccin_window_number_position "right"
   set -g @catppuccin_window_middle_separator "  █"
   set -g @catppuccin_window_default_fill "number"
   set -g @catppuccin_window_current_fill "number"
   set -g @catppuccin_window_current_text "#{pane_current_path}"
   ```
   - Configuración de la apariencia de las ventanas y paneles, incluidos separadores, posición de números y texto de la ventana actual.

4. **Configuración de la Barra de Estado (Status Bar):**
   ```bash
   set -g @catppuccin_status_modules "application session date_time"
   set -g @catppuccin_status_left_separator  ""
   set -g @catppuccin_status_right_separator " "
   set -g @catppuccin_status_right_separator_inverse "yes"
   set -g @catppuccin_status_fill "all"
   set -g @catppuccin_status_connect_separator "no"
   set -g @catppuccin_directory_text "#{pane_current_path}"
   ```
   - Configuración de módulos y apariencia de la barra de estado, incluidos separadores y texto del directorio actual.

5. **Inicialización del Tmux Plugin Manager (TPM):**
   ```bash
   run '~/.tmux/plugins/tpm/tpm'
   ```
   - Inicia el Tmux Plugin Manager. Este comando debe mantenerse al final del archivo de configuración de Tmux.

## Carpeta `GentlemanFish`

### Instalación de fish

#### HomeBrew (recomendado)
```brew install fish```

#### Ubuntu/Debian
```
sudo apt-get update
sudo apt-get install fish
```
#### Fedora
```sudo dnf install fish```

### Instalación de Oh My Fish 

```curl https://raw.githubusercontent.com/oh-my-fish/oh-my-fish/master/bin/install | fish```

### Transpaso de Configuraciones

```bash
git clone https://github.com/Gentleman-Programming/Gentleman.Dots
cp -r Gentleman.Dots/GentlemanFish/* ~/.config
```

### Configura path para carpetas de trabajo del plugin PJ de Oh My Fish

Ve al archivo `~/.config/fish/fish_variables` y cambia la siguiente variable por la ruta a tu carpeta de trabajo con tus projectos:

```SETUVAR --export PROJECT_PATHS: /TuRutaDeTrabajo```

¡Disfruta de tu nuevo entorno de desarrollo en Neovim!
