#!/usr/bin/env bash

set -euo pipefail

REPO=${1:-"$(cd "$(dirname "$0")/../.." && pwd)"}
FUNCTION="$REPO/fish/functions/fish-theme.fish"
FISH_BIN=$(command -v fish)
TEMP_ROOT=$(mktemp -d "${TMPDIR:-/tmp}/fish-theme-atuin-tests.XXXXXX")
trap 'rm -rf "$TEMP_ROOT"' EXIT

HOME_DIR="$TEMP_ROOT/home"
CONFIG_DIR="$HOME_DIR/.config"
BIN_DIR="$TEMP_ROOT/bin"
ATUIN_LOG="$TEMP_ROOT/atuin.log"

mkdir -p "$CONFIG_DIR/fish/themes" "$CONFIG_DIR/starship" "$CONFIG_DIR/atuin/themes" "$BIN_DIR"
cp "$REPO/fish/themes/gentleman.fish" "$CONFIG_DIR/fish/themes/gentleman.fish"
cp "$REPO/fish/themes/gentleman-cute.fish" "$CONFIG_DIR/fish/themes/gentleman-cute.fish"
printf '# Starship profile fixture\n' > "$CONFIG_DIR/starship/gentleman.toml"
printf '# Starship profile fixture\n' > "$CONFIG_DIR/starship/gentleman-cute.toml"
printf '[theme]\n[colors]\n' > "$CONFIG_DIR/atuin/themes/gentleman.toml"
printf '[theme]\n[colors]\n' > "$CONFIG_DIR/atuin/themes/gentleman-cute.toml"
cat > "$CONFIG_DIR/atuin/config.toml" <<'CONFIG'
enter_accept = true

[sync]
records = true
CONFIG

cat > "$BIN_DIR/fish" <<EOF
#!/usr/bin/env bash
exec "$FISH_BIN" "\$@"
EOF
chmod +x "$BIN_DIR/fish"

cat > "$BIN_DIR/atuin" <<'STUB'
#!/usr/bin/env bash
set -euo pipefail

printf '%s\n' "$*" >> "$ATUIN_LOG"
if [ "${ATUIN_STUB_EXIT:-0}" != '0' ]; then
  exit "$ATUIN_STUB_EXIT"
fi

if [ "$1" = config ] && [ "$2" = set ] && [ "$3" = theme.name ]; then
  config="$HOME/.config/atuin/config.toml"
  temporary="$config.tmp"
  awk -v name="$4" '
    BEGIN { in_theme = 0; found_theme = 0 }
    /^\[theme\]$/ { in_theme = 1; found_theme = 1; print; next }
    /^\[/ { in_theme = 0 }
    in_theme && /^name[[:space:]]*=/ { print "name = \"" name "\""; next }
    { print }
    END {
      if (!found_theme) {
        print ""
        print "[theme]"
        print "name = \"" name "\""
      }
    }
  ' "$config" > "$temporary"
  mv "$temporary" "$config"
fi
STUB
chmod +x "$BIN_DIR/atuin"

assert_eq() {
  local expected=$1
  local actual=$2
  local message=$3

  if [ "$expected" != "$actual" ]; then
    printf 'assertion failed: %s\nexpected: %s\nactual: %s\n' "$message" "$expected" "$actual" >&2
    exit 1
  fi
}

assert_file_contains() {
  local path=$1
  local expected=$2
  local message=$3

  if ! grep -Fq -- "$expected" "$path"; then
    printf 'assertion failed: %s\nmissing: %s\n' "$message" "$expected" >&2
    exit 1
  fi
}

assert_no_candidates() {
  if find "$CONFIG_DIR/fish" -maxdepth 1 -name '.gentleman-theme.*' -print -quit | grep -q .; then
    printf 'assertion failed: Fish selector candidate was not cleaned up\n' >&2
    exit 1
  fi
}

run_selector() {
  # Fish expands $argv inside the child process, not in this Bash test.
  # shellcheck disable=SC2016
  HOME="$HOME_DIR" PATH="$BIN_DIR:/usr/bin:/bin" ATUIN_LOG="$ATUIN_LOG" ATUIN_STUB_EXIT="${ATUIN_STUB_EXIT:-0}" "$FISH_BIN" -c 'source "$argv[1]"; fish-theme "$argv[2]"' "$FUNCTION" "$1"
}

# Managed assets use Atuin's complete 18.15.2 foreground-only theme schema and approved palettes.
assert_eq "$(cat <<'THEME'
[theme]
name = "gentleman"
parent = "default"

[colors]
AlertInfo = "#7FB4CA"
AlertWarn = "#FFE066"
AlertError = "#CB7C94"
Annotation = "#DEBA87"
Base = "#F3F6F9"
Guidance = "#7AA89F"
Important = "#FF8DD7"
Title = "#FF8DD7"
Muted = "#8394A3"
THEME
)" "$(cat "$REPO/atuin/themes/gentleman.toml")" 'Gentleman declares exact Atuin metadata and foreground semantic palette'

assert_eq "$(cat <<'THEME'
[theme]
name = "gentleman-cute"
parent = "default"

