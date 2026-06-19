#!/usr/bin/env bash
# agent-status — self-heal heartbeat (runs every status-interval, prints NOTHING).
#
# It lives in status-right only so tmux re-runs it on each refresh; the agent
# informativo is shown by the per-tab dots (window-status-format), so this script
# stays visually silent and the theme keeps the status-right (RAM/CPU/git).
#
# Self-heal: clears "blocked" on any pane the user is CURRENTLY VIEWING. No hook
# fires when an agent's question/permission is cancelled, so the blocked state would
# otherwise stick forever — once a blocked pane is on-screen, the alert is delivered.
set -uo pipefail
command -v tmux >/dev/null 2>&1 || exit 0

recompute_rollup() {
  local win="$1" worst=idle s
  while read -r s; do
    case "$s" in
      blocked) worst=blocked; break ;;
      working) [ "$worst" = idle ] && worst=working ;;
    esac
  done < <(tmux list-panes -t "$win" -F '#{@agent_state}' 2>/dev/null)
  tmux set -w -t "$win" @win_agent_state "$worst" 2>/dev/null
}

# --- self-heal: clear blocked on the pane(s) on screen right now ---
while IFS=$'\t' read -r pid st vis wid; do
  if [ "$st" = "blocked" ] && [ "$vis" = "1" ]; then
    tmux set -p -t "$pid" @agent_state idle 2>/dev/null
    recompute_rollup "$wid"
  fi
# note: substitute empty @agent_state with "-" so empty fields don't collapse
# (tmux IFS tab-splitting merges consecutive tabs).
done < <(tmux list-panes -a -F '#{pane_id}'$'\t''#{?#{@agent_state},#{@agent_state},-}'$'\t''#{&&:#{pane_active},#{&&:#{window_active},#{session_attached}}}'$'\t''#{window_id}' 2>/dev/null)

# No output: the per-tab dots are the informativo; the theme keeps the status-right.
exit 0
