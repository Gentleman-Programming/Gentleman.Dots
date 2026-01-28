#!/bin/bash

# Meeting - shows next calendar event (next 24 hours)
# IMPORTANT: Uses pgrep to check if Calendar is running (osascript can launch apps)

YELLOW=0xffffe066
DIM=0xff565f89

# Check if Calendar is running using pgrep (safe, won't launch the app)
if ! pgrep -x "Calendar" > /dev/null 2>&1; then
  sketchybar --set $NAME icon.color=$DIM label="--"
  exit 0
fi

# Get current hour and minute
CURRENT_TIME=$(date +%H%M)

# Get all events in next 24 hours (only if Calendar is running)
NEXT_EVENT=$(osascript -e '
tell application "Calendar"
    set now to current date
    set tomorrow to now + (24 * 60 * 60)

    set eventList to ""

    repeat with cal in calendars
        try
            set theEvents to (every event of cal whose start date >= now and start date <= tomorrow)
            repeat with evt in theEvents
                set evtStart to start date of evt
                set h to hours of evtStart
                set m to minutes of evtStart
                set hStr to h as text
                set mStr to m as text
                if h < 10 then set hStr to "0" & hStr
                if m < 10 then set mStr to "0" & mStr
                set eventList to eventList & hStr & mStr & "|" & hStr & ":" & mStr & " " & (summary of evt) & "\n"
            end repeat
        end try
    end repeat

    return eventList
end tell
' 2>/dev/null)

# Filter events that are in the future and get the earliest one
FILTERED=$(echo "$NEXT_EVENT" | while IFS='|' read -r time_num display; do
  if [ -n "$time_num" ] && [ "$time_num" -gt "$CURRENT_TIME" ]; then
    echo "$time_num|$display"
  fi
done | sort -n | head -1 | cut -d'|' -f2)

if [ -z "$FILTERED" ]; then
  sketchybar --set $NAME icon.color=$DIM label="--"
else
  sketchybar --set $NAME icon.color=$YELLOW label="$FILTERED"
fi
