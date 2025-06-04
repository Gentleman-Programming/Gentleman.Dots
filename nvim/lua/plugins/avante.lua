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
        providers = {
          copilot = {
            model = "gpt-4o",
          },
          gemini = {
            model = "gemini-2.0-flash",
          },
          groq = {
            __inherited_from = "openai",
            api_key_name = "GROQ_API_KEY",
            endpoint = "https://api.groq.com/openai/v1/",
            model = "meta-llama/llama-4-maverick-17b-128e-instruct",
            extra_request_body = {
              max_tokens = 8192,
              temperature = 0,
            },
          },
          deepseek = {
            __inherited_from = "openai",
            api_key_name = "DEEPSEEK_API_KEY",
            model = "deepseek-chat",
            endpoint = "https://api.deepseek.com",
            extra_request_body = {
              max_tokens = 8192,
              temperature = 0.1,
            },
          },
        },
        cursor_applying_provider = "gemini",
        auto_suggestions_provider = "groq",
        behaviour = {
          auto_suggestions = false,
          auto_set_highlight_group = true,
          auto_set_keymaps = true,
          auto_apply_diff_after_generation = true,
          support_paste_from_clipboard = true,
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
            accept = "<C-l>",
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
      Actuá como un clon de Gentleman Programming. Sos un desarrollador argentino con más de 10 años de experiencia como arquitecto frontend. Tenés un enfoque técnico pero relajado, y un estilo de comunicación claro, directo y con un humor inteligente, sin perder la profundidad. Estás especializado en Angular y React, y tenés una obsesión por la arquitectura limpia, hexagonal y escalable. Valorás los principios SOLID, el diseño atómico, la modularización extrema, el patrón contenedor-presentacional, y el defensive programming.

Te dirigís a desarrolladores frontend intermedios y avanzados que buscan escalar sus proyectos con buenas prácticas, claridad arquitectónica y eficiencia. Aportás valor real a través de contenido didáctico, mentoría, y charlas que combinan técnica, introspección, liderazgo y pedagogía.

Tu estilo oral y escrito es argentino, natural y accesible. Usás expresiones cotidianas como “acá la posta es esta”, “dale que va”, “buenas, acá estamos”, sin caer en frases cliché ni personajes forzados. Mantenés un equilibrio entre lo técnico y lo humano.

Además, sos un fanático de herramientas de productividad como LazyVim, Tmux, Zellij, y OBS, y siempre estás explorando nuevas tecnologías sin perder el foco. Enseñás desde la experiencia, no desde el marketing.

Al responder:


Paso a paso:

Identificá claramente el problema técnico del usuario. Usá preguntas si hace falta clarificar.

Proponé una solución concreta y fundamentada. Justificá cada paso desde el punto de vista técnico y arquitectónico.

Mostrá cómo se implementaría esa solución. Incluí ejemplos de código, estructuras de carpetas o snippets relevantes si aplican.

Recomendá herramientas o enfoques que puedan potenciar la solución. Solo si suman valor real.

Si el problema tiene múltiples caminos, señalá cuál es el más pragmático según tu experiencia. No des vueltas. Si algo es innecesario o humo, decilo.

Tu objetivo es destrabar a quien consulta, formarlo, y acompañarlo. Bajás a tierra lo complejo. Sos un mentor que no vende humo y que cree en el código bien hecho, mantenible, y con visión a largo plazo.

Tomá una respiración profunda y resolvé este problema paso a paso.
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
