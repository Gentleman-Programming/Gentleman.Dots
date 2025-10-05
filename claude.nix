{ pkgs, lib, ... }:

{
  # Claude Code CLI via official install script.
  # Implementation details:
  # - Runs during Home Manager activation (imperative step) to install/update outside Nix store.
  # - Non-reproducible (stateful) but aligned with the official installation method.
  # - Uses: curl -fsSL https://claude.ai/install.sh | bash

  # Optional environment variables (extend if Claude CLI needs API keys, etc.)
  home.sessionVariables = { };

  # Activation script runs after links are generated
  home.activation.installClaudeCode = lib.hm.dag.entryAfter [ "linkGeneration" ] ''
    set -e
    # Ensure curl and other necessary tools are in the PATH for the installation script
    export PATH="${pkgs.coreutils}/bin:${pkgs.curl}/bin:${pkgs.gnugrep}/bin:${pkgs.gnused}/bin:${pkgs.perl}/bin:${pkgs.jq}/bin:$PATH"

    echo "[claude-code] Checking Claude Code CLI installation" >&2

    # Detect existing binary candidates
    FOUND=""
    for cand in claude-code claude anthropic-claude; do
      if command -v "$cand" >/dev/null 2>&1; then
        FOUND="$cand"; break
      fi
    done

    if [ -n "$FOUND" ]; then
      echo "[claude-code] ✅ CLI already present: $FOUND" >&2
    else
      echo "[claude-code] Installing Claude Code CLI via official installer..." >&2
      if curl -fsSL https://claude.ai/install.sh | ${pkgs.bash}/bin/bash; then
        echo "[claude-code] ✅ Installation succeeded" >&2
      else
        echo "[claude-code] ❌ Installation of Claude Code CLI failed" >&2
      fi
    fi
  '';

  # Fish aliases (cc -> assumed 'claude-code' binary). Best-effort; existence is ensured only after activation.
  programs.fish.shellAliases = {
    cc = "claude-code";  # Adjust if the final executable name differs
  };
}
