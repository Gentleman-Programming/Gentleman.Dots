# Nushell Config File
#
# version = "0.99.1"

# For more information on defining custom themes, see
# https://www.nushell.sh/book/coloring_and_theming.html
# And here is the theme collection
# https://github.com/nushell/nu_scripts/tree/main/themes


# Sakura Theme
# let dark_theme = {
#    separator: "#786577"                # na: text (dark pink)
#    leading_trailing_space_bg: { attr: "n" }
#    header: "#878fb9_bold"              # va: blue + bold
#    empty: "#9e97d0"                    # ca: soft purple
#    bool: "#c58ea7"                     # ia: pink
#    int: "#786577"                      # na: text (dark pink)
#    filesize: "#878fb9"                 # va: blue
#    duration: "#786577"                 # na: text (dark pink)
#    date: "#9e97d0"                     # ca: purple
#    range: "#786577"                    # na: text (dark pink)
#    float: "#c58ea7"                    # ia: pink
#    string: "#786577"                   # na: text (dark pink)
#    nothing: "#878fb9"                  # va: blue
#    binary: "#786577"                   # na: text (dark pink)
#    cellpath: "#786577"                 # na: text (dark pink)
#    row_index: "#878fb9_bold"           # va: blue + bold
#    record: "#786577"                   # na: text (dark pink)
#    list: "#786577"                     # na: text (dark pink)
#    block: "#9e97d0_bold"               # ca: purple + bold
#    hints: "#3f3b3e"                    # gr: dark gray
#    search_result: { bg: "#c58ea7", fg: "#786577" } # ia/na
#    shape_and: "#9e97d0_bold"           # ca: purple + bold
#    shape_binary: "#9e97d0_bold"        # ca: purple + bold
#    shape_block: "#9e97d0_bold"         # ca: purple + bold
#    shape_bool: "#c58ea7"               # ia: pink
#    shape_closure: "#878fb9_bold"       # va: blue + bold
#    shape_custom: "#878fb9"             # va: blue
#    shape_datetime: "#878fb9_bold"      # va: blue + bold
#    shape_directory: "#9e97d0"          # ca: purple
#    shape_external: "#878fb9"           # va: blue
#    shape_externalarg: "#878fb9_bold"   # va: blue + bold
#    shape_external_resolved: "#9e97d0_bold" # ca: purple + bold
#    shape_filepath: "#878fb9"           # va: blue
#    shape_flag: "#9e97d0_bold"          # ca: purple + bold
#    shape_float: "#9e97d0_bold"         # ca: purple + bold
#    shape_garbage: { fg: "#786577", bg: "#c58ea7", attr: "b" } # na/ia
#    shape_glob_interpolation: "#878fb9_bold" # va: blue + bold
#    shape_globpattern: "#878fb9_bold"   # va: blue + bold
#    shape_int: "#9e97d0_bold"           # ca: purple + bold
#    shape_internalcall: "#878fb9_bold"  # va: blue + bold
#    shape_keyword: "#878fb9_bold"       # va: blue + bold
#    shape_list: "#878fb9_bold"          # va: blue + bold
#    shape_literal: "#9e97d0"            # ca: purple
#    shape_match_pattern: "#878fb9"      # va: blue
#    shape_matching_brackets: { attr: "u" }
#    shape_nothing: "#878fb9"            # va: blue
#    shape_operator: "#9e97d0"           # ca: purple (replaces yellow)
#    shape_or: "#9e97d0_bold"            # ca: purple + bold
#    shape_pipe: "#9e97d0_bold"          # ca: purple + bold
#    shape_range: "#9e97d0_bold"         # ca: purple + bold
#    shape_record: "#878fb9_bold"        # va: blue + bold
#    shape_redirection: "#9e97d0_bold"   # ca: purple + bold
#    shape_signature: "#878fb9_bold"     # va: blue + bold
#    shape_string: "#878fb9"             # va: blue
#    shape_string_interpolation: "#878fb9_bold" # va: blue + bold
#    shape_table: "#9e97d0_bold"         # ca: purple + bold
#    shape_variable: "#c58ea7"           # ia: pink
#    shape_vardecl: "#c58ea7"            # ia: pink
#    shape_raw_string: "#9e97d0"         # ca: purple
# 

