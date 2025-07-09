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

      echo "🚀 Installing latest OpenCode..."

      # Set PATH to include all required tools
      export PATH="${pkgs.curl}/bin:${pkgs.gawk}/bin:${pkgs.gnutar}/bin:${pkgs.gzip}/bin:${pkgs.coreutils}/bin:${pkgs.gh}/bin:$PATH"

      # Always install/update to latest version
      echo "Installing latest OpenCode version..."
      curl -fsSL https://raw.githubusercontent.com/opencode-ai/opencode/refs/heads/main/install | bash

      if [ -f "$OPENCODE_BIN" ]; then
        INSTALLED_VERSION=$("$OPENCODE_BIN" --version 2>/dev/null || echo "unknown")
        echo "✅ OpenCode v$INSTALLED_VERSION installed successfully!"
      else
        echo "❌ OpenCode installation failed"
        exit 1
      fi

      # Install GitHub Copilot extension if not present
      if ! gh extension list | grep -q "github/gh-copilot"; then
        echo "📦 Installing GitHub Copilot extension..."
        gh extension install github/gh-copilot
        echo "✅ GitHub Copilot extension installed!"
      else
        echo "✅ GitHub Copilot extension already installed"
      fi

      # Create default config if it doesn't exist
      if [ ! -f ~/.opencode.json ]; then
        echo "📝 Creating default OpenCode configuration..."
        cat > ~/.opencode.json << 'EOF'
{
  "providers": {
    "copilot": {
      "disabled": false
    }
  },
  "agents": {
    "coder": {
      "model": "copilot.claude-sonnet-4",
      "maxTokens": 16000,
      "reasoningEffort": ""
    },
    "task": {
      "model": "copilot.claude-sonnet-4",
      "maxTokens": 5000,
      "reasoningEffort": ""
    },
    "title": {
      "model": "copilot.gpt-4o-mini",
      "maxTokens": 80,
      "reasoningEffort": ""
    }
  },
  "tui": {
    "theme": "catppuccin"
  },
  "shell": {
    "path": "${pkgs.fish}/bin/fish",
    "args": ["-l"]
  },
  "autoCompact": true,
  "debug": false
}
EOF
        echo "✅ Default configuration created at ~/.opencode.json"
        echo "📖 You can edit the configuration with: opencode-config"
        echo "🤖 Available models: copilot.claude-sonnet-4, copilot.claude-3.5-sonnet, copilot.gpt-4o, copilot.gpt-4o-mini, copilot.gpt-4"
      else
        echo "✅ OpenCode configuration already exists at ~/.opencode.json"
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

    echo "🚀 Installing latest OpenCode..."

    # Set PATH to include all required tools
    export PATH="${pkgs.curl}/bin:${pkgs.gawk}/bin:${pkgs.gnutar}/bin:${pkgs.gzip}/bin:${pkgs.coreutils}/bin:${pkgs.gh}/bin:$PATH"

    # Always install/update to latest version
    echo "Installing latest OpenCode version..."
    curl -fsSL https://raw.githubusercontent.com/opencode-ai/opencode/refs/heads/main/install | bash

    if [ -f "$OPENCODE_BIN" ]; then
      INSTALLED_VERSION=$("$OPENCODE_BIN" --version 2>/dev/null || echo "unknown")
      echo "✅ OpenCode v$INSTALLED_VERSION installed successfully!"
    else
      echo "❌ OpenCode installation failed"
    fi

    # Install GitHub Copilot extension if not present
    if ! gh extension list | grep -q "github/gh-copilot"; then
      echo "📦 Installing GitHub Copilot extension..."
      gh extension install github/gh-copilot
      echo "✅ GitHub Copilot extension installed!"
    else
      echo "✅ GitHub Copilot extension already installed"
    fi

    # Create default config if it doesn't exist
    if [ ! -f ~/.opencode.json ]; then
      echo "📝 Creating default OpenCode configuration..."
      cat > ~/.opencode.json << 'EOF'
{
  "providers": {
    "copilot": {
      "disabled": false
    }
  },
  "agents": {
    "coder": {
      "model": "copilot.claude-sonnet-4",
      "maxTokens": 16000,
      "reasoningEffort": ""
    },
    "task": {
      "model": "copilot.claude-sonnet-4",
      "maxTokens": 5000,
      "reasoningEffort": ""
    },
    "title": {
      "model": "copilot.gpt-4o-mini",
      "maxTokens": 80,
      "reasoningEffort": ""
    }
  },
  "tui": {
    "theme": "catppuccin"
  },
  "shell": {
    "path": "${pkgs.fish}/bin/fish",
    "args": ["-l"]
  },
  "autoCompact": true,
  "debug": false
}
EOF
      echo "✅ Default configuration created at ~/.opencode.json"
      echo "📖 You can edit the configuration with: opencode-config"
      echo "🤖 Available models: copilot.claude-sonnet-4, copilot.claude-3.5-sonnet, copilot.gpt-4o, copilot.gpt-4o-mini, copilot.gpt-4"
    else
      echo "✅ OpenCode configuration already exists at ~/.opencode.json"
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
