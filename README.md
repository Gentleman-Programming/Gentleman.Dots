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

Este archivo configura el esquema de colores utilizando el plugin [nvim](https://github.com/catppuccin/nvim). Se elige el tema "kanagawa-dragon" con opciones específicas, también puedes elegir catppucin o modus. Para elegir un theme, solo cambiar la property ```colorscheme = "kanagawa-dragon"``` por el nombre del theme que quieras.
Si deseas tener un background transparente, haz ```:TransparentEnable``` y te quedará NVIM con una opacidad extra de acuerdo a la que pongas en tu terminal.

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
      -- For more plugins integrations please scroll down (https://github.com/catppuccin/nvim#integrations)
    },
  },
  {
    "miikanissi/modus-themes.nvim",
    name = "modus",
    priority = 1000,
  },
  {
    "rebelot/kanagawa.nvim",
    name = "kanagawa",
    opts = {
      transparent_background = true,
    },
    priority = 1000,
  },
  {
    "xiyaowong/transparent.nvim",
  },
  {
    "LazyVim/LazyVim",
    opts = {
      colorscheme = "kanagawa-dragon",
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

## Carpeta `GentlemanKitty`

### Archivo `kanagawa.nvim`

Este archivo configura el tema Kanagawa en la terminal Kitty, proporcionando ajustes visuales y atajos de teclado para la navegación entre pestañas.

```vim
# vim:fileencoding=utf-8:foldmethod=marker

#: Fuentes {{{

font_family      IosevkaTerm Nerd Font
font_size 14.0

#: Los colores de primer plano y fondo.

background_opacity 0.95
# background_blur 0

## nombre: Kanagawa
## licencia: MIT
## autor: Tommaso Laurenzi
## upstream: https://github.com/rebelot/kanagawa.nvim/

background #0d0c0c
foreground #DCD7BA
selection_background #2D4F67
selection_foreground #C8C093
url_color #72A7BC
cursor #C8C093

# Pestañas
active_tab_background #1F1F28
active_tab_foreground #C8C093
inactive_tab_background  #1F1F28
inactive_tab_foreground #727169
#tab_bar_background #15161E

# normal
color0 #16161D
color1 #C34043
color2 #76946A
color3 #C0A36E
color4 #7E9CD8
color5 #957FB8
color6 #6A9589
color7 #C8C093

# brillante
color8  #727169
color9  #E82424
color10 #98BB6C
color11 #E6C384
color12 #7FB4CA
color13 #938AA9
color14 #7AA89F
color15 #DCD7BA


# colores extendidos
color16 #FFA066
color17 #FF5D62
  

map cmd+1 goto_tab 1
map cmd+2 goto_tab 2
map cmd+3 goto_tab 3
map cmd+4 goto_tab 4
map cmd+5 goto_tab 5
map cmd+6 goto_tab 6
map cmd+7 goto_tab 7
map cmd+8 goto_tab 8
map cmd+9 goto_tab 9
```

Este archivo proporciona la configuración del tema Kanagawa en Neovim, utilizando la fuente IosevkaTerm Nerd Font con un tamaño de fuente de 14.0. Además, define una paleta de colores cuidadosamente seleccionada para mejorar la experiencia de codificación. Los ajustes incluyen el estilo de pestañas para pestañas activas e inactivas, junto con asignaciones de teclas para la navegación rápida entre pestañas.

### Transpaso de configuraciones
```bash
git clone https://github.com/Gentleman-Programming/Gentleman.Dots
cp -r Gentleman.Dots/GentlemanKitty/* ~/.config/kitty
```

**Detalles del Tema:**
- **Nombre:** Kanagawa
- **Autor:** Tommaso Laurenzi
- **Licencia:** MIT
- **Repositorio Upstream:** [Kanagawa.nvim](https://github.com/rebelot/kanagawa.nvim/)

**Nota:** Las asignaciones de teclas proporcionadas para la navegación entre pestañas están configuradas como `cmd+1` hasta `cmd+9`.

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

### Instalación de Oh My Fish 

```curl https://raw.githubusercontent.com/oh-my-fish/oh-my-fish/master/bin/install | fish```

### Transpaso de Configuraciones

```bash
git clone https://github.com/Gentleman-Programming/Gentleman.Dots
cp -r Gentleman.Dots/GentlemanFish/* ~/.config
```

### Configura path para carpetas de trabajo del plugin PJ de Oh My Fish

Ve al archivo `~/.config/fish/fish_variables` y cambia la siguiente variable por la ruta a tu carpeta de trabajo con tus projectos:

```SETUVAR --export PROJECT_PATHS: /TuRutaDeTrabajo```

### Elegir theme Kanagawa para Fish

Ejecuta: ```fish_config theme save Kanagawa```

Y cuando pregunte si deseas sobre escribir: ```Y``` y luego darle enter

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

### Cargamos los plugins de Tmux
```bash
<Ctrl-b> + I para cargar los plugins
```

### Si quieres que Tmux se ejecute de manera por defecto al abir la terminal

#### Abre `~/.config/fish/config.fish` y agrega la siguiente línea al final:

```bash
if status is-interactive
    and not set -q TMUX
    exec tmux
end
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
   set -g @plugin 'christoomey/vim-tmux-navigator'
   ```
   - Otros plugins utilizados, como el Plugin Manager de Tmux (`tpm`) y plugins sensibles por defecto.

   - Cabe destacar tmux-resurrect el cual guarda el estado de la session para que no lo perdamos, se utiliza mediante:
    ```bash
    <Ctrl-b> + <Ctrl-s> para guardar el estado
    <Ctrl-b> + <Ctrl-r> para recuperar el estado
    ```

   - Cabe destacar vim-tmux-navigator el cual permite cambiar entre splits de vim y tmux indistintivamente mediante la utilización de `<Ctrl-h/j/k/l>`:
    ```bash
    set -g @plugin 'christoomey/vim-tmux-navigator'
    # Smart pane switching with awareness of vim splits
    bind -n C-k run-shell 'tmux-vim-select-pane -U'
    bind -n C-j run-shell 'tmux-vim-select-pane -D'
    bind -n C-h run-shell 'tmux-vim-select-pane -L'
    bind -n C-l run-shell 'tmux-vim-select-pane -R'
    bind -n "C-\\" run-shell 'tmux-vim-select-pane -l'
    
    # Bring back clear screen under tmux prefix
    bind C-l send-keys 'C-l'

    ```
   - Configuración del tipo de terminal predeterminado y algunas configuraciones adicionales para la terminación.
   ```bash
   set -g default-terminal "tmux-256color"
   set-option -ga terminal-overrides ",xterm*:Tc"
  ``` 

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

¡Disfruta de tu nuevo entorno de desarrollo en Neovim!
