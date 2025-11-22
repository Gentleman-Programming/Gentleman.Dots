#!/bin/bash

# Script para instalar Neovim para todos los usuarios
# Requiere permisos de root/sudo

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

# URL y rutas
NVIM_URL="https://github.com/neovim/neovim/releases/latest/download/nvim-linux-x86_64.tar.gz"
NVIM_TAR="/tmp/nvim-linux-x86_64.tar.gz"
NVIM_PATH="/opt/nvim"
PROFILE_PATH="/etc/profile.d/nvim.sh"

# Métodos de mensajes estándar
info()    { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[OK]${NC} $1"; }
error()   { echo -e "${RED}[ERROR]${NC} $1"; }
warn()    { echo -e "${ORANGE}[WARN]${NC} $1"; }
bold()    { echo -e "${BOLD}$1${NC}"; }

# Función para obtener la versión más reciente disponible
get_latest_version() {
    local latest_version
    latest_version=$(curl -s https://api.github.com/repos/neovim/neovim/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$latest_version" ]; then
        return 1
    fi
    
    echo "$latest_version"
    return 0
}

# Función para obtener la versión instalada
get_installed_version() {
    if [ -x "$NVIM_PATH/bin/nvim" ]; then
        local installed_version
        installed_version=$($NVIM_PATH/bin/nvim --version 2>/dev/null | head -n1 | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+')
        echo "$installed_version"
        return 0
    fi
    return 1
}

# Función para comparar versiones
compare_versions() {
    local installed="$1"
    local latest="$2"
    
    # Remover la 'v' del inicio para comparación
    local installed_clean=$(echo "$installed" | sed 's/^v//')
    local latest_clean=$(echo "$latest" | sed 's/^v//')
    
    # Usar sort -V para comparación de versiones
    local higher_version=$(printf '%s\n%s\n' "$installed_clean" "$latest_clean" | sort -V | tail -n1)
    
    if [ "$higher_version" = "$latest_clean" ] && [ "$installed_clean" != "$latest_clean" ]; then
        return 0  # Hay una versión más nueva disponible
    else
        return 1  # Ya está actualizado
    fi
}

# Función para verificar si el comando se ejecutó correctamente
check_status() {
    if [ $? -eq 0 ]; then
        success "$1"
    else
        error "$2"
        exit 1
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

# Verificar si Neovim ya está instalado
check_existing_installation() {
    if [ -d "$NVIM_PATH" ] && [ -f "$PROFILE_PATH" ]; then
        warn "Neovim ya está instalado en el sistema"
        info "Ruta de instalación: ${BOLD}$NVIM_PATH${NC}"
        
        # Verificar versión actual
        if [ -x "$NVIM_PATH/bin/nvim" ]; then
            local current_version=$($NVIM_PATH/bin/nvim --version | head -n1)
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
                else
                    success "✅ Ya tienes la versión más reciente instalada"
                    info "No es necesario actualizar"
                    exit 0
                fi
            else
                warn "No se pudo verificar la versión más reciente"
                read -p "¿Desea reinstalar Neovim de todas formas? (y/N): " -n 1 -r
                echo
                if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                    info "Instalación cancelada por el usuario"
                    exit 0
                fi
                warn "Procediendo con la reinstalación..."
            fi
        else
            read -p "¿Desea reinstalar Neovim? (y/N): " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                info "Instalación cancelada por el usuario"
                exit 0
            fi
            warn "Procediendo con la reinstalación..."
        fi
    else
        # No está instalado, verificar la última versión disponible
    info "Neovim no está instalado en el sistema"
    info "Verificando la última versión disponible..."
        local latest_version=$(get_latest_version)
        
        if [ $? -eq 0 ] && [ -n "$latest_version" ]; then
            info "Se instalará la versión más reciente: ${BOLD}${GREEN}$latest_version${NC}"
        else
        warn "No se pudo verificar la versión más reciente, pero se procederá con la instalación"
        fi
    fi
}

# Limpiar instalación anterior si existe
cleanup_previous() {
    if [ -d "$NVIM_PATH" ]; then
    info "Removiendo instalación anterior de Neovim..."
        rm -rf "$NVIM_PATH"
        check_status "Instalación anterior removida correctamente" "Error al remover instalación anterior"
    fi
    
    # También limpiar la ruta anterior por si existía
    if [ -d "/opt/nvim-linux-x86_64" ]; then
    info "Removiendo instalación anterior en ruta legacy..."
        rm -rf "/opt/nvim-linux-x86_64"
        check_status "Instalación legacy removida correctamente" "Error al remover instalación legacy"
    fi
}

# Descargar Neovim
download_neovim() {
    if [ -f "$NVIM_TAR" ] && [ -s "$NVIM_TAR" ]; then
        info "Usando archivo de Neovim ya descargado en $NVIM_TAR."
    else
        info "Descargando Neovim desde GitHub..."
        bold "URL: $NVIM_URL"
        curl -L -o "$NVIM_TAR" "$NVIM_URL"
        check_status "Neovim descargado correctamente" "Error al descargar Neovim"
    fi
}

# Extraer e instalar Neovim
install_neovim() {
    info "Extrayendo Neovim a /opt..."
    
    # Crear directorio temporal para extracción
    local temp_dir="/tmp/nvim-extract"
    mkdir -p "$temp_dir"
    
    # Extraer a directorio temporal primero
    tar -C "$temp_dir" -xzf "$NVIM_TAR"
    check_status "Neovim extraído a directorio temporal" "Error al extraer Neovim"
    
    # Mover el contenido a /opt/nvim
    if [ -d "$temp_dir/nvim-linux-x86_64" ]; then
        mv "$temp_dir/nvim-linux-x86_64" "$NVIM_PATH"
    check_status "Neovim movido a $NVIM_PATH" "Error al mover Neovim a la ubicación final"
    else
    error "La estructura del archivo tar no es la esperada"
        rm -rf "$temp_dir"
        exit 1
    fi
    
    # Limpiar directorio temporal
    rm -rf "$temp_dir"
    
    # Verificar que la instalación fue exitosa
    if [ ! -d "$NVIM_PATH" ]; then
    error "El directorio de instalación no fue creado"
        exit 1
    fi
    
    if [ ! -x "$NVIM_PATH/bin/nvim" ]; then
    error "El ejecutable de Neovim no fue encontrado"
        exit 1
    fi
}

# Configurar PATH para todos los usuarios
setup_path() {
    info "Configurando PATH para todos los usuarios..."
    
    echo "export PATH=\"\$PATH:$NVIM_PATH/bin\"" > "$PROFILE_PATH"
    check_status "Archivo de perfil creado" "Error al crear archivo de perfil"
    
    chmod 644 "$PROFILE_PATH"
    check_status "Permisos configurados correctamente" "Error al configurar permisos"
}

# Verificar instalación
verify_installation() {
    info "Verificando instalación..."
    
    if [ -x "$NVIM_PATH/bin/nvim" ]; then
        local version=$($NVIM_PATH/bin/nvim --version | head -n1)
    success "Neovim instalado correctamente"
    bold "Versión instalada: $version"
    bold "Ubicación: $NVIM_PATH/bin/nvim"
    else
    error "La verificación de instalación falló"
        exit 1
    fi
}

# Recargar el entorno del shell
reload_shell_environment() {
    info "Recargando variables de entorno..."
    
    # Recargar el perfil de Neovim
    if [ -f "$PROFILE_PATH" ]; then
        source "$PROFILE_PATH"
        success "Variables de entorno de Neovim cargadas."
    fi
    
    # Verificar que Neovim esté disponible en el PATH actual
    if command -v nvim >/dev/null 2>&1; then
        local nvim_version=$(nvim --version | head -n1)
        success "Neovim está disponible: $nvim_version"
        
        # Verificar ruta del ejecutable
        local nvim_path=$(which nvim)
        info "Ejecutable encontrado en: $nvim_path"
    else
        warn "Neovim no está disponible en el PATH actual."
        info "Puedes ejecutar: ${YELLOW}${BOLD}source $PROFILE_PATH${NC}"
        info "O reinicia tu terminal para aplicar los cambios."
    fi
}

# Mostrar información post-instalación
show_post_install_info() {
    local installed_version=$(get_installed_version)
    
    echo
    info "Para usar Neovim en nuevas sesiones de terminal:"
    echo -e "  ${YELLOW}${BOLD}1.${NC} Las variables ya están configuradas globalmente"
    echo -e "  ${YELLOW}${BOLD}2.${NC} Reinicia tu terminal, o ejecuta: ${YELLOW}${BOLD}source $PROFILE_PATH${NC}"
    echo -e "  ${YELLOW}${BOLD}3.${NC} Verifica con: ${YELLOW}${BOLD}nvim --version${NC}"
    
    bold "\n=== INSTALACIÓN COMPLETADA ==="
    success "Neovim ha sido instalado correctamente para todos los usuarios."
}

# Función principal
main() {
    bold "=== INSTALADOR DE NEOVIM ==="
    info "Este script instalará Neovim para todos los usuarios del sistema"
    
    # Verificaciones iniciales
    check_root
    check_existing_installation
    
    # Proceso de instalación
    cleanup_previous
    download_neovim
    install_neovim
    setup_path
    verify_installation
    
    # Recargar el entorno para reconocer Neovim
    reload_shell_environment
    
    show_post_install_info
    
    success "\n¡Instalación completada exitosamente!"
}

# Ejecutar función principal
main "$@"
