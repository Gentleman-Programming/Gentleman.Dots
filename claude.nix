{ pkgs, lib, ... }:

{
  # Claude Code CLI via Bun global install.
  # Implementation details:
  # - Runs during Home Manager activation (imperative step) to install/update outside Nix store.
  # - Non-reproducible (stateful) but aligned with the requested behavior.
  # - Requested command: `bun install -g @anthropic-ai/claude-code` (current Bun docs generally use `bun add -g`).
  # - Falls back to `bun add -g` if `bun install -g` fails (covers both syntaxes across Bun versions).

  # Ensure Bun's global bin directory is on PATH for all sessions
  home.sessionPath = [ "$HOME/.bun/bin" ];

  # Optional environment variables (extend if Claude CLI needs API keys, etc.)
  home.sessionVariables = { };

  # Activation script runs after links are generated
  home.activation.installClaudeCode = lib.hm.dag.entryAfter [ "linkGeneration" ] ''
    set -e
    echo "[claude-code] Checking Claude Code CLI installation via Bun" >&2

    # Absolute Bun path from Nix (preferred if present in home.packages)
    BUN="${pkgs.bun}/bin/bun"

    if [ -x "$BUN" ]; then
      : # OK
    elif command -v bun >/dev/null 2>&1; then
      BUN="$(command -v bun)"
    else
      echo "[claude-code] ❌ Bun is not available (PATH=$PATH). Add 'bun' to home.packages first." >&2
      exit 0
    fi

    echo "[claude-code] Using bun: $BUN" >&2

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
      echo "[claude-code] Installing @anthropic-ai/claude-code globally..." >&2
      # First attempt with the requested syntax
      if "$BUN" install -g @anthropic-ai/claude-code 2>/dev/null; then
        echo "[claude-code] Installation succeeded (install -g)." >&2
      else
        echo "[claude-code] 'bun install -g' failed, retrying with 'bun add -g'..." >&2
        if "$BUN" add -g @anthropic-ai/claude-code; then
          echo "[claude-code] Installation succeeded (add -g)." >&2
        else
          echo "[claude-code] ❌ Global installation of @anthropic-ai/claude-code failed" >&2
        fi
      fi
    fi
  '';

  # Fish aliases (cc -> assumed 'claude-code' binary). Best-effort; existence is ensured only after activation.
  programs.fish.shellAliases = {
    cc = "claude-code";  # Adjust if the final executable name differs
  };
}
