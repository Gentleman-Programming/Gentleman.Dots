#!/bin/bash

# GitHub - displays unread notification count



COUNT=$(gh api notifications 2>/dev/null | jq 'length' 2>/dev/null)

if [ -z "$COUNT" ] || [ "$COUNT" = "null" ]; then
  COUNT=0
fi

if [ "$COUNT" -gt 0 ]; then
  sketchybar --set $NAME \
    icon="GH" \
    icon.color=0xffffe066 \
    label="$COUNT"
else
  sketchybar --set $NAME \
    icon="GH" \
    icon.color=0xff565f89 \
    label="0"
fi
