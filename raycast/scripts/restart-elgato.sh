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

# Kill logic lives in a temp helper so the privileged shell can exclude our
# own PID. Without this, `pkill -fi elgato` matches the script's path and
# SIGKILLs the script before the audio/video subsystems get restarted.
SELF_PID=$$
KILL_SCRIPT=$(mktemp -t restart-stack.XXXXXX)
trap 'rm -f "$KILL_SCRIPT"' EXIT

cat > "$KILL_SCRIPT" <<'KILLEOF'
#!/bin/bash
EXCLUDE_PID="$1"
for pattern in elgato "wave link" wavelink "stream deck" streamdeck "camera hub" camerahub crashpad_handler; do
  pids=$(pgrep -fi "$pattern" 2>/dev/null | grep -v "^${EXCLUDE_PID}$")
  [ -n "$pids" ] && echo "$pids" | xargs kill -9 2>/dev/null
done
# coreaudiod is SIP-protected on macOS Tahoe (26.x): `killall` is silently
# ignored and the same PID survives, so audio stays wedged. launchctl kickstart
# goes through launchd (which can restart it) and is what actually unwedges the
# stack when Wave Link / Camera Hub / OBS hang at launch with no window.
launchctl kickstart -k system/com.apple.audio.coreaudiod 2>/dev/null
killall VDCAssistant 2>/dev/null
killall AppleCameraAssistant 2>/dev/null
true
KILLEOF
chmod +x "$KILL_SCRIPT"

# osascript opens the native admin prompt (cachea pass ~5min)
osascript -e "do shell script \"'$KILL_SCRIPT' $SELF_PID\" with administrator privileges" >/dev/null 2>&1

sleep 2

# Relaunch apps that were running
for app in "${RUNNING_APPS[@]}"; do
  open -a "$app" 2>/dev/null
done

echo "Elgato stack restarted (${#RUNNING_APPS[@]} apps relaunched)"
