// tmux-agent-state — opencode adapter for the tmux agent-state notifier.
// Translates opencode's native events into the canonical vocabulary and forwards
// them to the single normalization core (~/.config/tmux/scripts/agent-report.sh).
// Inspired by herdr's plugin, but targets tmux ($TMUX_PANE) instead of the herdr socket.
//
// Active only when running inside a tmux pane (and NOT under herdr).

import { spawn } from "node:child_process";
import { homedir } from "node:os";
import { join } from "node:path";

const REPORT = join(homedir(), ".config", "tmux", "scripts", "agent-report.sh");

function enabled() {
  // tmux sets TMUX_PANE per pane. If we're under herdr, let herdr's plugin own it.
  return Boolean(process.env.TMUX_PANE) && process.env.HERDR_ENV !== "1";
}

function report(state, message) {
  const pane = process.env.TMUX_PANE;
  if (!pane) return Promise.resolve();
  return new Promise((resolve) => {
    const child = spawn("bash", [REPORT, pane, state, message ?? ""], {
      stdio: "ignore",
      detached: false,
    });
    child.on("error", () => resolve());
    child.on("close", () => resolve());
  });
}

function stateFromSessionStatus(status) {
  if (typeof status !== "string") return undefined;
  switch (status.toLowerCase()) {
    case "idle":
      return "idle";
    case "active":
    case "busy":
    case "pending":
    case "running":
    case "streaming":
    case "working":
      return "working";
    default:
      return undefined;
  }
}

export const TmuxAgentStatePlugin = async () => {
  if (!enabled()) return {};

  return {
    "chat.message": async () => {
      await report("working");
    },
    event: async ({ event }) => {
      const type = event?.type;
      const properties = event?.properties ?? {};

      switch (type) {
        case "session.status": {
          const state = stateFromSessionStatus(properties.status);
          if (state) await report(state);
          break;
        }
        case "tool.execute.before":
        case "tool.execute.after":
        case "permission.replied":
        case "question.replied":
        case "question.rejected":
        case "session.compacted":
          await report("working");
          break;
        case "permission.asked":
        case "question.asked":
        case "session.error":
          await report("blocked");
          break;
        case "session.idle":
          await report("idle");
          break;
        case "session.deleted":
          await report("idle");
          break;
        default:
          break;
      }
    },
  };
};
