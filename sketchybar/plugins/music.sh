#!/bin/bash

# Music - shows currently playing track (Spotify/Apple Music)
# IMPORTANT: Uses pgrep to check if apps are running (osascript can launch apps)

GREEN=0xffb7cc85
RED=0xffcb7c94
DIM=0xff565f89

# Check if Spotify is running using pgrep (safe, won't launch the app)
if pgrep -x "Spotify" > /dev/null 2>&1; then
  SPOTIFY_PLAYING=$(osascript -e 'tell application "Spotify" to player state' 2>/dev/null)
  if [ "$SPOTIFY_PLAYING" = "playing" ]; then
    TRACK=$(osascript -e 'tell application "Spotify" to name of current track' 2>/dev/null)
    ARTIST=$(osascript -e 'tell application "Spotify" to artist of current track' 2>/dev/null)
    sketchybar --set $NAME icon.color=$GREEN label="$ARTIST - $TRACK"
    exit 0
  fi
fi

# Check if Music is running using pgrep (safe, won't launch the app)
if pgrep -x "Music" > /dev/null 2>&1; then
  MUSIC_PLAYING=$(osascript -e 'tell application "Music" to player state' 2>/dev/null)
  if [ "$MUSIC_PLAYING" = "playing" ]; then
    TRACK=$(osascript -e 'tell application "Music" to name of current track' 2>/dev/null)
    ARTIST=$(osascript -e 'tell application "Music" to artist of current track' 2>/dev/null)
    sketchybar --set $NAME icon.color=$RED label="$ARTIST - $TRACK"
    exit 0
  fi
fi

# Nothing playing
sketchybar --set $NAME icon.color=$DIM label="--"
