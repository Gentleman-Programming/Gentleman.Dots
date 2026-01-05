#!/bin/bash

# Music - shows currently playing track (Spotify/Apple Music)

GREEN=0xffb7cc85
RED=0xffcb7c94
DIM=0xff565f89

SPOTIFY_PLAYING=$(osascript -e 'tell application "Spotify" to player state' 2>/dev/null)
MUSIC_PLAYING=$(osascript -e 'tell application "Music" to player state' 2>/dev/null)

if [ "$SPOTIFY_PLAYING" = "playing" ]; then
  TRACK=$(osascript -e 'tell application "Spotify" to name of current track')
  ARTIST=$(osascript -e 'tell application "Spotify" to artist of current track')
  sketchybar --set $NAME icon.color=$GREEN label="$ARTIST - $TRACK"
elif [ "$MUSIC_PLAYING" = "playing" ]; then
  TRACK=$(osascript -e 'tell application "Music" to name of current track')
  ARTIST=$(osascript -e 'tell application "Music" to artist of current track')
  sketchybar --set $NAME icon.color=$RED label="$ARTIST - $TRACK"
else
  sketchybar --set $NAME icon.color=$DIM label="--"
fi
