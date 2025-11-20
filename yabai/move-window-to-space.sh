#!/bin/bash
# Move window to space and follow
# Usage: ./move-window-to-space.sh <space_number>

SPACE=$1

# Move window using yabai
yabai -m window --space "$SPACE"

# Switch to that space using Mission Control shortcut
skhd -k "ctrl - $SPACE"
