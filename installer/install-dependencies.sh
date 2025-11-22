#!/bin/bash

set -e

YELLOW=$(tput setaf 221)
GREEN=$(tput setaf 114)
BOLD=$(tput bold)
NC=$(tput sgr0)

info()    { echo -e "${YELLOW}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[OK]${NC} $1"; }

info "Instalando dependencias esenciales (build-essential, curl, file, git, unzip)..."
apt update
apt install -y build-essential curl file git unzip
success "Dependencias instaladas correctamente."

# Mostrar cómo ejecutar este script directamente desde GitHub
echo -e "\n${YELLOW}${BOLD}Puedes ejecutar este instalador directamente con:${NC}"
echo -e "${YELLOW}${BOLD}curl -fsSL $REPO_URL | bash${NC}"

YELLOW=$(tput setaf 221)
BOLD=$(tput bold)
NC=$(tput sgr0)
# ----------------------------------------------------------------------------

# Ejecutar instaladores de Node.js y Neovim automáticamente
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
# ----------------------------------------------------------------------------
info "Instalando Node.js LTS..."
bash "$SCRIPT_DIR/install-node.sh" || { error "Error instalando Node.js"; exit 1; }
# ----------------------------------------------------------------------------
info "Instalando Neovim..."
bash "$SCRIPT_DIR/install-nvim.sh" || { error "Error instalando Neovim"; exit 1; }
# ----------------------------------------------------------------------------
info "Instalando Go..."
bash "$SCRIPT_DIR/install-go.sh" || { error "Error instalando Go"; exit 1; }
# ----------------------------------------------------------------------------
info "Instalando Obsidian..."
bash "$SCRIPT_DIR/install-obsidian.sh" || { error "Error instalando Obsidian"; exit 1; }
# ----------------------------------------------------------------------------
# === INSTALACIÓN DE AI CLI TOOLS ===
info "Verificando CLI de Claude Code..."
if ! command -v claude-code &>/dev/null; then
	info "Instalando Claude Code CLI..."
	CLAUDE_DEB_URL="https://github.com/anthropic-ai/claude-code/releases/latest/download/claude-code_linux_amd64.deb"
	curl -fsSL -o /tmp/claude-code.deb "$CLAUDE_DEB_URL"
	apt-get install -y /tmp/claude-code.deb && success "Claude Code CLI instalado."
	rm -f /tmp/claude-code.deb
else
	success "Claude Code CLI ya está instalado."
fi

info "Verificando CLI de OpenCode..."
if ! command -v opencode &>/dev/null; then
	info "Instalando OpenCode CLI..."
	curl -fsSL https://opencode.ai/install | bash && success "OpenCode CLI instalado."
else
	success "OpenCode CLI ya está instalado."
fi

info "Verificando Gemini CLI..."
if ! command -v gemini &>/dev/null; then
	info "Instalando Gemini CLI..."
	npm install -g @google/gemini-cli && success "Gemini CLI instalado."
else
	success "Gemini CLI ya está instalado."
fi
# ----------------------------------------------------------------------------
