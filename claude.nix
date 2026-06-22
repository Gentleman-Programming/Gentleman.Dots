{ pkgs, lib, ... }:

{
  # Required packages for statusline
  home.packages = [
    pkgs.jq
  ];

  # Merge MCP servers into ~/.claude.json on activation
  home.activation.installClaudeConfig = lib.hm.dag.entryAfter [ "linkGeneration" ] ''
    export PATH="${pkgs.coreutils}/bin:${pkgs.jq}/bin:$PATH"

    CLAUDE_JSON="$HOME/.claude.json"

    if [ -f "$CLAUDE_JSON" ]; then
      ${pkgs.jq}/bin/jq --argjson servers '{
        "context7":{"type":"http","url":"https://mcp.context7.com/mcp"},
        "engram":{"type":"stdio","command":"engram","args":["mcp"]},
        "notion":{"type":"http","url":"https://mcp.notion.com/mcp"}
      }' '.mcpServers = (.mcpServers // {}) + $servers' "$CLAUDE_JSON" > "$CLAUDE_JSON.tmp"
      mv "$CLAUDE_JSON.tmp" "$CLAUDE_JSON"
    else
      echo '{"mcpServers":{"context7":{"type":"http","url":"https://mcp.context7.com/mcp"},"engram":{"type":"stdio","command":"engram","args":["mcp"]},"notion":{"type":"http","url":"https://mcp.notion.com/mcp"}}}' > "$CLAUDE_JSON"
    fi
  '';

  programs.fish.shellAliases = {
    cc = "claude";
    claude-config = "nvim ~/.claude/settings.json";
  };

  programs.zsh.shellAliases = {
    cc = "claude";
    claude-config = "nvim ~/.claude/settings.json";
  };
}
