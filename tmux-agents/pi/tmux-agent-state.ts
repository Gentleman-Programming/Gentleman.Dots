// tmux-agent-state — pi adapter for the tmux agent-state notifier.
// Translates pi's native lifecycle events into the canonical vocabulary and
// forwards them to the single normalization core (~/.config/tmux/scripts/agent-report.sh).
// Inspired by herdr's pi extension, but targets tmux ($TMUX_PANE) instead of the herdr socket.
//
// Active only inside a tmux pane (and NOT under herdr).
// @ts-nocheck

import { spawn } from "node:child_process";
import { homedir } from "node:os";
import { join } from "node:path";

const REPORT = join(homedir(), ".config", "tmux", "scripts", "agent-report.sh");
const pane = process.env.TMUX_PANE;
const enabled = Boolean(pane) && process.env.HERDR_ENV !== "1";

function report(state) {
  if (!pane) return;
  const child = spawn("bash", [REPORT, pane, state], { stdio: "ignore" });
  child.on("error", () => {});
}

export default function (pi) {
  if (!enabled) return;

  let active = false;
  let idleTimer;

  const clearIdle = () => {
    if (idleTimer) clearTimeout(idleTimer);
    idleTimer = undefined;
  };

  pi.on?.("session_start", () => report("idle"));

  pi.on?.("agent_start", () => {
    active = true;
    clearIdle();
    report("working");
  });

  pi.on?.("agent_end", () => {
    active = false;
    clearIdle();
    // debounce so provider retries / quick back-to-back turns don't flash idle
    idleTimer = setTimeout(() => report("idle"), 250);
    idleTimer.unref?.();
  });

  // blocked = pi is waiting on the user. pi has NO native permission/question
  // event, but asking the user is implemented as a TOOL (ask_user_question).
  // So: that tool's execution start = waiting on you; its end = answered.
  const BLOCKING_TOOLS = new Set(["ask_user_question"]);

  pi.on?.("tool_execution_start", (ev) => {
    if (BLOCKING_TOOLS.has(ev?.toolName)) {
      clearIdle();
      report("blocked");
    }
  });

  pi.on?.("tool_execution_end", (ev) => {
    if (BLOCKING_TOOLS.has(ev?.toolName)) {
      report("working"); // answer received, agent keeps going
    }
  });

  pi.on?.("session_shutdown", () => {
    clearIdle();
    report("idle");
  });
}
