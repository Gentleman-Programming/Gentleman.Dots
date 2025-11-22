#!/bin/bash

# Script para instalar Obsidian para todos los usuarios
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

# Variables de configuración
OBSIDIAN_PATH="/opt/obsidian"
OBSIDIAN_APPIMAGE_PATH="$OBSIDIAN_PATH/Obsidian.AppImage"
OBSIDIAN_DESKTOP_PATH="/usr/share/applications/obsidian.desktop"
OBSIDIAN_ICON_PATH="/usr/share/pixmaps/obsidian.png"

# Métodos de mensajes estándar
info()    { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[OK]${NC} $1"; }
error()   { echo -e "${RED}[ERROR]${NC} $1"; }
warn()    { echo -e "${ORANGE}[WARN]${NC} $1"; }
bold()    { echo -e "${BOLD}$1${NC}"; }

die() {
    error "$1"
    exit 1
}

# Verificar si se está ejecutando como root
check_root() {
    if [ "$EUID" -ne 0 ]; then
        error "Este script debe ejecutarse como root o con sudo"
        info "Uso: ${BOLD}sudo $0${NC}"
        exit 1
    fi
}

# Función para obtener la versión más reciente disponible
get_latest_version() {
    local latest_version
    latest_version=$(curl -s https://api.github.com/repos/obsidianmd/obsidian-releases/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$latest_version" ]; then
        return 1
    fi
    
    echo "$latest_version"
    return 0
}

# Función para obtener la versión instalada
get_installed_version() {
    if [ -f "$OBSIDIAN_APPIMAGE_PATH" ]; then
        # Intentar obtener versión del AppImage (esto puede ser complicado)
        # Por simplicidad, asumimos que está instalado si el archivo existe
        echo "installed"
        return 0
    fi
    return 1
}

# Verificar si Obsidian ya está instalado
check_existing_installation() {
    if [ -f "$OBSIDIAN_APPIMAGE_PATH" ] && [ -f "$OBSIDIAN_DESKTOP_PATH" ]; then
        warn "Obsidian ya está instalado en el sistema"
        info "Ruta de instalación: ${BOLD}$OBSIDIAN_APPIMAGE_PATH${NC}"
        
        # Verificar versión disponible
        info "Verificando actualizaciones disponibles..."
        local latest_version=$(get_latest_version)
        
        if [ $? -eq 0 ] && [ -n "$latest_version" ]; then
            info "Última versión disponible: ${BOLD}$latest_version${NC}"
            
            echo ""
            read -p "¿Desea actualizar a la última versión? (Y/n): " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Nn]$ ]]; then
                info "Actualización cancelada por el usuario"
                exit 0
            fi
            success "Procediendo con la actualización..."
        else
            warn "No se pudo verificar la versión más reciente"
            read -p "¿Desea reinstalar Obsidian de todas formas? (y/N): " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                info "Instalación cancelada por el usuario"
                exit 0
            fi
            warn "Procediendo con la reinstalación..."
        fi
    else
        info "Obsidian no está instalado en el sistema"
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
    if [ -d "$OBSIDIAN_PATH" ]; then
        info "Removiendo instalación anterior de Obsidian..."
        rm -rf "$OBSIDIAN_PATH"
        [ $? -eq 0 ] && success "Instalación anterior removida correctamente" || die "Error al remover instalación anterior"
    fi
    
    # Limpiar archivos del sistema
    [ -f "$OBSIDIAN_DESKTOP_PATH" ] && rm -f "$OBSIDIAN_DESKTOP_PATH"
    [ -f "$OBSIDIAN_ICON_PATH" ] && rm -f "$OBSIDIAN_ICON_PATH"
}

