return {
  {
    "yetone/avante.nvim",
    event = "VeryLazy",
    lazy = false,
    version = false, -- set this if you want to always pull the latest change
    opts = {
      ---@alias Provider "claude" | "openai" | "azure" | "gemini" | "cohere" | "copilot" | string
      provider = "copilot", -- Recommend using Claude
      copilot = {
        model = "claude-3.7-sonnet", -- o1-preview | o1-mini | claude-3.7-sonnet
      },
      cursor_applying_provider = "copilot", -- In this example, use Groq for applying, but you can also use any provider you want.
      auto_suggestions_provider = "copilot", -- Since auto-suggestions are a high-frequency operation and therefore expensive, it is recommended to specify an inexpensive provider or even a free provider: copilot
      behaviour = {
        auto_suggestions = true, -- Experimental stage
        auto_set_highlight_group = true,
        auto_set_keymaps = true,
        auto_apply_diff_after_generation = false,
        support_paste_from_clipboard = false,
        enable_cursor_planning_mode = true, -- enable cursor planning mode!
      },
      -- File selector configuration
      --- @alias FileSelectorProvider "native" | "fzf" | "mini.pick" | "snacks" | "telescope" | string
      file_selector = {
        provider = "snacks", -- Avoid native provider issues
        provider_opts = {},
      },
      mappings = {
        --- @class AvanteConflictMappings
        diff = {
          ours = "co",
          theirs = "ct",
          all_theirs = "ca",
          both = "cb",
          cursor = "cc",
          next = "]x",
          prev = "[x",
        },
        suggestion = {
          accept = "<M-l>",
          next = "<M-]>",
          prev = "<M-[>",
          dismiss = "<C-]>",
        },
        jump = {
          next = "]]",
          prev = "[[",
        },
        submit = {
          normal = "<CR>",
          insert = "<C-s>",
        },
        sidebar = {
          apply_all = "A",
          apply_cursor = "a",
          switch_windows = "<Tab>",
          reverse_switch_windows = "<S-Tab>",
        },
      },
      hints = { enabled = false },
      windows = {
        ---@type "right" | "left" | "top" | "bottom" | "smart"
        position = "smart", -- the position of the sidebar
        wrap = true, -- similar to vim.o.wrap
        width = 30, -- default % based on available width
        sidebar_header = {
          enabled = true, -- true, false to enable/disable the header
          align = "center", -- left, center, right for title
          rounded = false,
        },
        input = {
          prefix = "> ",
          height = 8, -- Height of the input window in vertical layout
        },
        edit = {
          start_insert = true, -- Start insert mode when opening the edit window
        },
        ask = {
          floating = false, -- Open the 'AvanteAsk' prompt in a floating window
          start_insert = true, -- Start insert mode when opening the ask window
          ---@type "ours" | "theirs"
          focus_on_apply = "ours", -- which diff to focus after applying
        },
      },
      highlights = {
        ---@type AvanteConflictHighlights
        diff = {
          current = "DiffText",
          incoming = "DiffAdd",
        },
      },
      --- @class AvanteConflictUserConfig
      diff = {
        autojump = true,
        ---@type string | fun(): any
        list_opener = "copen",
        --- Override the 'timeoutlen' setting while hovering over a diff (see :help timeoutlen).
        --- Helps to avoid entering operator-pending mode with diff mappings starting with `c`.
        --- Disable by setting to -1.
        override_timeoutlen = 500,
      },
      system_prompt = [[
      Sos un clon de Gentleman Programming, un arquitecto frontend argentino con un enfoque técnico pero relajado. Tu estilo es claro, directo y con un toque de humor inteligente. Estás especializado en Angular y React, con obsesión por la arquitectura limpia, hexagonal y scalable, y fanático del patrón contenedor-presentacional, modularización, atomic design y defensive programming.

      Te dirigís a desarrolladores intermedios y avanzados, explicás conceptos complejos de forma clara y práctica, sin vueltas, con ejemplos útiles. Usás analogías del mundo de la construcción para ilustrar ideas difíciles. Tus charlas mezclan técnica con introspección, liderazgo y enseñanza. Tenés experiencia en mentoría, creación de contenido y comunidades tech.

      Hablas en tono argentino, natural y accesible. Usás expresiones como “buenas acá estamos”, “dale que va”, “acá la posta es esta”, pero sin caer en clichés forzados. Valorás las buenas prácticas, el testing, la productividad con herramientas como LazyVim, Tmux, Zellij y OBS, y la exploración de nuevas herramientas sin perder el foco.

      A la hora de responder:

        1. Identificás el problema técnico del usuario.
        2. Proponés una solución concreta con fundamentos.
        3. Dás ejemplos o snippets si aplican
        4. Recomendás herramientas si suman valor.

      Tu rol es acompañar, formar y destrabar nudos técnicos sin chamuyo. Si algo es complejo, lo bajás a tierra. Si algo es innecesario, lo decís. Tu estilo es: pragmático, apasionado, sin humo.]],  
    },
    -- if you want to build from source then do `make BUILD_FROM_SOURCE=true`
    build = "make",
    -- build = "powershell -ExecutionPolicy Bypass -File Build.ps1 -BuildFromSource false" -- for windows
    dependencies = {
      "nvim-treesitter/nvim-treesitter",
      "stevearc/dressing.nvim",
      "nvim-lua/plenary.nvim",
      "MunifTanjim/nui.nvim",
      --- The below dependencies are optional,
      "nvim-tree/nvim-web-devicons", -- or echasnovski/mini.icons
      {
        -- support for image pasting
        "HakonHarnes/img-clip.nvim",
        event = "VeryLazy",
        opts = {
          -- recommended settings
          default = {
            embed_image_as_base64 = false,
            prompt_for_file_name = false,
            drag_and_drop = {
              insert_mode = true,
            },
            -- required for Windows users
            use_absolute_path = true,
          },
        },
      },
    },
  },
}
