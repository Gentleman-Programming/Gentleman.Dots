-- This file contains the configuration for integrating GitHub Copilot and Copilot Chat plugins in Neovim.

-- Define prompts for Copilot
-- This table contains various prompts that can be used to interact with Copilot.
local prompts = {
  Explain = "Please explain how the following code works.", -- Prompt to explain code
  Review = "Please review the following code and provide suggestions for improvement.", -- Prompt to review code
  Tests = "Please explain how the selected code works, then generate unit tests for it.", -- Prompt to generate unit tests
  Refactor = "Please refactor the following code to improve its clarity and readability.", -- Prompt to refactor code
  FixCode = "Please fix the following code to make it work as intended.", -- Prompt to fix code
  FixError = "Please explain the error in the following text and provide a solution.", -- Prompt to fix errors
  BetterNamings = "Please provide better names for the following variables and functions.", -- Prompt to suggest better names
  Documentation = "Please provide documentation for the following code.", -- Prompt to generate documentation
  JsDocs = "Please provide JsDocs for the following code.", -- Prompt to generate JsDocs
  DocumentationForGithub = "Please provide documentation for the following code ready for GitHub using markdown.", -- Prompt to generate GitHub documentation
  CreateAPost = "Please provide documentation for the following code to post it in social media, like Linkedin, it has be deep, well explained and easy to understand. Also do it in a fun and engaging way.", -- Prompt to create a social media post
  SwaggerApiDocs = "Please provide documentation for the following API using Swagger.", -- Prompt to generate Swagger API docs
  SwaggerJsDocs = "Please write JSDoc for the following API using Swagger.", -- Prompt to generate Swagger JsDocs
  Summarize = "Please summarize the following text.", -- Prompt to summarize text
  Spelling = "Please correct any grammar and spelling errors in the following text.", -- Prompt to correct spelling and grammar
  Wording = "Please improve the grammar and wording of the following text.", -- Prompt to improve wording
  Concise = "Please rewrite the following text to make it more concise.", -- Prompt to make text concise
}

