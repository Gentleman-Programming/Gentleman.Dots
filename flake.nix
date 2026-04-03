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
      # ─── Cross-platform ─── Support macOS and Linux systems
      supportedSystems = [
        "x86_64-darwin"
        "aarch64-darwin"
        "x86_64-linux"
        "aarch64-linux"
      ];
      
      # ─── User Configuration ───
      # Change this to your macOS username
      username = "alanbuscaglia";

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

          # ─── Platform detection ───
          isDarwin = pkgs.stdenv.isDarwin;

          # ─── Cross-platform modules ───
          # These modules work on both macOS and Linux
          commonModules = [
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
            ./engram.nix  # Engram memory layer for AI agents
          ];

          # ─── macOS-only modules ───
          # These modules depend on macOS-specific tools and APIs
          darwinOnlyModules = [
            ./yabai.nix  # Yabai window manager (macOS only)
            ./skhd.nix  # Skhd hotkey daemon (macOS only)
            # ./simple-bar.nix  # simple-bar for Übersicht (disabled - using sketchybar)
            ./sketchybar.nix  # SketchyBar status bar (macOS only)
            ./raycast.nix  # Raycast scripts (macOS only)
          ];

          # ─── macOS-only packages ───
          darwinOnlyPackages = with pkgs; [
            yabai
            skhd
            unstablePkgs.sketchybar  # Use unstable for latest version
          ];
        in
        home-manager.lib.homeManagerConfiguration {
          inherit pkgs;
          
          # Pass extraSpecialArgs to make unstablePkgs available in modules
          extraSpecialArgs = {
            inherit unstablePkgs;
          };
          
          modules = commonModules
            ++ pkgs.lib.optionals isDarwin darwinOnlyModules
            ++ [
            {
              # Personal data
              home.username = username;
              # ─── Cross-platform ─── Conditional home directory
              home.homeDirectory =
                if isDarwin
                then "/Users/${username}"
                else "/home/${username}";
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
                unstablePkgs.neovim
                tree-sitter

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
              ]
              # ─── macOS-only ─── Window management packages
              ++ pkgs.lib.optionals isDarwin darwinOnlyPackages;

              # Enable programs explicitly (critical for binaries to appear)
              # All program enables are centralized here
              programs.neovim.enable = false;
              programs.fish.enable = true;
              programs.nushell.enable = true;
              programs.starship.enable = false;
              programs.zsh.enable = false;  # Managed via home.file in zsh.nix
              programs.git.enable = true;
              programs.gh.enable = true;  # GitHub CLI
              programs.home-manager.enable = true;
              # Note: tmux is configured via home.file in tmux.nix, not programs.tmux

              # NOTE: home.sessionVariables removed - it generates a recursive .zshenv bug
              # XDG_CONFIG_HOME is set in shell configs instead

              # Allow unfree packages
              nixpkgs.config.allowUnfree = true;
            }
          ];
        };
    in
    {
      # Home Manager configurations for each system
      homeConfigurations = {
        # ─── macOS system configurations ───
        "gentleman-macos-intel" = mkHomeConfiguration "x86_64-darwin";
        "gentleman-macos-arm" = mkHomeConfiguration "aarch64-darwin";
        
        # Default to Apple Silicon
        "gentleman" = mkHomeConfiguration "aarch64-darwin";

        # ─── Linux system configurations ───
        "gentleman-linux" = mkHomeConfiguration "x86_64-linux";
        "gentleman-linux-arm" = mkHomeConfiguration "aarch64-linux";
      };
    };
}
