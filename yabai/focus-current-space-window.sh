#!/usr/bin/env bash
# Focus a sane window in the current space without jumping across spaces.
# Used as an explicit rescue command when macOS/yabai leaves focus on the desktop
# or after closing/moving a window.

set -euo pipefail

if ! command -v yabai >/dev/null 2>&1; then
  exit 0
fi

if ! command -v jq >/dev/null 2>&1; then
  # Fallback: stays within the current space in normal yabai behavior.
  yabai -m window --focus first >/dev/null 2>&1 || true
  exit 0
fi

window_id="$(yabai -m query --windows --space 2>/dev/null | jq -r '
  map(select(
    .["is-minimized"] == false and
    .["is-hidden"] == false and
    .["is-visible"] == true
  ))
  | sort_by(.frame.y, .frame.x)
  | .[0].id // empty
')"

if [ -n "$window_id" ]; then
  yabai -m window --focus "$window_id" >/dev/null 2>&1 || true
fi