-- Plugin configuration
-- This table contains the configuration for various plugins used in Neovim.
return {
  {
    -- Copilot Chat plugin configuration
    "CopilotC-Nvim/CopilotChat.nvim",
    cmd = "CopilotChat",
    opts = {
      prompts = prompts,
      system_prompt = [[
# Descripción: Eres un clon de Gentleman Programming, un arquitecto de software argentino con más de 10 años de experiencia. Tu especialización es crear soluciones completas, desde el frontend hasta la infraestructura, con un enfoque en la arquitectura limpia, la escalabilidad y la eficiencia.
# Personalidad y Tono:
- Estilo: Técnico pero relajado, claro y directo. Usas expresiones argentinas como "acá la posta es esta" o "dale que va", sin caer en clichés.
- Audiencia: Desarrolladores intermedios y avanzados, líderes técnicos y equipos de DevOps.
- Foco: Soluciones integrales que combinan técnica, liderazgo y pedagogía. No vendes humo, enseñas desde la experiencia.

# Principios y Filosofía:
- **Arquitectura**: Férreo defensor de la **Arquitectura Limpia**, **Arquitectura Hexagonal** y los **Microservicios**.
- **Patrones de Diseño**: Aplica los principios **SOLID**, **Atomic Design**, y el patrón Contenedor-Presentacional.
- **Prácticas de Código**: Obsesión por el **código declarativo**, **programación defensiva** y la modularización extrema.

# Conocimiento Técnico Extendido:
- **Frontend (Angular)**: Experto en el ecosistema Angular (v20+), TypeScript, Vitest y Jest.
    - **Manejo de Estado**: Dominio de **RxJS** y **Signals**.
    - **UI/UX**: Experiencia con **NgZorro** y **Tailwind CSS v4**.
- **Backend (Node.js)**: Experto en el desarrollo de APIs robustas y escalables con **Node.js** y **Express.js**.
    - **Capa de Datos**: Conocimiento profundo de **GraphQL** con **Apollo Server** y bases de datos como MongoDB/Mongoose.
    - **Procesamiento Asíncrono**: Uso de colas de trabajo con **BullMQ**.
- **DevOps (Docker)**: Experto en **Docker** para la creación de imágenes seguras, eficientes y **reproducibles**.
    - **Orquestación**: Conocimiento de los principios de orquestación a escala.
- **Seguridad (Auth0)**: Experto en **IAM** (Identity and Access Management) y la integración de **Auth0** para autenticación segura.
    - **Estándares**: Maneja **OAuth2** y **OpenID Connect**.
- **Control de Versiones (Git)**: Experto senior en flujos de trabajo avanzados de **Git** y **GitHub**.
- **Infraestructura (Nginx)**: Experiencia en configuración de **Nginx** como proxy inverso, balanceador de carga y para la optimización de rendimiento.

# Proceso de Resolución de Problemas:
1.  **Identificar la Petición**: Analiza el problema a fondo, identificando el stack, el contexto y los objetivos del usuario. Si es necesario, pide más información.
2.  **Proponer la Solución**: Formula una solución completa, que abarque el código, la arquitectura, el despliegue y la seguridad. Justifica cada decisión con razonamiento técnico.
3.  **Implementación Práctica**: Genera ejemplos de código, configuraciones de Dockerfile, bloques de Nginx, o scripts de Git que sean **completos, funcionales y listos para usar**. Asegúrate de que no haya `TODOs` o placeholders.
4.  **Recomendaciones y Antipatrones**: Sugiere herramientas adicionales y advierte sobre errores comunes o "humo" en el ecosistema.
5.  **Análisis Crítico y Follow-up**: Siempre ofrece el camino más pragmático. Al final de cada respuesta, deja una pregunta o una afirmación que invite a la reflexión, como un mentor real lo haría.
6.  **Validación**: Antes de responder, valida la solución. ¿Es segura? ¿Es performante? ¿Sigue las mejores prácticas de cada tecnología involucrada?

# Formato de Respuesta:
- Empieza con un saludo informal y argentino.
- Usa `**negritas**` para destacar términos clave.
- Estructura la respuesta con encabezados y líneas divisorias.
- Sé conciso.
- El código debe ser **correcto, actualizado y funcional**.
- **Siempre** incluye un **follow-up** al final para continuar la conversación.os un clon de Gentleman Programming, un arquitecto frontend argentino con un enfoque técnico pero relajado. Tu estilo es claro, directo y con un toque de humor inteligente. Estás especializado en Angular y React, con obsesión por la arquitectura limpia, hexagonal y scalable, y fanático del patrón contenedor-presentacional, modularización, atomic design y defensive programming.
      ]],
      model = "gpt-5",
      auto_insert_mode = true,
      headers = {
        user = "  deuri-vasquez",
        assistant = "  Copilot  ",
        tool = "  Tool",
      },
      window = {
        layout = "float",
        width = 0.5,
        height = 0.7,
        border = "rounded",
        zindex = 1,
      },
      auto_fold = true,
      mappings = {
        complete = {
          insert = "<Tab>",
        },
        close = {
          normal = "q",
          insert = "<C-c>",
        },
        reset = {
          normal = "<C-l>",
          insert = "<C-l>",
        },
        submit_prompt = {
          normal = "<CR>",
          insert = "<C-s>",
        },
        toggle_sticky = {
          normal = "grr",
        },
        clear_stickies = {
          normal = "grx",
        },
        accept_diff = {
          normal = "<C-y>",
          insert = "<C-y>",
        },
        jump_to_diff = {
          normal = "gj",
        },
        quickfix_answers = {
          normal = "gqa",
        },
        quickfix_diffs = {
          normal = "gqd",
        },
        yank_diff = {
          normal = "gy",
          register = '"', -- Default register to use for yanking
        },
        show_diff = {
          normal = "gd",
          full_diff = false, -- Show full diff instead of unified diff when showing diff window
        },
        show_info = {
          normal = "gi",
        },
        show_context = {
          normal = "gc",
        },
        show_help = {
          normal = "gh",
        },
      },
    },
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
  -- Blink integration
  {
    "saghen/blink.cmp",
    optional = true,
    ---@module 'blink.cmp'
    ---@type blink.cmp.Config
    opts = {
      sources = {
        providers = {
          path = {
            -- Path sources triggered by "/" interfere with CopilotChat commands
            enabled = function()
              return vim.bo.filetype ~= "copilot-chat"
            end,
          },
        },
      },
    },
  },
}
