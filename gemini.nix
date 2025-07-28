{ pkgs, config, lib, ... }:

{
  # Install Bun for Gemini CLI
  home.packages = [
    pkgs.bun
  ];

  # Auto-install Gemini CLI on home-manager activation
  home.activation.installGeminiCLI = lib.hm.dag.entryAfter ["linkGeneration"] ''
    echo "🔧 Setting up Gemini CLI..."

    # Install Gemini CLI globally via bun
    echo "📦 Installing Gemini CLI..."
    ${pkgs.bun}/bin/bun install -g "@google/gemini-cli"
    echo "✅ Gemini CLI installed!"
    echo ""
    echo "🎉 Gemini CLI setup complete!"
    echo "Usage: gemini"
  '';
}
