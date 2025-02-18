{ config, pkgs, ... }:
let
  # Detección simple de si estamos en macOS o Linux
  systemType = if pkgs.stdenv.isDarwin then "mac" else "linux";

  nushellConfigText = ''
let dark_theme = {
    separator: "#C9C7CD"
    leading_trailing_space_bg: { attr: "n" }
    header: "#92A2D5_bold"
    empty: "#ACA1CF"
    bool: "#E29ECA"
    int: "#C9C7CD"
    filesize: "#85B5BA"
    duration: "#90B99F"
    date: "#E6B99D"
    range: "#C9C7CD"
    float: "#EA83A5"
    string: "#C9C7CD"
    nothing: "#92A2D5"
    binary: "#85B5BA"
    cellpath: "#90B99F"
    row_index: "#92A2D5_bold"
    record: "#ACA1CF"
    list: "#C9C7CD"
    block: "#ACA1CF_bold"
    hints: "#90B99F"
    search_result: { fg: "#000000", bg: "#F5A191" }
    shape_and: "#ACA1CF_bold"
    shape_binary: "#85B5BA_bold"
    shape_block: "#92A2D5"
    shape_bool: "#E29ECA"
    shape_closure: "#EA83A5"
    shape_custom: "#85B5BA"
    shape_datetime: "#E6B99D_bold"
    shape_directory: "#92A2D5"
    shape_external: "#85B5BA"
    shape_externalarg: "#ACA1CF_bold"
    shape_filepath: "#90B99F"
    shape_flag: "#92A2D5_bold"
    shape_float: "#EA83A5"
    shape_garbage: { fg: "#000000", bg: "#EA83A5", attr: "b" }
    shape_globpattern: "#85B5BA_bold"
    shape_int: "#ACA1CF"
    shape_internalcall: "#85B5BA_bold"
    shape_keyword: "#92A2D5"
    shape_literal: "#E6B99D"
    shape_operator: "#F5A191"
    shape_or: "#E29ECA_bold"
    shape_pipe: "#85B5BA"
    shape_string: "#90B99F"
    shape_variable: "#EA83A5"
}

$env.LS_COLORS = (
    "di=38;2;146;162;213:" +
    "fi=38;2;201;199;205:" +
    "ln=38;2;172;161;207:" +
    "ex=38;2;133;181;186:" +
    "or=38;2;234;131;165:" +
    "*.txt=38;2;201;199;205:" +
    "*.jpg=38;2;172;161;207:" +
    "*.png=38;2;172;161;207:" +
    "*.zip=38;2;133;181;186:" +
    "*.gz=38;2;133;181;186:" +
    "*.tar=38;2;133;181;186:" +
    "*.log=38;2;229;158;202:" +
    "*.md=38;2;229;158;202:" +
    "*.py=38;2;133;181;186:" +
    "*.rs=38;2;234;131;165:" +
    "*.sh=38;2;133;181;186:" +
    "*=38;2;201;199;205"
)

