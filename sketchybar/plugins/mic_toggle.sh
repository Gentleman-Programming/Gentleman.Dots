#!/bin/bash

# Toggle microphone mute/unmute with pulse animation

# Colors
RED=0xffcb7c94
GREEN=0xffb7cc85
DIM=0xff565f89
YELLOW=0xffffe066

# Animation settings
ANIM_DURATION=5
ANIM_CURVE="sin"

# Show loading with scale pulse
sketchybar --animate $ANIM_CURVE $ANIM_DURATION --set mic \
  icon="..." \
  icon.color=$YELLOW \
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
    icon.color=$DIM \
    label="OFF"
else
  sketchybar --animate $ANIM_CURVE $ANIM_DURATION --set mic \
    icon="MIC" \
    icon.color=$RED \
    label="ON"
fi
