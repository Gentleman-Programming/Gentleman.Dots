#!/usr/bin/env bash
# Move the focused window to a yabai space and follow it safely.
# Usage: move-window-to-space.sh <space_number>

set -euo pipefail

space="${1:-}"
if [[ -z "$space" || ! "$space" =~ ^[0-9]+$ ]]; then
	echo "Usage: $0 <space_number>" >&2
	exit 64
fi

if ! command -v yabai >/dev/null 2>&1; then
	echo "yabai not found" >&2
	exit 127
fi

if ! command -v jq >/dev/null 2>&1; then
	echo "jq not found" >&2
	exit 127
fi

window_json="$(yabai -m query --windows --window 2>/dev/null)" || {
	echo "No focused yabai window" >&2
	exit 1
}

window_id="$(jq -r '.id // empty' <<<"$window_json")"
can_move="$(jq -r '."can-move" // false' <<<"$window_json")"
current_space="$(jq -r '.space // empty' <<<"$window_json")"

if [[ -z "$window_id" || "$can_move" != "true" ]]; then
	echo "Focused window cannot be moved" >&2
	exit 1
fi

if ! yabai -m query --spaces | jq -e --arg space "$space" 'any(.[]; (.index | tostring) == $space)' >/dev/null; then
	echo "Space $space does not exist" >&2
	exit 1
fi

if [[ "$current_space" == "$space" ]]; then
	exit 0
fi

yabai -m window "$window_id" --space "$space"

moved_space=""
for _ in {1..20}; do
	moved_space="$(yabai -m query --windows --window "$window_id" 2>/dev/null | jq -r '.space // empty' || true)"
	[[ "$moved_space" == "$space" ]] && break
	sleep 0.05
done

if [[ "$moved_space" != "$space" ]]; then
	echo "Window $window_id did not move to space $space" >&2
	exit 1
fi

yabai -m space --focus "$space"
sleep 0.15

after_focus_space="$(yabai -m query --windows --window "$window_id" 2>/dev/null | jq -r '.space // empty' || true)"
if [[ "$after_focus_space" != "$space" ]]; then
	exit 1
fi

yabai -m window --focus "$window_id" 2>/dev/null || true
