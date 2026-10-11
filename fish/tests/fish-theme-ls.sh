#!/usr/bin/env bash
set -euo pipefail

REPO=${1:-"$(cd "$(dirname "$0")/../.." && pwd)"}
FISH_BIN=$(command -v fish)
TEMP_ROOT=$(mktemp -d "${TMPDIR:-/tmp}/fish-theme-ls-tests.XXXXXX")
trap 'rm -rf "$TEMP_ROOT"' EXIT

# No installed config or universal variables may influence these probes.
mkdir -p "$TEMP_ROOT/home"
export HOME="$TEMP_ROOT/home" XDG_CONFIG_HOME="$TEMP_ROOT/home/.config"
export LS_COLORS=stale LSCOLORS=stale

assert_eq() {
  if [ "$1" != "$2" ]; then
    printf 'assertion failed: %s\nexpected: %s\nactual: %s\n' "$3" "$1" "$2" >&2
    exit 1
  fi
}

GENTLEMAN='di=38;2;122;168;159:ln=38;2;255;141;215:ex=38;2;183;204;133:pi=38;2;255;224;102:so=38;2;163;181;214:bd=38;2;222;186;135:cd=38;2;222;186;135:or=38;2;203;124;148:mi=38;2;203;124;148'
CUTE='di=38;2;240;149;200:ln=38;2;169;199;238:ex=38;2;180;231;199:pi=38;2;224;194;122:so=38;2;201;106;162:bd=38;2;242;184;109:cd=38;2;242;184;109:or=38;2;255;113;143:mi=38;2;255;113;143'

# Child processes must receive the active palette, including after switching back.
# shellcheck disable=SC2016
actual=$("$FISH_BIN" --no-config -c '
    for profile in gentleman gentleman-cute gentleman
        source "$argv[1]/fish/themes/$profile.fish"; or exit
        command sh -c '\''printf "%s\n%s\n" "$LS_COLORS" "$LSCOLORS"'\''; or exit
    end
' "$REPO")
expected=$(printf '%s\n' "$GENTLEMAN" 'gxfxexdxcxdxdxabagacgc' "$CUTE" 'FxExfxdxcxdxdxabagacfc' "$GENTLEMAN" 'gxfxexdxcxdxdxabagacgc')
assert_eq "$expected" "$actual" 'Both listing palettes are exported and replaced on every profile activation'

# Existing syntax and fzf colors must not change when adding listing colors.
# shellcheck disable=SC2016
actual=$("$FISH_BIN" --no-config -c '
    source "$argv[1]/fish/themes/gentleman.fish"
    string join "," -- $fish_color_command
    string match -q "*prompt:#7AA89F,*" -- $FZF_DEFAULT_OPTS; or exit 1
    source "$argv[1]/fish/themes/gentleman-cute.fish"
    string join "," -- $fish_color_command
    string match -q "*prompt:#F095C8,*" -- $FZF_DEFAULT_OPTS; or exit 1
' "$REPO")
assert_eq $'7AA89F,cyan\n--bold,#A9C7EE' "$actual" 'Syntax and fzf retain their previous theme colors'

printf 'fish theme ls tests passed\n'
