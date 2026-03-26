#!/bin/bash

# Space/Workspace indicator script with fast animations

ACCENT=0xffe0c15a
WHITE=0xfff3f6f9
DIM=0xff565f89
ISLAND_BG=0xff121620
ISLAND_BORDER=0xff263356

# Animation settings - fast
ANIM_DURATION=4
ANIM_CURVE="tanh"

# Get space ID from item name (space.1 -> 1)
SPACE_ID=${NAME#space.}

# Check if this space has any windows (lightweight check)
HAS_WINDOWS=$(yabai -m query --windows --space $SPACE_ID 2>/dev/null | jq -r 'length' 2>/dev/null)

LABEL="$SPACE_ID"

# Check if this space is selected - animate the transition
if [ "$SELECTED" = "true" ]; then
  sketchybar --animate $ANIM_CURVE $ANIM_DURATION --set $NAME \
    label="$LABEL" \
    label.color=$ACCENT \
    label.font="IosevkaTerm NF:Bold:12.0" \
    background.border_color=$ACCENT
else
  if [ -n "$HAS_WINDOWS" ] && [ "$HAS_WINDOWS" -gt 0 ] 2>/dev/null; then
    sketchybar --animate $ANIM_CURVE $ANIM_DURATION --set $NAME \
      label="$LABEL" \
      label.color=$WHITE \
      label.font="IosevkaTerm NF:Regular:12.0" \
      background.border_color=$ISLAND_BORDER
  else
    sketchybar --animate $ANIM_CURVE $ANIM_DURATION --set $NAME \
      label="$LABEL" \
      label.color=$DIM \
      label.font="IosevkaTerm NF:Regular:12.0" \
      background.border_color=$ISLAND_BORDER
  fi
fi
