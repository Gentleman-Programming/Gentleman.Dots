{ ... }:

{
  # Atuin — Gentleman theme profiles for the shell history TUI.
  # The active theme is selected at runtime through ATUIN_THEME__NAME, which
  # fish exports from the persisted Gentleman theme marker (see fish.nix and
  # fish/functions/fish-theme.fish). Atuin resolves the name against these
  # files in ~/.config/atuin/themes/<name>.toml.
  home.file = {
    ".config/atuin/themes/gentleman.toml".source = ./atuin/themes/gentleman.toml;
    ".config/atuin/themes/gentleman-cute.toml".source = ./atuin/themes/gentleman-cute.toml;
  };
}
