#!/bin/bash

# Required parameters:
# @raycast.schemaVersion 1
# @raycast.title Restart Elgato
# @raycast.mode compact

# Optional parameters:
# @raycast.icon 🎙️
# @raycast.packageName System

# Apps to detect/relaunch
ELGATO_APPS=(
  "Elgato Stream Deck"
  "Elgato Wave Link"
  "Elgato Camera Hub"
  "Elgato Control Center"
)

# Detect running apps before killing
RUNNING_APPS=()
for app in "${ELGATO_APPS[@]}"; do
  pgrep -fi "$app" >/dev/null 2>&1 && RUNNING_APPS+=("$app")
done

# Kill -9 nuclear + restart audio/video subsystems via admin dialog
# (osascript abre el prompt nativo de macOS, cachea pass ~5min)
osascript -e 'do shell script "pkill -9 -fi elgato 2>/dev/null; pkill -9 -fi \"wave link\" 2>/dev/null; pkill -9 -fi wavelink 2>/dev/null; pkill -9 -fi \"stream deck\" 2>/dev/null; pkill -9 -fi streamdeck 2>/dev/null; pkill -9 -fi \"camera hub\" 2>/dev/null; pkill -9 -fi camerahub 2>/dev/null; pkill -9 -fi crashpad_handler 2>/dev/null; killall coreaudiod 2>/dev/null; killall VDCAssistant 2>/dev/null; killall AppleCameraAssistant 2>/dev/null; true" with administrator privileges' >/dev/null 2>&1

sleep 2

# Relaunch apps that were running
for app in "${RUNNING_APPS[@]}"; do
  open -a "$app" 2>/dev/null
done

echo "Elgato stack restarted (${#RUNNING_APPS[@]} apps relaunched)"
