#!/bin/bash

# CPU - displays current CPU usage percentage



CPU=$(top -l 1 -n 0 | grep "CPU usage" | awk '{print int($3)}')

sketchybar --set $NAME label="${CPU}%"
