return {
  {
    "yetone/avante.nvim",
    event = "VeryLazy",
    lazy = false,
    version = false,
    build = "make",
    opts = function()
      return {
        mode = "agentic",
        provider = "gemini",
        copilot = {
          model = "gpt-4o",
        },
        gemini = {
          endpoint = "https://generativelanguage.googleapis.com/v1beta/models",
          model = "gemini-2.0-flash",
          timeout = 30000,
          temperature = 0,
          max_tokens = 8192,
        },
        vendors = {
          groq = {
            __inherited_from = "openai",
            api_key_name = "GROQ_API_KEY",
            endpoint = "https://api.groq.com/openai/v1/",
            model = "meta-llama/llama-4-maverick-17b-128e-instruct",
            max_tokens = 8192,
            temperature = 0,
          },
          deepseek = {
            __inherited_from = "openai",
            api_key_name = "DEEPSEEK_API_KEY",
            endpoint = "https://api.deepseek.com",
            model = "deepseek-chat",
          },
        },
        cursor_applying_provider = "groq",
        auto_suggestions_provider = "gemini",
        behaviour = {
          auto_suggestions = true,
          auto_set_highlight_group = true,
          auto_set_keymaps = true,
          auto_apply_diff_after_generation = false,
          support_paste_from_clipboard = false,
          enable_cursor_planning_mode = true,
          enable_token_counting = false,
        },
        suggestion = {
          debounce = 400,
          throttle = 400,
        },
        file_selector = {
          provider = "snacks",
          provider_opts = {},
        },
        mappings = {
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
        windows = {
          position = "right",
          wrap = true,
          width = 30,
          sidebar_header = {
            enabled = true,
            align = "center",
            rounded = true,
          },
          input = {
            rounded = true,
            prefix = "> ",
            height = 8,
          },
          edit = {
            border = "rounded",
            start_insert = true,
          },
          ask = {
            border = "rounded",
            floating = false,
            start_insert = true,
            focus_on_apply = "ours",
          },
        },
        highlights = {
          diff = {
            current = "DiffText",
            incoming = "DiffAdd",
          },
        },
        diff = {
          autojump = true,
          list_opener = "copen",
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

      Tu rol es acompañar, formar y destrabar nudos técnicos sin chamuyo. Si algo es complejo, lo bajás a tierra. Si algo es innecesario, lo decís. Tu estilo es: pragmático, apasionado, sin humo.
]],
      }
    end,
    dependencies = {
      "nvim-treesitter/nvim-treesitter",
      "stevearc/dressing.nvim",
      "nvim-lua/plenary.nvim",
      "MunifTanjim/nui.nvim",
      "nvim-tree/nvim-web-devicons",
      {
        "HakonHarnes/img-clip.nvim",
        event = "VeryLazy",
        opts = {
          default = {
            embed_image_as_base64 = false,
            prompt_for_file_name = false,
            drag_and_drop = {
              insert_mode = true,
            },
            use_absolute_path = true,
          },
        },
      },
    },
  },
}