# Descargar Obsidian AppImage
download_obsidian() {
    info "Descargando Obsidian AppImage..."
    
    # Obtener URL de descarga del AppImage
    local obsidian_url=$(curl -s https://api.github.com/repos/obsidianmd/obsidian-releases/releases/latest | grep "browser_download_url.*AppImage" | cut -d '"' -f 4)
    
    if [ -z "$obsidian_url" ]; then
        return 1
    fi
    
    bold "URL: $obsidian_url"
    curl -fsSL -o "$OBSIDIAN_APPIMAGE_PATH" "$obsidian_url" || return 1
    success "Obsidian AppImage descargado correctamente"
    return 0
}

# Descargar Obsidian .deb como fallback
download_obsidian_deb() {
    info "Intentando instalación via paquete .deb..."
    
    # Obtener URL de descarga del .deb
    local obsidian_deb_url=$(curl -s https://api.github.com/repos/obsidianmd/obsidian-releases/releases/latest | grep "browser_download_url.*amd64.deb" | cut -d '"' -f 4)
    
    if [ -z "$obsidian_deb_url" ]; then
        return 1
    fi
    
    bold "URL: $obsidian_deb_url"
    curl -fsSL -o /tmp/obsidian.deb "$obsidian_deb_url" || return 1
    apt-get install -y /tmp/obsidian.deb || return 1
    rm -f /tmp/obsidian.deb
    success "Obsidian instalado via paquete .deb"
    return 0
}

# Instalar Obsidian
install_obsidian() {
    # Crear directorio para Obsidian
    mkdir -p "$OBSIDIAN_PATH"
    
    # Intentar descargar AppImage primero
    if download_obsidian; then
        # Dar permisos de ejecución al AppImage
        chmod +x "$OBSIDIAN_APPIMAGE_PATH"
        [ $? -eq 0 ] && success "Permisos de ejecución configurados" || die "Error al configurar permisos"
        
        # Crear archivo .desktop
        setup_desktop_entry
        
        # Descargar ícono
        download_icon
        
        success "Obsidian instalado correctamente como AppImage"
    elif download_obsidian_deb; then
        success "Obsidian instalado correctamente via paquete .deb"
    else
        die "No se pudo instalar Obsidian automáticamente. Puedes descargarlo manualmente desde https://obsidian.md/"
    fi
}

# Configurar entrada en el menú de aplicaciones
setup_desktop_entry() {
    info "Configurando entrada en el menú de aplicaciones..."
    
    cat > "$OBSIDIAN_DESKTOP_PATH" << EOF
[Desktop Entry]
Name=Obsidian
Comment=A powerful knowledge base that works on top of a local folder of plain text Markdown files
Exec=$OBSIDIAN_APPIMAGE_PATH %U
Terminal=false
Type=Application
Icon=obsidian
Categories=Office;TextEditor;Utility;
MimeType=x-scheme-handler/obsidian;
StartupWMClass=obsidian
EOF
    
    [ $? -eq 0 ] && success "Archivo .desktop creado correctamente" || warn "Error al crear archivo .desktop"
}

# Descargar ícono de Obsidian
download_icon() {
    info "Descargando ícono de Obsidian..."
    
    # Intentar descargar el ícono oficial (SVG convertido a PNG puede no funcionar directamente)
    # Como fallback, usar un ícono genérico o crear uno simple
    curl -fsSL -o "$OBSIDIAN_ICON_PATH" "https://obsidian.md/images/obsidian-logo-gradient.svg" 2>/dev/null || {
        warn "No se pudo descargar el ícono oficial, usando ícono genérico"
        # Crear un ícono simple de texto como fallback
        echo "Obsidian" > "$OBSIDIAN_ICON_PATH.txt"
    }
}

# Verificar instalación
verify_installation() {
    info "Verificando instalación..."
    
    if [ -f "$OBSIDIAN_APPIMAGE_PATH" ] && [ -x "$OBSIDIAN_APPIMAGE_PATH" ]; then
        success "Obsidian AppImage instalado correctamente"
        bold "Ubicación: $OBSIDIAN_APPIMAGE_PATH"
    elif command -v obsidian >/dev/null 2>&1; then
        success "Obsidian instalado correctamente via paquete .deb"
        local obsidian_path=$(which obsidian)
        bold "Ubicación: $obsidian_path"
    else
        die "La verificación de instalación falló"
    fi
    
    if [ -f "$OBSIDIAN_DESKTOP_PATH" ]; then
        success "Entrada del menú configurada correctamente"
    else
        warn "La entrada del menú puede no estar disponible"
    fi
}

# Crear directorio vault para compatibilidad con obsidian.nvim
setup_obsidian_vault() {
    info "Configurando vault de Obsidian para integración con Neovim..."
    
    # Obtener usuarios del sistema (excluyendo usuarios del sistema)
    local users=$(getent passwd | grep -E ":[0-9]{4}:" | cut -d: -f1)
    
    for user in $users; do
        local user_home=$(getent passwd "$user" | cut -d: -f6)
        local vault_dir="$user_home/.config/obsidian"
        
        if [ -d "$user_home" ] && [ "$user_home" != "/root" ]; then
            info "Configurando vault para usuario: $user"
            
            # Crear directorio vault si no existe
            if [ ! -d "$vault_dir" ]; then
                sudo -u "$user" mkdir -p "$vault_dir"
                [ $? -eq 0 ] && success "Directorio vault creado: $vault_dir" || warn "No se pudo crear el directorio vault para $user"
            else
                success "Directorio vault ya existe: $vault_dir"
            fi
            
            # Crear archivo de bienvenida si no existe
            local welcome_file="$vault_dir/Bienvenida.md"
            if [ ! -f "$welcome_file" ]; then
                sudo -u "$user" cat > "$welcome_file" << 'EOF'
# 🎉 Bienvenido a tu Obsidian Vault

Este es tu vault de Obsidian, configurado para trabajar tanto con:
- **Obsidian (aplicación)**: Para navegación visual y grafo de conexiones
- **obsidian.nvim**: Para edición rápida desde Neovim

## 🚀 Primeros Pasos

1. Crea nuevas notas usando `[[Nombre de la Nota]]`
2. Conecta ideas entre notas
3. Explora el grafo de conexiones
4. Usa templates en la carpeta `templates/`

## 📁 Estructura recomendada

```
~/.config/obsidian/
├── Inbox/          # Notas rápidas
├── Projects/       # Proyectos específicos
├── Resources/      # Referencias y recursos
├── Archive/        # Notas archivadas
└── templates/      # Plantillas
```

## 🔗 Enlaces útiles

- [[Índice de Proyectos]]
- [[Ideas Rápidas]]
- [[Recursos de Desarrollo]]

¡Feliz escritura! ✨
EOF
                [ $? -eq 0 ] && success "Archivo de bienvenida creado para $user" || warn "No se pudo crear el archivo de bienvenida para $user"
            fi
            
            # Crear directorio de templates
            local templates_dir="$vault_dir/templates"
            if [ ! -d "$templates_dir" ]; then
                sudo -u "$user" mkdir -p "$templates_dir"
                [ $? -eq 0 ] && success "Directorio de templates creado para $user"
            fi
        fi
    done
}

# Mostrar información post-instalación
show_post_install_info() {
    echo
    info "Obsidian está listo para usar:"
    echo -e "  ${YELLOW}${BOLD}1.${NC} Busca 'Obsidian' en el menú de aplicaciones"
    echo -e "  ${YELLOW}${BOLD}2.${NC} O ejecuta directamente: ${YELLOW}${BOLD}$OBSIDIAN_APPIMAGE_PATH${NC}"
    echo -e "  ${YELLOW}${BOLD}3.${NC} Abre el vault en: ${YELLOW}${BOLD}~/.config/obsidian${NC}"
    
    echo
    info "Integración con Neovim:"
    echo -e "  ${BOLD}•${NC} El vault está configurado para obsidian.nvim"
    echo -e "  ${BOLD}•${NC} Edita notas desde Neovim usando comandos :ObsidianOpen"
    echo -e "  ${BOLD}•${NC} Ambas herramientas comparten la misma carpeta de notas"
    
    echo
    info "Consejos para empezar:"
    echo -e "  ${BOLD}•${NC} Lee el archivo Bienvenida.md en tu vault"
    echo -e "  ${BOLD}•${NC} Usa [[Nombre de Nota]] para crear enlaces entre notas"
    echo -e "  ${BOLD}•${NC} Explora el grafo de conexiones con Ctrl+G"
    echo -e "  ${BOLD}•${NC} Usa templates en la carpeta templates/"
    
    bold "\n=== INSTALACIÓN COMPLETADA ==="
    success "Obsidian ha sido instalado correctamente."
}

# Función principal
main() {
    bold "=== INSTALADOR DE OBSIDIAN ==="
    info "Este script instalará Obsidian para todos los usuarios del sistema"
    
    # Verificaciones iniciales
    check_root
    check_existing_installation
    
    # Proceso de instalación
    cleanup_previous
    install_obsidian
    verify_installation
    
    # Configurar vault para integración con Neovim
    setup_obsidian_vault
    
    show_post_install_info
    
    success "\n¡Instalación completada exitosamente!"
}

# Ejecutar función principal
main "$@"