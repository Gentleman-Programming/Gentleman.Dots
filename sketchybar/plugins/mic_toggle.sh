#!/bin/bash

# Toggle microphone mute/unmute with pulse animation

# CONFIG_DIR is supplied by SketchyBar at runtime.
# shellcheck disable=SC1091
source "$CONFIG_DIR/theme.sh"

# Animation settings
ANIM_DURATION=5
ANIM_CURVE="sin"

# Show loading with scale pulse
sketchybar --animate $ANIM_CURVE $ANIM_DURATION --set mic \
  icon="..." \
  icon.color="$WARNING" \
  background.y_offset=2 background.y_offset=0

# Get current state
MIC_VOLUME=$(osascript -e "input volume of (get volume settings)")

if [ "$MIC_VOLUME" -eq 0 ]; then
  # Unmute - set to 100%
  osascript -e "set volume input volume 100"
else
  # Mute - set to 0%
  osascript -e "set volume input volume 0"
fi

# Small delay
sleep 0.15

# Update with animation
MIC_VOLUME=$(osascript -e "input volume of (get volume settings)")

if [ "$MIC_VOLUME" -eq 0 ]; then
  sketchybar --animate $ANIM_CURVE $ANIM_DURATION --set mic \
    icon="MIC" \
    icon.color="$MUTED" \
    label="OFF"
else
  sketchybar --animate $ANIM_CURVE $ANIM_DURATION --set mic \
    icon="MIC" \
    icon.color="$SOFT_ROSE" \
    label="ON"
fi
