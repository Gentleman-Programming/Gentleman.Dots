#!/bin/bash

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

# Directorio de instaladores
INSTALLER_DIR="$(pwd)/installer"

# Verificar que existe el directorio installer
check_installer_dir() {
    if [ ! -d "$INSTALLER_DIR" ]; then
        echo "${RED}❌ Error: El directorio 'installer' no existe${NC}"
        echo "${BLUE}📁 Debe ejecutar este script desde el directorio raíz del proyecto${NC}"
        exit 1
    fi
}

# Verificar que existe el script de instalación
check_installer_script() {
    local script_name="$1"
    if [ ! -f "$INSTALLER_DIR/$script_name" ]; then
        echo "${RED}❌ Error: El script '$script_name' no existe en el directorio installer${NC}"
        echo "${YELLOW}📋 Scripts disponibles:${NC}"
        ls -1 "$INSTALLER_DIR"/*.sh 2>/dev/null || echo "   No hay scripts disponibles"
        exit 1
    fi
}

# Función para probar con Ubuntu
test_ubuntu() {
    local script_name="${1:-install-nvim.sh}"
    
    check_installer_dir
    check_installer_script "$script_name"
    
    echo "${BLUE}🐳 Iniciando contenedor Ubuntu${NC}..."
    docker run -it --rm \
        -e TERM=xterm-256color \
        -v "$(pwd):/root/.marckv.dots" \
        ubuntu:latest bash -c '
                export TERM=xterm-256color
                tput_color() { tput setaf "$1" 2>/dev/null || echo ""; }
                tput_bold() { tput bold 2>/dev/null || echo ""; }
                tput_reset() { tput sgr0 2>/dev/null || echo ""; }
                PINK=$(tput_color 204)
                PURPLE=$(tput_color 141)
                GREEN=$(tput_color 114)
                ORANGE=$(tput_color 208)
                BLUE=$(tput_color 75)
                YELLOW=$(tput_color 221)
                RED=$(tput_color 196)
                BOLD=$(tput_bold)
                NC=$(tput_reset)
                apt update -qq && apt install -y curl > /dev/null 2>&1
                chmod +x /root/.marckv.dots/installer/*.sh 2>/dev/null
                clear
                echo "${GREEN}=== CONTENEDOR UBUNTU LISTO ===${NC}"
                echo "${BLUE}Scripts disponibles en /root/.marckv.dots/installer:${NC}"
                for f in /root/.marckv.dots/installer/*.sh; do
                    [ -e "$f" ] && echo "  ${YELLOW}$(basename "$f")${NC}";
                done
                echo "${YELLOW}Para ejecutar: ./<script>.sh${NC} (ejemplo: /root/.marckv.dots/installer/install-nvim.sh)"
                echo "${YELLOW}Para salir del contenedor: exit${NC}"
                echo ""
                bash'
}

# Función para probar con Debian
test_debian() {
    local script_name="${1:-install-nvim.sh}"
    
    check_installer_dir
    check_installer_script "$script_name"
    
    echo "${BLUE}🐳 Iniciando contenedor Debian${NC}..."
    docker run -it --rm \
        -e TERM=xterm-256color \
        -v "$(pwd):/root/.marckv.dots" \
        debian:latest bash -c '
                export TERM=xterm-256color
                tput_color() { tput setaf "$1" 2>/dev/null || echo ""; }
                tput_bold() { tput bold 2>/dev/null || echo ""; }
                tput_reset() { tput sgr0 2>/dev/null || echo ""; }
                PINK=$(tput_color 204)
                PURPLE=$(tput_color 141)
                GREEN=$(tput_color 114)
                ORANGE=$(tput_color 208)
                BLUE=$(tput_color 75)
                YELLOW=$(tput_color 221)
                RED=$(tput_color 196)
                BOLD=$(tput_bold)
                NC=$(tput_reset)
                apt update -qq && apt install -y curl > /dev/null 2>&1
                chmod +x /root/.marckv.dots/installer/*.sh 2>/dev/null
                clear
                echo "${GREEN}=== CONTENEDOR DEBIAN LISTO ===${NC}"
                echo "${BLUE}Scripts disponibles en /root/.marckv.dots/installer:${NC}"
                for f in /root/.marckv.dots/installer/*.sh; do
                    [ -e "$f" ] && echo "  ${YELLOW}$(basename "$f")${NC}";
                done
                echo "${YELLOW}Para ejecutar: ./<script>.sh${NC} (ejemplo: /root/.marckv.dots/installer/install-nvim.sh)"
                echo "${YELLOW}Para salir del contenedor: exit${NC}"
                echo ""
                bash'
}

# Función para listar scripts disponibles
list_scripts() {
    check_installer_dir
    
    echo "${PURPLE}📋 Scripts disponibles en el directorio installer:${NC}"
    echo ""
    
    if ls "$INSTALLER_DIR"/*.sh >/dev/null 2>&1; then
        for script in "$INSTALLER_DIR"/*.sh; do
            script_name=$(basename "$script")
            echo "  ${GREEN}•${NC} ${YELLOW}$script_name${NC}"
        done
    else
        echo "  ${RED}❌ No hay scripts disponibles${NC}"
    fi
    
    echo ""
}

# Menú principal
case "$1" in
    ubuntu)
        test_ubuntu "$2"
        ;;
    debian)
        test_debian "$2"
        ;;
    list|ls)
        list_scripts
        ;;
    *)
        echo "${BOLD}🚀 Script para probar instaladores en contenedores Docker${NC}"
        echo ""
        echo "${BLUE}📁 Los scripts se cargan desde el directorio: ${YELLOW}./installer/${NC}"
        echo ""
        echo "${PURPLE}Uso:${NC}"
        echo "  ${GREEN}$0 ubuntu [script]${NC}  - Probar en Ubuntu (default: install-nvim.sh)"
        echo "  ${GREEN}$0 debian [script]${NC}  - Probar en Debian (default: install-nvim.sh)"
        echo "  ${GREEN}$0 list${NC}            - Listar scripts disponibles"
        echo ""
        echo "${ORANGE}Ejemplos:${NC}"
        echo "  ${YELLOW}$0 ubuntu${NC}                    # Probar install-nvim.sh en Ubuntu"
        echo "  ${YELLOW}$0 debian install-nvim.sh${NC}   # Probar script específico en Debian"
        echo ""
        echo "${PINK}📋 El contenedor se eliminará automáticamente al salir.${NC}"
        echo ""
        list_scripts
        ;;
esac
