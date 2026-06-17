#!/bin/bash

# Required parameters:
# @raycast.schemaVersion 1
# @raycast.title Restart Nehir
# @raycast.mode compact

# Optional parameters:
# @raycast.icon 🪟
# @raycast.packageName System

set -euo pipefail

# Nehir is a regular .app (Homebrew cask), not a LaunchAgent — quit and relaunch
# to re-tile windows cleanly when the layout gets stuck.
/usr/bin/osascript -e 'quit app "Nehir"' >/dev/null 2>&1 || true

# Wait for it to exit; force-kill if it lingers.
for _ in $(seq 1 20); do
  /usr/bin/pgrep -x Nehir >/dev/null 2>&1 || break
  sleep 0.1
done
/usr/bin/pkill -x Nehir >/dev/null 2>&1 || true

sleep 0.5
/usr/bin/open -a Nehir

echo "Nehir restarted"
