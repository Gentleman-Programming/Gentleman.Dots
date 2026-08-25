#!/bin/bash

set -euo pipefail

REPO=${1:-"$(cd "$(dirname "$0")/../.." && pwd)"}
SOURCE_CONFIG="$REPO/sketchybar"
TEMP_ROOT=$(mktemp -d "${TMPDIR:-/tmp}/sketchybar-theme-tests.XXXXXX")
trap 'rm -rf "$TEMP_ROOT"' EXIT

HOME_DIR="$TEMP_ROOT/home"
CONFIG_DIR="$HOME_DIR/.config/sketchybar"
MARKER_DIR="$HOME_DIR/.config/gentleman"
BIN_DIR="$TEMP_ROOT/bin"
LOG_FILE="$TEMP_ROOT/sketchybar.log"

mkdir -p "$MARKER_DIR" "$BIN_DIR"
cp -R "$SOURCE_CONFIG" "$CONFIG_DIR"

cat > "$BIN_DIR/sketchybar" <<'STUB'
#!/bin/bash
printf '%s\n' "$*" >> "$SKETCHYBAR_LOG"
exit "${SKETCHYBAR_STUB_EXIT:-0}"
STUB
chmod +x "$BIN_DIR/sketchybar"

assert_eq() {
  local expected=$1
  local actual=$2
  local message=$3

  if [ "$expected" != "$actual" ]; then
    printf 'assertion failed: %s\nexpected: %s\nactual: %s\n' "$message" "$expected" "$actual" >&2
    exit 1
  fi
}

assert_contains() {
  local needle=$1
  local haystack=$2
  local message=$3

  if [[ "$haystack" != *"$needle"* ]]; then
    printf 'assertion failed: %s\nmissing: %s\nactual: %s\n' "$message" "$needle" "$haystack" >&2
    exit 1
  fi
}

icon_color_for_item() {
  local item=$1
  local command_stream=$2

  printf '%s\n' "$command_stream" | grep -E -- "(--add item $item |--set $item )" | sed -n 's/.*icon\.color=\([^ ]*\).*/\1/p'
}

assert_no_candidates() {
  if /usr/bin/find "$MARKER_DIR" -maxdepth 1 -name '.sketchybar-theme.*' -print -quit | grep -q .; then
    printf 'assertion failed: selector candidate was not cleaned up\n' >&2
    exit 1
  fi
}

resolve_value() {
  local variable=$1

  THEME_VARIABLE="$variable" HOME="$HOME_DIR" CONFIG_DIR="$CONFIG_DIR" /bin/bash -c 'source "$CONFIG_DIR/theme.sh"; printf "%s" "${!THEME_VARIABLE}"'
}

run_selector() {
  HOME="$HOME_DIR" PATH="$BIN_DIR:$PATH" SKETCHYBAR_LOG="$LOG_FILE" /bin/bash "$CONFIG_DIR/sketchybar-theme" "$@"
}

run_plugin() {
  local plugin=$1
  shift
  : > "$LOG_FILE"
  env HOME="$HOME_DIR" CONFIG_DIR="$CONFIG_DIR" PATH="$BIN_DIR:$PATH" SKETCHYBAR_LOG="$LOG_FILE" NAME=test-item "$@" "$CONFIG_DIR/plugins/$plugin"
  cat "$LOG_FILE"
}

run_config() {
  : > "$LOG_FILE"
  env HOME="$HOME_DIR" CONFIG_DIR="$CONFIG_DIR" PATH="$BIN_DIR:$PATH" SKETCHYBAR_LOG="$LOG_FILE" /bin/bash "$CONFIG_DIR/sketchybarrc"
  cat "$LOG_FILE"
}

# Resolver defaults to the preserved Gentleman palette without a marker.
assert_eq '0xff121620' "$(resolve_value ISLAND_BG)" "missing marker defaults to Gentleman island background"
assert_eq '0xfff3f6f9' "$(resolve_value TEXT)" "Gentleman text preserves the current white"
assert_eq '0xff565f89' "$(resolve_value DIM)" "Gentleman dim preserves the current dim"
assert_eq '0xffe0c15a' "$(resolve_value ACCENT)" "Gentleman accent preserves the current accent"
assert_eq '0xff263356' "$(resolve_value ISLAND_BORDER)" "Gentleman border preserves the current island border"
assert_eq '0xff121620' "$(resolve_value SELECTED_BG)" "Gentleman selected workspace preserves the current island background"
assert_eq '0xff7aa89f' "$(resolve_value CPU_NORMAL)" "Gentleman normal CPU preserves the current cyan"
assert_eq '0xffcb7c94' "$(resolve_value CPU_INITIAL)" "Gentleman initial CPU preserves the old red"
assert_eq '0xffe0c15a' "$(resolve_value CHARGING)" "Gentleman charging preserves the existing accent"
assert_eq '0xfffff7b1' "$(resolve_value GPU_NORMAL)" "Gentleman normal GPU preserves the current orange"
assert_eq '0xffff8dd7' "$(resolve_value RAM_NORMAL)" "Gentleman normal RAM preserves the current magenta"
assert_eq '0xff7fb4ca' "$(resolve_value VOLUME_NORMAL)" "Gentleman normal volume preserves the current blue"