$env.config = {
    show_banner: false
    ls: {
        use_ls_colors: true
        clickable_links: true
    }
    rm: { always_trash: false }
    table: {
        mode: rounded
        index_mode: always
        show_empty: true
        padding: { left: 1, right: 1 }
        trim: {
            methodology: wrapping
            wrapping_try_keep_words: true
            truncating_suffix: "..."
        }
        header_on_separator: false
    }
    error_style: "fancy"
    display_errors: {
        exit_code: false
        termination_signal: true
    }
    datetime_format: {}
    explore: {
        status_bar_background: { fg: "#1D1F21", bg: "#C4C9C6" }
        command_bar_text: { fg: "#C4C9C6" }
        highlight: { fg: "black", bg: "yellow" }
        status: {
            error: { fg: "white", bg: "red" }
            warn: {}
            info: {}
        }
        selected_cell: { bg: light_blue }
    }
    history: {
        max_size: 100_000
        sync_on_enter: true
        file_format: "plaintext"
        isolation: false
    }
    completions: {
        case_sensitive: false
        quick: true
        partial: true
        algorithm: "prefix"
        sort: "smart"
        external: {
            enable: true
            max_results: 100
            completer: null
        }
        use_ls_colors: true
    }
    filesize: {
        unit: "MB"
    }
    cursor_shape: {
        emacs: line
        vi_insert: block
        vi_normal: underscore
    }
    color_config: $dark_theme
    footer_mode: 25
    float_precision: 2
    buffer_editor: null
    use_ansi_coloring: true
    bracketed_paste: true
    edit_mode: vi
    shell_integration: {
        osc2: true
        osc7: true
        osc8: true
        osc9_9: false
        osc133: false
        osc633: true
        reset_application_mode: true
    }
    render_right_prompt_on_last_line: false
    use_kitty_protocol: false
    highlight_resolved_externals: false
    recursion_limit: 50
    plugins: {}
    plugin_gc: {
        default: {
            enabled: true
            stop_after: 10sec
        }
        plugins: {}
    }
    hooks: {
        pre_prompt: [{ null }]
        pre_execution: [{ null }]
        env_change: {
            PWD: [{|before, after| null }]
        }
        display_output: "if (term size).columns >= 100 { table -e } else { table }"
        command_not_found: { null }
    }
    menus: []
    keybindings: []
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

let MULTIPLEXER = "zellij"
let MULTIPLEXER_ENV_PREFIX = "ZELLIJ"

def start_multiplexer [] {
  if $MULTIPLEXER_ENV_PREFIX not-in ($env | columns) {
    run-external $MULTIPLEXER
  }
}

start_multiplexer
  '';

  nushellEnvText = ''
def create_left_prompt [] {
    let dir = match (do --ignore-shell-errors { $env.PWD | path relative-to $nu.home-path }) {
        null => $env.PWD
        "" => "~"
        $relative_pwd => ([~ $relative_pwd] | path join)
    }

    let path_color = (if (is-admin) { ansi red_bold } else { ansi green_bold })
    let separator_color = (if (is-admin) { ansi light_red_bold } else { ansi light_green_bold })
    let path_segment = $"($path_color)($dir)(ansi reset)"

    $path_segment | str replace --all (char path_sep) $"(ansi reset)$separator_color(char path_sep)$path_color"
}

def create_right_prompt [] {
    let time_segment = ([
        (ansi reset)
        (ansi magenta)
        (date now | format date '%x %X')
    ] | str join | 
    str replace --regex --all "([/:])" { |r | (ansi green) + (r.captures.1 | to-string) + (ansi magenta) } | str replace --regex --all "([AP]M)" { |r | (ansi magenta_underline) + (r.captures.1 | to-string)
})

    let last_exit_code = if ($env.LAST_EXIT_CODE != 0) {
      ([
        (ansi rb)
        ($env.LAST_EXIT_CODE | to-string)
      ] | str join)
    } else {
      ""
    }

    ([$last_exit_code, (char space), $time_segment] | str join)
}

$env.PROMPT_COMMAND = {|| create_left_prompt }
$env.PROMPT_COMMAND_RIGHT = {|| create_right_prompt }

$env.PROMPT_INDICATOR = {|| "> " }
$env.PROMPT_INDICATOR_VI_INSERT = {|| ": " }
$env.PROMPT_INDICATOR_VI_NORMAL = {|| "> " }
$env.PROMPT_MULTILINE_INDICATOR = {|| "::: " }

$env.ENV_CONVERSIONS = {
    "PATH": {
        from_string: { |s| $s | split row (char esep) | path expand --no-symlink }
        to_string: { |v| $v | path expand --no-symlink | str join (char esep) }
    }
    "Path": {
        from_string: { |s| $s | split row (char esep) | path expand --no-symlink }
        to_string: { |v| $v | path expand --no-symlink | str join (char esep) }
    }
}

$env.NU_LIB_DIRS = [
    ($nu.default-config-dir | path join 'scripts')
    ($nu.data-dir | path join 'completions')
]

$env.NU_PLUGIN_DIRS = [
    ($nu.default-config-dir | path join 'plugins')
]

$env.EDITOR = "nvim"

$env.PATH = (
    $env.PATH
    | split row (char esep)
    | prepend '/opt/homebrew/bin'
    | prepend '/nix/var/nix/profiles/per-user/($env.USER)/home-manager/bin'
    | prepend ($env.HOME | path join ".volta/bin")
    | prepend ($env.HOME | path join ".bun/bin")
    | prepend '/nix/var/nix/profiles/default/bin'
    | prepend '/Users/var/nix/profiles/default/bin'
    | append '/usr/local/bin'
    | append ($env.HOME | path join ".config")
    | append ($env.HOME | path join ".cargo/bin")
    | append '/usr/local/lib/*'
)

mkdir ~/.cache/starship
mkdir ~/.cache/carapace
mkdir ~/.local/share/atuin

$env.STARSHIP_CONFIG = $env.HOME | path join ".config/starship.toml"
$env.CARAPACE_BRIDGES = 'zsh,fish,bash,inshellisense'

starship init nu | save -f ~/.cache/starship/init.nu
zoxide init nushell | save -f ~/.zoxide.nu
atuin init nu | save -f ~/.local/share/atuin/init.nu
carapace _carapace nushell | save --force ~/.cache/carapace/init.nu
  '';
in {
  programs.nushell.enable = true;

  home.file = if systemType == "mac" then {
    "Library/Application Support/nushell/config.nu" = { text = nushellConfigText; };
    "Library/Application Support/nushell/env.nu" = { text = nushellEnvText; };
  } else {
    ".config/nushell/config.nu" = { text = nushellConfigText; };
    ".config/nushell/env.nu" = { text = nushellEnvText; };
  };
}
