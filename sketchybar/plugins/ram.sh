#!/bin/bash

# RAM - displays current memory usage percentage with color based on load

RED=0xffcb7c94
YELLOW=0xffffe066
MAGENTA=0xffff8dd7

MEMORY_PRESSURE=$(memory_pressure | grep "System-wide memory free percentage:" | awk '{print 100-$5}')

# Color based on usage
if [ "$MEMORY_PRESSURE" -ge 80 ]; then
  COLOR=$RED
elif [ "$MEMORY_PRESSURE" -ge 60 ]; then
  COLOR=$YELLOW
else
  COLOR=$MAGENTA
fi

sketchybar --set $NAME icon.color="$COLOR" label="${MEMORY_PRESSURE}%"
