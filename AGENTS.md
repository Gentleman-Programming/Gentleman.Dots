# Agents Guide for Gentleman.Dots

## Build/Test Commands

- `nix run github:nix-community/home-manager -- switch --flake .#gentleman -b backup` - Apply configuration
- No traditional test commands - this is a dotfiles configuration repository
- Validate Nix syntax: `nix flake check`

## Repository Structure

This is a **dotfiles repository** using Nix flakes for declarative system configuration management. The main entry point is `flake.nix` which orchestrates all configurations.

## Code Style Guidelines

### Nix Files (.nix)

- Use kebab-case for file names (`nushell.nix`, `fish.nix`)
- 2-space indentation
- Use `lib.mkIf` for conditional logic
- Platform detection with `pkgs.stdenv.isDarwin`
- File paths should use `./relative/path` syntax
- Always include proper error handling for cross-platform compatibility

### Lua Files (Neovim config)

- Use snake_case for variables and functions
- 2-space indentation
- Prefer `require()` over `vim.cmd()`
- Comment complex configurations
- Use lazy loading for plugins

### Configuration Files

- Follow original tool conventions (Fish shell, Nushell, etc.)
- Use appropriate comment syntax for each language
- Maintain existing color schemes and themes
- Keep platform-specific configurations clearly separated

### Error Handling

- Always check for tool availability before configuration
- Use conditional blocks for OS-specific features
- Provide fallbacks for missing dependencies

