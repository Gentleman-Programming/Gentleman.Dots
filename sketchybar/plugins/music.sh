#!/bin/bash

# Music - single osascript call instead of 3 separate ones

# CONFIG_DIR is supplied by SketchyBar at runtime.
# shellcheck disable=SC1091
source "$CONFIG_DIR/theme.sh"

if pgrep -x "Spotify" > /dev/null 2>&1; then
  INFO=$(osascript -e 'tell application "Spotify" to if player state is playing then return artist of current track & " - " & name of current track' 2>/dev/null)
  if [ -n "$INFO" ]; then
    sketchybar --set "$NAME" icon.color="$MEDIA_SPOTIFY" label="$INFO"
    exit 0
  fi
fi

if pgrep -x "Music" > /dev/null 2>&1; then
  INFO=$(osascript -e 'tell application "Music" to if player state is playing then return artist of current track & " - " & name of current track' 2>/dev/null)
  if [ -n "$INFO" ]; then
    sketchybar --set "$NAME" icon.color="$MEDIA_MUSIC" label="$INFO"
    exit 0
  fi
fi

sketchybar --set "$NAME" icon.color="$MUTED" label="--"
