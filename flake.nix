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

  outputs = { nixpkgs, nixpkgs-unstable, home-manager, flake-utils, ... }:
    let
      # Support macOS systems only
      supportedSystems = [ "x86_64-darwin" "aarch64-darwin" ];
      
      # Function to create home configuration for a specific system
      mkHomeConfiguration = system:
        let
          pkgs = import nixpkgs {
            inherit system;
            config.allowUnfree = true;
          };
          
          unstablePkgs = import nixpkgs-unstable {
            inherit system;
            config.allowUnfree = true;
          };
        in
        home-manager.lib.homeManagerConfiguration {
          inherit pkgs;
          
          # Pass extraSpecialArgs to make unstablePkgs available in modules
          extraSpecialArgs = {
            inherit unstablePkgs;
          };
          
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
            ./claude.nix  # Claude Code CLI configuration
            {
              # Personal data
              home.username = "YourUser";  # Replace with your username
              home.homeDirectory = "/Users/YourUser";  # macOS home directory
              home.stateVersion = "24.11";  # State version

              # Base packages that should be available everywhere
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

              # Enable programs explicitly (critical for binaries to appear)
              # All program enables are centralized here
              programs.neovim.enable = true;
              programs.fish.enable = true;
              programs.nushell.enable = true;
              programs.starship.enable = true;
              programs.zsh.enable = true;
              programs.git.enable = true;
              programs.gh.enable = true;  # GitHub CLI
              programs.home-manager.enable = true;
              # Note: tmux is configured via home.file in tmux.nix, not programs.tmux

              # Allow unfree packages
              nixpkgs.config.allowUnfree = true;
            }
          ];
        };
    in
    {
      # Home Manager configurations for each system
      homeConfigurations = {
        # macOS system configurations
        "gentleman-macos-intel" = mkHomeConfiguration "x86_64-darwin";
        "gentleman-macos-arm" = mkHomeConfiguration "aarch64-darwin";
        
        # Default to Apple Silicon
        "gentleman" = mkHomeConfiguration "aarch64-darwin";
      };
    };
}
