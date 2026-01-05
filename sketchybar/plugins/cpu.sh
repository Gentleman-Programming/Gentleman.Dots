#!/bin/bash

# CPU - displays current CPU usage percentage with color based on load

RED=0xffcb7c94
YELLOW=0xffffe066
CYAN=0xff7aa89f

CPU=$(top -l 1 -n 0 | grep "CPU usage" | awk '{print int($3)}')

# Color based on usage
if [ "$CPU" -ge 80 ]; then
  COLOR=$RED
elif [ "$CPU" -ge 50 ]; then
  COLOR=$YELLOW
else
  COLOR=$CYAN
fi

sketchybar --set $NAME icon.color="$COLOR" label="${CPU}%"
