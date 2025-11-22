#!/bin/bash

# Script para instalar Go siguiendo las instrucciones oficiales de go.dev
# Basado en: https://go.dev/doc/install

# === VARIABLES CONFIGURABLES ===
GO_VERSION=""  # Se detectará automáticamente
GO_ARCH=""     # Se detectará automáticamente
GO_INSTALL_DIR="/usr/local"
GO_TAR=""
GO_DOWNLOAD_URL=""

# Colores para output
GREEN=$(tput setaf 114)
ORANGE=$(tput setaf 208)
BLUE=$(tput setaf 75)
RED=$(tput setaf 196)
BOLD=$(tput bold)
NC=$(tput sgr0)

# Funciones de mensaje
info()    { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[OK]${NC} $1"; }
error()   { echo -e "${RED}[ERROR]${NC} $1"; }
warn()    { echo -e "${ORANGE}[WARN]${NC} $1"; }

die() {
    error "$1"
    exit 1
}

# Verificar permisos de root
check_root() {
    if [ "$EUID" -ne 0 ]; then
        error "Este script debe ejecutarse como root o con sudo"
        info "Uso: sudo $0"
        exit 1
    fi
}

# Detectar arquitectura del sistema
detect_architecture() {
    local os_type
    local arch_type
    
    # Este script está diseñado para Linux/WSL
    if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]]; then
        error "Este script está diseñado para Linux/WSL, no para Git Bash en Windows"
        info "Por favor ejecuta este script en WSL o un sistema Linux"
        exit 1
    fi
    
    os_type="linux"
    
    # Detectar arquitectura
    arch_type=$(uname -m)
    case $arch_type in
        x86_64|amd64)
            arch_type="amd64"
            ;;
        aarch64|arm64)
            arch_type="arm64"
            ;;
        armv6l)
            arch_type="armv6l"
            ;;
        i386|i686)
            arch_type="386"
            ;;
        *)
            arch_type="amd64"  # Fallback
            ;;
    esac
    
    echo "${os_type}-${arch_type}"
}

