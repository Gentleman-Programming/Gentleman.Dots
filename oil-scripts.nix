{ lib, pkgs, ... }:
{
  home.packages = with pkgs; [
    # Oil.nvim helper scripts
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
      cd "$DIR"
      exec ${pkgs.neovim}/bin/nvim \
          --cmd "set noswapfile" \
          -c "lua require('oil').open()" \
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
      cd "$DIR"
      exec ${pkgs.neovim}/bin/nvim \
          --cmd "set noswapfile" \
          --cmd "lua vim.g.disable_obsidian = true" \
          --cmd "lua vim.g.disable_copilot = true" \
          -c "lua require('oil').open()" \
          -c "autocmd FileType oil nnoremap <buffer><silent> q :qa!<CR>" \
          -c "autocmd FileType oil nnoremap <buffer><silent> <Esc> :qa!<CR>"
    '')

    # Special script for Zed integration
    (writeShellScriptBin "oil-zed" ''
      #!/usr/bin/env bash
      # oil-zed: Launch Oil.nvim with Zed integration using minimal config
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

      # Use minimal nvim config to avoid plugin conflicts
      OIL_MINIMAL_DIR="$HOME/.config/nvim-oil-minimal"

      # Check if minimal config exists
      if [[ ! -f "$OIL_MINIMAL_DIR/init.lua" ]]; then
          echo "Error: Minimal Oil.nvim configuration not found at $OIL_MINIMAL_DIR" >&2
          echo "Run 'home-manager switch' to ensure the minimal config is installed" >&2
          exit 1
      fi

      # Launch Neovim with minimal config and Zed integration
      cd "$DIR"
      exec ${pkgs.neovim}/bin/nvim \
          --cmd "set noswapfile" \
          -u "$OIL_MINIMAL_DIR/init.lua" \
          -c "lua vim.g.oil_open_in_zed = true" \
          -c "lua require('oil').open()" \
          -c "autocmd FileType oil nnoremap <buffer><silent> q :qa!<CR>" \
          -c "autocmd FileType oil nnoremap <buffer><silent> <Esc> :qa!<CR>"
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
      # oil-float: Launch Oil.nvim in floating window mode with minimal config
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

      # Use minimal nvim config to avoid plugin conflicts
      OIL_MINIMAL_DIR="$HOME/.config/nvim-oil-minimal"

      # Check if minimal config exists
      if [[ ! -f "$OIL_MINIMAL_DIR/init.lua" ]]; then
          echo "Error: Minimal Oil.nvim configuration not found at $OIL_MINIMAL_DIR" >&2
          echo "Run 'home-manager switch' to ensure the minimal config is installed" >&2
          exit 1
      fi

      # Launch Neovim with minimal config and Oil in floating mode
      cd "$DIR"
      exec ${pkgs.neovim}/bin/nvim \
          --cmd "set noswapfile" \
          -u "$OIL_MINIMAL_DIR/init.lua" \
          -c "lua vim.g.oil_open_in_zed = true" \
          -c "lua require('oil').open_float()" \
          -c "autocmd FileType oil nnoremap <buffer><silent> q :qa!<CR>" \
          -c "autocmd FileType oil nnoremap <buffer><silent> <Esc> :qa!<CR>"
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

    # Dynamic script: open Oil in the current file's directory
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

    # Dynamic script: open Oil Float in the current file's directory
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

  # Copy minimal Oil configuration
  home.activation.copyOilMinimal = lib.hm.dag.entryAfter ["writeBoundary"] ''
    echo "Copying Oil minimal configuration..."

    # Create directory for minimal configuration
    OIL_MINIMAL_DIR="$HOME/.config/nvim-oil-minimal"
    rm -rf "$OIL_MINIMAL_DIR"
    mkdir -p "$OIL_MINIMAL_DIR"

    # Copy minimal configuration
    cp -r ${toString ./nvim-oil-minimal}/* "$OIL_MINIMAL_DIR/" 2>/dev/null || true
    chmod -R u+w "$OIL_MINIMAL_DIR" 2>/dev/null || true

    echo "Oil minimal configuration copied to $OIL_MINIMAL_DIR"
  '';

  # Shell configuration to ensure scripts are in PATH
  programs.fish = {
    enable = true;
    shellAliases = {
      "o" = "oil";
      "oo" = "oil .";
      "of" = "oil-float";
      "oz" = "oil-zed";
    };
    shellInit = ''
      # Ensure nix-profile bin is in PATH
      if not contains ~/.nix-profile/bin $PATH
        set -gx PATH ~/.nix-profile/bin $PATH
      end
    '';
  };

  # Zsh aliases handled via home.file in zsh.nix - don't enable programs.zsh here

  # Extend Nushell configuration with Oil aliases
  programs.nushell.extraConfig = ''
    # Ensure nix-profile bin is in PATH
    $env.PATH = ($env.PATH | split row (char esep) | prepend $"($env.HOME)/.nix-profile/bin")

    # Aliases for Oil
    alias o = oil
    alias oo = oil .
    alias of = oil-float
    alias oz = oil-zed
  '';
}
