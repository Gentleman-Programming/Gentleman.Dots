#!/bin/bash

# Gentleman theme colors (ANSI 256)
PRIMARY='\033[38;5;110m'      # #7FB4CA azul claro
ACCENT='\033[38;5;179m'       # #E0C15A dorado
SECONDARY='\033[38;5;146m'    # #A3B5D6 azul gris
MUTED='\033[38;5;242m'        # #5C6170 gris
SUCCESS='\033[38;5;150m'      # #B7CC85 verde
ERROR='\033[38;5;174m'        # #CB7C94 rosa/rojo
PURPLE='\033[38;5;183m'       # #C99AD6 púrpura
BOLD='\033[1m'
STRIKE='\033[9m'
NC='\033[0m'

# Cache for MCP (don't call every 300ms)
MCP_CACHE_FILE="/tmp/claude_mcp_cache"
MCP_CACHE_TTL=120  # 2 minutes

# Read JSON from stdin
input=$(cat)

# Parse basic fields
MODEL=$(echo "$input" | jq -r '.model.display_name // "Claude"')
DIR=$(echo "$input" | jq -r '.workspace.current_dir // "~"')
ADDED=$(echo "$input" | jq -r '.cost.total_lines_added // 0')
REMOVED=$(echo "$input" | jq -r '.cost.total_lines_removed // 0')

# Context window (before compression)
CTX_SIZE=$(echo "$input" | jq -r '.context_window.context_window_size // 200000')
INPUT_TOKENS=$(echo "$input" | jq -r '.context_window.current_usage.input_tokens // 0')
CACHE_CREATE=$(echo "$input" | jq -r '.context_window.current_usage.cache_creation_input_tokens // 0')
CACHE_READ=$(echo "$input" | jq -r '.context_window.current_usage.cache_read_input_tokens // 0')

TOTAL_USED=$((INPUT_TOKENS + CACHE_CREATE + CACHE_READ))
if [ "$CTX_SIZE" -gt 0 ] 2>/dev/null; then
  CTX_PERCENT=$((TOTAL_USED * 100 / CTX_SIZE))
else
  CTX_PERCENT=0
fi
[ "$CTX_PERCENT" -gt 100 ] && CTX_PERCENT=100
[ "$CTX_PERCENT" -lt 0 ] && CTX_PERCENT=0

# Function to get MCP servers from config
get_mcp_servers() {
  # Check cache first
  if [ -f "$MCP_CACHE_FILE" ]; then
    CACHE_AGE=$(($(date +%s) - $(stat -f %m "$MCP_CACHE_FILE" 2>/dev/null || echo 0)))
    if [ "$CACHE_AGE" -lt "$MCP_CACHE_TTL" ]; then
      cat "$MCP_CACHE_FILE"
      return
    fi
  fi

  # Read MCP servers from ~/.claude.json config
  local CURRENT_DIR
  CURRENT_DIR=$(echo "$input" | jq -r '.workspace.current_dir // ""')

  # Get servers from current project first, fallback to home
  local SERVERS=""

  if [ -n "$CURRENT_DIR" ]; then
    SERVERS=$(jq -r ".projects[\"$CURRENT_DIR\"].mcpServers // {} | keys[]" ~/.claude.json 2>/dev/null | tr '\n' ',' | sed 's/,$//')
  fi

  # Fallback to home config if current project has no MCP
  if [ -z "$SERVERS" ]; then
    SERVERS=$(jq -r '.projects["/Users/alanbuscaglia"].mcpServers // {} | keys[]' ~/.claude.json 2>/dev/null | tr '\n' ',' | sed 's/,$//')
  fi

  local ALL_SERVERS="$SERVERS"

  # Save to cache (all as "configured" - we can't check connection status easily)
  echo "$ALL_SERVERS|" > "$MCP_CACHE_FILE"
  echo "$ALL_SERVERS|"
}

# Get MCP status
MCP_DATA=$(get_mcp_servers)
MCP_CONNECTED=$(echo "$MCP_DATA" | cut -d'|' -f1)
MCP_DISCONNECTED=$(echo "$MCP_DATA" | cut -d'|' -f2)

# Format MCP display
format_mcp() {
  local result=""

  # Connected servers (green)
  if [ -n "$MCP_CONNECTED" ]; then
    IFS=',' read -ra SERVERS <<< "$MCP_CONNECTED"
    for srv in "${SERVERS[@]}"; do
      if [ -n "$result" ]; then
        result+=" "
      fi
      result+="${SUCCESS}${srv}${NC}"
    done
  fi

  # Disconnected servers (red + strikethrough)
  if [ -n "$MCP_DISCONNECTED" ]; then
    IFS=',' read -ra SERVERS <<< "$MCP_DISCONNECTED"
    for srv in "${SERVERS[@]}"; do
      if [ -n "$result" ]; then
        result+=" "
      fi
      result+="${ERROR}${STRIKE}${srv}${NC}"
    done
  fi

  if [ -z "$result" ]; then
    echo "${MUTED}no mcp${NC}"
  else
    echo "$result"
  fi
}

MCP_DISPLAY=$(format_mcp)

# Directory name
DIR_NAME=$(basename "$DIR")

# Git info
BRANCH=""
GIT_DIRTY=""
if git rev-parse --git-dir > /dev/null 2>&1; then
  BRANCH=$(git branch --show-current 2>/dev/null)
  if [[ -n $(git status --porcelain 2>/dev/null) ]]; then
    GIT_DIRTY="*"
  fi
fi

# Model icon
MODEL_ICON="🤖"
case "$MODEL" in
  *Opus*) MODEL_ICON="🎭" ;;
  *Sonnet*) MODEL_ICON="📝" ;;
  *Haiku*) MODEL_ICON="🍃" ;;
esac

# Progress bar
BAR_WIDTH=8
FILLED=$((CTX_PERCENT * BAR_WIDTH / 100))
EMPTY=$((BAR_WIDTH - FILLED))

if [ "$CTX_PERCENT" -ge 80 ]; then
  BAR_COLOR="$ERROR"
elif [ "$CTX_PERCENT" -ge 50 ]; then
  BAR_COLOR="$ACCENT"
else
  BAR_COLOR="$SUCCESS"
fi

BAR="${BAR_COLOR}"
for ((i=0; i<FILLED; i++)); do BAR+="█"; done
BAR+="${MUTED}"
for ((i=0; i<EMPTY; i++)); do BAR+="░"; done
BAR+="${NC}"

# Build status line
SEP="${MUTED}  ${NC}"

LINE="${BOLD}${PURPLE}${MODEL_ICON} ${MODEL}${NC}"
LINE+="${SEP}"
LINE+="${ACCENT}󰉋 ${DIR_NAME}${NC}"

if [ -n "$BRANCH" ]; then
  LINE+="${SEP}"
  LINE+="${SECONDARY} ${BRANCH}${GIT_DIRTY}${NC}"
fi

LINE+="${SEP}"
LINE+="${SUCCESS}+${ADDED}${NC} ${ERROR}-${REMOVED}${NC}"

LINE+="${SEP}"
LINE+="${MUTED}ctx${NC} ${BAR} ${MUTED}${CTX_PERCENT}%${NC}"

echo -e "$LINE"
