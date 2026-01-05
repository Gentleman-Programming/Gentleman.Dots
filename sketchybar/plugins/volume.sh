#!/bin/bash

# Volume - displays current volume level



VOLUME=$(osascript -e "output volume of (get volume settings)")
MUTED=$(osascript -e "output muted of (get volume settings)")

if [ "$MUTED" = "true" ] || [ "$VOLUME" -eq 0 ]; then
  ICON="VOL"
  COLOR=0xffcb7c94
else
  ICON="VOL"
  COLOR=0xff7fb4ca
fi

sketchybar --set $NAME \
  icon="$ICON" \
  icon.color="$COLOR" \
  label="${VOLUME}%"
