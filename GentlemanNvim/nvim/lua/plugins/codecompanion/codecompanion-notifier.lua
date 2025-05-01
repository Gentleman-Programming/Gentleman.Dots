local M = {}

local uv = vim.uv
local api = vim.api

local active = {} -- dictionary of active requests

local S = {
  frames = { "⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏" },
  speed = 80, -- milliseconds per frame
}

local function spinner_frame()
  local time = math.floor(uv.hrtime() / (1e6 * S.speed))
  local idx = time % #S.frames + 1
  local frame = S.frames[idx]

  return frame
end

local function refresh_notifications(key)
  return function()
    local req = active[key]
    if not req then
      return
    end

    vim.notify(req.msg, vim.log.levels.INFO, {
      id = "cc_progress",
      title = req.adapter,
      opts = function(notif)
        local icon = " "
        if not req.done then
          icon = spinner_frame()
        end

        notif.icon = icon
      end,
    })
  end
end

local function request_key(data)
  local adapter = data.adapter or {}
  local name = adapter.formatted_name or adapter.name or "unknown"
  return string.format("%s:%s", name, data.id or "???")
end

local function start(ev)
  local data = ev.data or {}
  local key = request_key(data)
  local adapter = data.adapter and data.adapter.name or "CodeCompanion"
  local refresh = refresh_notifications(key)

  local timer = uv.new_timer()
  local req = {
    adapter = adapter,
    done = false,
    msg = "Thinking...",
    refresh = refresh,
    timer = timer,
  }

  active[key] = req

  timer:start(0, 150, vim.schedule_wrap(refresh))
  refresh()
end

local function finished(ev)
  local data = ev.data or {}
  local key = request_key(data)
  local req = active[key]

  if not req then
    return
  end

  req.done = true

  if data.status == "success" then
    req.msg = "Done."
  elseif data.status == "error" then
    req.msg = "Error!"
  else
    req.msg = "Cancelled."
  end

  req.refresh()

  -- clear the finished request
  active[key] = nil
  if req.timer then
    req.timer:stop()
    req.timer:close()
  end
end

function M.setup()
  local group = vim.api.nvim_create_augroup("CodeCompanionSnacks", { clear = true })

  vim.api.nvim_create_autocmd("User", {
    pattern = "CodeCompanionRequestStarted",
    group = group,
    callback = start,
  })

  vim.api.nvim_create_autocmd("User", {
    pattern = "CodeCompanionRequestFinished",
    group = group,
    callback = finished,
  })
end

function M.init()
  M.setup()
end

return M
