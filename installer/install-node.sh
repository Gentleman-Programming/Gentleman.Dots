#!/bin/bash

# === VARIABLES PARAMETRIZABLES ===
NODE_VERSION="22.0.0"
NODE_DISTRO="linux-x64"
NODE_BASE_URL="https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-$NODE_DISTRO.tar.xz"
NODE_TAR="/tmp/node-v$NODE_VERSION-$NODE_DISTRO.tar.xz"
NODE_DIR="/opt/nodejs"
NODE_PROFILE="/etc/profile.d/node.sh"

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

info()    { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[OK]${NC} $1"; }
error()   { echo -e "${RED}[ERROR]${NC} $1"; }
warn()    { echo -e "${ORANGE}[WARN]${NC} $1"; }
bold()    { echo -e "${BOLD}$1${NC}"; }

die() {
    error "$1"
    exit 1
}

# Función para obtener la versión más reciente de Node.js LTS
get_latest_node_version() {
    local latest_version
    latest_version=$(curl -s https://nodejs.org/dist/index.json | grep -o '"version":"[^"]*"' | head -1 | cut -d'"' -f4 | sed 's/v//')
    
    if [ -z "$latest_version" ]; then
        return 1
    fi
    
    echo "$latest_version"
    return 0
}

# Función para obtener la versión instalada
get_installed_node_version() {
    if command -v node >/dev/null 2>&1; then
        local installed_version
        installed_version=$(node --version 2>/dev/null | sed 's/v//')
        if [ -n "$installed_version" ]; then
            echo "$installed_version"
            return 0
        fi
    fi
    return 1
}

# Función para comparar versiones
compare_node_versions() {
    local installed="$1"
    local latest="$2"
    
    # Usar sort -V para comparación de versiones
    local higher_version=$(printf '%s\n%s\n' "$installed" "$latest" | sort -V | tail -n1)
    
    if [ "$higher_version" = "$latest" ] && [ "$installed" != "$latest" ]; then
        return 0  # Hay una versión más nueva disponible
    else
        return 1  # Ya está actualizado
    fi
}

# Verificar si Node.js ya está instalado
check_existing_node_installation() {
    if command -v node >/dev/null 2>&1 && [ -f "$NODE_PROFILE" ]; then
        warn "Node.js ya está instalado en el sistema"
        local current_version=$(node --version)
        local installed_version=$(get_installed_node_version)
        info "Versión actual: ${BOLD}$current_version${NC}"
        
        # Verificar si hay una nueva versión disponible
        info "Verificando actualizaciones disponibles..."
        local latest_version=$(get_latest_node_version)
        
        if [ $? -eq 0 ] && [ -n "$latest_version" ]; then
            info "Última versión disponible: ${BOLD}v$latest_version${NC}"
            
            if compare_node_versions "$installed_version" "$latest_version"; then
                bold "\n🚀 ¡NUEVA VERSIÓN DISPONIBLE!"
                info "Versión instalada: ${YELLOW}v$installed_version${NC}"
                info "Versión disponible: ${GREEN}v$latest_version${NC}"
                warn "Se recomienda actualizar para obtener las últimas mejoras y correcciones"
                
                echo ""
                read -p "¿Desea actualizar a la última versión? (Y/n): " -n 1 -r
                echo
                if [[ $REPLY =~ ^[Nn]$ ]]; then
                    info "Actualización cancelada por el usuario"
                    exit 0
                fi
                success "Procediendo con la actualización..."
                # Actualizar variables para usar la última versión
                NODE_VERSION="$latest_version"
                NODE_BASE_URL="https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-$NODE_DISTRO.tar.xz"
                NODE_TAR="/tmp/node-v$NODE_VERSION-$NODE_DISTRO.tar.xz"
            else
                success "✅ Ya tienes la versión más reciente instalada"
                info "No es necesario actualizar"
                exit 0
            fi
        else
            warn "No se pudo verificar la versión más reciente"
            read -p "¿Desea reinstalar Node.js de todas formas? (y/N): " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                info "Instalación cancelada por el usuario"
                exit 0
            fi
            warn "Procediendo con la reinstalación..."
        fi
    else
        info "Node.js no está instalado en el sistema"
        info "Verificando la última versión disponible..."
        local latest_version=$(get_latest_node_version)
        
        if [ $? -eq 0 ] && [ -n "$latest_version" ]; then
            info "Se instalará la versión más reciente: ${BOLD}${GREEN}v$latest_version${NC}"
            # Actualizar variables para usar la última versión
            NODE_VERSION="$latest_version"
            NODE_BASE_URL="https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-$NODE_DISTRO.tar.xz"
            NODE_TAR="/tmp/node-v$NODE_VERSION-$NODE_DISTRO.tar.xz"
        else
            warn "No se pudo verificar la versión más reciente, se procederá con la versión predefinida: $NODE_VERSION"
        fi
    fi
}

install_node() {
    if [ ! -f "$NODE_TAR" ] || [ ! -s "$NODE_TAR" ]; then
        info "Descargando Node.js $NODE_VERSION en $NODE_TAR..."
        curl -L -o "$NODE_TAR" "$NODE_BASE_URL" || die "No se pudo descargar Node.js."
    else
        info "Usando archivo Node.js ya descargado en $NODE_TAR."
    fi
    # Verificar e instalar xz-utils si es necesario
    if ! command -v xz >/dev/null 2>&1; then
        info "Instalando xz-utils para descomprimir .tar.xz..."
        if command -v apt-get >/dev/null 2>&1; then
            apt-get update && apt-get install -y xz-utils || die "No se pudo instalar xz-utils."
        elif command -v dnf >/dev/null 2>&1; then
            dnf install -y xz || die "No se pudo instalar xz."
        elif command -v yum >/dev/null 2>&1; then
            yum install -y xz || die "No se pudo instalar xz."
        elif command -v pacman >/dev/null 2>&1; then
            pacman -Sy --noconfirm xz || die "No se pudo instalar xz."
        else
            die "No se pudo encontrar un gestor de paquetes compatible para instalar xz-utils."
        fi
    fi
    info "Extrayendo Node.js en $NODE_DIR..."
    rm -rf "$NODE_DIR"
    mkdir -p "$NODE_DIR"
    tar -xf "$NODE_TAR" -C "$NODE_DIR" --strip-components=1 || die "No se pudo extraer Node.js."
    info "Configurando PATH global para Node.js..."
    echo "export PATH=\"\$PATH:$NODE_DIR/bin\"" > "$NODE_PROFILE"
    chmod 644 "$NODE_PROFILE"
    success "Node.js $NODE_VERSION instalado correctamente en $NODE_DIR."
}

reload_shell_environment() {
    info "Recargando variables de entorno..."
    
    # Recargar el perfil de Node.js
    if [ -f "$NODE_PROFILE" ]; then
        source "$NODE_PROFILE"
        success "Variables de entorno de Node.js cargadas."
    fi
    
    # Verificar que Node.js esté disponible
    if command -v node >/dev/null 2>&1; then
        local node_version=$(node --version)
        success "Node.js está disponible: $node_version"
        
        # Verificar npm también
        if command -v npm >/dev/null 2>&1; then
            local npm_version=$(npm --version)
            success "npm está disponible: v$npm_version"
        fi
    else
        warn "Node.js no está disponible en el PATH actual."
        info "Puedes ejecutar: ${YELLOW}${BOLD}source $NODE_PROFILE${NC}"
        info "O reinicia tu terminal para aplicar los cambios."
    fi
}

bold "=== Instalador de Node.js $NODE_VERSION ==="

# Verificar instalación existente
check_existing_node_installation

install_node

# Recargar el entorno para reconocer Node.js
reload_shell_environment

echo
info "Para usar Node.js en nuevas sesiones de terminal:"
echo -e "  ${YELLOW}${BOLD}1.${NC} Las variables ya están configuradas globalmente"
echo -e "  ${YELLOW}${BOLD}2.${NC} Reinicia tu terminal, o ejecuta: ${YELLOW}${BOLD}source $NODE_PROFILE${NC}"
echo -e "  ${YELLOW}${BOLD}3.${NC} Verifica con: ${YELLOW}${BOLD}node --version${NC}"

success "Todo listo: Node.js $NODE_VERSION instalado y configurado."
