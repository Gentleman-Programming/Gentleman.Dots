{ config, pkgs, ... }:
{
  home.file = {
    ".tmux.conf" = {
      text = ''
set -g prefix C-a
unbind C-b
bind C-a send-prefix

bind -n M-Up resize-pane -U 5
bind -n M-Down resize-pane -D 5
bind -n M-Left resize-pane -L 5
bind -n M-Right resize-pane -R 5

bind-key -n M-1 select-window -t 1
bind-key -n M-2 select-window -t 2
bind-key -n M-3 select-window -t 3
bind-key -n M-4 select-window -t 4
bind-key -n M-5 select-window -t 5
bind-key -n M-6 select-window -t 6
bind-key -n M-7 select-window -t 7
bind-key -n M-8 select-window -t 8
bind-key -n M-9 select-window -t 9

bind '"' split-window -c "#{pane_current_path}"
bind % split-window -h -c "#{pane_current_path}"
bind c new-window -c "#{pane_current_path}"

bind-key -n M-g if-shell -F '#{==:#{session_name},scratch}' {
  detach-client
} {
  display-popup -d "#{pane_current_path}" -E "tmux new-session -A -s scratch"
}

set -g @plugin 'Nybkox/tmux-kanagawa'
set -g @kanagawa-theme "dragon"
set -g @kanagawa-plugins "cpu-usage ram-usage time"
set -g @kanagawa-show-powerline true
set -g @kanagawa-show-timezone false
set -g @kanagawa-ignore-window-colors true

set -g @plugin 'tmux-plugins/tpm'
set -g @plugin 'tmux-plugins/tmux-sensible'
set -g @plugin 'tmux-plugins/tmux-resurrect'
set -g @plugin 'christoomey/vim-tmux-navigator'
set-option -ga terminal-overrides ",xterm*:Tc"
set -g default-terminal "tmux-256color"
set -g default-terminal "screen-256color"

set -sg escape-time 0 
set -g status-interval 0
set -g status-position top
set -g mode-keys vi

if-shell 'uname | grep -q Darwin' 'bind-key -T copy-mode-vi y send-keys -X copy-pipe-and-cancel "pbcopy"' 'bind-key -T copy-mode-vi y send-keys -X copy-pipe-and-cancel "clip"'

run '~/.tmux/plugins/tpm/tpm'
      '';
    };
  };
}
