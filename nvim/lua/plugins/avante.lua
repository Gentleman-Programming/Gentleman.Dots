return {
  {
    "yetone/avante.nvim",
    -- if you want to build from source then do `make BUILD_FROM_SOURCE=true`
    -- ⚠️ must add this setting! ! !
    build = function()
      -- conditionally use the correct build system for the current OS
      if vim.fn.has("win32") == 1 then
        return "powershell -ExecutionPolicy Bypass -File Build.ps1 -BuildFromSource false"
      else
        return "make"
      end
    end,
    event = "VeryLazy",
    version = false, -- Never set this value to "*"! Never!
    ---@module 'avante'
    ---@type avante.Config
    opts = function(_, opts)
      -- Track avante's internal state during resize
      local in_resize = false
      local original_cursor_win = nil
      local avante_filetypes = { "Avante", "AvanteInput", "AvanteAsk", "AvanteSelectedFiles" }

      -- Check if current window is avante
      local function is_in_avante_window()
        local win = vim.api.nvim_get_current_win()
        local buf = vim.api.nvim_win_get_buf(win)
        local ft = vim.api.nvim_buf_get_option(buf, "filetype")

        for _, avante_ft in ipairs(avante_filetypes) do
          if ft == avante_ft then
            return true, win, ft
          end
        end
        return false
      end

      -- Temporarily move cursor away from avante during resize
      local function temporarily_leave_avante()
        local is_avante, avante_win, avante_ft = is_in_avante_window()
        if is_avante and not in_resize then
          in_resize = true
          original_cursor_win = avante_win

          -- Find a non-avante window to switch to
          local target_win = nil
          for _, win in ipairs(vim.api.nvim_list_wins()) do
            local buf = vim.api.nvim_win_get_buf(win)
            local ft = vim.api.nvim_buf_get_option(buf, "filetype")

            local is_avante_ft = false
            for _, aft in ipairs(avante_filetypes) do
              if ft == aft then
                is_avante_ft = true
                break
              end
            end

            if not is_avante_ft and vim.api.nvim_win_is_valid(win) then
              target_win = win
              break
            end
          end

          -- Switch to non-avante window if found
          if target_win then
            vim.api.nvim_set_current_win(target_win)
            return true
          end
        end
        return false
      end

      -- Restore cursor to original avante window
      local function restore_cursor_to_avante()
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
      local function cleanup_duplicate_avante_windows()
        local seen_filetypes = {}
        local windows_to_close = {}

        for _, win in ipairs(vim.api.nvim_list_wins()) do
          local buf = vim.api.nvim_win_get_buf(win)
          local ft = vim.api.nvim_buf_get_option(buf, "filetype")

          -- Special handling for Ask and Select Files panels
          if ft == "AvanteAsk" or ft == "AvanteSelectedFiles" then
            if seen_filetypes[ft] then
              -- Found duplicate, mark for closing
              table.insert(windows_to_close, win)
            else
              seen_filetypes[ft] = win
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
      vim.api.nvim_create_augroup("AvanteResizeFix", { clear = true })

      -- Main resize handler for Resize
      vim.api.nvim_create_autocmd({ "VimResized" }, {
        group = "AvanteResizeFix",
        callback = function()
          -- Move cursor away from avante before resize processing
          local moved = temporarily_leave_avante()

          if moved then
            -- Let resize happen, then restore cursor
            vim.defer_fn(function()
              restore_cursor_to_avante()
              -- Force a clean redraw
              vim.cmd("redraw!")
            end, 100)
          end

          -- Cleanup duplicates after resize completes
          vim.defer_fn(cleanup_duplicate_avante_windows, 150)
        end,
      })

      -- Prevent avante from responding to scroll/resize events during resize
      vim.api.nvim_create_autocmd({ "WinScrolled", "WinResized" }, {
        group = "AvanteResizeFix",
        pattern = "*",
        callback = function(args)
          local buf = args.buf
          if buf and vim.api.nvim_buf_is_valid(buf) then
            local ft = vim.api.nvim_buf_get_option(buf, "filetype")

            for _, avante_ft in ipairs(avante_filetypes) do
              if ft == avante_ft then
                -- Prevent event propagation for avante buffers during resize
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
        group = "AvanteResizeFix",
        callback = function()
          -- Reset resize state on focus gain
          in_resize = false
          original_cursor_win = nil
          -- Clean up any duplicate windows
          vim.defer_fn(cleanup_duplicate_avante_windows, 100)
        end,
      })

      return {
        -- add any opts here
        -- for example
        provider = "copilot",
        providers = {
          copilot = {
            model = "claude-sonnet-4",
          },
        },
        web_search_engine = {
          provider = "google", -- tavily, serpapi, searchapi, google, kagi, brave, or searxng
          proxy = nil, -- proxy support, e.g., http://127.0.0.1:7890
        },
        cursor_applying_provider = "copilot",
        auto_suggestions_provider = "copilot",
        behaviour = {
          enable_cursor_planning_mode = true,
          enable_token_counting = false,
        },
        -- File selector configuration
        --- @alias FileSelectorProvider "native" | "fzf" | "mini.pick" | "snacks" | "telescope" | string
        file_selector = {
          provider = "snacks", -- Avoid native provider issues
          provider_opts = {},
        },
        windows = {
          ---@type "right" | "left" | "top" | "bottom" | "smart"
          position = "right", -- the position of the sidebar
          wrap = true, -- similar to vim.o.wrap
          width = 30, -- default % based on available width
          sidebar_header = {
            enabled = true, -- true, false to enable/disable the header
            align = "center", -- left, center, right for title
            rounded = true,
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
- **Siempre** incluye un **follow-up** al final para continuar la conversación.os un clon de Gentleman Programming, un arquitecto frontend argentino con un enfoque técnico pero relajado. Tu estilo es claro, directo y con un toque de humor inteligente. Estás especializado en Angular y React, con obsesión por la arquitectura limpia, hexagonal y scalable, y fanático del patrón contenedor-presentacional, modularización, atomic design y defensive programming.,
          ]],
      }
    end,
    dependencies = {
      "MunifTanjim/nui.nvim",
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
