-- prompts/copilot.lua
local prompts = {
  Explain = "Please explain how the following code works.",
  Review = "Please review the following code and provide suggestions for improvement.",
  Tests = "Please explain how the selected code works, then generate unit tests for it.",
  Refactor = "Please refactor the following code to improve its clarity and readability.",
  FixCode = "Please fix the following code to make it work as intended.",
  FixError = "Please explain the error in the following text and provide a solution.",
  BetterNamings = "Please provide better names for the following variables and functions.",
  Documentation = "Please provide documentation for the following code.",
  JsDocs = "Please provide JsDocs for the following code.",
  DocumentationForGithub = "Please provide documentation for the following code ready for GitHub using markdown.",
  CreateAPost = "Please write a social media post (e.g. LinkedIn) explaining the following code in a fun, deep, and engaging way.",
  SwaggerApiDocs = "Please provide documentation for the following API using Swagger.",
  SwaggerJsDocs = "Please write JSDoc for the following API using Swagger.",
  Summarize = "Please summarize the following text.",
  Spelling = "Please correct grammar and spelling errors in the following text.",
  Wording = "Please improve grammar and wording of the following text.",
  Concise = "Please rewrite the following text to be more concise.",
}

return {
  {
    "CopilotC-Nvim/CopilotChat.nvim",
    branch = "main",
    cmd = "CopilotChat",
    dependencies = {
      "github/copilot.vim", -- o zbirenbaum/copilot.lua
      "nvim-lua/plenary.nvim",
    },
    opts = function()
      return {
        prompts = prompts,
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
        model = "gpt-4o",
        answer_header = "󱗞  The Gentleman  ",
        question_header = "  deuri-vasquez  ",
        auto_insert_mode = "insert",
        window = {
          layout = "float",
          width = 0.5,
          height = 0.6,
          relative = "editor",
          border = "rounded",
          title = "Copilot Chat",
          zindex = 1,
        },
        mappings = {
          complete = { insert = "<Tab>" },
          close = { normal = "q", insert = "<C-c>" },
          reset = { normal = "<C-l>", insert = "<C-l>" },
          submit_prompt = { normal = "<CR>", insert = "<C-s>" },
          toggle_sticky = { normal = "grr" },
          clear_stickies = { normal = "grx" },
          accept_diff = { normal = "<C-y>", insert = "<C-y>" },
          jump_to_diff = { normal = "gj" },
          quickfix_answers = { normal = "gqa" },
          quickfix_diffs = { normal = "gqd" },
          yank_diff = { normal = "gy", register = '"' },
          show_diff = { normal = "gd", full_diff = false },
          show_info = { normal = "gi" },
          show_context = { normal = "gc" },
          show_help = { normal = "gh" },
        },
      }
    end,
    config = function(_, opts)
      local chat = require("CopilotChat")

      vim.api.nvim_create_autocmd("BufEnter", {
        pattern = "copilot-chat",
        callback = function()
          vim.opt_local.relativenumber = true
          vim.opt_local.number = false
        end,
      })

      chat.setup(opts)
    end,
  },

  -- Blink integration to avoid interference with path completion in chat buffer
  {
    "saghen/blink.cmp",
    optional = true,
    opts = {
      sources = {
        providers = {
          path = {
            enabled = function()
              return vim.bo.filetype ~= "copilot-chat"
            end,
          },
        },
      },
    },
  },
}
