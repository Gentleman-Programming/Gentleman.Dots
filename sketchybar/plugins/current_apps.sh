#!/bin/bash

# Current Apps - shows all apps in the current workspace
# The focused app is highlighted with brackets [App]

# CONFIG_DIR is supplied by SketchyBar at runtime.
# shellcheck disable=SC1091
source "$CONFIG_DIR/theme.sh"

# Get current space
CURRENT_SPACE=$(yabai -m query --spaces --space 2>/dev/null | jq -r '.index')

# Get focused window app name
FOCUSED_APP=$(yabai -m query --windows --window 2>/dev/null | jq -r '.app' 2>/dev/null)

# Get all apps in current space
APPS=$(yabai -m query --windows --space "$CURRENT_SPACE" 2>/dev/null | jq -r '.[].app' 2>/dev/null | sort -u)

if [ -z "$APPS" ]; then
  sketchybar --set "$NAME" label="--"
  exit 0
fi

# Build label with focused app highlighted
LABEL=""
while IFS= read -r app; do
  if [ -n "$app" ]; then
    if [ "$app" = "$FOCUSED_APP" ]; then
      # Focused app in brackets
      if [ -z "$LABEL" ]; then
        LABEL="[$app]"
      else
        LABEL="$LABEL | [$app]"
      fi
    else
      if [ -z "$LABEL" ]; then
        LABEL="$app"
      else
        LABEL="$LABEL | $app"
      fi
    fi
  fi
done <<< "$APPS"

sketchybar --set "$NAME" label="$LABEL"
