{ pkgs, lib, ... }:
let
  systemType = if pkgs.stdenv.isDarwin then "mac" else "linux";
in {
  # Use programs.nushell instead of home.file to avoid conflicts
  programs.nushell = {
    enable = true;
    configFile.source = ./nushell/config.nu;
    envFile.source = ./nushell/env.nu;
  };

  # Additional bash environment files
  home.file = {
    ".config/bash-env-json" = { source = ./nushell/bash-env-json; };
    ".config/bash-env.nu" = { source = ./nushell/bash-env.nu; };
  };
}