cat > "$BIN_DIR/pmset" <<'STUB'
#!/bin/bash
printf "Now drawing from 'AC Power'\n -InternalBattery-0 (id=1234567)\t50%%; charging; 0:30 remaining present: true\n"
STUB
chmod +x "$BIN_DIR/pmset"

# Charging must emit the selected profile's semantic color, not pmset output.
for profile in gentleman gentleman-cute; do
  printf '%s\n' "$profile" > "$MARKER_DIR/sketchybar-theme"
  battery_output=$(run_plugin battery.sh env)
  assert_eq "$(resolve_value CHARGING)" "$(icon_color_for_item test-item "$battery_output")" "$profile charging battery color uses CHARGING"
done

printf 'not-a-theme\n' > "$MARKER_DIR/sketchybar-theme"
assert_eq '0xff121620' "$(resolve_value ISLAND_BG)" "invalid marker defaults to Gentleman"

printf 'gentleman\n' > "$MARKER_DIR/sketchybar-theme"
assert_eq '0xffcb7c94' "$(resolve_value ERROR)" "Gentleman error preserves the current red"
assert_eq '0xffffe066' "$(resolve_value WARNING)" "Gentleman warning preserves the current yellow"

printf 'gentleman-cute\n' > "$MARKER_DIR/sketchybar-theme"
assert_eq '0xff1a1218' "$(resolve_value PANEL_BG)" "Cute panel background"
assert_eq '0xff241822' "$(resolve_value ISLAND_BG)" "Cute island background"
assert_eq '0xff342230' "$(resolve_value SELECTED_BG)" "Cute selected workspace background"
assert_eq '0xfff6eff3' "$(resolve_value TEXT)" "Cute text"
assert_eq '0xffa78e9b' "$(resolve_value MUTED)" "Cute muted"
assert_eq '0xff76616b' "$(resolve_value DIM)" "Cute dim"
assert_eq '0xfff095c8' "$(resolve_value ACCENT)" "Cute accent"
assert_eq '0xffc96aa2' "$(resolve_value DEEP_ROSE)" "Cute deep rose"
assert_eq '0xffffb1dd' "$(resolve_value ACTIVE_ROSE)" "Cute active rose"
assert_eq '0xffe0c27a' "$(resolve_value CHAMPAGNE)" "Cute champagne"
assert_eq '0xffe0c27a' "$(resolve_value CHARGING)" "Cute charging uses champagne"
assert_eq '0xffd2cbd0' "$(resolve_value PEARL)" "Cute pearl"
assert_eq '0xffd7a0b8' "$(resolve_value SOFT_ROSE)" "Cute soft rose"
assert_eq '0xffa9c7ee' "$(resolve_value POWDER_BLUE)" "Cute powder blue"
assert_eq '0xfff2b86d' "$(resolve_value WARNING)" "Cute warning"
assert_eq '0xffff718f' "$(resolve_value ERROR)" "Cute error"
assert_eq '0xffffb1dd' "$(resolve_value WORKSPACE_ACTIVE)" "Cute selected workspace uses active rose"
assert_eq '0xffd7a0b8' "$(resolve_value CPU_NORMAL)" "Cute normal CPU uses soft rose"
assert_eq '0xffd7a0b8' "$(resolve_value CPU_INITIAL)" "Cute initial CPU uses soft rose"
assert_eq "$(resolve_value CPU_NORMAL)" "$(resolve_value CPU_INITIAL)" "Cute initial CPU aliases normal CPU"
assert_eq '0xffd7a0b8' "$(resolve_value NETWORK_NORMAL)" "Cute network uses soft rose"
assert_eq '0xfff095c8' "$(resolve_value VOLUME_NORMAL)" "Cute normal volume uses primary rose"
assert_eq '0xffc96aa2' "$(resolve_value RAM_NORMAL)" "Cute RAM uses deep rose"

# The selector accepts only the two profile names and publishes the marker before reload.
printf 'gentleman\n' > "$MARKER_DIR/sketchybar-theme"
: > "$LOG_FILE"
run_selector gentleman
assert_eq 'gentleman' "$(cat "$MARKER_DIR/sketchybar-theme")" "selector publishes Gentleman"
assert_eq '--reload' "$(cat "$LOG_FILE")" "selector reloads after publishing Gentleman"
assert_no_candidates

