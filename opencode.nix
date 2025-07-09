{ pkgs, config, lib, ... }:

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

    # Set PATH to include all required tools
    export PATH="${pkgs.unzip}/bin:${pkgs.curl}/bin:${pkgs.gawk}/bin:${pkgs.gnutar}/bin:${pkgs.gzip}/bin:${pkgs.coreutils}/bin:${pkgs.gh}/bin:$PATH"

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
  programs.fish.shellAliases.opencode-config = "nvim ~/.opencode.json";
  programs.zsh.shellAliases.opencode-config = "nvim ~/.opencode.json";
  programs.nushell.shellAliases.opencode-config = "nvim ~/.opencode.json";
}
