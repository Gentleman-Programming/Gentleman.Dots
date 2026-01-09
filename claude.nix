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
  home.activation.installClaudeConfig = lib.hm.dag.entryAfter [ "linkGeneration" ] ''
    export PATH="${pkgs.coreutils}/bin:${pkgs.jq}/bin:$PATH"

    echo "🔧 Setting up Claude Code..."

    CLAUDE_SRC="${toString ./claude}"
    CLAUDE_DST="$HOME/.claude"

    mkdir -p "$CLAUDE_DST/output-styles"
    mkdir -p "$CLAUDE_DST/skills"

    # Copy CLAUDE.md (global instructions)
    if [ -f "$CLAUDE_SRC/CLAUDE.md" ]; then
      cp -f "$CLAUDE_SRC/CLAUDE.md" "$CLAUDE_DST/"
      echo "⚙️ Copied CLAUDE.md"
    fi

    # Copy statusline script
    if [ -f "$CLAUDE_SRC/statusline.sh" ]; then
      cp -f "$CLAUDE_SRC/statusline.sh" "$CLAUDE_DST/"
      chmod +x "$CLAUDE_DST/statusline.sh"
      echo "📊 Copied statusline.sh"
    fi

    # Copy settings.json (merge with existing to preserve plugins, etc.)
    if [ -f "$CLAUDE_SRC/settings.json" ]; then
      if [ -f "$CLAUDE_DST/settings.json" ]; then
        # Merge: keep existing plugins, override permissions/outputStyle/statusLine
        ${pkgs.jq}/bin/jq -s '.[0] * .[1]' "$CLAUDE_DST/settings.json" "$CLAUDE_SRC/settings.json" > "$CLAUDE_DST/settings.json.tmp"
        mv "$CLAUDE_DST/settings.json.tmp" "$CLAUDE_DST/settings.json"
        echo "⚙️ Merged settings.json"
      else
        cp -f "$CLAUDE_SRC/settings.json" "$CLAUDE_DST/"
        echo "⚙️ Copied settings.json"
      fi
    fi

    # Copy output styles
    if [ -d "$CLAUDE_SRC/output-styles" ]; then
      cp -f "$CLAUDE_SRC/output-styles"/* "$CLAUDE_DST/output-styles/" 2>/dev/null || true
      echo "🎨 Copied output styles"
    fi

    # Copy tweakcc theme (visual colors for Claude Code)
    if [ -f "$CLAUDE_SRC/tweakcc-theme.json" ]; then
      cp -f "$CLAUDE_SRC/tweakcc-theme.json" "$CLAUDE_DST/"
      echo "🎨 Copied tweakcc theme (run 'npx tweakcc --apply' to enable)"
    fi

    # Copy skills
    if [ -d "$CLAUDE_SRC/skills" ]; then
      cp -rf "$CLAUDE_SRC/skills"/* "$CLAUDE_DST/skills/" 2>/dev/null || true
      echo "🧠 Copied skills"
    fi

    # Merge MCP servers into ~/.claude.json (the actual config file)
    CLAUDE_JSON="$HOME/.claude.json"
    if [ -f "$CLAUDE_SRC/mcp-servers.template.json" ]; then
      # Keep template as reference for other servers (Jira, Figma need manual tokens)
      cp -f "$CLAUDE_SRC/mcp-servers.template.json" "$CLAUDE_DST/"

      if [ -f "$CLAUDE_JSON" ]; then
        # Merge only context7 into existing ~/.claude.json (safe - no tokens needed)
        ${pkgs.jq}/bin/jq --argjson ctx7 '{"context7":{"type":"http","url":"https://mcp.context7.com/mcp"}}' \
          '.mcpServers = (.mcpServers // {}) + $ctx7' "$CLAUDE_JSON" > "$CLAUDE_JSON.tmp"
        mv "$CLAUDE_JSON.tmp" "$CLAUDE_JSON"
        echo "📡 Merged context7 MCP server into ~/.claude.json"
      else
        # Create new ~/.claude.json with context7
        echo '{"mcpServers":{"context7":{"type":"http","url":"https://mcp.context7.com/mcp"}}}' > "$CLAUDE_JSON"
        echo "📡 Created ~/.claude.json with context7 MCP server"
      fi
      echo "💡 Other MCP servers (Jira, Figma) need tokens - see ~/.claude/mcp-servers.template.json"
    fi

    echo ""
    echo "🎉 Claude Code setup complete!"
    echo "Usage: cc | claude-config"
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
