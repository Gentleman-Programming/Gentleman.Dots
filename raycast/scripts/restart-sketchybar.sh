#!/bin/bash

# Required parameters:
# @raycast.schemaVersion 1
# @raycast.title Restart SketchyBar
# @raycast.mode silent

# Optional parameters:
# @raycast.icon 🔄
# @raycast.packageName System

export PATH="$HOME/.local/state/nix/profiles/home-manager/home-path/bin:$PATH"
pkill sketchybar
sleep 1
sketchybar &

echo "Sketchybar restarted"
