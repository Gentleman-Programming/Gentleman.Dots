#!/bin/bash

# Media - shows currently playing track (Spotify/Apple Music)

# CONFIG_DIR is supplied by SketchyBar at runtime.
# shellcheck disable=SC1091
source "$CONFIG_DIR/theme.sh"

SPOTIFY_PLAYING=$(osascript -e 'tell application "Spotify" to player state' 2>/dev/null)
MUSIC_PLAYING=$(osascript -e 'tell application "Music" to player state' 2>/dev/null)

if [ "$SPOTIFY_PLAYING" = "playing" ]; then
  TRACK=$(osascript -e 'tell application "Spotify" to name of current track')
  ARTIST=$(osascript -e 'tell application "Spotify" to artist of current track')
  sketchybar --set "$NAME" \
    icon="" \
    icon.color="$MEDIA_SPOTIFY" \
    label="$ARTIST - $TRACK" \
    drawing=on
elif [ "$MUSIC_PLAYING" = "playing" ]; then
  TRACK=$(osascript -e 'tell application "Music" to name of current track')
  ARTIST=$(osascript -e 'tell application "Music" to artist of current track')
  sketchybar --set "$NAME" \
    icon="" \
    icon.color="$MEDIA_MUSIC" \
    label="$ARTIST - $TRACK" \
    drawing=on
else
  sketchybar --set "$NAME" drawing=off
fi
