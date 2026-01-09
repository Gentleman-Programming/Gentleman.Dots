{ pkgs, lib, ... }:

{
  # Claude Code configuration files
  # Note: CLI installation is handled separately (brew, official installer, etc.)
  # This module only manages config files: CLAUDE.md, settings.json, statusline, skills, output-styles

  # Required packages for statusline
  home.packages = [
    pkgs.jq      # For statusline JSON parsing
  ];

  # Activation script to copy config files
  home.activation.setupClaudeConfig = lib.hm.dag.entryAfter [ "linkGeneration" ] ''
    set -e
    export PATH="${pkgs.coreutils}/bin:${pkgs.jq}/bin:$PATH"

    echo "[claude-code] Setting up Claude Code config..." >&2

    CLAUDE_SRC="${toString ./claude}"
    CLAUDE_DST="$HOME/.claude"

    mkdir -p "$CLAUDE_DST/output-styles"
    mkdir -p "$CLAUDE_DST/skills"

    # Copy CLAUDE.md (global instructions)
    if [ -f "$CLAUDE_SRC/CLAUDE.md" ]; then
      cp -f "$CLAUDE_SRC/CLAUDE.md" "$CLAUDE_DST/"
      echo "[claude-code] ⚙️ Copied CLAUDE.md" >&2
    fi

    # Copy statusline script
    if [ -f "$CLAUDE_SRC/statusline.sh" ]; then
      cp -f "$CLAUDE_SRC/statusline.sh" "$CLAUDE_DST/"
      chmod +x "$CLAUDE_DST/statusline.sh"
      echo "[claude-code] 📊 Copied statusline.sh" >&2
    fi

    # Copy settings.json (merge with existing to preserve plugins, etc.)
    if [ -f "$CLAUDE_SRC/settings.json" ]; then
      if [ -f "$CLAUDE_DST/settings.json" ]; then
        # Merge: keep existing plugins, override permissions/outputStyle/statusLine
        ${pkgs.jq}/bin/jq -s '.[0] * .[1]' "$CLAUDE_DST/settings.json" "$CLAUDE_SRC/settings.json" > "$CLAUDE_DST/settings.json.tmp"
        mv "$CLAUDE_DST/settings.json.tmp" "$CLAUDE_DST/settings.json"
        echo "[claude-code] ⚙️ Merged settings.json" >&2
      else
        cp -f "$CLAUDE_SRC/settings.json" "$CLAUDE_DST/"
        echo "[claude-code] ⚙️ Copied settings.json" >&2
      fi
    fi

    # Copy output styles
    if [ -d "$CLAUDE_SRC/output-styles" ]; then
      cp -f "$CLAUDE_SRC/output-styles"/* "$CLAUDE_DST/output-styles/" 2>/dev/null || true
      echo "[claude-code] 🎨 Copied output styles" >&2
    fi

    # Copy skills
    if [ -d "$CLAUDE_SRC/skills" ]; then
      cp -rf "$CLAUDE_SRC/skills"/* "$CLAUDE_DST/skills/" 2>/dev/null || true
      echo "[claude-code] 🧠 Copied skills" >&2
    fi

    # Show MCP template info
    if [ -f "$CLAUDE_SRC/mcp-servers.template.json" ]; then
      cp -f "$CLAUDE_SRC/mcp-servers.template.json" "$CLAUDE_DST/"
      echo "[claude-code] 📡 MCP template copied to $CLAUDE_DST/mcp-servers.template.json" >&2
      echo "[claude-code] 💡 To add MCP servers, run: claude mcp add <name> or edit ~/.claude.json" >&2
    fi

    echo "[claude-code] 🎉 Claude Code config setup complete!" >&2
  '';

  # Shell aliases
  programs.fish.shellAliases = {
    cc = "claude";
    claude-config = "nvim ~/.claude/settings.json";
  };

  programs.zsh.shellAliases = {
    cc = "claude";
    claude-config = "nvim ~/.claude/settings.json";
  };
}
