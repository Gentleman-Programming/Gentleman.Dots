{ pkgs, lib, ... }:

{
  # Install GitHub CLI and tools for Copilot authentication and OpenCode installation
  home.packages = [
    pkgs.gh
    pkgs.curl
    pkgs.gawk
    pkgs.gnutar
    pkgs.gzip
    pkgs.coreutils
    pkgs.unzip
  ];

  # Setup GitHub CLI
  programs.gh = {
    enable = true;
    settings = {
      version = "1";
    };
  };

  # Add OpenCode to PATH for all shells
  home.sessionPath = [
    "$HOME/.opencode/bin"
  ];

  # Create OpenCode installation script for manual use
  home.file."bin/install-opencode" = {
    text = ''
      #!/usr/bin/env bash
      set -e

      OPENCODE_DIR="$HOME/.opencode"
      OPENCODE_BIN="$OPENCODE_DIR/bin/opencode"

      # Create required cache directories
      mkdir -p "$HOME/.cache/nvim/opencode"

      # Set PATH to include all required tools
      export PATH="${pkgs.unzip}/bin:${pkgs.curl}/bin:${pkgs.gawk}/bin:${pkgs.gnutar}/bin:${pkgs.gzip}/bin:${pkgs.coreutils}/bin:${pkgs.gh}/bin:$PATH"

      # Check if OpenCode is already installed and working
      if [ -f "$OPENCODE_BIN" ] && "$OPENCODE_BIN" --version &>/dev/null; then
        INSTALLED_VERSION=$("$OPENCODE_BIN" --version 2>/dev/null | head -n1 || echo "unknown")
        echo "✅ OpenCode already installed: $INSTALLED_VERSION"
        read -p "Do you want to reinstall/update? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
          echo "Skipping OpenCode installation"
        else
          echo "🚀 Reinstalling OpenCode..."
          curl -fsSL https://opencode.ai/install | bash 
          if [ -f "$OPENCODE_BIN" ]; then
            INSTALLED_VERSION=$("$OPENCODE_BIN" --version 2>/dev/null | head -n1 || echo "unknown")
            echo "✅ OpenCode v$INSTALLED_VERSION installed successfully!"
          else
            echo "❌ OpenCode installation failed"
            exit 1
          fi
        fi
      else
        echo "🚀 Installing latest OpenCode..."
        curl -fsSL https://opencode.ai/install | bash 
        if [ -f "$OPENCODE_BIN" ]; then
          INSTALLED_VERSION=$("$OPENCODE_BIN" --version 2>/dev/null | head -n1 || echo "unknown")
          echo "✅ OpenCode v$INSTALLED_VERSION installed successfully!"
        else
          echo "❌ OpenCode installation failed"
          exit 1
        fi
      fi

      # Install GitHub Copilot extension if not present
      if ! gh extension list | grep -q "github/gh-copilot"; then
        echo "📦 Installing GitHub Copilot extension..."
        gh extension install github/gh-copilot
        echo "✅ GitHub Copilot extension installed!"
      else
        echo "✅ GitHub Copilot extension already installed"
      fi
      echo ""
      echo "🎉 OpenCode setup complete!"
      echo "Usage: opencode | opencode-config | gh auth status"
    '';
    executable = true;
  };

  # Auto-install OpenCode on home-manager activation
  home.activation.installOpenCode = lib.hm.dag.entryAfter ["linkGeneration"] ''
    echo "🔧 Setting up OpenCode..."

    OPENCODE_DIR="$HOME/.opencode"
    OPENCODE_BIN="$OPENCODE_DIR/bin/opencode"

    # Create required cache directories
    mkdir -p "$HOME/.cache/nvim/opencode"

    # Set PATH to include all required tools
    export PATH="${pkgs.unzip}/bin:${pkgs.curl}/bin:${pkgs.gawk}/bin:${pkgs.gnutar}/bin:${pkgs.gzip}/bin:${pkgs.coreutils}/bin:${pkgs.gh}/bin:$PATH"

    # Copy bundled config and themes into user config
    OPENCODE_SRC="${toString ./opencode}"
    OPENCODE_DST="$HOME/.config/opencode"
    mkdir -p "$OPENCODE_DST/themes"
    
    # Copy main config file
    if [ -f "$OPENCODE_SRC/opencode.json" ]; then
      cp -f "$OPENCODE_SRC/opencode.json" "$OPENCODE_DST/" 2>/dev/null || true
      echo "⚙️ Copied OpenCode config to $OPENCODE_DST"
    else
      echo "⚠️ Config source not found: $OPENCODE_SRC/opencode.json"
    fi
    
    # Copy themes
    if [ -d "$OPENCODE_SRC/themes" ]; then
      cp -f "$OPENCODE_SRC/themes"/* "$OPENCODE_DST/themes/" 2>/dev/null || true
      echo "🎨 Copied OpenCode themes to $OPENCODE_DST/themes"
    else
      echo "⚠️ Themes source not found: $OPENCODE_SRC/themes"
    fi

    # Copy AGENTS.md (referenced by agents via {file:./AGENTS.md})
    if [ -f "$OPENCODE_SRC/AGENTS.md" ]; then
      cp -f "$OPENCODE_SRC/AGENTS.md" "$OPENCODE_DST/" 2>/dev/null || true
      echo "📋 Copied AGENTS.md to $OPENCODE_DST"
    fi

    # Copy skills
    if [ -d "$OPENCODE_SRC/skills" ]; then
      cp -rf "$OPENCODE_SRC/skills" "$OPENCODE_DST/" 2>/dev/null || true
      echo "🧠 Copied OpenCode skills to $OPENCODE_DST/skills"
    fi

    # Copy commands
    if [ -d "$OPENCODE_SRC/commands" ]; then
      mkdir -p "$OPENCODE_DST/commands"
      cp -f "$OPENCODE_SRC/commands"/* "$OPENCODE_DST/commands/" 2>/dev/null || true
      echo "⚡ Copied OpenCode commands to $OPENCODE_DST/commands"
    else
      echo "⚠️ Commands source not found: $OPENCODE_SRC/commands"
    fi

    # Copy plugins
    if [ -d "$OPENCODE_SRC/plugins" ]; then
      mkdir -p "$OPENCODE_DST/plugins"
      cp -f "$OPENCODE_SRC/plugins"/* "$OPENCODE_DST/plugins/" 2>/dev/null || true
      echo "🔌 Copied OpenCode plugins to $OPENCODE_DST/plugins"
    else
      echo "⚠️ Plugins source not found: $OPENCODE_SRC/plugins"
    fi

    # Check if OpenCode is already installed and working
    if [ -f "$OPENCODE_BIN" ] && "$OPENCODE_BIN" --version &>/dev/null; then
      INSTALLED_VERSION=$("$OPENCODE_BIN" --version 2>/dev/null | head -n1 || echo "unknown")
      echo "✅ OpenCode already installed: $INSTALLED_VERSION"
    else
      echo "🚀 Installing latest OpenCode..."
      curl -fsSL https://opencode.ai/install | bash 

      if [ -f "$OPENCODE_BIN" ]; then
        INSTALLED_VERSION=$("$OPENCODE_BIN" --version 2>/dev/null | head -n1 || echo "unknown")
        echo "✅ OpenCode v$INSTALLED_VERSION installed successfully!"
      else
        echo "❌ OpenCode installation failed"
      fi
    fi

    # Install GitHub Copilot extension if not present
    if ! gh extension list | grep -q "github/gh-copilot"; then
      echo "📦 Installing GitHub Copilot extension..."
      gh extension install github/gh-copilot
      echo "✅ GitHub Copilot extension installed!"
    else
      echo "✅ GitHub Copilot extension already installed"
    fi
    echo ""
    echo "🎉 OpenCode setup complete!"
    echo "Usage: opencode | opencode-config | gh auth status"
  '';

  # Add aliases for all configured shells
  programs.fish.shellAliases.opencode-config = "nvim ~/.config/opencode/opencode.json";
  programs.zsh.shellAliases.opencode-config = "nvim ~/.config/opencode/opencode.json";
  programs.nushell.shellAliases.opencode-config = "nvim ~/.config/opencode/opencode.json";
}
