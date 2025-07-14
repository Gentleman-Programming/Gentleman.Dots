{ pkgs, ... }:

{
  # Instalar TPM (Tmux Plugin Manager)
  home.activation.installTpm = ''
    if [ ! -d ~/.tmux/plugins/tpm ]; then
      ${pkgs.git}/bin/git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm
    fi
  '';

  home.file = {
    ".config/tmux/tmux.conf" = {
      text = ''
# Carga TPM
set -g @plugin 'tmux-plugins/tpm'

# Opciones ya probadas que limpian problemas de TMUX
set -g @plugin 'tmux-plugins/tmux-sensible'

# Tmux Navigation
set -g @plugin 'christoomey/vim-tmux-navigator'

# Tmux Resurrect
set -g @plugin 'tmux-plugins/tmux-resurrect'

# Which Key
set -g @plugin 'alexwforsythe/tmux-which-key'

# Floating window
bind-key -n M-g if-shell -F '#{==:#{session_name},scratch}' {
detach-client
} {
# open in the same directory of the current pane
display-popup -d "#{pane_current_path}" -E "tmux new-session -A -s scratch"
}

# Tema Kanagawa
set -g @plugin 'Nybkox/tmux-kanagawa'
set -g @kanagawa-theme 'Kanagawa'
set -g @kanagawa-plugins "git cpu-usage ram-usage"
set -g @kanagawa-ignore-window-colors true

# Fix colors for the terminal
set -g default-terminal 'tmux-256color'
set -ga terminal-overrides ",xterm-256color:Tc"

# Modo vim
set -g mode-keys vi
if-shell 'uname | grep -q Darwin' 'bind-key -T copy-mode-vi y send-keys -X copy-pipe-and-cancel "pbcopy"' 'bind-key -T copy-mode-vi y send-keys -X copy-pipe-and-cancel "clip"'

# Keymaps
unbind C-b
set -g prefix C-a
bind C-a send-prefix

unbind %
unbind '"'
bind v split-window -h -c "#{pane_current_path}"
bind d split-window -v -c "#{pane_current_path}"

# Mouse support
set -g mouse on

# Status bar position
set -g status-position top

# Fix index
set -g base-index 1
setw -g pane-base-index 1

run '~/.tmux/plugins/tpm/tpm'
      '';
    };
  };
}
