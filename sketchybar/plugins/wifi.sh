#!/bin/bash

# Network - displays connection status

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
  sketchybar --set $NAME \
    icon="NET" \
    icon.color=0xffcb7c94 \
    label="Off"
else
  sketchybar --set $NAME \
    icon="NET" \
    icon.color=0xff7aa89f \
    label="$WIFI"
fi