: > "$LOG_FILE"
run_selector gentleman-cute
assert_eq 'gentleman-cute' "$(cat "$MARKER_DIR/sketchybar-theme")" "selector publishes Cute"
assert_eq '--reload' "$(cat "$LOG_FILE")" "selector reloads after publishing Cute"
assert_no_candidates

for invalid_args in '' 'unknown-theme'; do
  printf 'gentleman\n' > "$MARKER_DIR/sketchybar-theme"
  : > "$LOG_FILE"
  set +e
  if [ -n "$invalid_args" ]; then
    run_selector "$invalid_args" >/dev/null 2>&1
  else
    run_selector >/dev/null 2>&1
  fi
  selector_status=$?
  set -e
  assert_eq '2' "$selector_status" "invalid selector input exits 2"
  assert_eq 'gentleman' "$(cat "$MARKER_DIR/sketchybar-theme")" "invalid selector input does not mutate the marker"
  assert_eq '' "$(cat "$LOG_FILE")" "invalid selector input does not reload"
done

# Installed resolver/profile validation happens before publication and leaves no candidate behind.
printf 'gentleman\n' > "$MARKER_DIR/sketchybar-theme"
rm "$CONFIG_DIR/themes/gentleman-cute.sh"
set +e
run_selector gentleman-cute >/dev/null 2>&1
selector_status=$?
set -e
if [ "$selector_status" -eq 0 ]; then
  printf 'assertion failed: missing profile validation unexpectedly succeeded\n' >&2
  exit 1
fi
assert_eq 'gentleman' "$(cat "$MARKER_DIR/sketchybar-theme")" "validation failure does not publish a marker"
assert_no_candidates
cp "$SOURCE_CONFIG/themes/gentleman-cute.sh" "$CONFIG_DIR/themes/gentleman-cute.sh"

printf 'if then\n' > "$CONFIG_DIR/themes/gentleman-cute.sh"
set +e
run_selector gentleman-cute >/dev/null 2>&1
selector_status=$?
set -e
if [ "$selector_status" -eq 0 ]; then
  printf 'assertion failed: invalid profile syntax unexpectedly succeeded\n' >&2
  exit 1
fi
assert_eq 'gentleman' "$(cat "$MARKER_DIR/sketchybar-theme")" "syntax validation failure does not publish a marker"
assert_no_candidates
cp "$SOURCE_CONFIG/themes/gentleman-cute.sh" "$CONFIG_DIR/themes/gentleman-cute.sh"

# A pre-publication move failure cleans its candidate and retains the previous marker.
cat > "$BIN_DIR/mv" <<'STUB'
#!/bin/bash
if [ "${MV_STUB_FAIL:-0}" = '1' ]; then
  exit 1
fi
exec /bin/mv "$@"
STUB
chmod +x "$BIN_DIR/mv"
printf 'gentleman\n' > "$MARKER_DIR/sketchybar-theme"
export MV_STUB_FAIL=1
set +e
run_selector gentleman-cute >/dev/null 2>&1
selector_status=$?
set -e
unset MV_STUB_FAIL
if [ "$selector_status" -eq 0 ]; then
  printf 'assertion failed: marker move failure unexpectedly succeeded\n' >&2
  exit 1
fi
assert_eq 'gentleman' "$(cat "$MARKER_DIR/sketchybar-theme")" "move failure does not publish a marker"
assert_no_candidates

# A reload failure retains the valid, published marker for the next launch.
printf 'gentleman\n' > "$MARKER_DIR/sketchybar-theme"
: > "$LOG_FILE"
set +e
reload_output=$(SKETCHYBAR_STUB_EXIT=41 run_selector gentleman-cute 2>&1)
selector_status=$?
set -e
if [ "$selector_status" -eq 0 ]; then
  printf 'assertion failed: reload failure unexpectedly succeeded\n' >&2
  exit 1
fi
assert_eq 'gentleman-cute' "$(cat "$MARKER_DIR/sketchybar-theme")" "reload failure retains the published marker"
assert_contains 'reload failed' "$reload_output" "reload failure is reported clearly"
assert_eq '--reload' "$(cat "$LOG_FILE")" "reload failure still invokes SketchyBar reload"
assert_no_candidates

# Every color-bearing plugin resolves profile tokens rather than retaining a local palette.
for plugin in battery.sh cpu.sh current_apps.sh github.sh gpu.sh media.sh meeting.sh mic.sh mic_toggle.sh music.sh ram.sh space.sh volume.sh wifi.sh; do
  grep -Fq "source \"\$CONFIG_DIR/theme.sh\"" "$CONFIG_DIR/plugins/$plugin"
  if grep -Eq '0x[0-9A-Fa-f]+' "$CONFIG_DIR/plugins/$plugin"; then
    printf 'assertion failed: %s retains a hard-coded color\n' "$plugin" >&2
    exit 1
  fi
