{ lib, pkgs, ... }:
{
  home.packages = with pkgs; [
    # Scripts para Oil.nvim
    (writeShellScriptBin "oil" ''
      #!/usr/bin/env bash
      # oil: Launch Neovim with Oil.nvim using full LazyVim config
      # Usage: oil [directory]
      #
      # This script opens Neovim with your complete configuration and Oil
      # Press 'q' or ':qa' to exit back to terminal

      set -euo pipefail

      # Default to current directory if no argument provided
      DIR="''${1:-.}"

      # Resolve to absolute path
      DIR=$(realpath "$DIR")

      # Check if directory exists
      if [[ ! -d "$DIR" ]]; then
          echo "Error: Directory '$DIR' does not exist" >&2
          exit 1
      fi

      # Use full nvim config
      NVIM_CONFIG="$HOME/.config/nvim"

      # Check if oil.nvim is available in the config
      if [[ ! -f "$NVIM_CONFIG/init.lua" ]] && [[ ! -f "$NVIM_CONFIG/lua/plugins/oil.lua" ]]; then
          echo "Error: Oil.nvim configuration not found in $NVIM_CONFIG" >&2
          echo "Make sure you have oil.nvim installed in your Neovim configuration" >&2
          exit 1
      fi


      # Launch Neovim with full config and Oil
      exec ${pkgs.neovim}/bin/nvim \
          --cmd "set noswapfile" \
          -c "lua require('oil').open('$DIR')" \
          -c "autocmd FileType oil nnoremap <buffer><silent> q :qa!<CR>" \
          -c "autocmd FileType oil nnoremap <buffer><silent> <Esc> :qa!<CR>"
    '')

    (writeShellScriptBin "oil-simple" ''
      #!/usr/bin/env bash
      # oil-simple: Launch Neovim with Oil.nvim using full LazyVim config but disable problematic plugins
      # Usage: oil-simple [directory]
      #
      # This script opens Neovim with your configuration but disables plugins that cause issues in Zed context

      set -euo pipefail

      # Default to current directory if no argument provided
      DIR="''${1:-.}"

      # Resolve to absolute path
      DIR=$(realpath "$DIR")

      # Check if directory exists
      if [[ ! -d "$DIR" ]]; then
          echo "Error: Directory '$DIR' does not exist" >&2
          exit 1
      fi

      # Use full nvim config but disable problematic plugins
      NVIM_CONFIG="$HOME/.config/nvim"

      # Check if oil.nvim is available in the config
      if [[ ! -f "$NVIM_CONFIG/init.lua" ]] && [[ ! -f "$NVIM_CONFIG/lua/plugins/oil.lua" ]]; then
          echo "Error: Oil.nvim configuration not found in $NVIM_CONFIG" >&2
          echo "Make sure you have oil.nvim installed in your Neovim configuration" >&2
          exit 1
      fi

      # Launch Neovim with full config but disable problematic plugins for Zed context
      exec ${pkgs.neovim}/bin/nvim \
          --cmd "set noswapfile" \
          --cmd "lua vim.g.disable_obsidian = true" \
          --cmd "lua vim.g.disable_copilot = true" \
          -c "lua require('oil').open('$DIR')" \
          -c "autocmd FileType oil nnoremap <buffer><silent> q :qa!<CR>" \
          -c "autocmd FileType oil nnoremap <buffer><silent> <Esc> :qa!<CR>"
    '')

    # Special script for Zed integration
    (writeShellScriptBin "oil-zed" ''
      #!/usr/bin/env bash
      # oil-zed: Launch Oil.nvim with Zed integration using full config
      # Usage: oil-zed [directory]
      #
      # This script opens Oil configured to open files in Zed when selected

      set -euo pipefail

      # Default to current directory if no argument provided
      DIR="''${1:-.}"

      # Resolve to absolute path
      DIR=$(realpath "$DIR")

      # Check if directory exists
      if [[ ! -d "$DIR" ]]; then
          echo "Error: Directory '$DIR' does not exist" >&2
          exit 1
      fi

      # Use full nvim config but disable problematic plugins
      NVIM_CONFIG="$HOME/.config/nvim"

      # Check if oil.nvim is available in the config
      if [[ ! -f "$NVIM_CONFIG/init.lua" ]] && [[ ! -f "$NVIM_CONFIG/lua/plugins/oil.lua" ]]; then
          echo "Error: Oil.nvim configuration not found in $NVIM_CONFIG" >&2
          echo "Make sure you have oil.nvim installed in your Neovim configuration" >&2
          exit 1
      fi

      # Launch Neovim with full config but disable problematic plugins for Zed context
      exec ${pkgs.neovim}/bin/nvim \
          --cmd "set noswapfile" \
          --cmd "lua vim.g.disable_obsidian = true" \
          --cmd "lua vim.g.disable_copilot = true" \
          --cmd "lua vim.g.oil_open_in_zed = true" \
          --cmd "lua function ZedOilOpen() local oil = require('oil'); local entry = oil.get_cursor_entry(); local dir = oil.get_current_dir(); if entry and entry.type == 'file' and dir then vim.fn.jobstart({'zed', dir .. entry.name}, {detach = true}); vim.cmd('qa!'); else require('oil.actions').select.callback(); end end" \
          -c "lua require('oil').open('$DIR')" \
          -c "autocmd FileType oil nnoremap <buffer><silent> q :qa!<CR>" \
          -c "autocmd FileType oil nnoremap <buffer><silent> <Esc> :qa!<CR>" \
          -c "autocmd FileType oil nnoremap <buffer><silent> <CR> :lua ZedOilOpen()<CR>"
    '')

    (writeShellScriptBin "ocd" ''
      #!/usr/bin/env bash
      # ocd: Change directory and launch Oil.nvim with full LazyVim config
      # Usage: ocd [directory]
      #
      # This script combines 'cd' and 'oil' - it changes to the specified directory
      # and then launches Oil.nvim with your complete configuration

      set -euo pipefail

      # Default to current directory if no argument provided
      DIR="''${1:-.}"

      # Resolve to absolute path
      DIR=$(realpath "$DIR")

      # Check if directory exists
      if [[ ! -d "$DIR" ]]; then
          echo "Error: Directory '$DIR' does not exist" >&2
          exit 1
      fi

      # Change to the directory
      cd "$DIR"

      # Display current location
      echo "Changed to: $(pwd)"

      # Use full nvim config
      NVIM_CONFIG="$HOME/.config/nvim"

      # Check if oil.nvim is available in the config
      if [[ ! -f "$NVIM_CONFIG/init.lua" ]] && [[ ! -f "$NVIM_CONFIG/lua/plugins/oil.lua" ]]; then
          echo "Error: Oil.nvim configuration not found in $NVIM_CONFIG" >&2
          echo "Make sure you have oil.nvim installed in your Neovim configuration" >&2
          exit 1
      fi


      # Launch Neovim with full config and Oil
      exec ${pkgs.neovim}/bin/nvim \
          --cmd "set noswapfile" \
          -c "lua require('oil').open('.')" \
          -c "autocmd FileType oil nnoremap <buffer><silent> q :qa!<CR>" \
          -c "autocmd FileType oil nnoremap <buffer><silent> <Esc> :qa!<CR>"
    '')

    (writeShellScriptBin "oil-float" ''
      #!/usr/bin/env bash
      # oil-float: Launch Oil.nvim in floating window mode with full config
      # Usage: oil-float [directory]
      #
      # This script opens Oil.nvim in a floating window for quick file browsing

      set -euo pipefail

      # Default to current directory if no argument provided
      DIR="''${1:-.}"

      # Resolve to absolute path
      DIR=$(realpath "$DIR")

      # Check if directory exists
      if [[ ! -d "$DIR" ]]; then
          echo "Error: Directory '$DIR' does not exist" >&2
          exit 1
      fi

      # Use full nvim config but disable problematic plugins
      NVIM_CONFIG="$HOME/.config/nvim"

      # Check if oil.nvim is available in the config
      if [[ ! -f "$NVIM_CONFIG/init.lua" ]] && [[ ! -f "$NVIM_CONFIG/lua/plugins/oil.lua" ]]; then
          echo "Error: Oil.nvim configuration not found in $NVIM_CONFIG" >&2
          echo "Make sure you have oil.nvim installed in your Neovim configuration" >&2
          exit 1
      fi

      # Launch Neovim with full config but disable problematic plugins for Zed context
      exec ${pkgs.neovim}/bin/nvim \
          --cmd "set noswapfile" \
          --cmd "lua vim.g.disable_obsidian = true" \
          --cmd "lua vim.g.disable_copilot = true" \
          --cmd "lua vim.g.oil_open_in_zed = true" \
          --cmd "lua function ZedOilOpen() local oil = require('oil'); local entry = oil.get_cursor_entry(); local dir = oil.get_current_dir(); if entry and entry.type == 'file' and dir then vim.fn.jobstart({'zed', dir .. entry.name}, {detach = true}); vim.cmd('qa!'); else require('oil.actions').select.callback(); end end" \
          -c "lua require('oil').open_float('$DIR')" \
          -c "autocmd FileType oil nnoremap <buffer><silent> q :qa!<CR>" \
          -c "autocmd FileType oil nnoremap <buffer><silent> <Esc> :qa!<CR>" \
          -c "autocmd FileType oil nnoremap <buffer><silent> <CR> :lua ZedOilOpen()<CR>"
    '')

    # Zed-specific helper scripts
    (writeShellScriptBin "zed-oil" ''
      #!/usr/bin/env bash
      # zed-oil: Helper script for Zed Oil task
      # Usage: zed-oil

      if [[ -n "$ZED_FILE" ]]; then
          DIR=$(dirname "$ZED_FILE")
      else
          DIR="$ZED_WORKTREE_ROOT"
      fi

      exec oil-zed "$DIR"
    '')

    (writeShellScriptBin "zed-oil-float" ''
      #!/usr/bin/env bash
      # zed-oil-float: Helper script for Zed Oil Float task
      # Usage: zed-oil-float

      if [[ -n "$ZED_FILE" ]]; then
          DIR=$(dirname "$ZED_FILE")
      else
          DIR="$ZED_WORKTREE_ROOT"
      fi

      exec oil-float "$DIR"
    '')

    (writeShellScriptBin "zed-ocd" ''
      #!/usr/bin/env bash
      # zed-ocd: Helper script for Zed Oil CD task
      # Usage: zed-ocd

      if [[ -n "$ZED_FILE" ]]; then
          DIR=$(dirname "$ZED_FILE")
      else
          DIR="$ZED_WORKTREE_ROOT"
      fi

      exec ocd "$DIR"
    '')

    # Script dinámico para Oil en directorio del archivo actual
    (writeShellScriptBin "oil-file-dir" ''
      #!/usr/bin/env bash
      # oil-file-dir: Open Oil in the directory of the current file
      # Usage: oil-file-dir [fallback_dir]

      FALLBACK_DIR="''${1:-$PWD}"

      if [[ -n "$ZED_FILE" && -f "$ZED_FILE" ]]; then
          DIR=$(dirname "$ZED_FILE")
          echo "Opening Oil in file directory: $DIR"
      else
          DIR="$FALLBACK_DIR"
          echo "Opening Oil in fallback directory: $DIR"
      fi

      exec oil-zed "$DIR"
    '')

    # Script dinámico para Oil Float en directorio del archivo actual
    (writeShellScriptBin "oil-float-file-dir" ''
      #!/usr/bin/env bash
      # oil-float-file-dir: Open Oil Float in the directory of the current file
      # Usage: oil-float-file-dir [fallback_dir]

      FALLBACK_DIR="''${1:-$PWD}"

      if [[ -n "$ZED_FILE" && -f "$ZED_FILE" ]]; then
          DIR=$(dirname "$ZED_FILE")
          echo "Opening Oil Float in file directory: $DIR"
      else
          DIR="$FALLBACK_DIR"
          echo "Opening Oil Float in fallback directory: $DIR"
      fi

      exec oil-float "$DIR"
    '')
  ];

  # Copiar la configuración mínima de Oil
  home.activation.copyOilMinimal = lib.hm.dag.entryAfter ["writeBoundary"] ''
    echo "Copying Oil minimal configuration..."

    # Crear directorio para la configuración mínima
    OIL_MINIMAL_DIR="$HOME/.config/nvim-oil-minimal"
    rm -rf "$OIL_MINIMAL_DIR"
    mkdir -p "$OIL_MINIMAL_DIR"

    # Copiar configuración mínima
    cp -r ${toString ./nvim-oil-minimal}/* "$OIL_MINIMAL_DIR/" 2>/dev/null || true
    chmod -R u+w "$OIL_MINIMAL_DIR" 2>/dev/null || true

    echo "Oil minimal configuration copied to $OIL_MINIMAL_DIR"
  '';

  # Configuración de shells para asegurar que los scripts estén en PATH
  programs.fish = {
    enable = true;
    shellAliases = {
      "o" = "oil";
      "oo" = "oil .";
      "of" = "oil-float";
      "oz" = "oil-zed";
    };
    shellInit = ''
      # Asegurar que nix-profile esté en PATH
      if not contains ~/.nix-profile/bin $PATH
        set -gx PATH ~/.nix-profile/bin $PATH
      end
    '';
  };

  programs.zsh = {
    enable = true;
    shellAliases = {
      "o" = "oil";
      "oo" = "oil .";
      "of" = "oil-float";
      "oz" = "oil-zed";
    };
    initExtra = ''
      # Asegurar que nix-profile esté en PATH
      export PATH="$HOME/.nix-profile/bin:$PATH"
    '';
  };

  programs.nushell = {
    enable = true;
    extraConfig = ''
      # Asegurar que nix-profile esté en PATH
      $env.PATH = ($env.PATH | split row (char esep) | prepend $"($env.HOME)/.nix-profile/bin")

      # Alias para Oil
      alias o = oil
      alias oo = oil .
      alias of = oil-float
      alias oz = oil-zed
    '';
  };
}
