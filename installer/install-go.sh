#!/bin/bash

# === VARIABLES PARAMETRIZABLES ===
GO_VERSION="1.21.5"
GO_ARCH="linux-amd64"
GO_BASE_URL="https://golang.org/dl/go$GO_VERSION.$GO_ARCH.tar.gz"
GO_TAR="/tmp/go$GO_VERSION.$GO_ARCH.tar.gz"
GO_PATH="/opt/go"
GO_PROFILE="/etc/profile.d/go.sh"

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

# Función para obtener la versión más reciente de Go
get_latest_version() {
    local latest_version
    latest_version=$(curl -s https://golang.org/VERSION?m=text | head -n1 | sed 's/go//')
    
    if [ -z "$latest_version" ]; then
        return 1
    fi
    
    echo "$latest_version"
    return 0
}

# Función para obtener la versión instalada
get_installed_version() {
    if [ -x "$GO_PATH/bin/go" ]; then
        local installed_version
        installed_version=$($GO_PATH/bin/go version 2>/dev/null | grep -oE 'go[0-9]+\.[0-9]+(\.[0-9]+)?' | sed 's/go//')
        echo "$installed_version"
        return 0
    fi
    return 1
}

# Función para comparar versiones
compare_versions() {
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

# Verificar si se está ejecutando como root
check_root() {
    if [ "$EUID" -ne 0 ]; then
        error "Este script debe ejecutarse como root o con sudo"
        info "Uso: ${BOLD}sudo $0${NC}"
        exit 1
    fi
}

# Verificar si Go ya está instalado
check_existing_installation() {
    if [ -d "$GO_PATH" ] && [ -f "$GO_PROFILE" ]; then
        warn "Go ya está instalado en el sistema"
        info "Ruta de instalación: ${BOLD}$GO_PATH${NC}"
        
        # Verificar versión actual
        if [ -x "$GO_PATH/bin/go" ]; then
            local current_version=$($GO_PATH/bin/go version)
            local installed_version=$(get_installed_version)
            info "Versión actual: ${BOLD}$current_version${NC}"
            
            # Verificar si hay una nueva versión disponible
            info "Verificando actualizaciones disponibles..."
            local latest_version=$(get_latest_version)
            
            if [ $? -eq 0 ] && [ -n "$latest_version" ]; then
                info "Última versión disponible: ${BOLD}$latest_version${NC}"
                
                if compare_versions "$installed_version" "$latest_version"; then
                    bold "\n🚀 ¡NUEVA VERSIÓN DISPONIBLE!"
                    info "Versión instalada: ${YELLOW}$installed_version${NC}"
                    info "Versión disponible: ${GREEN}$latest_version${NC}"
                    warn "Se recomienda actualizar para obtener las últimas mejoras y correcciones"
                    
                    echo ""
                    read -p "¿Desea actualizar a la última versión? (Y/n): " -n 1 -r
                    echo
                    if [[ $REPLY =~ ^[Nn]$ ]]; then
                        info "Actualización cancelada por el usuario"
                        exit 0
                    fi
                    success "Procediendo con la actualización..."
                    # Actualizar la versión para descargar
                    GO_VERSION="$latest_version"
                    GO_BASE_URL="https://golang.org/dl/go$GO_VERSION.$GO_ARCH.tar.gz"
                    GO_TAR="/tmp/go$GO_VERSION.$GO_ARCH.tar.gz"
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
                warn "Procediendo con la reinstalación..."
            fi
        else
            read -p "¿Desea reinstalar Go? (y/N): " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                info "Instalación cancelada por el usuario"
                exit 0
            fi
            warn "Procediendo con la reinstalación..."
        fi
    else
        # No está instalado, verificar la última versión disponible
        info "Go no está instalado en el sistema"
        info "Verificando la última versión disponible..."
        local latest_version=$(get_latest_version)
        
        if [ $? -eq 0 ] && [ -n "$latest_version" ]; then
            info "Se instalará la versión más reciente: ${BOLD}${GREEN}$latest_version${NC}"
            # Actualizar variables para usar la última versión
            GO_VERSION="$latest_version"
            GO_BASE_URL="https://golang.org/dl/go$GO_VERSION.$GO_ARCH.tar.gz"
            GO_TAR="/tmp/go$GO_VERSION.$GO_ARCH.tar.gz"
        else
            warn "No se pudo verificar la versión más reciente, pero se procederá con la instalación de la versión predefinida: $GO_VERSION"
        fi
    fi
}

# Limpiar instalación anterior si existe
cleanup_previous() {
    if [ -d "$GO_PATH" ]; then
        info "Removiendo instalación anterior de Go..."
        rm -rf "$GO_PATH"
        [ $? -eq 0 ] && success "Instalación anterior removida correctamente" || die "Error al remover instalación anterior"
    fi
}

# Descargar Go
download_go() {
    if [ -f "$GO_TAR" ] && [ -s "$GO_TAR" ]; then
        info "Usando archivo de Go ya descargado en $GO_TAR."
    else
        info "Descargando Go $GO_VERSION desde el sitio oficial..."
        bold "URL: $GO_BASE_URL"
        curl -L -o "$GO_TAR" "$GO_BASE_URL" || die "No se pudo descargar Go."
        success "Go descargado correctamente"
    fi
}

# Extraer e instalar Go
install_go() {
    info "Extrayendo Go a /opt..."
    
    # Extraer directamente a /opt
    tar -C /opt -xzf "$GO_TAR" || die "No se pudo extraer Go."
    
    # Verificar que la instalación fue exitosa
    if [ ! -d "$GO_PATH" ]; then
        die "El directorio de instalación no fue creado"
    fi
    
    if [ ! -x "$GO_PATH/bin/go" ]; then
        die "El ejecutable de Go no fue encontrado"
    fi
    
    success "Go extraído e instalado correctamente en $GO_PATH"
}

# Configurar variables de entorno para todos los usuarios
setup_environment() {
    info "Configurando variables de entorno para todos los usuarios..."
    
    cat > "$GO_PROFILE" << 'EOF'
export GOROOT=/opt/go
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
EOF
    
    [ $? -eq 0 ] && success "Archivo de perfil creado" || die "Error al crear archivo de perfil"
    
    chmod 644 "$GO_PROFILE"
    [ $? -eq 0 ] && success "Permisos configurados correctamente" || die "Error al configurar permisos"
}

# Verificar instalación
verify_installation() {
    info "Verificando instalación..."
    
    if [ -x "$GO_PATH/bin/go" ]; then
        local version=$($GO_PATH/bin/go version)
        success "Go instalado correctamente"
        bold "Versión instalada: $version"
        bold "Ubicación: $GO_PATH/bin/go"
    else
        die "La verificación de instalación falló"
    fi
}

# Recargar el entorno del shell
reload_shell_environment() {
    info "Recargando variables de entorno..."
    
    # Recargar el perfil de Go
    if [ -f "$GO_PROFILE" ]; then
        source "$GO_PROFILE"
        success "Variables de entorno de Go cargadas."
    fi
    
    # Verificar que Go esté disponible en el PATH actual
    if command -v go >/dev/null 2>&1; then
        local go_version=$(go version)
        success "Go está disponible: $go_version"
        
        # Verificar variables de entorno importantes
        if [ -n "$GOROOT" ]; then
            info "GOROOT configurado en: $GOROOT"
        fi
        
        if [ -n "$GOPATH" ]; then
            info "GOPATH configurado en: $GOPATH"
        else
            info "GOPATH se configurará automáticamente en: \$HOME/go"
        fi
        
        # Verificar ruta del ejecutable
        local go_path=$(which go)
        info "Ejecutable encontrado en: $go_path"
    else
        warn "Go no está disponible en el PATH actual."
        info "Puedes ejecutar: ${YELLOW}${BOLD}source $GO_PROFILE${NC}"
        info "O reinicia tu terminal para aplicar los cambios."
    fi
}

# Mostrar información post-instalación
show_post_install_info() {
    echo
    info "Para usar Go en nuevas sesiones de terminal:"
    echo -e "  ${YELLOW}${BOLD}1.${NC} Las variables ya están configuradas globalmente"
    echo -e "  ${YELLOW}${BOLD}2.${NC} Reinicia tu terminal, o ejecuta: ${YELLOW}${BOLD}source $GO_PROFILE${NC}"
    echo -e "  ${YELLOW}${BOLD}3.${NC} Verifica con: ${YELLOW}${BOLD}go version${NC}"
    echo -e "  ${YELLOW}${BOLD}4.${NC} Tu workspace de Go estará en: ${YELLOW}${BOLD}\$HOME/go${NC}"
    
    echo
    info "Variables de entorno configuradas:"
    echo -e "  ${BOLD}GOROOT${NC}: $GO_PATH (instalación de Go)"
    echo -e "  ${BOLD}GOPATH${NC}: \$HOME/go (tu workspace)"
    echo -e "  ${BOLD}PATH${NC}: incluye \$GOROOT/bin y \$GOPATH/bin"
    
    bold "\n=== INSTALACIÓN COMPLETADA ==="
    success "Go ha sido instalado correctamente para todos los usuarios."
}

# Función principal
main() {
    bold "=== INSTALADOR DE GO $GO_VERSION ==="
    info "Este script instalará Go para todos los usuarios del sistema"
    
    # Verificaciones iniciales
    check_root
    check_existing_installation
    
    # Proceso de instalación
    cleanup_previous
    download_go
    install_go
    setup_environment
    verify_installation
    
    # Recargar el entorno para reconocer Go
    reload_shell_environment
    
    show_post_install_info
    
    success "\n¡Instalación completada exitosamente!"
}

# Ejecutar función principal
main "$@"
