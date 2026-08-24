#!/bin/bash

THEME_NAME="gentleman"
MARKER_FILE="$HOME/.config/gentleman/sketchybar-theme"

if [ -r "$MARKER_FILE" ]; then
  REQUESTED_THEME=$(cat "$MARKER_FILE" 2>/dev/null || true)
  case "$REQUESTED_THEME" in
    gentleman|gentleman-cute)
      THEME_NAME="$REQUESTED_THEME"
      ;;
  esac
fi

# THEME_NAME is selected only from the whitelist above.
# shellcheck disable=SC1090
source "$CONFIG_DIR/themes/$THEME_NAME.sh"
