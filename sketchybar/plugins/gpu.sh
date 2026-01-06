#!/bin/bash

# GPU - displays GPU usage for Apple Silicon (M1/M2/M3/M4)
# Uses IOKit to get Device Utilization % from the GPU driver

RED=0xffcb7c94
YELLOW=0xffffe066
ORANGE=0xfffff7b1

# Get GPU utilization from IOAccelerator
GPU=$(ioreg -r -d 1 -c IOAccelerator 2>/dev/null | grep -o '"Device Utilization %"=[0-9]*' | awk -F'=' '{print $2}' | head -1)

# Fallback if empty
if [ -z "$GPU" ]; then
  GPU="--"
fi

# Color based on usage
if [ "$GPU" = "--" ]; then
  COLOR=$ORANGE
elif [ "$GPU" -ge 80 ]; then
  COLOR=$RED
elif [ "$GPU" -ge 50 ]; then
  COLOR=$YELLOW
else
  COLOR=$ORANGE
fi

if [ "$GPU" = "--" ]; then
  sketchybar --set $NAME icon.color="$ORANGE" label="--"
else
  sketchybar --set $NAME icon.color="$COLOR" label="${GPU}%"
fi