let dark_theme = {
    # --- base elements ---
    separator: "#54546D"                     # table borders (wave gray)
    leading_trailing_space_bg: { attr: "n" } # spaces without background
    header: "#7E9CD8_bold"                   # headers (wave blue + bold)
    empty: "#957FB8"                         # empty elements (wave purple)
    bool: "#D27E99"                          # booleans (wave pink)
    int: "#54546D"                           # integers (wave gray)
    filesize: "#6A9589"                      # file sizes (wave teal)
    duration: "#98BB6C"                      # duration (wave green)
    date: "#E6C384"                          # dates (wave beige)
    range: "#54546D"                         # ranges (wave gray)
    float: "#DCA561"                         # floats (wave orange)
    string: "#54546D"                        # general text (wave gray)
    nothing: "#7E9CD8"                       # null values (wave blue)
    binary: "#6A9589"                        # binaries (wave teal)
    cellpath: "#98BB6C"                      # cell paths (wave green)
    row_index: "#7E9CD8_bold"                # row indices (wave blue + bold)
    record: "#957FB8"                        # records (wave purple)
    list: "#54546D"                          # lists (wave gray)
    block: "#957FB8_bold"                    # blocks (wave purple + bold)
    hints: "#98BB6C"                         # hints (wave green)
    search_result: { fg: "#1F1F28", bg: "#E46876" } # search result (wave red background)

    # --- syntax elements/commands ---
    shape_and: "#957FB8_bold"                # AND operator (wave purple + bold)
    shape_binary: "#6A9589_bold"             # binaries (wave teal + bold)
    shape_block: "#7E9CD8"                   # blocks (wave blue)
    shape_bool: "#D27E99"                    # booleans (wave pink)
    shape_closure: "#DCA561"                 # closures (wave orange)
    shape_custom: "#6A9589"                  # custom commands (wave teal)
    shape_datetime: "#E6C384_bold"           # dates (wave beige + bold)
    shape_directory: "#7E9CD8"               # directories (wave blue)
    shape_external: "#6A9589"                # external commands (wave teal)
    shape_externalarg: "#957FB8_bold"        # external arguments (wave purple + bold)
    shape_filepath: "#98BB6C"                # file paths (wave green)
    shape_flag: "#7E9CD8_bold"               # flags (wave blue + bold)
    shape_float: "#DCA561"                   # floats (wave orange)
    shape_garbage: { fg: "#1F1F28", bg: "#DCA561", attr: "b" } # error (wave orange background)
    shape_globpattern: "#6A9589_bold"        # glob patterns (wave teal + bold)
    shape_int: "#957FB8"                     # integers (wave purple)
    shape_internalcall: "#6A9589_bold"       # internal calls (wave teal + bold)
    shape_keyword: "#7E9CD8"                 # keywords (wave blue)
    shape_literal: "#E6C384"                 # literals (wave beige)
    shape_operator: "#E46876"                # operators (wave red)
    shape_or: "#D27E99_bold"                 # OR operator (wave pink + bold)
    shape_pipe: "#6A9589"                    # pipes (wave teal)
    shape_string: "#98BB6C"                  # strings (wave green)
    shape_variable: "#DCA561"                # variables (wave orange)
}

# $env.LS_COLORS = (
#    "di=38;2;197;142;167:" +       # Directories: intense pink (#C58EA7)
#    "fi=38;2;197;163;169:" +       # Regular files: light pink (#C5A3A9)
#    "ln=38;2;158;151;208:" +       # Symbolic links: soft purple (#9E97D0)
#    "ex=38;2;135;143;185:" +       # Executable files: light blue (#878FB9)
#    "or=38;2;197;142;167:" +       # Broken symbolic links: intense pink (#C58EA7)
#    "*.txt=38;2;197;163;169:" +    # .txt files: light pink (#C5A3A9)
#    "*.jpg=38;2;158;151;208:" +    # .jpg files: soft purple (#9E97D0)
#    "*.png=38;2;158;151;208:" +    # .png files: soft purple (#9E97D0)
#    "*.zip=38;2;135;143;185:" +    # .zip files: light blue (#878FB9)
#    "*.gz=38;2;135;143;185:" +     # .gz files: light blue (#878FB9)
#    "*.tar=38;2;135;143;185:" +    # .tar files: light blue (#878FB9)
#    "*.log=38;2;63;59;62:" +       # .log files: dark gray (#3F3B3E)
#    "*.md=38;2;197;163;169:" +     # .md files: light pink (#C5A3A9)
#    "*.py=38;2;135;143;185:" +     # .py files: light blue (#878FB9)
#    "*.rs=38;2;197;142;167:" +     # .rs files: intense pink (#C58EA7)
#    "*.sh=38;2;135;143;185:"       # .sh files: light blue (#878FB9)
# )

