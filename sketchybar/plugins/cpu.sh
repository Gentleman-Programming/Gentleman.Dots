#!/bin/bash

# CPU - ps-based (19ms) instead of top (500ms)

# CONFIG_DIR is supplied by SketchyBar at runtime.
# shellcheck disable=SC1091
source "$CONFIG_DIR/theme.sh"

NCPU=$(sysctl -n hw.ncpu)
CPU=$(ps -A -o %cpu | awk -v n="$NCPU" '{s+=$1} END {v=s/n; if(v>100)v=100; printf "%d",v}')

if [ "$CPU" -ge 80 ]; then
  COLOR=$ERROR
elif [ "$CPU" -ge 50 ]; then
  COLOR=$WARNING
else
  COLOR=$CPU_NORMAL
fi

sketchybar --set "$NAME" icon.color="$COLOR" label="${CPU}%"
