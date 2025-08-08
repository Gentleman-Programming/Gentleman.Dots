{ pkgs, lib, ... }:

{
  # Claude Code CLI configuration
  # Note: Claude Code must be installed via Homebrew:
  # brew install --cask claude-code
  
  # Add any Claude Code related configuration here
  # For example, you could add shell aliases or environment variables
  home.sessionVariables = {
    # Add any Claude Code related environment variables here if needed
  };
  
  programs.fish.shellAliases = {
    # Add any Claude Code related aliases here if needed
    # cc = "claude-code";
  };
}