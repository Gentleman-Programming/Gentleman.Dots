#!/bin/bash

# RAM - vm_stat (7ms) instead of memory_pressure

RED=0xffcb7c94
YELLOW=0xffffe066
MAGENTA=0xffff8dd7

RAM=$(vm_stat | awk '/Pages active/ {a=$3} /Pages wired/ {w=$3} /Pages free/ {f=$3} /Pages inactive/ {i=$3} END {used=a+w; total=used+f+i; printf "%d", (used*100/total)}')

if [ "$RAM" -ge 80 ]; then
  COLOR=$RED
elif [ "$RAM" -ge 60 ]; then
  COLOR=$YELLOW
else
  COLOR=$MAGENTA
fi

sketchybar --set $NAME icon.color="$COLOR" label="${RAM}%"
