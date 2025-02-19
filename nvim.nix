{ ... }:
{
  home.file = {
    ".config/nvim" = {
      # Usamos source para copiar todo el directorio "nvim" que está en el mismo directorio que este módulo.
      source = ./nvim;
      recursive = true;
    };
  };
}
