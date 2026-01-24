# Raycast Scripts

Custom Raycast scripts for macOS automation.

## Setup

1. Open Raycast Settings → Extensions → Script Commands
2. Add the `scripts/` directory from this folder
3. Scripts will be available in Raycast search

## Scripts

| Script | Description |
|--------|-------------|
| `set-4k.sh` | Set LG TV to 4K resolution with multi-monitor arrangement |
| `set-1080p.sh` | Set LG TV to 1080p resolution with multi-monitor arrangement |
| `reset-display-placement.sh` | Auto-detect current resolution and reset display arrangement |
| `restart-sketchybar.sh` | Kill and restart sketchybar |

## Requirements

- [displayplacer](https://github.com/jakehilborn/displayplacer) - `brew install displayplacer`
- Raycast with Script Commands enabled

## Note

The display scripts use hardcoded monitor IDs. Run `displayplacer list` to get your monitor IDs and update the scripts accordingly.
