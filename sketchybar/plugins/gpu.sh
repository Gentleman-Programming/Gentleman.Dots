#!/bin/bash

# GPU - displays GPU usage for Apple Silicon (M1/M2/M3/M4)
# Uses IOKit to get Device Utilization % from the GPU driver

# CONFIG_DIR is supplied by SketchyBar at runtime.
# shellcheck disable=SC1091
source "$CONFIG_DIR/theme.sh"

# Get GPU utilization from IOAccelerator
GPU=$(ioreg -r -d 1 -c IOAccelerator 2>/dev/null | grep -o '"Device Utilization %"=[0-9]*' | awk -F'=' '{print $2}' | head -1)

# Fallback if empty
if [ -z "$GPU" ]; then
  GPU="--"
fi

# Color based on usage
if [ "$GPU" = "--" ]; then
  COLOR=$GPU_NORMAL
elif [ "$GPU" -ge 80 ]; then
  COLOR=$ERROR
elif [ "$GPU" -ge 50 ]; then
  COLOR=$WARNING
else
  COLOR=$GPU_NORMAL
fi

if [ "$GPU" = "--" ]; then
  sketchybar --set "$NAME" icon.color="$GPU_NORMAL" label="--"
else
  sketchybar --set "$NAME" icon.color="$COLOR" label="${GPU}%"
fi
