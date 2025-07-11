{
  description = "Gentleman: Single config for all systems in one go";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    nixpkgs-unstable.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    home-manager = {
      url = "github:nix-community/home-manager";  # Home Manager repository
      inputs.nixpkgs.follows = "nixpkgs";  # Follow nixpkgs input
    };
    flake-utils.url = "github:numtide/flake-utils";  # Flake utilities
  };

  outputs = { nixpkgs, nixpkgs-unstable, home-manager, ... }:
    let
      unstablePkgs = import nixpkgs-unstable {
        system = "aarch64-darwin";  # Make sure this matches your system
        config.allowUnfree = true;
      };
      system = "aarch64-darwin";  # Make sure this matches your system
      pkgs = import nixpkgs { inherit system; };  # Import nixpkgs for the specified system
    in {
      homeConfigurations = {
        "gentleman" =
          home-manager.lib.homeManagerConfiguration {
            inherit pkgs;
            modules = [
              ./nushell.nix  # Nushell configuration
              ./ghostty.nix  # Ghostty configuration
              ./zed.nix  # Zed configuration
              ./television.nix  # Television configuration
              ./wezterm.nix  # WezTerm configuration
              # ./zellij.nix  # Zellij configuration (commented out)
              ./tmux.nix  # Tmux configuration
              ./fish.nix  # Fish shell configuration
              ./starship.nix  # Starship prompt configuration
              ./nvim.nix  # Neovim configuration
              ./zsh.nix  # Zsh configuration
              ./oil-scripts.nix  # Oil.nvim scripts configuration
              ./opencode.nix  # OpenCode AI assistant configuration
              {
                # Personal data
                home.username = "YourUser";  # Replace with your username
                home.homeDirectory = "/Users/YourUser/";  # Replace with your home directory
                home.stateVersion = "24.11";  # State version

                home.packages = with pkgs; [
                  # ─── Terminals and utilities ───
                  # zellij
                  tmux
                  fish
                  zsh
                  nushell

                  # ─── Development tools ───
                  volta
                  carapace
                  zoxide
                  atuin
                  jq
                  bash
                  starship
                  fzf
                  unstablePkgs.neovim
                  nodejs
                  bun
                  cargo
                  go
                  nil
                  unstablePkgs.nixd


                  # ─── Compilers and system utilities ───
                  gcc
                  fd
                  ripgrep
                  coreutils
                  unzip
                  bat
                  lazygit
                  yazi
                  television

                  # ─── Nerd Fonts ───
                  nerd-fonts.iosevka-term
                ];
              }
            ];
          };
      };
    };
}
