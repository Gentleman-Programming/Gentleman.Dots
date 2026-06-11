#!/bin/bash

# Required parameters:
# @raycast.schemaVersion 1
# @raycast.title Restart Yabai
# @raycast.mode compact

# Optional parameters:
# @raycast.icon 🔄
# @raycast.packageName System

set -euo pipefail

uid="$(/usr/bin/id -u)"
label="com.koekeishiya.yabai"
domain="gui/${uid}"
plist="$HOME/Library/LaunchAgents/${label}.plist"

if [ ! -f "$plist" ]; then
  echo "Missing yabai LaunchAgent: $plist"
  exit 1
fi

if /bin/launchctl print "${domain}/${label}" >/dev/null 2>&1; then
  /bin/launchctl kickstart -k "${domain}/${label}"
else
  /bin/launchctl bootstrap "$domain" "$plist"
fi

echo "Yabai restarted"
