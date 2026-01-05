#!/bin/bash

# Toggle microphone mute/unmute with loading state

# Show loading spinner
sketchybar --set mic icon="..." icon.color=0xffffe066

# Get current state
MIC_VOLUME=$(osascript -e "input volume of (get volume settings)")

if [ "$MIC_VOLUME" -eq 0 ]; then
  # Unmute - set to 100%
  osascript -e "set volume input volume 100"
else
  # Mute - set to 0%
  osascript -e "set volume input volume 0"
fi

# Small delay to show the spinner
sleep 0.2

# Update the sketchybar item with new state
MIC_VOLUME=$(osascript -e "input volume of (get volume settings)")

if [ "$MIC_VOLUME" -eq 0 ]; then
  sketchybar --set mic \
    icon="MIC" \
    icon.color=0xff565f89 \
    label="OFF"
else
  sketchybar --set mic \
    icon="MIC" \
    icon.color=0xffcb7c94 \
    label="ON"
fi
