return {
  "sudo-tee/opencode.nvim",
  config = function()
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
        },
      },
      ui = {
        floating = true, -- Use floating windows for input and output
        window_width = 0.40, -- Width as percentage of editor width
        input_height = 0.15, -- Input height as percentage of window height
        fullscreen = false, -- Start in fullscreen mode (default: false)
        layout = "center", -- Options: "center" or "right"
        floating_height = 0.8, -- Height as percentage of editor height for "center" layout
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
