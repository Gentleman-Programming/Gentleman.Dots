# OpenCode Configuration for Gentleman.Dots

## Build/Lint/Test Commands
- **Apply configuration:** `nix run github:nix-community/home-manager -- switch --flake .#gentleman -b backup`
- **Build flake:** `nix build .#homeConfigurations.gentleman.activationPackage`
- **Check flake:** `nix flake check`
- **Update dependencies:** `nix flake update`
- **Lua formatting:** `stylua nvim/lua/` (follows stylua.toml: 2 spaces, 120 columns)
- **Test Neovim config:** `nvim --headless -c "checkhealth" -c "qa"`

## Code Style Guidelines
- **Nix files:** Use 2-space indentation, snake_case for variables, kebab-case for packages
- **Lua (Neovim):** 2-space indentation, camelCase for functions, snake_case for variables
- **Configuration files:** Follow existing patterns - declarative modules in .nix files
- **Comments:** Minimal in code, extensive in configuration explanations
- **Imports:** Group by source (builtin → external → local), sort alphabetically within groups
- **Error handling:** Use safe defaults, graceful fallbacks for cross-platform compatibility
- **File structure:** One concern per .nix module, mirror target system paths when possible
- **Platform detection:** Use `pkgs.stdenv.isDarwin` for macOS-specific config
- **Dependencies:** Declare all tools in flake.nix packages, avoid external downloads in configs
- **Theming:** Consistent color schemes across all tools (Catppuccin Mocha, Kanagawa variants)