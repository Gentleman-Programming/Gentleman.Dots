{ pkgs, lib, ... }: 
{
  programs.fish = {
    enable = true;
    
    interactiveShellInit = ''
      # Set Catppuccin Mocha as default theme on startup
      if status is-interactive
        fish_config theme choose "Catppuccin Mocha"
      end
    '';
    
    plugins = [
      {
        name = "fisher";
        src = pkgs.fetchFromGitHub {
          owner = "jorgebucaran";
          repo = "fisher";
          rev = "4.4.4";
          sha256 = "sha256-k8aBgZuKPB784qa9vZJe1E8bLqXXGFRi6xfRiR3yJ5c=";
        };
      }
      {
        name = "catppuccin";
        src = pkgs.fetchFromGitHub {
          owner = "catppuccin";
          repo = "fish";
          rev = "0ce27b518e8ead555dec34dd8be3df5bd75cff8e";
          sha256 = "sha256-Dc/zdxfzAUM5NX8PxzfljRbYvO9f9syuLJ4WuCGgzQE=";
        };
      }
    ];
  };

  home.activation.copyFish = lib.hm.dag.entryAfter ["writeBoundary"] ''
    echo "Copying fish configuration..."
    rm -rf "$HOME/.config/fish"

    cp -r ${toString ./fish} "$HOME/.config/fish"
    chmod -R u+w "$HOME/.config/fish"
    
    # Ensure directories exist for Fisher
    mkdir -p "$HOME/.config/fish/fisher/functions"
    mkdir -p "$HOME/.config/fish/fisher/completions"
    mkdir -p "$HOME/.config/fish/fisher/conf.d"
  '';
}
