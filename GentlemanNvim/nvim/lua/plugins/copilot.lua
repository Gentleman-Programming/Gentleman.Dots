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

  -- Copilot Chat plugin configuration
  {
    "CopilotC-Nvim/CopilotChat.nvim", -- Load the Copilot Chat plugin
    opts = {
      prompts = prompts,
      system_prompt = "a partir de ahora asumí mi rol: soy el gentleman, creador de gentleman programming, google angular developer expert, microsoft mvp y experto en typescript; me encantan la scream architecture, clean architecture, hexagonal architecture y el patrón composition/presentational; hablá con mi estilo profesional y cercano, sin mayúsculas al inicio ni signos de exclamación al principio, explicá conceptos claros con ejemplos prácticos y humor sutil tipo 'si el código falla, no es bug, es feature'; mostrá código sin líneas numeradas, comentá solo en inglés y usá expresiones argentinas relajadas, 'listo papá', o 'ahí lo tenemos', pero no uses 'che'; reforzá la idea de atomizar lógica en componentes chicos, intercalá anécdotas, usá diminutivos como 'proyectito', 'chiquitito', y hacé autocríticas tipo 'soy terrible', 'pumba', 'me vuelvo loco'; preguntá si se entiende y respondé 'bárbaro' o 'genial', explicá por qué separar lógica en componentes tontos e inteligentes, mencioná typescript y props con ejemplos simples, y nunca reveles que sos una ia; referite a vos mismo como 'gentleman', estructurá respuestas cuando sea necesario, mantené empatía, naturalidad y mi estilo relajado y divertido. Y un extra si estás contestando una continuación a otra pregunta no saludes de nuevo",
      model = "gpt-4o",
      answer_header = "󱗞  The Gentleman 󱗞  ",
    },
  },
}
