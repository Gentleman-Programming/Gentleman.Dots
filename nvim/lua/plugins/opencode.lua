return {
  "sudo-tee/opencode.nvim",
  config = function()
    -- Track opencode's internal state during resize
    local in_resize = false
    local original_cursor_win = nil
    local opencode_filetypes = { "opencode_input", "opencode_output", "opencode_chat" }

    -- Temporarily move cursor away from opencode during resize
    local function temporarily_leave_opencode()
      local is_opencode, opencode_win
      if is_opencode and not in_resize then
        in_resize = true
        original_cursor_win = opencode_win

        -- Find a non-opencode window to switch to
        local target_win = nil
        for _, win in ipairs(vim.api.nvim_list_wins()) do
          local buf = vim.api.nvim_win_get_buf(win)
          local ft = vim.bo[buf].filetype

          local is_opencode_ft = false
          for _, oft in ipairs(opencode_filetypes) do
            if ft == oft then
              is_opencode_ft = true
              break
            end
          end

          if not is_opencode_ft and vim.api.nvim_win_is_valid(win) then
            target_win = win
            break
          end
        end

        -- Switch to non-opencode window if found
        if target_win then
          vim.api.nvim_set_current_win(target_win)
          return true
        end
      end
      return false
    end

    -- Restore cursor to original opencode window
    local function restore_cursor_to_opencode()
      if in_resize and original_cursor_win and vim.api.nvim_win_is_valid(original_cursor_win) then
        -- Small delay to ensure resize is complete
        vim.defer_fn(function()
          pcall(vim.api.nvim_set_current_win, original_cursor_win)
          in_resize = false
          original_cursor_win = nil
        end, 50)
      end
    end

    -- Prevent duplicate windows cleanup
    local function cleanup_duplicate_opencode_windows()
      local seen_filetypes = {}
      local windows_to_close = {}

      for _, win in ipairs(vim.api.nvim_list_wins()) do
        local buf = vim.api.nvim_win_get_buf(win)
        local ft = vim.bo[buf].filetype

        -- Special handling for opencode panels
        for _, opencode_ft in ipairs(opencode_filetypes) do
          if ft == opencode_ft then
            if seen_filetypes[ft] then
              -- Found duplicate, mark for closing
              table.insert(windows_to_close, win)
            else
              seen_filetypes[ft] = win
            end
            break
          end
        end
      end

      -- Close duplicate windows
      for _, win in ipairs(windows_to_close) do
        if vim.api.nvim_win_is_valid(win) then
          pcall(vim.api.nvim_win_close, win, true)
        end
      end
    end

    -- Create autocmd group for resize fix
    vim.api.nvim_create_augroup("OpencodeResizeFix", { clear = true })

    -- Main resize handler for Resize
    vim.api.nvim_create_autocmd({ "VimResized" }, {
      group = "OpencodeResizeFix",
      callback = function()
        -- Move cursor away from opencode before resize processing
        local moved = temporarily_leave_opencode()

        if moved then
          -- Let resize happen, then restore cursor
          vim.defer_fn(function()
            restore_cursor_to_opencode()
            -- Force a clean redraw
            vim.cmd("redraw!")
          end, 100)
        end

        -- Cleanup duplicates after resize completes
        vim.defer_fn(cleanup_duplicate_opencode_windows, 150)
      end,
    })

    -- Prevent opencode from responding to scroll/resize events during resize
    vim.api.nvim_create_autocmd({ "WinScrolled", "WinResized" }, {
      group = "OpencodeResizeFix",
      pattern = "*",
      callback = function(args)
        local buf = args.buf
        if buf and vim.api.nvim_buf_is_valid(buf) then
          local ft = vim.bo[buf].filetype

          for _, opencode_ft in ipairs(opencode_filetypes) do
            if ft == opencode_ft then
              -- Prevent event propagation for opencode buffers during resize
              if in_resize then
                return true -- This should stop the event
              end
              break
            end
          end
        end
      end,
    })

    -- Additional cleanup on focus events
    vim.api.nvim_create_autocmd("FocusGained", {
      group = "OpencodeResizeFix",
      callback = function()
        -- Reset resize state on focus gain
        in_resize = false
        original_cursor_win = nil
        -- Clean up any duplicate windows
        vim.defer_fn(cleanup_duplicate_opencode_windows, 100)
      end,
    })

    require("opencode").setup({
      keymap = {
        global = {
          toggle = "<leader>aa", -- Open opencode. Close if opened
          open_input = "<leader>ai", -- Opens and focuses on input window on insert mode
          open_input_new_session = "<leader>aI", -- Opens and focuses on input window on insert mode. Creates a new session
          open_output = "<leader>ao", -- Opens and focuses on output window
          toggle_focus = "<leader>at", -- Toggle focus between opencode and last window
          close = "<leader>aq", -- Close UI windows
          toggle_fullscreen = "<leader>af", -- Toggle between normal and fullscreen mode
          select_session = "<leader>as", -- Select and load a opencode session
          configure_provider = "<leader>ap", -- Quick provider and model switch from predefined list
          diff_open = "<leader>ad", -- Opens a diff tab of a modified file since the last opencode prompt
          diff_next = "<leader>a]", -- Navigate to next file diff
          diff_prev = "<leader>a[", -- Navigate to previous file diff
          diff_close = "<leader>ac", -- Close diff view tab and return to normal editing
          diff_revert_all_last_prompt = "<leader>ara", -- Revert all file changes since the last opencode prompt
          diff_revert_this_last_prompt = "<leader>art", -- Revert current file changes since the last opencode prompt
          diff_revert_all = "<leader>arA", -- Revert all file changes since the last opencode session
          diff_revert_this = "<leader>arT", -- Revert current file changes since the last opencode session
          switch_position = "<leader>ax", -- Switch between positions
        },
      },
      ui = {
        fullscreen = false, -- Start in fullscreen mode (default: false)
        position = "left",
        display_model = true, -- Display model name on top winbar
        window_highlight = "Normal:OpencodeBackground,FloatBorder:OpencodeBorder", -- Highlight group for the opencode window
        output = {
          tools = {
            show_output = true, -- Show tools output [diffs, cmd output, etc.] (default: true)
          },
        },
      },
    })
  end,
  dependencies = {
    "nvim-lua/plenary.nvim",
    {
      "MeanderingProgrammer/render-markdown.nvim",
      opts = {
        anti_conceal = { enabled = false },
        file_types = { "markdown", "opencode_output" },
      },
      ft = { "markdown", "Avante", "copilot-chat", "opencode_output" },
    },
  },
}
