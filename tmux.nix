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

# Tested options for TMUX compatibility
set -g @plugin 'tmux-plugins/tmux-sensible'

# Clipboard management
set -g @plugin 'tmux-plugins/tmux-yank'

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
  display-popup -d "#{pane_current_path}" -E -k "tmux new-session -A -s scratch"
}

# Tema Kanagawa
set -g @plugin 'Nybkox/tmux-kanagawa'
set -g @kanagawa-theme 'Kanagawa'
set -g @kanagawa-plugins "git cpu-usage ram-usage"
set -g @kanagawa-ignore-window-colors true

# --- terminal & key handling ---
set -g default-terminal "tmux-256color"
set -as terminal-features ",*:RGB"
set -as terminal-features ",*:usstyle"
set -as terminal-features ",*:hyperlinks"

set -s extended-keys on
set -s extended-keys-format csi-u
set -as terminal-features 'xterm*:extkeys'
set -sg escape-time 10

# Modo vim
set -g mode-keys vi
set -g set-clipboard on
bind-key -T copy-mode-vi y send-keys -X copy-pipe-and-cancel

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

# Kill all sessions except current
bind K confirm-before -p "Kill all other sessions? (y/n)" "kill-session -a"

# Fix index
set -g base-index 1
setw -g pane-base-index 1

# TPM init
run '~/.tmux/plugins/tpm/tpm'
      '';
    };
  };
}