# Obtener la última versión de Go
get_latest_go_version() {
    local latest_version
    
    # Usar scraping de la página oficial de descargas
    latest_version=$(curl -s https://go.dev/dl/ | grep -oP 'go\K[0-9]+\.[0-9]+\.[0-9]+' | head -n1)
    
    # Fallback: versión conocida estable
    if [ -z "$latest_version" ] || [[ "$latest_version" =~ [^0-9.] ]]; then
        latest_version="1.23.3"
    fi
    
    echo "$latest_version"
}

# Comparar versiones semánticamente
compare_versions() {
    local installed="$1"
    local available="$2"
    
    if [ -z "$installed" ] || [ -z "$available" ]; then
        return 1
    fi
    
    # Usar sort -V para comparación semántica
    local higher=$(printf '%s\n%s\n' "$installed" "$available" | sort -V | tail -n1)
    
    if [ "$higher" = "$available" ] && [ "$installed" != "$available" ]; then
        return 0  # Hay una versión más nueva
    else
        return 1  # Ya está actualizada
    fi
}

# Verificar si Go ya está instalado
check_existing_installation() {
    # Detectar arquitectura del sistema
    GO_ARCH=$(detect_architecture)
    info "Arquitectura detectada: $GO_ARCH"
    
    # Verificar en múltiples ubicaciones porque sudo puede cambiar el PATH
    local go_binary=""
    local installed_version=""
    
    # 1. Verificar en PATH actual
    if command -v go >/dev/null 2>&1; then
        go_binary=$(command -v go)
        installed_version=$(go version 2>/dev/null | awk '{print $3}' | sed 's/go//')
    # 2. Verificar en ubicación estándar de instalación
    elif [ -x "/usr/local/go/bin/go" ]; then
        go_binary="/usr/local/go/bin/go"
        installed_version=$(/usr/local/go/bin/go version 2>/dev/null | awk '{print $3}' | sed 's/go//')
    # 3. Verificar en /usr/bin
    elif [ -x "/usr/bin/go" ]; then
        go_binary="/usr/bin/go"
        installed_version=$(/usr/bin/go version 2>/dev/null | awk '{print $3}' | sed 's/go//')
    fi
    
    if [ -n "$go_binary" ] && [ -n "$installed_version" ]; then
        warn "Go ya está instalado en el sistema"
        info "Versión instalada: go$installed_version"
        info "Ubicación: $go_binary"
        
        # Verificar si hay una versión más nueva
        info "Verificando actualizaciones disponibles..."
        local latest_version
        latest_version=$(get_latest_go_version)
        
        if [ -n "$latest_version" ]; then
            info "Última versión disponible: go$latest_version"
            
            if compare_versions "$installed_version" "$latest_version"; then
                success "🚀 ¡Nueva versión disponible!"
                info "Versión instalada: go$installed_version"
                info "Versión disponible: go$latest_version"
                
                read -p "¿Desea actualizar a la última versión? (Y/n): " -n 1 -r
                echo
                if [[ $REPLY =~ ^[Nn]$ ]]; then
                    info "Actualización cancelada por el usuario"
                    exit 0
                fi
                
                # Configurar variables para la nueva versión
                GO_VERSION="$latest_version"
            else
                success "✅ Ya tienes la versión más reciente instalada"
                info "No es necesario actualizar"
                exit 0
            fi
        else
            warn "No se pudo verificar la versión más reciente"
            read -p "¿Desea reinstalar Go de todas formas? (y/N): " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                info "Instalación cancelada por el usuario"
                exit 0
            fi
            GO_VERSION="1.23.3"  # Versión fallback
        fi
    else
        # Go no está instalado, obtener la última versión
        info "Go no está instalado en el sistema"
        GO_VERSION=$(get_latest_go_version)
        info "Se instalará la última versión: go$GO_VERSION"
    fi
    
    # Configurar URLs después de determinar la versión
    GO_TAR="/tmp/go${GO_VERSION}.${GO_ARCH}.tar.gz"
    GO_DOWNLOAD_URL="https://go.dev/dl/go${GO_VERSION}.${GO_ARCH}.tar.gz"
    
    # Mostrar información de lo que se va a instalar
    echo ""
    info "📦 Preparando instalación de Go $GO_VERSION"
}

# Descargar Go
download_go() {
    info "Descargando Go $GO_VERSION..."
    
    if curl -L --progress-bar -o "$GO_TAR" "$GO_DOWNLOAD_URL"; then
        success "Go descargado correctamente"
    else
        die "Error al descargar Go desde: $GO_DOWNLOAD_URL"
    fi
}

# Instalar Go
install_go() {
    info "Removiendo instalación anterior de Go (si existe)..."
    rm -rf "$GO_INSTALL_DIR/go"
    
    info "Instalando Go en $GO_INSTALL_DIR..."
    if tar -C "$GO_INSTALL_DIR" -xzf "$GO_TAR"; then
        success "Go instalado correctamente"
    else
        die "Error al extraer Go"
    fi
    
    # Limpiar archivo temporal
    rm -f "$GO_TAR"
}

# Configurar PATH
setup_environment() {
    info "Configurando variables de entorno..."
    
    # Crear archivo de configuración en /etc/profile.d/
    cat > /etc/profile.d/go.sh << 'EOF'
# Go programming language configuration
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$GOROOT/bin:$GOPATH/bin:$PATH
EOF

    chmod +x /etc/profile.d/go.sh
    success "Variables de entorno configuradas en /etc/profile.d/go.sh"
    
    # Aplicar configuración en la sesión actual
    export GOROOT=/usr/local/go
    export GOPATH=$HOME/go
    export PATH=$GOROOT/bin:$GOPATH/bin:$PATH
}

# Verificar instalación
verify_installation() {
    info "Verificando instalación..."
    
    if [ ! -x "$GO_INSTALL_DIR/go/bin/go" ]; then
        die "Binario de Go no encontrado en $GO_INSTALL_DIR/go/bin/go"
    fi
    
    local go_version
    go_version=$("$GO_INSTALL_DIR/go/bin/go" version 2>/dev/null | awk '{print $3}')
    
    if [ -n "$go_version" ]; then
        success "Go instalado correctamente: $go_version"
        info "Ubicación: $GO_INSTALL_DIR/go"
        info "Variables configuradas en: /etc/profile.d/go.sh"
    else
        die "Error: Go no responde correctamente"
    fi
}

# Mostrar información final
show_final_info() {
    echo ""
    success "🎉 ¡Instalación de Go completada!"
    info "Para usar Go en la sesión actual, ejecuta:"
    echo "   source /etc/profile.d/go.sh"
    info "O reinicia tu terminal/sesión"
    echo ""
    info "Comandos útiles:"
    echo "   go version      # Ver versión instalada"
    echo "   go mod init     # Inicializar un módulo"
    echo "   go run main.go  # Ejecutar un programa"
    echo ""
}

# === FUNCIÓN PRINCIPAL ===
main() {
    info "🚀 Iniciando instalador de Go"
    
    # Validar entorno antes que nada
    if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]]; then
        error "Este script debe ejecutarse en Linux o WSL, no en Git Bash"
        info "Opciones:"
        info "  1. Usar WSL: wsl -d Ubuntu sudo bash install-go.sh"
        info "  2. Ejecutar en un sistema Linux real"
        exit 1
    fi
    
    check_root
    check_existing_installation
    download_go
    install_go
    setup_environment
    verify_installation
    show_final_info
}

# Ejecutar función principal
main "$@"
