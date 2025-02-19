{
  description = "Gentleman: Single config for all systems in one go";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    home-manager = {
      url = "github:nix-community/home-manager";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { nixpkgs, home-manager, ... }:
    let
      system = "aarch64-darwin";  # Make sure this matches your system
      pkgs = import nixpkgs { inherit system; };
    in {
      homeConfigurations = {
        "gentleman" =
          home-manager.lib.homeManagerConfiguration {
            inherit pkgs;
            modules = [
              ./nushell.nix  
              ./ghostty.nix  
              ./wezterm.nix  
              ./zellij.nix   
              ./starship.nix 
              ./nvim.nix     
              {
                # Personal data
                home.username = "YourUser";  # Change this to your username
                home.homeDirectory = "/Users/YourUser";  # On macOS; on Linux it would be "/home/yourUser"
                home.stateVersion = "24.11";

                # Group packages by categories to keep everything organized
                home.packages = with pkgs; [
                  # ─── Terminals and window managers ──────────────────────────────
                  zellij
                  nushell

                  # ─── Development tools and utilities ─────────────────────────
                  volta
                  carapace
                  zoxide
                  atuin
                  jq
                  bash
                  starship
                  fzf
                  neovim
                  nodejs

                  # ─── Compilers, search tools, and system utilities ─────────────
                  gcc
                  fd
                  ripgrep
                  coreutils
                  bat
                  lazygit

                  # ─── Nerd Fonts ────────────────────────────────────────────────────
                  # Adding IosevkaTerm NF to improve terminal look
                  nerd-fonts.iosevka-term
                ];

                # Enable specific programs
                programs.nushell.enable = true;
                programs.starship.enable = true;

                # Custom activation: create directories for Obsidian
                home.activation.createObsidianDirs = ''
                  mkdir -p "$HOME/.config/obsidian/templates"
                '';
              }
            ];
          };
      };
    };
}
