#!/bin/bash

# RAM - vm_stat (7ms) instead of memory_pressure

# CONFIG_DIR is supplied by SketchyBar at runtime.
# shellcheck disable=SC1091
source "$CONFIG_DIR/theme.sh"

RAM=$(vm_stat | awk '/Pages active/ {a=$3} /Pages wired/ {w=$3} /Pages free/ {f=$3} /Pages inactive/ {i=$3} END {used=a+w; total=used+f+i; printf "%d", (used*100/total)}')

if [ "$RAM" -ge 80 ]; then
  COLOR=$ERROR
elif [ "$RAM" -ge 60 ]; then
  COLOR=$WARNING
else
  COLOR=$RAM_NORMAL
fi

sketchybar --set "$NAME" icon.color="$COLOR" label="${RAM}%"
