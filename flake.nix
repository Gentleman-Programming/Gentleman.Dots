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
      system = "aarch64-darwin";  # make sure this matches your system
      pkgs = import nixpkgs { inherit system; };
    in {
      homeConfigurations = {
        "gentleman" =
          home-manager.lib.homeManagerConfiguration {
            inherit pkgs;
            modules = [
              ./nushell.nix  
              ./wezterm.nix  
              ./zellij.nix   
              ./starship.nix 
              ./nvim.nix     
              {
                home.username = "YourUser";  # ensure this is your username
                home.homeDirectory = "/Users/YourUser"; # ensure this is your home directory
                home.stateVersion = "24.11";  # use a valid version
                home.packages = [
                  pkgs.wezterm
                  pkgs.zellij
                  pkgs.nushell
                  pkgs.volta
                  pkgs.carapace
                  pkgs.zoxide
                  pkgs.atuin
                  pkgs.jq
                  pkgs.bash
                  pkgs.starship
                  pkgs.fzf
                  pkgs.neovim
                  pkgs.nodejs  # npm is included with nodejs
                  pkgs.gcc
                  pkgs.fd
                  pkgs.ripgrep
                  pkgs.coreutils
                  pkgs.bat
                  pkgs.lazygit
                ];
                programs.nushell.enable = true;
                programs.starship.enable = true;

                home.activation.createObsidianDirs = ''
                  mkdir -p "$HOME/.config/obsidian/templates"
                '';

              }
            ];
          };
      };
    };
}
