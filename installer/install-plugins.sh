#!/bin/bash

# === VARIABLES PARAMETRIZABLES ===
# URLs y rutas principales
DOTFILES_REPO="https://github.com/villcabo/marckv.dots.git"
DOTFILES_BRANCH="main"
TMP_DOTFILES="$HOME/.tmp-dotfiles"
NVIM_SRC_DIR="$TMP_DOTFILES/GentlemanNvim/nvim"
NVIM_DEST_DIR="$HOME/.config/nvim"

# Colores para output usando tput (256 colores)
PINK=$(tput setaf 204)
PURPLE=$(tput setaf 141)
GREEN=$(tput setaf 114)
ORANGE=$(tput setaf 208)
BLUE=$(tput setaf 75)
YELLOW=$(tput setaf 221)
RED=$(tput setaf 196)
BOLD=$(tput bold)
NC=$(tput sgr0) # No Color

# Función para imprimir mensajes
info()    { echo "${BLUE}[INFO]${NC} $1"; }
success() { echo "${GREEN}[OK]${NC} $1"; }
error()   { echo "${RED}[ERROR]${NC} $1"; }
warn()    { echo "${ORANGE}[WARN]${NC} $1"; }
bold()    { echo "${BOLD}$1${NC}"; }

# Función para controlar errores
die() {
    error "$1"
    exit 1
}

# Clonar dotfiles y preparar configuración de Neovim
prepare_nvim_config() {
    local DOTS_DIR="$HOME/.marckv.dots"
    local NVIM_SRC="$DOTS_DIR/GentlemanNvim/nvim"
    local NVIM_DEST="$HOME/.config/nvim"
    if [ ! -d "$DOTS_DIR" ]; then
        info "Clonando dotfiles en $DOTS_DIR..."
        git clone --depth=1 --branch "$DOTFILES_BRANCH" "$DOTFILES_REPO" "$DOTS_DIR" || die "No se pudo clonar el repositorio de dotfiles."
    else
        info "Dotfiles ya existen en $DOTS_DIR."
    fi
    if [ -d "$NVIM_SRC" ]; then
        if [ -L "$NVIM_DEST" ] || [ -d "$NVIM_DEST" ]; then
            info "Eliminando configuración previa de Neovim en $NVIM_DEST..."
            rm -rf "$NVIM_DEST"
        fi
        # Crear ~/.config si no existe
        mkdir -p "$(dirname "$NVIM_DEST")"
        info "Enlazando $NVIM_SRC a $NVIM_DEST..."
        ln -s "$NVIM_SRC" "$NVIM_DEST" || die "No se pudo crear el enlace simbólico para Neovim."
        success "Configuración de Neovim lista en ${YELLOW}${BOLD}$NVIM_DEST${NC}"
    else
        die "No se encontró $NVIM_SRC en el repositorio de dotfiles."
    fi
}

# Instalar plugins de Neovim
install_nvim_plugins() {
    if ! command -v nvim >/dev/null 2>&1; then
        die "Neovim no está instalado. Ejecute primero el instalador de Neovim."
    fi
    info "Instalando plugins de Neovim..."
    nvim --headless "+Lazy! sync" +qa || die "Error al instalar plugins de Neovim (Lazy.nvim)."
    success "Plugins de Neovim instalados correctamente."
}

bold "=== Instalador de configuración/plugins de Neovim ==="
prepare_nvim_config
install_nvim_plugins
