{
  description = "Gentleman.Dots: Complete configuration using local modules";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    home-manager = {
      url = "https://github.com/nix-community/home-manager/archive/master.tar.gz";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, home-manager, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        lib = pkgs.lib;

        # Base options used by default
        baseOptions = {
          terminal = "wezterm";
          shell = "nushell";
          windowManager = "zellij";
          installNeovim = true;
          osType = if system == "x86_64-darwin" then "mac" else "linux";
          starship = true;
          powerlevel10k = false;
          useTmux = false;
        };

        # Extra module to install dependencies and perform replacements
        hmExtraModule = { config, pkgs, ... }: {
          home.packages = lib.concatLists [
            [ pkgs.git pkgs.curl pkgs.rustc pkgs.cargo ]
            (if config.gentlemanOptions.terminal == "alacritty" then [ pkgs.alacritty ] else [])
            (if config.gentlemanOptions.terminal == "wezterm" then [ pkgs.wezterm ] else [])
            (if config.gentlemanOptions.terminal == "kitty" then [ pkgs.kitty ] else [])
            (if config.gentlemanOptions.shell == "zsh" then [ pkgs.zsh pkgs.volta pkgs.carapace pkgs.zoxide pkgs.atuin pkgs.fzf pkgs.zsh-autosuggestions pkgs.zsh-syntax-highlighting pkgs.zsh-autocomplete ] else [])
            (if config.gentlemanOptions.shell == "fish" then [ pkgs.fish pkgs.volta pkgs.carapace pkgs.zoxide pkgs.atuin pkgs.starship pkgs.fzf ] else [])
            (if config.gentlemanOptions.shell == "nushell" then [ pkgs.nushell pkgs.volta pkgs.carapace pkgs.zoxide pkgs.atuin pkgs.jq pkgs.bash pkgs.starship pkgs.fzf ] else [])
            (if (config.gentlemanOptions.windowManager == "tmux") || config.gentlemanOptions.useTmux then [ pkgs.tmux ] else [])
            (if config.gentlemanOptions.windowManager == "zellij" then [ pkgs.zellij ] else [])
            (if config.gentlemanOptions.installNeovim then [ pkgs.neovim pkgs.nodejs pkgs.npm pkgs.gcc pkgs.fzf pkgs.fd pkgs.ripgrep pkgs.coreutils pkgs.bat pkgs.lazygit ] else [])
          ];

          home.activation.replaceZshConfig = lib.mkAfter ''
            if [ "${toString config.gentlemanOptions.shell}" = "zsh" ] && [ "${toString config.gentlemanOptions.windowManager}" = "zellij" ]; then
              echo "Replacing placeholders in .zshrc for Zellij..."
              sed -i.bak 's/TMUX/WM_VAR="\/$ZELLIJ"/g; s/tmux/WM_CMD="zellij"/g' ~/.zshrc
            fi
          '';

          home.activation.replaceFishConfig = lib.mkAfter ''
            if [ "${toString config.gentlemanOptions.shell}" = "fish" ] && [ "${toString config.gentlemanOptions.windowManager}" = "zellij" ]; then
              echo "Replacing placeholders in config.fish for Zellij..."
              sed -i.bak 's/TMUX/if not set -q ZELLIJ/g; s/tmux/zellij/g' ~/.config/fish/config.fish
            fi
          '';

          home.activation.replaceNushellConfig = lib.mkAfter ''
            if [ "${toString config.gentlemanOptions.shell}" = "nushell" ] && [ "${toString config.gentlemanOptions.windowManager}" = "zellij" ]; then
              echo "Replacing placeholders in config.nu for Zellij..."
              sed -i.bak 's/"tmux"/let MULTIPLEXER = "zellij"/g; s/"TMUX"/let MULTIPLEXER_ENV_PREFIX = "ZELLIJ"/g' ~/.config/nushell/config.nu
            fi
          '';
        };

        hmConfig = home-manager.lib.homeManagerConfiguration { inherit pkgs; };
      in {
        # Preset 1: Zellij with Nushell and Starship
        homeConfigurations."zellij-nushell-starship" = hmConfig {
          modules = [
            ({ config, ... }: {
              gentlemanOptions = baseOptions // {
                shell = "nushell";
                windowManager = "zellij";
                starship = true;
                powerlevel10k = false;
                useTmux = false;
              };
            })
            hmExtraModule
          ];
          home.username = "Gentleman";
        };

        # Preset 2: Zellij with Fish and Starship
        homeConfigurations."zellij-fish-starship" = hmConfig {
          modules = [
            ({ config, ... }: {
              gentlemanOptions = baseOptions // {
                shell = "fish";
                windowManager = "zellij";
                starship = true;
                powerlevel10k = false;
                useTmux = false;
              };
            })
            hmExtraModule
          ];
          home.username = "Gentleman";
        };

        # Preset 3: Zellij with Zsh and Powerlevel10k
        homeConfigurations."zellij-zsh-power10k" = hmConfig {
          modules = [
            ({ config, ... }: {
              gentlemanOptions = baseOptions // {
                shell = "zsh";
                windowManager = "zellij";
                starship = false;
                powerlevel10k = true;
                useTmux = false;
              };
            })
            hmExtraModule
          ];
          home.username = "Gentleman";
        };

        # Preset 4: Tmux with Nushell and Starship
        homeConfigurations."tmux-nushell-starship" = hmConfig {
          modules = [
            ({ config, ... }: {
              gentlemanOptions = baseOptions // {
                shell = "nushell";
                windowManager = "tmux";
                starship = true;
                powerlevel10k = false;
                useTmux = true;
              };
            })
            hmExtraModule
          ];
          home.username = "Gentleman";
        };

        # Preset 5: Tmux with Fish and Starship
        homeConfigurations."tmux-fish-starship" = hmConfig {
          modules = [
            ({ config, ... }: {
              gentlemanOptions = baseOptions // {
                shell = "fish";
                windowManager = "tmux";
                starship = true;
                powerlevel10k = false;
                useTmux = true;
              };
            })
            hmExtraModule
          ];
          home.username = "Gentleman";
        };

        # Preset 6: Tmux with Zsh and Powerlevel10k
        homeConfigurations."tmux-zsh-power10k" = hmConfig {
          modules = [
            ({ config, ... }: {
              gentlemanOptions = baseOptions // {
                shell = "zsh";
                windowManager = "tmux";
                starship = false;
                powerlevel10k = true;
                useTmux = true;
              };
            })
            hmExtraModule
          ];
          home.username = "Gentleman";
        };
      }
    );
}