done

# A controlled SketchyBar run proves the Gentleman layout renders the current colors.
printf 'gentleman\n' > "$MARKER_DIR/sketchybar-theme"
main_output=$(run_config)
assert_contains "color=0x00000000" "$main_output" "bar preserves its transparent base"
assert_contains "icon.color=0xfff3f6f9" "$main_output" "default text preserves its current white"
assert_contains "label.color=0xff565f89" "$main_output" "workspace labels preserve their dim color"
assert_contains "icon.color=0xffe0c15a" "$main_output" "separator preserves the current accent"
assert_contains "label.color=0xff06080f" "$main_output" "front app text preserves its current black"
assert_contains "background.color=0xffb7cc85" "$main_output" "front app preserves its current green island"
assert_contains "icon.color=0xff7fb4ca" "$main_output" "volume preserves its initial blue"
assert_eq '0xffcb7c94' "$(icon_color_for_item cpu "$main_output")" "Gentleman initial CPU renders the old red"
assert_contains "icon.color=0xfffff7b1" "$main_output" "GPU preserves its initial orange"
assert_contains "icon.color=0xffff8dd7" "$main_output" "RAM preserves its initial magenta"

# Controlled plugin runs prove semantic tokens produce the preserved Gentleman outcomes.
space_selected=$(run_plugin space.sh env SELECTED=true)
assert_contains 'label.color=0xffe0c15a' "$space_selected" "selected workspace preserves accent"
assert_contains 'background.color=0xff121620' "$space_selected" "selected workspace preserves island background"
assert_contains 'background.border_color=0xffe0c15a' "$space_selected" "selected workspace preserves accent border"

printf 'gentleman-cute\n' > "$MARKER_DIR/sketchybar-theme"
cute_main_output=$(run_config)
assert_eq "$(resolve_value CPU_INITIAL)" "$(icon_color_for_item cpu "$cute_main_output")" "Cute initial CPU renders its profile token"
assert_eq '0xffd7a0b8' "$(icon_color_for_item cpu "$cute_main_output")" "Cute initial CPU renders soft rose"
cute_space_selected=$(run_plugin space.sh env SELECTED=true)
assert_contains 'label.color=0xffffb1dd' "$cute_space_selected" "Cute selected workspace uses active rose"
assert_contains 'background.color=0xff342230' "$cute_space_selected" "Cute selected workspace uses selected background"
assert_contains 'background.border_color=0xffffb1dd' "$cute_space_selected" "Cute selected workspace uses active rose border"

printf 'gentleman\n' > "$MARKER_DIR/sketchybar-theme"
cat > "$BIN_DIR/osascript" <<'STUB'
#!/bin/bash
case "$*" in
  *'output volume'*) printf '50\n' ;;
  *'output muted'*) printf 'false\n' ;;
  *) exit 1 ;;
esac
STUB
chmod +x "$BIN_DIR/osascript"
volume_output=$(run_plugin volume.sh env)
assert_contains 'icon.color=0xff7fb4ca' "$volume_output" "normal volume preserves blue"

cat > "$BIN_DIR/sysctl" <<'STUB'
#!/bin/bash
printf '4\n'
STUB
cat > "$BIN_DIR/ps" <<'STUB'
#!/bin/bash
printf '40\n'
STUB
chmod +x "$BIN_DIR/sysctl" "$BIN_DIR/ps"
cpu_output=$(run_plugin cpu.sh env)
assert_contains 'icon.color=0xff7aa89f' "$cpu_output" "normal CPU preserves cyan"

cat > "$BIN_DIR/ioreg" <<'STUB'
#!/bin/bash
printf '"Device Utilization %%"=20\n'
STUB
chmod +x "$BIN_DIR/ioreg"
gpu_output=$(run_plugin gpu.sh env)
assert_contains 'icon.color=0xfffff7b1' "$gpu_output" "normal GPU preserves orange"

cat > "$BIN_DIR/vm_stat" <<'STUB'
#!/bin/bash
cat <<'OUTPUT'
Pages active: 1000.
Pages wired down: 1000.
Pages free: 8000.
Pages inactive: 0.
OUTPUT
STUB
chmod +x "$BIN_DIR/vm_stat"
ram_output=$(run_plugin ram.sh env)
assert_contains 'icon.color=0xffff8dd7' "$ram_output" "normal RAM preserves magenta"

printf 'theme profile tests passed\n'