$env.LS_COLORS = (
    # --- Directories and file types ---
    "di=38;2;146;162;213:" +       # Directories: lavender blue (#92A2D5)
    "fi=38;2;201;199;205:" +       # Regular files: light gray (#C9C7CD)
    "ln=38;2;172;161;207:" +       # Symbolic links: lilac gray (#ACA1CF)
    "ex=38;2;133;181;186:" +       # Executable files: blue-green (#85B5BA)
    "or=38;2;234;131;165:" +       # Broken links: intense pink (#EA83A5)

    # --- Specific extensions ---
    "*.txt=38;2;201;199;205:" +    # .txt: light gray (#C9C7CD)
    "*.jpg=38;2;172;161;207:" +    # .jpg: lilac gray (#ACA1CF)
    "*.png=38;2;172;161;207:" +    # .png: lilac gray (#ACA1CF)
    "*.zip=38;2;133;181;186:" +    # .zip: blue-green (#85B5BA)
    "*.gz=38;2;133;181;186:" +     # .gz: blue-green (#85B5BA)
    "*.tar=38;2;133;181;186:" +    # .tar: blue-green (#85B5BA)
    "*.log=38;2;229;158;202:" +    # .log: soft pink (#E29ECA)
    "*.md=38;2;229;158;202:" +     # .md: soft pink (#E29ECA)
    "*.py=38;2;133;181;186:" +     # .py: blue-green (#85B5BA)
    "*.rs=38;2;234;131;165:" +     # .rs: intense pink (#EA83A5)
    "*.sh=38;2;133;181;186:" +     # .sh: blue-green (#85B5BA)

    # --- Default color for other files ---
    "*=38;2;201;199;205"           # Default: light gray (#C9C7CD)
)

let light_theme = {
    # color for nushell primitives
    separator: dark_gray
    leading_trailing_space_bg: { attr: n } # no fg, no bg, attr none effectively turns this off
    header: green_bold
    empty: blue
    # Closures can be used to choose colors for specific values.
    # The value (in this case, a bool) is piped into the closure.
    # eg) {|| if $in { 'dark_cyan' } else { 'dark_gray' } }
    bool: dark_cyan
    int: dark_gray
    filesize: cyan_bold
    duration: dark_gray
    date: purple
    range: dark_gray
    float: dark_gray
    string: dark_gray
    nothing: dark_gray
    binary: dark_gray
    cell-path: dark_gray
    row_index: green_bold
    record: dark_gray
    list: dark_gray
    block: dark_gray
    hints: dark_gray
    search_result: { fg: white bg: red }
    shape_and: purple_bold
    shape_binary: purple_bold
    shape_block: blue_bold
    shape_bool: light_cyan
    shape_closure: green_bold
    shape_custom: green
    shape_datetime: cyan_bold
    shape_directory: cyan
    shape_external: cyan
    shape_externalarg: green_bold
    shape_external_resolved: light_purple_bold
    shape_filepath: cyan
    shape_flag: blue_bold
    shape_float: purple_bold
    # shapes are used to change the cli syntax highlighting
    shape_garbage: { fg: white bg: red attr: b }
    shape_glob_interpolation: cyan_bold
    shape_globpattern: cyan_bold
    shape_int: purple_bold
    shape_internalcall: cyan_bold
    shape_keyword: cyan_bold
    shape_list: cyan_bold
    shape_literal: blue
    shape_match_pattern: green
    shape_matching_brackets: { attr: u }
    shape_nothing: light_cyan
    shape_operator: yellow
    shape_or: purple_bold
    shape_pipe: purple_bold
    shape_range: yellow_bold
    shape_record: cyan_bold
    shape_redirection: purple_bold
    shape_signature: green_bold
    shape_string: green
    shape_string_interpolation: cyan_bold
    shape_table: blue_bold
    shape_variable: purple
    shape_vardecl: purple
    shape_raw_string: light_purple
}

