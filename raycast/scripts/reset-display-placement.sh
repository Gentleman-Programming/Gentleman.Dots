#!/bin/bash

# Required parameters:
# @raycast.schemaVersion 1
# @raycast.title Reset Display Placement
# @raycast.mode compact

# Optional parameters:
# @raycast.icon 🖥️

sleep 1  # Give displays time to settle after resolution change

# Monitor IDs
LG_TV="C41C0C23-D683-42F0-BDD8-453B6D927811"
ELGATO="9FE7C924-AEEE-4338-B16E-F5B29054CC26"
ARZOPA="2CC539C8-E8F3-44C6-925E-AE98DE7F6EFB"

# Get currently connected displays
CONNECTED=$(/opt/homebrew/bin/displayplacer list 2>/dev/null)

# Get current resolution of LG TV
CURRENT_RES=$(echo "$CONNECTED" | grep -A5 "Persistent screen id: $LG_TV" | grep "Resolution:" | awk '{print $2}')

# Determine which config to use based on LG resolution
if [ "$CURRENT_RES" = "1920x1080" ]; then
    # 1080p config
    /opt/homebrew/bin/displayplacer \
        "id:$LG_TV res:1920x1080 hz:120 color_depth:8 enabled:true scaling:on origin:(0,0) degree:0" \
        "id:$ELGATO res:1280x800 hz:60 color_depth:8 enabled:true scaling:on origin:(372,1080) degree:0" \
        "id:$ARZOPA res:1024x600 hz:60 color_depth:4 enabled:true scaling:off origin:(-55,-600) degree:0"
    echo "Applied 1080p arrangement"
elif [ "$CURRENT_RES" = "3008x1692" ]; then
    # 4K/More Space config
    /opt/homebrew/bin/displayplacer \
        "id:$LG_TV res:3008x1692 hz:120 color_depth:8 enabled:true scaling:on origin:(0,0) degree:0" \
        "id:$ELGATO res:1280x800 hz:60 color_depth:8 enabled:true scaling:on origin:(704,1692) degree:0" \
        "id:$ARZOPA res:1024x600 hz:60 color_depth:4 enabled:true scaling:off origin:(0,-600) degree:0"
    echo "Applied 4K arrangement"
else
    echo "Unknown resolution: $CURRENT_RES"
fi

# Reset Raycast window position
sleep 2
open raycast://extensions/raycast/raycast/reset-window-position
