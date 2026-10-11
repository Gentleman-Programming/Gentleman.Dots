#!/usr/bin/env bash
set -euo pipefail
REPO=${1:-"$(cd "$(dirname "$0")/../.." && pwd)"}
FISH_BIN=$(command -v fish)
TEMP_ROOT=$(mktemp -d "${TMPDIR:-/tmp}/legacy-fish-ls-tests.XXXXXX")
trap 'rm -rf "$TEMP_ROOT"' EXIT
mkdir -p "$TEMP_ROOT/home"
export HOME="$TEMP_ROOT/home" XDG_CONFIG_HOME="$TEMP_ROOT/home/.config"
# The Termux startup branch avoids executing the host's Homebrew shellenv.
# Alias platform selection still uses the mocked uname below.
export LS_COLORS=stale LSCOLORS=stale TERMUX_VERSION=fixture PREFIX="$TEMP_ROOT"

# Source the real configuration, isolating unrelated integrations and platform probes.
cat > "$TEMP_ROOT/probe.fish" <<'FISH'
set -g fixture_os $argv[2]
function uname
    echo $fixture_os
end
for tool in starship zoxide atuin fzf carapace clear
    function $tool
        return 0
    end
end
source "$argv[1]/GentlemanFish/fish/config.fish"; or exit 1
command sh -c 'printf "%s\n%s\n" "$LS_COLORS" "$LSCOLORS"'
string join ',' -- $fish_color_command
functions ls
FISH

GNU_COLORS='di=38;2;122;168;159:ln=38;2;255;141;215:ex=38;2;183;204;133:pi=38;2;255;224;102:so=38;2;163;181;214:bd=38;2;222;186;135:cd=38;2;222;186;135:or=38;2;203;124;148:mi=38;2;203;124;148'
for os in Darwin Linux
 do
    output=$("$FISH_BIN" --no-config "$TEMP_ROOT/probe.fish" "$REPO" "$os")
    expected=$(printf '%s\n%s' "$GNU_COLORS" 'gxfxexdxcxdxdxabagacgc')
    if [ "$(printf '%s\n' "$output" | head -n 2)" != "$expected" ]; then
        printf 'assertion failed: %s must export the Gentleman listing palette\n%s\n' "$os" "$output" >&2
        exit 1
    fi
    printf '%s\n' "$output" | grep -Fxq '7AA89F,cyan'
    if [ "$os" = Darwin ]; then
        printf '%s\n' "$output" | grep -Fq 'ls -G $argv'
    else
        printf '%s\n' "$output" | grep -Fq 'ls --color=auto $argv'
    fi
 done
printf 'legacy Fish ls theme tests passed\n'
