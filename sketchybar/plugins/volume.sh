#!/bin/bash

# Volume - displays current volume level

RAW_VOLUME=$(osascript -e "output volume of (get volume settings)" 2>/dev/null)
MUTED=$(osascript -e "output muted of (get volume settings)" 2>/dev/null)

if [[ "$RAW_VOLUME" =~ ^[0-9]+$ ]]; then
  VOLUME="$RAW_VOLUME"
  LABEL="${VOLUME}%"
else
  VOLUME="0"
  LABEL="--"
fi

if [ "$MUTED" = "true" ] || [ "$VOLUME" -eq 0 ]; then
  COLOR=0xffcb7c94
else
  COLOR=0xff7fb4ca
fi

sketchybar --set $NAME icon.color="$COLOR" label="$LABEL"
