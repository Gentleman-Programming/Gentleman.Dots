{ lib, pkgs, ... }:

let
  script = builtins.readFile ./scripts/omniwm/activate.sh;
  runScript = args: ''
    ${pkgs.bash}/bin/bash -s -- ${args} <<'OMNIWM_ACTIVATION'
${script}
OMNIWM_ACTIVATION
  '';
in
{
  # Reject unsupported hosts before Home Manager's write boundary.
  home.activation.checkOmniWM = lib.hm.dag.entryBefore [ "writeBoundary" ]
    (runScript "--check-platform");

  # Do not manage the GUI-owned settings or register an extra launch agent.
  home.activation.setupOmniWM = lib.hm.dag.entryAfter [ "linkGeneration" ] ''
    if [ -z "''${DRY_RUN_CMD:-}" ]; then
      ${runScript ""}
    fi
  '';
}
