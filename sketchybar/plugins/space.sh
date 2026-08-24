#!/bin/bash

# Space/Workspace indicator - zero yabai queries, pure sketchybar state

# CONFIG_DIR is supplied by SketchyBar at runtime.
# shellcheck disable=SC1091
source "$CONFIG_DIR/theme.sh"

if [ "$SELECTED" = "true" ]; then
  sketchybar --set "$NAME" \
    label.color="$WORKSPACE_ACTIVE" \
    label.font="IosevkaTerm NF:Bold:12.0" \
    background.color="$SELECTED_BG" \
    background.border_color="$WORKSPACE_ACTIVE"
else
  sketchybar --set "$NAME" \
    label.color="$DIM" \
    label.font="IosevkaTerm NF:Regular:12.0" \
    background.color="$ISLAND_BG" \
    background.border_color="$ISLAND_BORDER"
fi
