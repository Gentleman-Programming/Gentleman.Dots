#!/bin/bash

# GitHub - displays unread notification count

# CONFIG_DIR is supplied by SketchyBar at runtime.
# shellcheck disable=SC1091
source "$CONFIG_DIR/theme.sh"

COUNT=$(gh api notifications 2>/dev/null | jq 'length' 2>/dev/null)

if [ -z "$COUNT" ] || [ "$COUNT" = "null" ]; then
  COUNT=0
fi

if [ "$COUNT" -gt 0 ]; then
  sketchybar --set "$NAME" icon.color="$WARNING" label="$COUNT"
else
  sketchybar --set "$NAME" icon.color="$MUTED" label="0"
fi
