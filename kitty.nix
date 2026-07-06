{ config, ... }:
{
  home.file.".config/kitty/kitty.conf".text = ''
    # ╔══════════════════════════════════════════════════════════════════════════════╗
    # ║                          GENTLEMAN DOTS - KITTY                              ║
    # ║                Local config matched to current Ghostty setup                  ║
    # ╚══════════════════════════════════════════════════════════════════════════════╝
    # vim:fileencoding=utf-8:foldmethod=marker


    # ┌──────────────────────────────────────────────────────────────────────────────┐
    # │                                   FONT                                       │
    # └──────────────────────────────────────────────────────────────────────────────┘

    font_family      IosevkaTerm Nerd Font Mono
    font_size        14.0


    # ┌──────────────────────────────────────────────────────────────────────────────┐
    # │                                  WINDOW                                      │
    # └──────────────────────────────────────────────────────────────────────────────┘

    background_opacity 0.95
    # Kitty does not expose Ghostty's background-blur-radius directly.

    window_padding_width 0
    placement_strategy center
    remember_window_size no
    initial_window_width 100c
    initial_window_height 100c

    # Match Ghostty's hidden tab UI preference.
    tab_bar_style hidden


    # ┌──────────────────────────────────────────────────────────────────────────────┐
    # │                                  SHELL                                       │
    # └──────────────────────────────────────────────────────────────────────────────┘

    # Force the same Fish shell used by the rest of this machine. Fish auto-starts
    # Herdr for fresh interactive sessions when HERDR_ENV/TMUX/ZELLIJ are not set.
    shell ${config.home.homeDirectory}/.local/state/nix/profiles/home-manager/home-path/bin/fish --login

    # When Kitty is launched from an existing Herdr session, macOS can inherit
    # HERDR_ENV. Remove multiplexer markers so a fresh Kitty window can start Herdr.
    env HERDR_ENV
    env TMUX
    env ZELLIJ


    # ┌──────────────────────────────────────────────────────────────────────────────┐
    # │                                  CURSOR                                      │
    # └──────────────────────────────────────────────────────────────────────────────┘

    cursor_shape          block
    cursor_blink_interval 0.5
    cursor_stop_blinking_after 0
    shell_integration     no-cursor


    # ┌──────────────────────────────────────────────────────────────────────────────┐
    # │                            NEOVIM OPTIMIZATIONS                              │
    # └──────────────────────────────────────────────────────────────────────────────┘

    term xterm-kitty
    undercurl_style thin-sparse
    scrollback_lines 10000
    repaint_delay    10
    input_delay      3
    sync_to_monitor  yes
    allow_remote_control yes
    listen_on unix:/tmp/kitty
    enable_audio_bell no


    # ┌──────────────────────────────────────────────────────────────────────────────┐
    # │                            INPUT / KEYBINDINGS                               │
    # └──────────────────────────────────────────────────────────────────────────────┘

    # macOS Alt key fix, matching Ghostty's left Option behavior.
    macos_option_as_alt left
    map alt+left no_op
    map alt+right no_op

    # Multiplexer-style shortcuts are intentionally not bound in Kitty.
    # Herdr owns Ctrl+a and pane/tab/workspace navigation inside the terminal.

    map cmd+k clear_terminal reset active
    map shift+enter send_text all \x1b\r

    # macOS-style command delete behavior.
    # On Apple keyboards the Backspace key is labeled Delete, so bind both names.
    map cmd+backspace send_text all \x15
    map cmd+delete send_text all \x0b

    # Ghostty's write_screen_file/write_scrollback_file actions do not have a direct
    # Kitty equivalent. Use Kitty scrollback tools manually when needed.


    # ┌──────────────────────────────────────────────────────────────────────────────┐
    # │                           GENTLEMAN THEME                                    │
    # └──────────────────────────────────────────────────────────────────────────────┘

    # --- Base Colors ---
    background            #06080f
    foreground            #f3f6f9
    cursor                #e0c15a
    selection_background  #263356
    selection_foreground  #f3f6f9
    url_color             #7fb4ca

    # --- Tabs ---
    active_tab_background   #263356
    active_tab_foreground   #f3f6f9
    inactive_tab_background #06080f
    inactive_tab_foreground #8a8fa3

    # --- Normal Colors ---
    color0  #06080f
    color1  #cb7c94
    color2  #b7cc85
    color3  #ffe066
    color4  #7fb4ca
    color5  #ff8dd7
    color6  #7aa89f
    color7  #f3f6f9

    # --- Bright Colors ---
    color8  #8a8fa3
    color9  #de8fa8
    color10 #d1e8a9
    color11 #fff7b1
    color12 #a3d4d5
    color13 #ffaeea
    color14 #7fb4ca
    color15 #f3f6f9
  '';
}
