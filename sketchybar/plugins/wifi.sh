#!/bin/bash

# Network - displays connection status

# CONFIG_DIR is supplied by SketchyBar at runtime.
# shellcheck disable=SC1091
source "$CONFIG_DIR/theme.sh"

# Try to get WiFi SSID
WIFI=$(/System/Library/PrivateFrameworks/Apple80211.framework/Resources/airport -I 2>/dev/null | awk -F': ' '/^ *SSID/{print $2}')

# If no SSID, check if we have an IP (connected via ethernet or hidden SSID)
if [ -z "$WIFI" ]; then
  IP=$(ipconfig getifaddr en0 2>/dev/null)
  if [ -n "$IP" ]; then
    WIFI="Online"
  fi
fi

if [ -z "$WIFI" ]; then
  sketchybar --set "$NAME" \
    icon="NET" \
    icon.color="$ERROR" \
    label="Off"
else
  sketchybar --set "$NAME" \
    icon="NET" \
    icon.color="$NETWORK_NORMAL" \
    label="$WIFI"
fi
