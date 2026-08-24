#!/bin/bash

# Battery - displays battery percentage with dynamic icon and color

# CONFIG_DIR is supplied by SketchyBar at runtime.
# shellcheck disable=SC1091
source "$CONFIG_DIR/theme.sh"

PERCENTAGE=$(pmset -g batt | grep -Eo "\d+%" | cut -d% -f1)
ON_AC_POWER=$(pmset -g batt | grep 'AC Power')

if [ -z "$PERCENTAGE" ]; then
  sketchybar --set "$NAME" drawing=off
  exit 0
fi

# Determine icon and color based on level
if [ -n "$ON_AC_POWER" ]; then
  ICON="󰂄"
  COLOR=$CHARGING
elif [ "$PERCENTAGE" -ge 80 ]; then
  ICON="󰁹"
  COLOR=$BATTERY_HEALTHY
elif [ "$PERCENTAGE" -ge 60 ]; then
  ICON="󰂁"
  COLOR=$BATTERY_HEALTHY
elif [ "$PERCENTAGE" -ge 40 ]; then
  ICON="󰁿"
  COLOR=$WARNING
elif [ "$PERCENTAGE" -ge 20 ]; then
  ICON="󰁼"
  COLOR=$ERROR
else
  ICON="󰂃"
  COLOR=$ERROR
fi

sketchybar --set "$NAME" icon="$ICON" icon.color="$COLOR" label="${PERCENTAGE}%"
