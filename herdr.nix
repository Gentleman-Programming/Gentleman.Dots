{ pkgs, lib, ... }:

{
  # Herdr — agent multiplexer that lives in your terminal (https://herdr.dev)
  # Binary installed via Homebrew (homebrew-core). Config is managed declaratively below.

  # Auto-install herdr on home-manager activation if it is missing.
  # Guarded so a missing/failed brew never breaks the activation (same approach as engram.nix).
  home.activation.installHerdr = lib.hm.dag.entryAfter [ "linkGeneration" ] ''
    echo "🔧 Setting up Herdr..."

    if command -v herdr &>/dev/null; then
      echo "✅ Herdr already installed"
    elif command -v brew &>/dev/null; then
      echo "🚀 Installing Herdr via Homebrew..."
      brew install herdr || echo "❌ Herdr installation failed (run 'brew install herdr' manually)"
    else
      echo "⚠️  Homebrew not found — install Herdr manually: brew install herdr"
    fi
  '';

  home.file = {
    # Custom keybindings are MERGED on top of herdr's built-in v2 defaults.
    # A partial config only overrides what it declares; `herdr config reset-keys`
    # backs this up and restores defaults. Apply changes live with `herdr server reload-config`.
    ".config/herdr/config.toml" = {
      text = ''
[keys]
# Jump straight to agent N (1..9) as listed in the sidebar.
# Independent from prefix+1..9 (tabs) and prefix+w (workspaces).
focus_agent = "prefix+shift+1..9"
      '';
    };
  };
}
