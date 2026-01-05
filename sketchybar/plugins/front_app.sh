#!/bin/bash

# Front App - shows the currently focused application with icon + smooth fade

# Animation settings
ANIM_DURATION=6
ANIM_CURVE="tanh"

if [ "$SENDER" = "front_app_switched" ]; then
  sketchybar --animate $ANIM_CURVE $ANIM_DURATION --set $NAME \
    label="$INFO" \
    icon.background.image="app.$INFO"
fi
