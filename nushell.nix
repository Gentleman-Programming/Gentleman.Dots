{ pkgs, lib, ... }:
let
  systemType = if pkgs.stdenv.isDarwin then "mac" else "linux";
in {
  home.file = 
    if systemType == "mac" then { 
      "Library/Application Support/nushell/config.nu" = { source = ./nushell/config.nu; };
      "Library/Application Support/nushell/env.nu" = { source = ./nushell/env.nu; };
      ".config/bash-env-json" = { source = ./nushell/bash-env-json; };
      ".config/bash-env.nu" = { source = ./nushell/bash-env.nu; };
    } else {
      ".config/nushell/config.nu" = { source = ./nushell/config.nu; };
      ".config/nushell/env.nu" = { source = ./nushell/env.nu; };
      ".config/bash-env-json" = { source = ./nushell/bash-env-json; };
      ".config/bash-env.nu" = { source = ./nushell/bash-env.nu; };
    };
}
