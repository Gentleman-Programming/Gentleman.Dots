#!/bin/bash

# Space/Workspace indicator script

ACCENT=0xffe0c15a
WHITE=0xfff3f6f9
DIM=0xff565f89
ISLAND_BG=0xff121620
ISLAND_BORDER=0xff263356

# Get space ID from item name (space.1 -> 1)
SPACE_ID=${NAME#space.}

# Get apps in this space from yabai, join with " | "
APPS=$(yabai -m query --windows --space $SPACE_ID 2>/dev/null | jq -r '.[].app' 2>/dev/null | sort -u | head -3 | paste -sd '|' - | sed 's/|/ | /g')

# Build label with app names
if [ -n "$APPS" ]; then
  LABEL="$SPACE_ID · $APPS"
else
  LABEL="$SPACE_ID"
fi

# Check if this space is selected
if [ "$SELECTED" = "true" ]; then
  sketchybar --set $NAME \
    label="$LABEL" \
    label.color=$ACCENT \
    label.font="IosevkaTerm NF:Bold:12.0" \
    background.border_color=$ACCENT
else
  if [ -n "$APPS" ]; then
    sketchybar --set $NAME \
      label="$LABEL" \
      label.color=$WHITE \
      label.font="IosevkaTerm NF:Regular:12.0" \
      background.border_color=$ISLAND_BORDER
  else
    sketchybar --set $NAME \
      label="$LABEL" \
      label.color=$DIM \
      label.font="IosevkaTerm NF:Regular:12.0" \
      background.border_color=$ISLAND_BORDER
  fi
fi
