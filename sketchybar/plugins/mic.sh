#!/bin/bash

# Mic - shows microphone status (on/off)

# CONFIG_DIR is supplied by SketchyBar at runtime.
# shellcheck disable=SC1091
source "$CONFIG_DIR/theme.sh"

# Get mic input volume (0 = muted)
MIC_VOLUME=$(osascript -e "input volume of (get volume settings)")

if [ "$MIC_VOLUME" -eq 0 ]; then
  sketchybar --set "$NAME" icon="MIC" icon.color="$MUTED" label="OFF"
else
  sketchybar --set "$NAME" icon="MIC" icon.color="$SOFT_ROSE" label="ON"
fi