[colors]
AlertInfo = "#D2CBD0"
AlertWarn = "#E0C27A"
AlertError = "#FF718F"
Annotation = "#A78E9B"
Base = "#F6EFF3"
Guidance = "#D2CBD0"
Important = "#F095C8"
Title = "#FFB1DD"
Muted = "#76616B"
THEME
)" "$(cat "$REPO/atuin/themes/gentleman-cute.toml")" 'Cute declares exact Atuin metadata and foreground semantic palette'

assert_file_contains "$REPO/fish.nix" '".config/atuin/themes/gentleman.toml".source = ./atuin/themes/gentleman.toml;' 'Home Manager installs the Gentleman Atuin asset'
assert_file_contains "$REPO/fish.nix" '".config/atuin/themes/gentleman-cute.toml".source = ./atuin/themes/gentleman-cute.toml;' 'Home Manager installs the Cute Atuin asset'

# Both supported profiles invoke Atuin's config setter and publish the shared marker.
: > "$ATUIN_LOG"
run_selector gentleman
assert_eq 'gentleman' "$(cat "$CONFIG_DIR/fish/gentleman-theme")" 'Gentleman publishes the Fish marker'
assert_file_contains "$ATUIN_LOG" 'config set theme.name gentleman' 'Gentleman persists the Atuin profile'
assert_file_contains "$CONFIG_DIR/atuin/config.toml" 'enter_accept = true' 'Gentleman retains unrelated Atuin config'
assert_file_contains "$CONFIG_DIR/atuin/config.toml" 'records = true' 'Gentleman retains Atuin sync configuration'
assert_file_contains "$CONFIG_DIR/atuin/config.toml" 'name = "gentleman"' 'Gentleman activates the matching Atuin theme'
assert_no_candidates

run_selector gentleman-cute
assert_eq 'gentleman-cute' "$(cat "$CONFIG_DIR/fish/gentleman-theme")" 'Cute publishes the Fish marker'
assert_file_contains "$ATUIN_LOG" 'config set theme.name gentleman-cute' 'Cute persists the Atuin profile'
assert_file_contains "$CONFIG_DIR/atuin/config.toml" 'enter_accept = true' 'Cute retains unrelated Atuin config'
assert_file_contains "$CONFIG_DIR/atuin/config.toml" 'records = true' 'Cute retains Atuin sync configuration'
assert_file_contains "$CONFIG_DIR/atuin/config.toml" 'name = "gentleman-cute"' 'Cute activates the matching Atuin theme'
assert_no_candidates

# A missing selected theme is rejected before either Atuin or the Fish marker changes.
printf 'gentleman\n' > "$CONFIG_DIR/fish/gentleman-theme"
mv "$CONFIG_DIR/atuin/themes/gentleman-cute.toml" "$CONFIG_DIR/atuin/themes/gentleman-cute.toml.unavailable"
: > "$ATUIN_LOG"
set +e
run_selector gentleman-cute >/dev/null 2>&1
selector_status=$?
set -e
assert_eq '1' "$selector_status" 'Missing Atuin theme exits with failure'
assert_eq 'gentleman' "$(cat "$CONFIG_DIR/fish/gentleman-theme")" 'Missing Atuin theme retains the previous Fish marker'
assert_eq '' "$(cat "$ATUIN_LOG")" 'Missing Atuin theme does not invoke Atuin'
assert_no_candidates
mv "$CONFIG_DIR/atuin/themes/gentleman-cute.toml.unavailable" "$CONFIG_DIR/atuin/themes/gentleman-cute.toml"

# A missing Atuin executable is rejected before the Fish marker changes.
printf 'gentleman\n' > "$CONFIG_DIR/fish/gentleman-theme"
mv "$BIN_DIR/atuin" "$BIN_DIR/atuin.unavailable"
set +e
run_selector gentleman-cute >/dev/null 2>&1
selector_status=$?
set -e
assert_eq '1' "$selector_status" 'Missing Atuin executable exits with failure'
assert_eq 'gentleman' "$(cat "$CONFIG_DIR/fish/gentleman-theme")" 'Missing Atuin executable retains the previous Fish marker'
assert_no_candidates
mv "$BIN_DIR/atuin.unavailable" "$BIN_DIR/atuin"

# A failed Atuin update does not report success or claim a fully applied Fish profile.
printf 'gentleman\n' > "$CONFIG_DIR/fish/gentleman-theme"
: > "$ATUIN_LOG"
set +e
ATUIN_STUB_EXIT=23 run_selector gentleman-cute >/dev/null 2>&1
selector_status=$?
set -e
assert_eq '23' "$selector_status" 'Failed Atuin update returns Atuin failure'
assert_eq 'gentleman' "$(cat "$CONFIG_DIR/fish/gentleman-theme")" 'Failed Atuin update retains the previous Fish marker'
assert_file_contains "$ATUIN_LOG" 'config set theme.name gentleman-cute' 'Failed update attempted the selected Atuin profile'
assert_no_candidates

printf 'fish theme Atuin tests passed\n'
