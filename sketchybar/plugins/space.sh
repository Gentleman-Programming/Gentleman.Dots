#!/bin/bash

# Space/Workspace indicator - zero yabai queries, pure sketchybar state

ACCENT=0xffe0c15a
DIM=0xff565f89
ISLAND_BORDER=0xff263356

if [ "$SELECTED" = "true" ]; then
  sketchybar --set $NAME \
    label.color=$ACCENT \
    label.font="IosevkaTerm NF:Bold:12.0" \
    background.border_color=$ACCENT
else
  sketchybar --set $NAME \
    label.color=$DIM \
    label.font="IosevkaTerm NF:Regular:12.0" \
    background.border_color=$ISLAND_BORDER
fi
