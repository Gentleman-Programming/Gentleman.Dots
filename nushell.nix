{ pkgs, ... }:
let
  systemType = if pkgs.stdenv.isDarwin then "mac" else "linux";
in {
  programs.nushell.enable = true;
  
  home.file = 
    if systemType == "mac" then {
      "Library/Application Support/nushell/config.nu" = { source = ./nushell/config.nu; };
      "Library/Application Support/nushell/env.nu" = { source = ./nushell/env.nu; };
    } else {
      ".config/nushell/config.nu" = { source = ./nushell/config.nu; };
      ".config/nushell/env.nu" = { source = ./nushell/env.nu; };
    };
}