# External completer example
# let carapace_completer = {|spans|
#     carapace $spans.0 nushell ...$spans | from json
# }

# The default config record. This is where much of your global configuration is setup.
$env.config = {
    show_banner: false # true or false to enable or disable the welcome banner at startup

    ls: {
        use_ls_colors: true # use the LS_COLORS environment variable to colorize output
        clickable_links: true # enable or disable clickable links. Your terminal has to support links.
    }

    rm: {
        always_trash: false # always act as if -t was given. Can be overridden with -p
    }

    table: {
        mode: rounded # basic, compact, compact_double, light, thin, with_love, rounded, reinforced, heavy, none, other
        index_mode: always # "always" show indexes, "never" show indexes, "auto" = show indexes when a table has "index" column
        show_empty: true # show 'empty list' and 'empty record' placeholders for command output
        padding: { left: 1, right: 1 } # a left right padding of each column in a table
        trim: {
            methodology: wrapping # wrapping or truncating
            wrapping_try_keep_words: true # A strategy used by the 'wrapping' methodology
            truncating_suffix: "..." # A suffix used by the 'truncating' methodology
        }
        header_on_separator: false # show header text on separator/border line
        # abbreviated_row_count: 10 # limit data rows from top and bottom after reaching a set point
    }

    error_style: "fancy" # "fancy" or "plain" for screen reader-friendly error messages

    # Whether an error message should be printed if an error of a certain kind is triggered.
    display_errors: {
        exit_code: false # assume the external command prints an error message
        # Core dump errors are always printed, and SIGPIPE never triggers an error.
        # The setting below controls message printing for termination by all other signals.
        termination_signal: true
    }

    # datetime_format determines what a datetime rendered in the shell would look like.
    # Behavior without this configuration point will be to "humanize" the datetime display,
    # showing something like "a day ago."
    datetime_format: {
        # normal: '%a, %d %b %Y %H:%M:%S %z'    # shows up in displays of variables or other datetime's outside of tables
        # table: '%m/%d/%y %I:%M:%S%p'          # generally shows up in tabular outputs such as ls. commenting this out will change it to the default human readable datetime format
    }

    explore: {
        status_bar_background: { fg: "#1D1F21", bg: "#C4C9C6" },
        command_bar_text: { fg: "#C4C9C6" },
        highlight: { fg: "black", bg: "yellow" },
        status: {
            error: { fg: "white", bg: "red" },
            warn: {}
            info: {}
        },
        selected_cell: { bg: light_blue },
    }

    history: {
        max_size: 100_000 # Session has to be reloaded for this to take effect
        sync_on_enter: true # Enable to share history between multiple sessions, else you have to close the session to write history to file
        file_format: "plaintext" # "sqlite" or "plaintext"
        isolation: false # only available with sqlite file_format. true enables history isolation, false disables it. true will allow the history to be isolated to the current session using up/down arrows. false will allow the history to be shared across all sessions.
    }

    completions: {
        case_sensitive: false # set to true to enable case-sensitive completions
        quick: true    # set this to false to prevent auto-selecting completions when only one remains
        partial: true    # set this to false to prevent partial filling of the prompt
        algorithm: "prefix"    # prefix or fuzzy
        sort: "smart" # "smart" (alphabetical for prefix matching, fuzzy score for fuzzy matching) or "alphabetical"
        external: {
            enable: true # set to false to prevent nushell looking into $env.PATH to find more suggestions, `false` recommended for WSL users as this look up may be very slow
            max_results: 100 # setting it lower can improve completion performance at the cost of omitting some options
            completer: null # check 'carapace_completer' above as an example
        }
        use_ls_colors: true # set this to true to enable file/path/directory completions using LS_COLORS
    }

    filesize: {
        unit: "MB" # b, kb, kib, mb, mib, gb, gib, tb, tib, pb, pib, eb, eib, auto
    }

    cursor_shape: {
        emacs: line # block, underscore, line, blink_block, blink_underscore, blink_line, inherit to skip setting cursor shape (line is the default)
        vi_insert: block # block, underscore, line, blink_block, blink_underscore, blink_line, inherit to skip setting cursor shape (block is the default)
        vi_normal: underscore # block, underscore, line, blink_block, blink_underscore, blink_line, inherit to skip setting cursor shape (underscore is the default)
    }

    color_config: $dark_theme # if you want a more interesting theme, you can replace the empty record with `$dark_theme`, `$light_theme` or another custom record
    footer_mode: 25 # always, never, number_of_rows, auto
    float_precision: 2 # the precision for displaying floats in tables
    buffer_editor: null # command that will be used to edit the current line buffer with ctrl+o, if unset fallback to $env.EDITOR and $env.VISUAL
    use_ansi_coloring: true
    bracketed_paste: true # enable bracketed paste, currently useless on windows
    edit_mode: vi # emacs, vi
    shell_integration: {
        # osc2 abbreviates the path if in the home_dir, sets the tab/window title, shows the running command in the tab/window title
        osc2: true
        # osc7 is a way to communicate the path to the terminal, this is helpful for spawning new tabs in the same directory
        osc7: true
        # osc8 is also implemented as the deprecated setting ls.show_clickable_links, it shows clickable links in ls output if your terminal supports it. show_clickable_links is deprecated in favor of osc8
        osc8: true
        # osc9_9 is from ConEmu and is starting to get wider support. It's similar to osc7 in that it communicates the path to the terminal
        osc9_9: false
        # osc133 is several escapes invented by Final Term which include the supported ones below.
        # 133;A - Mark prompt start
        # 133;B - Mark prompt end
        # 133;C - Mark pre-execution
        # 133;D;exit - Mark execution finished with exit code
        # This is used to enable terminals to know where the prompt is, the command is, where the command finishes, and where the output of the command is
        osc133: false
        # osc633 is closely related to osc133 but only exists in visual studio code (vscode) and supports their shell integration features
        # 633;A - Mark prompt start
        # 633;B - Mark prompt end
        # 633;C - Mark pre-execution
        # 633;D;exit - Mark execution finished with exit code
        # 633;E - Explicitly set the command line with an optional nonce
        # 633;P;Cwd=<path> - Mark the current working directory and communicate it to the terminal
        # and also helps with the run recent menu in vscode
        osc633: true
        # reset_application_mode is escape \x1b[?1l and was added to help ssh work better
        reset_application_mode: true
    }
    render_right_prompt_on_last_line: false # true or false to enable or disable right prompt to be rendered on last line of the prompt.
    use_kitty_protocol: false # enables keyboard enhancement protocol implemented by kitty console, only if your terminal support this.
    highlight_resolved_externals: false # true enables highlighting of external commands in the repl resolved by which.
    recursion_limit: 50 # the maximum number of times nushell allows recursion before stopping it

    plugins: {} # Per-plugin configuration. See https://www.nushell.sh/contributor-book/plugins.html#configuration.

    plugin_gc: {
        # Configuration for plugin garbage collection
        default: {
            enabled: true # true to enable stopping of inactive plugins
            stop_after: 10sec # how long to wait after a plugin is inactive to stop it
        }
        plugins: {
            # alternate configuration for specific plugins, by name, for example:
            #
            # gstat: {
            #     enabled: false
            # }
        }
    }

    hooks: {
        pre_prompt: [{ null }] # run before the prompt is shown
        pre_execution: [{ null }] # run before the repl input is run
        env_change: {
            PWD: [{|before, after| null }] # run if the PWD environment is different since the last repl input
        }
        display_output: "if (term size).columns >= 100 { table -e } else { table }" # run to display the output of a pipeline
        command_not_found: { null } # return an error message when a command is not found
    }

    menus: [
        # Configuration for default nushell menus
        # Note the lack of source parameter
        {
            name: completion_menu
            only_buffer_difference: false
            marker: "| "
            type: {
                layout: columnar
                columns: 4
                col_width: 20     # Optional value. If missing all the screen width is used to calculate column width
                col_padding: 2
            }
            style: {
                text: green
                selected_text: { attr: r }
                description_text: yellow
                match_text: { attr: u }
                selected_match_text: { attr: ur }
            }
        }
        {
            name: ide_completion_menu
            only_buffer_difference: false
            marker: "| "
            type: {
                layout: ide
                min_completion_width: 0,
                max_completion_width: 50,
                max_completion_height: 10, # will be limited by the available lines in the terminal
                padding: 0,
                border: true,
                cursor_offset: 0,
                description_mode: "prefer_right"
                min_description_width: 0
                max_description_width: 50
                max_description_height: 10
                description_offset: 1
                # If true, the cursor pos will be corrected, so the suggestions match up with the typed text
                #
                # C:\> str
                #      str join
                #      str trim
                #      str split
                correct_cursor_pos: false
            }
            style: {
                text: green
                selected_text: { attr: r }
                description_text: yellow
                match_text: { attr: u }
                selected_match_text: { attr: ur }
            }
        }
        {
            name: history_menu
            only_buffer_difference: true
            marker: "? "
            type: {
                layout: list
                page_size: 10
            }
            style: {
                text: green
                selected_text: green_reverse
                description_text: yellow
            }
        }
        {
            name: help_menu
            only_buffer_difference: true
            marker: "? "
            type: {
                layout: description
                columns: 4
                col_width: 20     # Optional value. If missing all the screen width is used to calculate column width
                col_padding: 2
                selection_rows: 4
                description_rows: 10
            }
            style: {
                text: green
                selected_text: green_reverse
                description_text: yellow
            }
        }
    ]

    keybindings: [
        {
            name: completion_menu
            modifier: none
            keycode: tab
            mode: [emacs vi_normal vi_insert]
            event: {
                until: [
                    { send: menu name: completion_menu }
                    { send: menunext }
                    { edit: complete }
                ]
            }
        }
        {
            name: completion_previous_menu
            modifier: shift
            keycode: backtab
            mode: [emacs, vi_normal, vi_insert]
            event: { send: menuprevious }
        }
        {
            name: ide_completion_menu
            modifier: control
            keycode: space
            mode: [emacs vi_normal vi_insert]
            event: {
                until: [
                    { send: menu name: ide_completion_menu }
                    { send: menunext }
                    { edit: complete }
                ]
            }
        }
        {
            name: history_menu
            modifier: control
            keycode: char_r
            mode: [emacs, vi_insert, vi_normal]
            event: { send: menu name: history_menu }
        }
        {
            name: help_menu
            modifier: none
            keycode: f1
            mode: [emacs, vi_insert, vi_normal]
            event: { send: menu name: help_menu }
        }
        {
            name: next_page_menu
            modifier: control
            keycode: char_x
            mode: emacs
            event: { send: menupagenext }
        }
        {
            name: undo_or_previous_page_menu
            modifier: control
            keycode: char_z
            mode: emacs
            event: {
                until: [
                    { send: menupageprevious }
                    { edit: undo }
                ]
            }
        }
        {
            name: escape
            modifier: none
            keycode: escape
            mode: [emacs, vi_normal, vi_insert]
            event: { send: esc }    # NOTE: does not appear to work
        }
        {
            name: cancel_command
            modifier: control
            keycode: char_c
            mode: [emacs, vi_normal, vi_insert]
            event: { send: ctrlc }
        }
        {
            name: quit_shell
            modifier: control
            keycode: char_d
            mode: [emacs, vi_normal, vi_insert]
            event: { send: ctrld }
        }
        {
            name: clear_screen
            modifier: control
            keycode: char_l
            mode: [emacs, vi_normal, vi_insert]
            event: { send: clearscreen }
        }
        {
            name: search_history
            modifier: control
            keycode: char_q
            mode: [emacs, vi_normal, vi_insert]
            event: { send: searchhistory }
        }
        {
            name: open_command_editor
            modifier: control
            keycode: char_o
            mode: [emacs, vi_normal, vi_insert]
            event: { send: openeditor }
        }
        {
            name: move_up
            modifier: none
            keycode: up
            mode: [emacs, vi_normal, vi_insert]
            event: {
                until: [
                    { send: menuup }
                    { send: up }
                ]
            }
        }
        {
            name: move_down
            modifier: none
            keycode: down
            mode: [emacs, vi_normal, vi_insert]
            event: {
                until: [
                    { send: menudown }
                    { send: down }
                ]
            }
        }
        {
            name: move_left
            modifier: none
            keycode: left
            mode: [emacs, vi_normal, vi_insert]
            event: {
                until: [
                    { send: menuleft }
                    { send: left }
                ]
            }
        }
        {
            name: move_right_or_take_history_hint
            modifier: none
            keycode: right
            mode: [emacs, vi_normal, vi_insert]
            event: {
                until: [
                    { send: historyhintcomplete }
                    { send: menuright }
                    { send: right }
                ]
            }
        }
        {
            name: move_one_word_left
            modifier: control
            keycode: left
            mode: [emacs, vi_normal, vi_insert]
            event: { edit: movewordleft }
        }
        {
            name: move_one_word_right_or_take_history_hint
            modifier: control
            keycode: right
            mode: [emacs, vi_normal, vi_insert]
            event: {
                until: [
                    { send: historyhintwordcomplete }
                    { edit: movewordright }
                ]
            }
        }
        {
            name: move_to_line_start
            modifier: none
            keycode: home
            mode: [emacs, vi_normal, vi_insert]
            event: { edit: movetolinestart }
        }
        {
            name: move_to_line_start
            modifier: control
            keycode: char_a
            mode: [emacs, vi_normal, vi_insert]
            event: { edit: movetolinestart }
        }
        {
            name: move_to_line_end_or_take_history_hint
            modifier: none
            keycode: end
            mode: [emacs, vi_normal, vi_insert]
            event: {
                until: [
                    { send: historyhintcomplete }
                    { edit: movetolineend }
                ]
            }
        }
        {
            name: move_to_line_end_or_take_history_hint
            modifier: control
            keycode: char_e
            mode: [emacs, vi_normal, vi_insert]
            event: {
                until: [
                    { send: historyhintcomplete }
                    { edit: movetolineend }
                ]
            }
        }
        {
            name: move_to_line_start
            modifier: control
            keycode: home
            mode: [emacs, vi_normal, vi_insert]
            event: { edit: movetolinestart }
        }
        {
            name: move_to_line_end
            modifier: control
            keycode: end
            mode: [emacs, vi_normal, vi_insert]
            event: { edit: movetolineend }
        }
        {
            name: move_down
            modifier: control
            keycode: char_n
            mode: [emacs, vi_normal, vi_insert]
            event: {
                until: [
                    { send: menudown }
                    { send: down }
                ]
            }
        }
        {
            name: move_up
            modifier: control
            keycode: char_p
            mode: [emacs, vi_normal, vi_insert]
            event: {
                until: [
                    { send: menuup }
                    { send: up }
                ]
            }
        }
        {
            name: delete_one_character_backward
            modifier: none
            keycode: backspace
            mode: [emacs, vi_insert]
            event: { edit: backspace }
        }
        {
            name: delete_one_word_backward
            modifier: control
            keycode: backspace
            mode: [emacs, vi_insert]
            event: { edit: backspaceword }
        }
        {
            name: delete_one_character_forward
            modifier: none
            keycode: delete
            mode: [emacs, vi_insert]
            event: { edit: delete }
        }
        {
            name: delete_one_character_forward
            modifier: control
            keycode: delete
            mode: [emacs, vi_insert]
            event: { edit: delete }
        }
        {
            name: delete_one_character_backward
            modifier: control
            keycode: char_h
            mode: [emacs, vi_insert]
            event: { edit: backspace }
        }
        {
            name: delete_one_word_backward
            modifier: control
            keycode: char_w
            mode: [emacs, vi_insert]
            event: { edit: backspaceword }
        }
        {
            name: move_left
            modifier: none
            keycode: backspace
            mode: vi_normal
            event: { edit: moveleft }
        }
        {
            name: newline_or_run_command
            modifier: none
            keycode: enter
            mode: emacs
            event: { send: enter }
        }
        {
            name: move_left
            modifier: control
            keycode: char_b
            mode: emacs
            event: {
                until: [
                    { send: menuleft }
                    { send: left }
                ]
            }
        }
        {
            name: move_right_or_take_history_hint
            modifier: control
            keycode: char_f
            mode: emacs
            event: {
                until: [
                    { send: historyhintcomplete }
                    { send: menuright }
                    { send: right }
                ]
            }
        }
        {
            name: redo_change
            modifier: control
            keycode: char_g
            mode: emacs
            event: { edit: redo }
        }
        {
            name: undo_change
            modifier: control
            keycode: char_z
            mode: emacs
            event: { edit: undo }
        }
        {
            name: paste_before
            modifier: control
            keycode: char_y
            mode: emacs
            event: { edit: pastecutbufferbefore }
        }
        {
            name: cut_word_left
            modifier: control
            keycode: char_w
            mode: emacs
            event: { edit: cutwordleft }
        }
        {
            name: cut_line_to_end
            modifier: control
            keycode: char_k
            mode: emacs
            event: { edit: cuttolineend }
        }
        {
            name: cut_line_from_start
            modifier: control
            keycode: char_u
            mode: emacs
            event: { edit: cutfromstart }
        }
        {
            name: swap_graphemes
            modifier: control
            keycode: char_t
            mode: emacs
            event: { edit: swapgraphemes }
        }
        {
            name: move_one_word_left
            modifier: alt
            keycode: left
            mode: emacs
            event: { edit: movewordleft }
        }
        {
            name: move_one_word_right_or_take_history_hint
            modifier: alt
            keycode: right
            mode: emacs
            event: {
                until: [
                    { send: historyhintwordcomplete }
                    { edit: movewordright }
                ]
            }
        }
        {
            name: move_one_word_left
            modifier: alt
            keycode: char_b
            mode: emacs
            event: { edit: movewordleft }
        }
        {
            name: move_one_word_right_or_take_history_hint
            modifier: alt
            keycode: char_f
            mode: emacs
            event: {
                until: [
                    { send: historyhintwordcomplete }
                    { edit: movewordright }
                ]
            }
        }
        {
            name: delete_one_word_forward
            modifier: alt
            keycode: delete
            mode: emacs
            event: { edit: deleteword }
        }
        {
            name: delete_one_word_backward
            modifier: alt
            keycode: backspace
            mode: emacs
            event: { edit: backspaceword }
        }
        {
            name: delete_one_word_backward
            modifier: alt
            keycode: char_m
            mode: emacs
            event: { edit: backspaceword }
        }
        {
            name: cut_word_to_right
            modifier: alt
            keycode: char_d
            mode: emacs
            event: { edit: cutwordright }
        }
        {
            name: upper_case_word
            modifier: alt
            keycode: char_u
            mode: emacs
            event: { edit: uppercaseword }
        }
        {
            name: lower_case_word
            modifier: alt
            keycode: char_l
            mode: emacs
            event: { edit: lowercaseword }
        }
        {
            name: capitalize_char
            modifier: alt
            keycode: char_c
            mode: emacs
            event: { edit: capitalizechar }
        }
        # The following bindings with `*system` events require that Nushell has
        # been compiled with the `system-clipboard` feature.
        # If you want to use the system clipboard for visual selection or to
        # paste directly, uncomment the respective lines and replace the version
        # using the internal clipboard.
        {
            name: copy_selection
            modifier: control_shift
            keycode: char_c
            mode: emacs
            event: { edit: copyselection }
            # event: { edit: copyselectionsystem }
        }
        {
            name: cut_selection
            modifier: control_shift
            keycode: char_x
            mode: emacs
            event: { edit: cutselection }
            # event: { edit: cutselectionsystem }
        }
        # {
        #     name: paste_system
        #     modifier: control_shift
        #     keycode: char_v
        #     mode: emacs
        #     event: { edit: pastesystem }
        # }
        {
            name: select_all
            modifier: control_shift
            keycode: char_a
            mode: emacs
            event: { edit: selectall }
        }
    ]
}

def fzfbat [] {
  fzf --preview "bat --theme=gruvbox-dark --color=always {}" 
}

def fzfnvim [] {
  nvim (fzf --preview "bat --theme=gruvbox-dark --color=always {}")
}

 source ~/.zoxide.nu
 source ~/.cache/carapace/init.nu
 source ~/.local/share/atuin/init.nu
 use ~/.cache/starship/init.nu
 use ~/.config/bash-env.nu

let MULTIPLEXER = "tmux" 
let MULTIPLEXER_ENV_PREFIX = "TMUX"

def start_multiplexer [] {
  if $MULTIPLEXER_ENV_PREFIX not-in ($env | columns) {
    run-external $MULTIPLEXER
  }
}

start_multiplexer
