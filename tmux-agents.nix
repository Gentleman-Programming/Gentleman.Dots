{ pkgs, lib, ... }:

{
  # tmux agent-state notifier — surfaces AI agent state (working / blocked / idle)
  # into tmux: per-tab dot, status-right informativo, and audible+visual alerts.
  #
  # Architecture (anti-corruption layer): every agent adapter normalizes its native
  # events into one canonical vocabulary and calls a single core script. tmux only
  # ever sees the core; it never knows which agent ran.
  #
  #   adapters (opencode / pi / claude / codex)  ->  agent-report.sh  ->  tmux
  #
  # Sources live in ./tmux-agents and are copied into place on activation, mirroring
  # the opencode.nix / claude.nix convention.

  home.packages = [ pkgs.jq ];

  home.activation.installTmuxAgents = lib.hm.dag.entryAfter [ "linkGeneration" ] ''
    export PATH="${pkgs.jq}/bin:${pkgs.coreutils}/bin:$PATH"

    echo "🔧 Setting up tmux agent-state notifier..."

    SRC="${toString ./tmux-agents}"

    # --- core scripts + tmux display layer ---
    mkdir -p "$HOME/.config/tmux/scripts"
    cp -f "$SRC/scripts/"*.sh "$HOME/.config/tmux/scripts/" 2>/dev/null || true
    chmod +x "$HOME/.config/tmux/scripts/"*.sh 2>/dev/null || true
    cp -f "$SRC/agents.conf" "$HOME/.config/tmux/agents.conf" 2>/dev/null || true
    echo "  ⚙️  scripts + agents.conf installed"

    # --- opencode adapter (plugin) ---
    mkdir -p "$HOME/.config/opencode/plugins"
    cp -f "$SRC/opencode/tmux-agent-state.js" "$HOME/.config/opencode/plugins/" 2>/dev/null || true

    # --- pi adapter (extension) ---
    mkdir -p "$HOME/.pi/agent/extensions"
    cp -f "$SRC/pi/tmux-agent-state.ts" "$HOME/.pi/agent/extensions/" 2>/dev/null || true
    echo "  🔌 opencode + pi adapters installed"

    # --- claude / codex adapters: idempotent hook merge into their config files ---
    # Both use a claude-style hooks schema. We APPEND our hook (never replace arrays),
    # only if not already present, so we never clobber herdr/gentle hooks.
    HOOK="bash $HOME/.config/tmux/scripts/hook-adapter.sh"

    merge_hooks() {
      # $1 = target json file, remaining args = "Event:matcher" pairs
      local file="$1"; shift
      [ -f "$file" ] || echo '{"hooks":{}}' > "$file"
      local filter='.hooks = (.hooks // {})'
      local pair ev matcher
      for pair in "$@"; do
        ev="''${pair%%:*}"; matcher="''${pair#*:}"
        filter="$filter | addhook(\"$ev\"; \"$matcher\")"
      done
      jq --arg base "$HOOK" "
        def addhook(\$ev; \$matcher):
          (\$base + \" \" + \$ev) as \$cmd
          | .hooks[\$ev] = ((.hooks[\$ev] // [])
            | if any(.[].hooks[]?; .command == \$cmd) then .
              else . + [{matcher: \$matcher, hooks: [{type: \"command\", command: \$cmd, timeout: 5}]}] end);
        $filter
      " "$file" > "$file.tmp" && mv "$file.tmp" "$file"
    }

    CLAUDE_SETTINGS="$HOME/.claude/settings.json"
    if [ -f "$CLAUDE_SETTINGS" ]; then
      merge_hooks "$CLAUDE_SETTINGS" "UserPromptSubmit:" "PreToolUse:*" "Notification:" "Stop:"
      echo "  🪝 claude hooks merged"
    else
      echo "  ⏭️  claude settings.json not found yet — run Claude once, then re-switch"
    fi

    CODEX_HOOKS="$HOME/.codex/hooks.json"
    if [ -f "$CODEX_HOOKS" ]; then
      merge_hooks "$CODEX_HOOKS" "UserPromptSubmit:" "PreToolUse:*" "Notification:" "Stop:" "TurnComplete:"
      echo "  🪝 codex hooks merged (approve the new hooks on next codex start)"
    else
      echo "  ⏭️  codex hooks.json not found yet — run Codex once, then re-switch"
    fi

    echo "🎉 tmux agent-state notifier ready (open a fresh tmux + restart your agents)"
  '';
}
