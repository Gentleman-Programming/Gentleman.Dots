# Arch Package Install Specification

## Purpose

Defines which packages the installer's Arch/pacman path installs, and requires that no unresolvable package name can abort the entire `pacman -S` transaction and silently skip unrelated packages (zoxide, atuin, starship, neovim, zellij).

## Requirements

### Requirement: Arch package lists MUST only contain names resolvable in the official Arch repositories

The `Arch:` package list for each shell profile (fish, zsh, nushell) MUST NOT include a package name that is unavailable in the official Arch repositories (`core`, `extra`) and requires AUR. Per the official Arch package database, `carapace` and `zsh-theme-powerlevel10k` are AUR-only and MUST be removed from every list that currently contains them.

#### Scenario: Fish profile Arch list has no AUR-only names

- GIVEN the fish shell profile's Arch package list
- WHEN the list is inspected
- THEN it MUST NOT contain `carapace`
- AND every remaining name MUST resolve via `pacman -Si <name>` on a vanilla Arch install

#### Scenario: Zsh profile Arch list has no AUR-only names

- GIVEN the zsh shell profile's Arch package list
- WHEN the list is inspected
- THEN it MUST NOT contain `carapace` or `zsh-theme-powerlevel10k`
- AND the remaining names MUST be exactly `zsh`, `zoxide`, `atuin`, `zsh-autosuggestions`, `zsh-syntax-highlighting`, `zsh-autocomplete`, each of which resolves on vanilla Arch

#### Scenario: Nushell profile Arch list has no AUR-only names

- GIVEN the nushell shell profile's Arch package list
- WHEN the list is inspected
- THEN it MUST NOT contain `carapace`
- AND every remaining name MUST resolve on vanilla Arch

### Requirement: A single unresolvable pacman package name MUST NOT block installation of the remaining packages

Because `pacman -S` is atomic, the installer's Arch package installation step MUST be resilient to a package name that pacman cannot resolve: it MUST NOT let one bad name prevent installation of otherwise-installable packages in the same profile, and MUST NOT silently report success when packages were skipped.

#### Scenario: All listed packages resolve (nominal path)

- GIVEN the corrected Arch package list for a shell profile
- WHEN the installer runs the Arch install step on a vanilla Arch host with no extra repos
- THEN `pacman -S --needed --noconfirm` MUST complete successfully
- AND zoxide, atuin, starship, neovim, and zellij MUST be installed

#### Scenario: Homebrew fallback is not required for the corrected list

- GIVEN the corrected Arch package list
- WHEN the installer runs on a host where `hasBrew` is false
- THEN the install step MUST still succeed without falling back to Homebrew

## Non-Goals

- AUR fallback (`carapace-bin` via `paru`/`yay`) is explicitly out of scope; losing the `carapace` completion feature on Arch is an accepted, documented tradeoff. The install summary SHOULD note the omission.
