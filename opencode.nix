{ pkgs, lib, ... }:

{
  home.packages = [
    pkgs.gh
    pkgs.curl
    pkgs.gawk
    pkgs.gnutar
    pkgs.gzip
    pkgs.coreutils
    pkgs.unzip
  ];

  programs.gh = {
    enable = true;
    settings = {
      version = "1";
    };
  };

  home.sessionPath = [
    "$HOME/.opencode/bin"
  ];

  # Install OpenCode binary and GitHub Copilot extension on activation
  home.activation.installOpenCode = lib.hm.dag.entryAfter ["linkGeneration"] ''
    export PATH="${pkgs.unzip}/bin:${pkgs.curl}/bin:${pkgs.gawk}/bin:${pkgs.gnutar}/bin:${pkgs.gzip}/bin:${pkgs.coreutils}/bin:${pkgs.gh}/bin:$PATH"

    OPENCODE_BIN="$HOME/.opencode/bin/opencode"

    if [ ! -f "$OPENCODE_BIN" ] || ! "$OPENCODE_BIN" --version &>/dev/null; then
      curl -fsSL https://opencode.ai/install | bash
    fi

    if gh copilot --help >/dev/null 2>&1; then
      :
    elif ! gh extension list | grep -q "github/gh-copilot"; then
      gh extension install github/gh-copilot
    fi
  '';

  programs.fish.shellAliases.opencode-config = "nvim ~/.config/opencode/opencode.json";
  programs.zsh.shellAliases.opencode-config = "nvim ~/.config/opencode/opencode.json";
  programs.nushell.shellAliases.opencode-config = "nvim ~/.config/opencode/opencode.json";
}
