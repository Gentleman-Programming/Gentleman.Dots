{

  home.file = {
    ".config/zellij/plugins" = { source = ./zellij; };
  };

  home.file = {
    ".config/zellij/config.kdl" = {
      text = ''
keybinds clear-defaults=true {
    locked {
        bind "Ctrl g" { SwitchToMode "normal"; }
    }
    pane {
        bind "left" { MoveFocus "left"; }
        bind "down" { MoveFocus "down"; }
        bind "up" { MoveFocus "up"; }
        bind "right" { MoveFocus "right"; }
        bind "c" { SwitchToMode "renamepane"; PaneNameInput 0; }
        bind "d" { NewPane "down"; SwitchToMode "locked"; }
        bind "e" { TogglePaneEmbedOrFloating; SwitchToMode "locked"; }
        bind "f" { ToggleFocusFullscreen; SwitchToMode "locked"; }
        bind "h" { MoveFocus "left"; }
        bind "j" { MoveFocus "down"; }
        bind "k" { MoveFocus "up"; }
        bind "l" { MoveFocus "right"; }
        bind "n" { NewPane; SwitchToMode "locked"; }
        bind "p" { SwitchToMode "normal"; }
        bind "r" { NewPane "right"; SwitchToMode "locked"; }
        bind "w" { ToggleFloatingPanes; SwitchToMode "locked"; }
        bind "x" { CloseFocus; SwitchToMode "locked"; }
        bind "z" { TogglePaneFrames; SwitchToMode "locked"; }
        bind "tab" { SwitchFocus; }
    }
    tab {
        bind "left" { GoToPreviousTab; }
        bind "down" { GoToNextTab; }
        bind "up" { GoToPreviousTab; }
        bind "right" { GoToNextTab; }
        bind "1" { GoToTab 1; SwitchToMode "locked"; }
        bind "2" { GoToTab 2; SwitchToMode "locked"; }
        bind "3" { GoToTab 3; SwitchToMode "locked"; }
        bind "4" { GoToTab 4; SwitchToMode "locked"; }
        bind "5" { GoToTab 5; SwitchToMode "locked"; }
        bind "6" { GoToTab 6; SwitchToMode "locked"; }
        bind "7" { GoToTab 7; SwitchToMode "locked"; }
        bind "8" { GoToTab 8; SwitchToMode "locked"; }
        bind "9" { GoToTab 9; SwitchToMode "locked"; }
        bind "[" { BreakPaneLeft; SwitchToMode "locked"; }
        bind "]" { BreakPaneRight; SwitchToMode "locked"; }
        bind "b" { BreakPane; SwitchToMode "locked"; }
        bind "h" { GoToPreviousTab; }
        bind "j" { GoToNextTab; }
        bind "k" { GoToPreviousTab; }
        bind "l" { GoToNextTab; }
        bind "n" { NewTab; SwitchToMode "locked"; }
        bind "r" { SwitchToMode "renametab"; TabNameInput 0; }
        bind "s" { ToggleActiveSyncTab; SwitchToMode "locked"; }
        bind "t" { SwitchToMode "normal"; }
        bind "x" { CloseTab; SwitchToMode "locked"; }
        bind "tab" { ToggleTab; }
    }
    resize {
        bind "left" { Resize "Increase left"; }
        bind "down" { Resize "Increase down"; }
        bind "up" { Resize "Increase up"; }
        bind "right" { Resize "Increase right"; }
        bind "+" { Resize "Increase"; }
        bind "-" { Resize "Decrease"; }
        bind "=" { Resize "Increase"; }
        bind "H" { Resize "Decrease left"; }
        bind "J" { Resize "Decrease down"; }
        bind "K" { Resize "Decrease up"; }
        bind "L" { Resize "Decrease right"; }
        bind "h" { Resize "Increase left"; }
        bind "j" { Resize "Increase down"; }
        bind "k" { Resize "Increase up"; }
        bind "l" { Resize "Increase right"; }
        bind "r" { SwitchToMode "normal"; }
    }
    move {
        bind "left" { MovePane "left"; }
        bind "down" { MovePane "down"; }
        bind "up" { MovePane "up"; }
        bind "right" { MovePane "right"; }
        bind "h" { MovePane "left"; }
        bind "j" { MovePane "down"; }
        bind "k" { MovePane "up"; }
        bind "l" { MovePane "right"; }
        bind "m" { SwitchToMode "normal"; }
        bind "n" { MovePane; }
        bind "p" { MovePaneBackwards; }
        bind "tab" { MovePane; }
    }
    scroll {
        bind "Alt left" { MoveFocusOrTab "left"; SwitchToMode "locked"; }
        bind "Alt down" { MoveFocus "down"; SwitchToMode "locked"; }
        bind "Alt up" { MoveFocus "up"; SwitchToMode "locked"; }
        bind "Alt right" { MoveFocusOrTab "right"; SwitchToMode "locked"; }
        bind "e" { EditScrollback; SwitchToMode "locked"; }
        bind "f" { SwitchToMode "entersearch"; SearchInput 0; }
        bind "Alt h" { MoveFocusOrTab "left"; SwitchToMode "locked"; }
        bind "Alt j" { MoveFocus "down"; SwitchToMode "locked"; }
        bind "Alt k" { MoveFocus "up"; SwitchToMode "locked"; }
        bind "Alt l" { MoveFocusOrTab "right"; SwitchToMode "locked"; }
        bind "s" { SwitchToMode "normal"; }
    }
    search {
        bind "c" { SearchToggleOption "CaseSensitivity"; }
        bind "n" { Search "down"; }
        bind "o" { SearchToggleOption "WholeWord"; }
        bind "p" { Search "up"; }
        bind "w" { SearchToggleOption "Wrap"; }
    }
    session {
        bind "c" {
            LaunchOrFocusPlugin "configuration" {
                floating true
                move_to_focused_tab true
            }
            SwitchToMode "locked"
        }
        bind "d" { Detach; }
        bind "o" { SwitchToMode "normal"; }
        bind "p" {
            LaunchOrFocusPlugin "plugin-manager" {
                floating true
                move_to_focused_tab true
            }
            SwitchToMode "locked"
        }
        bind "w" {
            LaunchOrFocusPlugin "session-manager" {
                floating true
                move_to_focused_tab true
            }
            SwitchToMode "locked"
        }
    }
    shared_among "normal" "locked" {
        bind "Alt y" {
            LaunchOrFocusPlugin "file:~/.config/zellij/plugins/zellij_forgot.wasm" {
            "lock"                  "ctrl + g"
            "unlock"                "ctrl + g"
            "new pane"              "Alt + n || ctrl + g + p + n"
            "change focus of pane"  "Alt + arrow key || ctrl + g + p + arrow key"
            "close pane"            "ctrl + g + p + x"
            "rename pane"           "ctrl + g + p + c"
            "toggle fullscreen"     "ctrl + g + p + f"
            "toggle floating pane"  "Alt + f || ctrl + g + p + w"
            "toggle embed pane"     "ctrl + g + p + e"
            "choose right pane"     "ctrl + g + p + l"
            "choose left pane"      "ctrl + g + p + r"
            "choose upper pane"     "ctrl + g + p + k"
            "choose lower pane"     "ctrl + g + p + j"
            "new tab"               "ctrl + g + t + n"
            "close tab"             "ctrl + g + t + x"
            "change focus of tab"   "ctrl + g + t + arrow key"
            "rename tab"            "ctrl + g + t + r"
            "sync tab"              "ctrl + g + t + s"
            "brake pane to new tab" "ctrl + g + t + b"
            "brake pane left"       "ctrl + g + t + ["
            "brake pane right"      "ctrl + g + t + ]"
            "toggle tab"            "ctrl + g + t + tab"
            "increase pane size"    "ctrl + g + n + +"
            "decrease pane size"    "ctrl + g + n + -"
            "increase pane top"     "ctrl + g + n + k"
            "increase pane right"   "ctrl + g + n + l"
            "increase pane bottom"  "ctrl + g + n + j"
            "increase pane left"    "ctrl + g + n + h"
            "decrease pane top"     "ctrl + g + n + K"
            "decrease pane right"   "ctrl + g + n + L"
            "decrease pane bottom"  "ctrl + g + n + J"
            "decrease pane left"    "ctrl + g + n + H"
            "move pane to top"      "ctrl + g + h + k"
            "move pane to right"    "ctrl + g + h + l"
            "move pane to bottom"   "ctrl + g + h + j"
            "move pane to left"     "ctrl + g + h + h"
            "search"                "ctrl + g + s + s"
            "go into edit mode"     "ctrl + g + s + e"
            "detach session"        "ctrl + g + o + w"
            "open session manager"  "ctrl + g + o + w"
            "quit zellij"           "ctrl + g + q"
            floating true
            }
        }
        bind "Alt left" { MoveFocusOrTab "left"; }
        bind "Alt down" { MoveFocus "down"; }
        bind "Alt up" { MoveFocus "up"; }
        bind "Alt right" { MoveFocusOrTab "right"; }
        bind "Alt +" { Resize "Increase"; }
        bind "Alt -" { Resize "Decrease"; }
        bind "Alt =" { Resize "Increase"; }
        bind "Alt [" { PreviousSwapLayout; }
        bind "Alt ]" { NextSwapLayout; }
        bind "Alt f" { ToggleFloatingPanes; }
        bind "Alt h" { MoveFocusOrTab "left"; }
        bind "Alt i" { MoveTab "left"; }
        bind "Alt j" { MoveFocus "down"; }
        bind "Alt k" { MoveFocus "up"; }
        bind "Alt l" { MoveFocusOrTab "right"; }
        bind "Alt n" { NewPane; }
        bind "Alt o" { MoveTab "right"; }
    }
    shared_except "locked" "renametab" "renamepane" {
        bind "Ctrl g" { SwitchToMode "locked"; }
        bind "Ctrl q" { Quit; }
    }
    shared_except "locked" "entersearch" {
        bind "enter" { SwitchToMode "locked"; }
    }
    shared_except "locked" "entersearch" "renametab" "renamepane" {
        bind "esc" { SwitchToMode "locked"; }
    }
    shared_except "locked" "entersearch" "renametab" "renamepane" "move" {
        bind "m" { SwitchToMode "move"; }
    }
    shared_except "locked" "tab" "entersearch" "renametab" "renamepane" {
        bind "t" { SwitchToMode "tab"; }
    }
    shared_except "locked" "tab" "scroll" "entersearch" "renametab" "renamepane" {
        bind "s" { SwitchToMode "scroll"; }
    }
    shared_among "normal" "resize" "tab" "scroll" "prompt" "tmux" {
        bind "p" { SwitchToMode "pane"; }
    }
    shared_among "normal" "resize" "tab" "scroll" "prompt" "tmux" {
        bind "o" { SwitchToMode "session"; }
    }
    shared_except "locked" "resize" "pane" "tab" "entersearch" "renametab" "renamepane" {
        bind "r" { SwitchToMode "resize"; }
    }
    shared_among "scroll" "search" {
        bind "PageDown" { PageScrollDown; }
        bind "PageUp" { PageScrollUp; }
        bind "left" { PageScrollUp; }
        bind "down" { ScrollDown; }
        bind "up" { ScrollUp; }
        bind "right" { PageScrollDown; }
        bind "Ctrl b" { PageScrollUp; }
        bind "Ctrl c" { ScrollToBottom; SwitchToMode "locked"; }
        bind "d" { HalfPageScrollDown; }
        bind "Ctrl f" { PageScrollDown; }
        bind "h" { PageScrollUp; }
        bind "j" { ScrollDown; }
        bind "k" { ScrollUp; }
        bind "l" { PageScrollDown; }
        bind "u" { HalfPageScrollUp; }
    }
    entersearch {
        bind "Ctrl c" { SwitchToMode "scroll"; }
        bind "esc" { SwitchToMode "scroll"; }
        bind "enter" { SwitchToMode "search"; }
    }
    renametab {
        bind "esc" { UndoRenameTab; SwitchToMode "tab"; }
    }
    shared_among "renametab" "renamepane" {
        bind "Ctrl c" { SwitchToMode "locked"; }
    }
    renamepane {
        bind "esc" { UndoRenamePane; SwitchToMode "pane"; }
    }
}

plugins {
    compact-bar location="zellij:compact-bar"
    configuration location="zellij:configuration"
    filepicker location="zellij:strider" {
        cwd "/"
    }
    plugin-manager location="zellij:plugin-manager"
    session-manager location="zellij:session-manager"
    status-bar location="zellij:status-bar"
    strider location="zellij:strider"
    tab-bar location="zellij:tab-bar"
    welcome-screen location="zellij:session-manager" {
        welcome_screen true
    }
}

themes {
    kanagawa_dragon {
        fg "#dcdccc"
        bg "#282828"
        red "#C34043"
        green "#728F66"
        yellow "#9D8F6F"
        blue "#8BA4B0"
        magenta "#8BA4B0"
        cyan "#8be9fd"
        orange "#ffb86c"
        black "#1a1a1a"
        white "#8D909D"
    } 

    everforest {
        fg "#d3c6aa"
        bg "#282828"
        red "#e67e80"
        green "#a7c080"
        yellow "#dbbc7f"
        blue "#7fbbb3"
        magenta "#d699b6"
        cyan "#83c092"
        orange "#e69875"
        black "#1a1a1a"
        white "#8D908D"
    }

    rose_pine_moon {
        fg "#e0def4"
        bg "#191724"
        red "#eb6f92"
        green "#31748f"
        yellow "#f6c177"
        blue "#9ccfd8"
        magenta "#c4a7e7"
        cyan "#9ccfd8"
        orange "#f6c177"
        black "#191724"
        white "#e0def4"
    }

    sakura {
      fg "#c5a3a9"
      bg "#1c1a1c"
      red "#2B1720"
      green "#878fb9"
      yellow "#9e97d0"
      blue "#878fb9"
      magenta "#9e97d0"
      cyan "#878fb9"
      orange "#9e97d0"
      black "#1c1a1c"
      white "#c5a3a9"
    }

    oldWorld {
      fg "#C9C7CD"
      bg "#000000"
      red "#EA83A5"
      green "#90B99F"
      yellow "#E6B99D"
      blue "#85B5BA"
      magenta "#92A2D5"
      cyan "#85B5BA"
      orange "#F5A191"
      black "#000000"
      white "#C9C7CD"
    }
}
theme "catppuccin-mocha"
default_mode "locked"
scrollback_editor "nvim"
default_layout "work_oldWorld"
      '';
    };
    ".config/zellij/layouts/work_oldWorld.kdl" = {
      text = ''
layout {
    tab name="nvim" focus=true {
        pane
    }

    tab name="shell" {
        pane
    }

    default_tab_template {
        pane size=1 borderless=true {
            plugin location="file:~/.config/zellij/plugins/zjstatus.wasm" {
                format_left   "{mode} #[fg=#E29ECA,bold]{session}{tabs}"
                format_right  "{command_git_branch} {datetime}"
                format_space  ""

                border_enabled  "false"
                border_char     "─"
                border_format   "#[fg=#161617]{char}"
                border_position "top"

                hide_frame_for_single_pane "true"
                mode_normal  "#[bg=#85B5BA] "
                mode_tmux    "#[bg=#EA83A5] "

                tab_normal   "#[fg=#C9C7CD] {name} "
                tab_active   "#[fg=#92A2D5,bold,italic] {name} "

                command_git_branch_command     "git rev-parse --abbrev-ref HEAD"
                command_git_branch_format      "#[fg=#85B5BA] {stdout} "
                command_git_branch_interval    "10"
                command_git_branch_rendermode  "static"

                datetime        "#[fg=#C9C7CD,bold] {format} "
                datetime_format "%A, %d %b %Y %H:%M"
                datetime_timezone "Europe/Berlin"
            }
        }
        children
        pane size=1 borderless=true  {
            plugin location="zellij:status-bar"
        }
    }
}
      '';
    };
  };
}
