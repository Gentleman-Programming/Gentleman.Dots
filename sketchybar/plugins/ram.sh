#!/bin/bash

# RAM - displays current memory usage percentage



MEMORY_PRESSURE=$(memory_pressure | grep "System-wide memory free percentage:" | awk '{print 100-$5}')

sketchybar --set $NAME label="${MEMORY_PRESSURE}%"
