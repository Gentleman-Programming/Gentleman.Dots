#!/usr/bin/env bash
# hook-adapter — shared adapter for claude-code-style hook systems (claude, codex).
# The hook registration passes the event name as $1; the hook's JSON payload comes
# on stdin (carries tool_name for tool events). We map both to the canonical
# vocabulary and forward to the normalization core.
#
# usage (from a hook command): bash hook-adapter.sh <EventName>
set -uo pipefail

event="${1:-}"
pane="${TMUX_PANE:-}"
input="$(cat 2>/dev/null || true)"   # capture the hook's JSON stdin

[ -n "$pane" ] && [ -n "$event" ] || exit 0
[ "${HERDR_ENV:-}" = "1" ] && exit 0   # under herdr, let herdr's integration own it

# Tools that mean "the agent is waiting on the user" (asking a question).
# PreToolUse on these = blocked, not working. Extend as needed per agent.
is_blocking_tool() {
  case "$1" in
    AskUserQuestion|ask_user_question|request_user_input|ExitPlanMode) return 0 ;;
    *) return 1 ;;
  esac
}

tool=""
msg=""
case "$event" in
  PreToolUse|PostToolUse)
    tool="$(printf '%s' "$input" | python3 -c 'import sys,json
try: print(json.load(sys.stdin).get("tool_name","") or json.load(sys.stdin).get("toolName",""))
except Exception: print("")' 2>/dev/null)"
    ;;
  Notification)
    msg="$(printf '%s' "$input" | python3 -c 'import sys,json
try: print(json.load(sys.stdin).get("message","") or "")
except Exception: print("")' 2>/dev/null)"
    ;;
esac

case "$event" in
  PreToolUse)
    if is_blocking_tool "$tool"; then state=blocked; else state=working; fi ;;
  PostToolUse)
    state=working ;;   # tool finished (incl. question answered) → agent continues
  UserPromptSubmit|PreCompact)
    state=working ;;
  Notification)
    # Claude fires Notification for two unrelated cases. Only the permission case
    # is a real "needs you"; the idle "waiting for your input" one is NOT blocked
    # (it fires ~60s after Stop and would otherwise re-stick the pane to blocked).
    case "$msg" in
      *permission*|*"needs your"*|*approve*) state=blocked ;;
      *) state=idle; msg="" ;;
    esac ;;
  Stop|SessionStart|SessionEnd|TurnComplete)
    state=idle ;;
  *) exit 0 ;;   # SubagentStop and unknown events: ignore
esac

exec bash "$HOME/.config/tmux/scripts/agent-report.sh" "$pane" "$state" "$msg"
