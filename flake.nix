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
      system = "aarch64-darwin";  # Asegurate de que coincida con tu sistema
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
                # Datos personales
                home.username = "YourUser";
                home.homeDirectory = "/Users/YourUser/"; 
                home.stateVersion = "24.11";

                home.packages = with pkgs; [
                  # ─── Terminals y utilidades ───
                  zellij
                  nushell

                  # ─── Herramientas de desarrollo ───
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
                  lazygit
                  bun

                  # ─── Compiladores y utilidades de sistema ───
                  gcc
                  fd
                  ripgrep
                  coreutils
                  bat
                  lazygit

                  # ─── Nerd Fonts ───
                  nerd-fonts.iosevka-term
                ];

                programs.nushell.enable = true;
                programs.starship.enable = true;
              }
            ];
          };
      };
    };
}
