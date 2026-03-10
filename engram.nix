{ pkgs, lib, ... }:

{
  # Engram — memory layer for AI agents
  # Binary installed via `go install`, requires `go` in home.packages (see flake.nix)

  # Ensure $HOME/go/bin is in PATH so engram is available in all shells
  home.sessionPath = [
    "$HOME/go/bin"
  ];

  # Auto-install engram on home-manager activation
  home.activation.installEngram = lib.hm.dag.entryAfter [ "linkGeneration" ] ''
    echo "🔧 Setting up Engram..."

    export PATH="${pkgs.go}/bin:$HOME/go/bin:$PATH"
    export GOPATH="$HOME/go"

    # Check if engram is already installed and working
    if command -v engram &>/dev/null; then
      echo "✅ Engram already installed"
    else
      echo "🚀 Installing Engram via go install..."
      ${pkgs.go}/bin/go install github.com/Gentleman-Programming/engram/cmd/engram@latest

      if command -v engram &>/dev/null; then
        echo "✅ Engram installed successfully!"
      else
        echo "❌ Engram installation failed"
      fi
    fi

    echo ""
    echo "🎉 Engram setup complete!"
  '';
}
