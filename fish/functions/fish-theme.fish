function fish-theme --description "Select the Gentleman Fish, Starship, and Atuin theme"
    if test (count $argv) -ne 1
        echo "Usage: fish-theme <gentleman|gentleman-cute>" >&2
        return 2
    end

    set -l theme "$argv[1]"
    switch $theme
        case gentleman gentleman-cute
        case '*'
            echo "Unknown theme: $theme" >&2
            echo "Usage: fish-theme <gentleman|gentleman-cute>" >&2
            return 2
    end

    set -l config_home "$HOME/.config"

    set -l theme_dir "$config_home/fish/themes"
    set -l profile "$theme_dir/$theme.fish"
    set -l marker_dir "$config_home/fish"
    set -l marker "$marker_dir/gentleman-theme"
    set -l starship_profile "$config_home/starship/$theme.toml"
    set -l atuin_theme "$config_home/atuin/themes/$theme.toml"

    if not test -r "$profile"
        echo "Fish profile is unavailable: $profile" >&2
        return 1
    end

    if not test -r "$starship_profile"
        echo "Starship profile is unavailable: $starship_profile" >&2
        return 1
    end

    if not test -r "$atuin_theme"
        echo "Atuin theme is unavailable: $atuin_theme" >&2
        return 1
    end

    if not command -q atuin
        echo "Atuin executable is unavailable." >&2
        return 1
    end

    fish --no-execute "$profile"
    or return 1

    mkdir -p "$marker_dir"
    or begin
        echo "Unable to create Fish configuration directory: $marker_dir" >&2
        return 1
    end

    set -l candidate (mktemp "$marker_dir/.gentleman-theme.XXXXXX")
    if test $status -ne 0
        echo "Unable to create Fish theme marker candidate." >&2
        return 1
    end

    if not printf '%s\n' "$theme" > "$candidate"
        command rm -f -- "$candidate"
        echo "Unable to write Fish theme marker candidate." >&2
        return 1
    end

    chmod u+w "$candidate"
    or begin
        command rm -f -- "$candidate"
        echo "Unable to make Fish theme marker writable." >&2
        return 1
    end

    command atuin config set theme.name "$theme"
    set -l atuin_status $status
    if test $atuin_status -ne 0
        command rm -f -- "$candidate"
        echo "Unable to activate Atuin theme: $theme" >&2
        return $atuin_status
    end

    mv -f "$candidate" "$marker"
    or begin
        command rm -f -- "$candidate"
        echo "Unable to activate Fish theme marker." >&2
        return 1
    end

    source "$profile"
    or return 1
    set -gx STARSHIP_CONFIG "$starship_profile"

    if status is-interactive
        commandline -f repaint
    end

    printf 'Fish, Starship, and Atuin theme selected: %s\n' "$theme"
end
