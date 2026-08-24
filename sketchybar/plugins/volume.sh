#!/bin/bash

# Volume - displays current volume level

# CONFIG_DIR is supplied by SketchyBar at runtime.
# shellcheck disable=SC1091
source "$CONFIG_DIR/theme.sh"

RAW_VOLUME=$(osascript -e "output volume of (get volume settings)" 2>/dev/null)
MUTED=$(osascript -e "output muted of (get volume settings)" 2>/dev/null)

if [[ "$RAW_VOLUME" =~ ^[0-9]+$ ]]; then
  VOLUME="$RAW_VOLUME"
  LABEL="${VOLUME}%"
else
  VOLUME="0"
  LABEL="--"
fi

if [ "$MUTED" = "true" ] || [ "$VOLUME" -eq 0 ]; then
  COLOR=$ERROR
else
  COLOR=$VOLUME_NORMAL
fi

sketchybar --set "$NAME" icon.color="$COLOR" label="$LABEL"
