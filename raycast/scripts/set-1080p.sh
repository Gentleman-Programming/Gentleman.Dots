#!/bin/bash

# Required parameters:
# @raycast.schemaVersion 1
# @raycast.title Set 1080p
# @raycast.mode compact

# Optional parameters:
# @raycast.icon 🖥️

# Monitor IDs
LG_TV="C41C0C23-D683-42F0-BDD8-453B6D927811"
ELGATO="9FE7C924-AEEE-4338-B16E-F5B29054CC26"
ARZOPA="2CC539C8-E8F3-44C6-925E-AE98DE7F6EFB"

# Set 1080p with correct arrangement
/opt/homebrew/bin/displayplacer \
    "id:$LG_TV res:1920x1080 hz:120 color_depth:8 enabled:true scaling:on origin:(0,0) degree:0" \
    "id:$ELGATO res:1280x800 hz:60 color_depth:8 enabled:true scaling:on origin:(372,1080) degree:0" \
    "id:$ARZOPA res:1024x600 hz:60 color_depth:4 enabled:true scaling:off origin:(-55,-600) degree:0"

# Reset Raycast window position (use -g to not focus Raycast)
sleep 2
open -g raycast://extensions/raycast/raycast/reset-raycast-window-position

echo "1080p applied"
