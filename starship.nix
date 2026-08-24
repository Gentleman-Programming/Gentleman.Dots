{ pkgs, ... }:

let
  toml = pkgs.formats.toml { };

  catppuccinMocha = {
    rosewater = "#f5e0dc";
    flamingo = "#f2cdcd";
    pink = "#f5c2e7";
    mauve = "#cba6f7";
    red = "#f38ba8";
    maroon = "#eba0ac";
    peach = "#fab387";
    yellow = "#f9e2af";
    green = "#a6e3a1";
    teal = "#94e2d5";
    sky = "#89dceb";
    sapphire = "#74c7ec";
    blue = "#89b4fa";
    lavender = "#b4befe";
    text = "#cdd6f4";
    subtext1 = "#bac2de";
    subtext0 = "#a6adc8";
    overlay2 = "#9399b2";
    overlay1 = "#7f849c";
    overlay0 = "#6c7086";
    surface2 = "#585b70";
    surface1 = "#45475a";
    surface0 = "#313244";
    base = "#1e1e2e";
    mantle = "#181825";
    crust = "#11111b";
  };

  gentlemanPalette = {
    text = "#F3F6F9";
    red = "#CB7C94";
    green = "#B7CC85";
    yellow = "#FFE066";
    blue = "#7FB4CA";
    mauve = "#A3B5D6";
    pink = "#FF8DD7";
    teal = "#7AA89F";
    peach = "#DEBA87";
    subtext0 = "#5C6170";
    overlay0 = "#232A40";
    rosewater = "#E0C15A";
    flamingo = "#FF8DD7";
    maroon = "#C4746E";
    lavender = "#B99BF2";
    subtext1 = "#8A8FA3";
    overlay2 = "#313342";
    overlay1 = "#191E28";
    surface2 = "#27345C";
    surface1 = "#232A40";
    surface0 = "#191E28";
    base = "none";
    mantle = "#06080f";
    crust = "#06080f";
  };

  gentlemanCutePalette = {
    text = "#F6EFF3";
    red = "#FF718F";
    green = "#D2CBD0";
    yellow = "#E0C27A";
    blue = "#F095C8";
    mauve = "#C49BFF";
    pink = "#FFB1DD";
    teal = "#D7A0B8";
    peach = "#F2B86D";
    subtext0 = "#76616B";
    overlay0 = "#342230";
    rosewater = "#E0C27A";
    flamingo = "#FFB1DD";
    maroon = "#C96AA2";
    lavender = "#C49BFF";
    subtext1 = "#A78E9B";
    overlay2 = "#D7A0B8";
    overlay1 = "#B99AA9";
    surface2 = "#342230";
    surface1 = "#241822";
    surface0 = "#20161E";
    base = "#1A1218";
    mantle = "#1A1218";
    crust = "#1A1218";
  };

  commonConfig = {
    format = "($directory)$os$git_branch$fill$nodejs$rust$golang$php$bun$java$c$conda$zig$cmd_duration$time\n$character";
    add_newline = true;
    command_timeout = 3600000;

    fill = {
      symbol = " ";
    };

    character = {
      success_symbol = "[󱗞 ](fg:green)";
      error_symbol = "[󱗞 ](fg:red)";
      vimcmd_symbol = "[N](bold red)";
      vimcmd_replace_one_symbol = "[R](bold peach)";
      vimcmd_visual_symbol = "[V](bold mauve)";
    };

    username = {
      style_user = "bold blue";
      style_root = "bold red";
      format = "[󱗞 $user](fg:$style) ";
      disabled = false;
      show_always = true;
    };

    directory = {
      format = "[$path](bold $style)[$read_only]($read_only_style) ";
      truncation_length = 2;
      style = "fg:blue";
      read_only_style = "fg:blue";
      before_repo_root_style = "fg:blue";
      truncation_symbol = "…/";
      truncate_to_repo = true;
      read_only = "  ";
      substitutions = {
        Documents = "󰈙 ";
        Downloads = " ";
        Music = " ";
        Pictures = " ";
      };
    };

    cmd_duration = {
      format = " took [ $duration]($style) ";
      style = "bold fg:yellow";
      min_time = 500;
    };

    git_branch = {
      format = "-> [$symbol$branch]($style) ";
      style = "bold fg:mauve";
      symbol = " ";
    };

    git_status = {
      format = "[$all_status$ahead_behind ]($style)";
      style = "fg:text bg:pink";
    };

    docker_context = {
      disabled = true;
      symbol = " ";
    };

    python = {
      disabled = false;
      format = "[$symbol$pyenv_prefix($version)( $virtualenv)](fg:peach)";
      symbol = " ";
      version_format = "$raw";
    };

    java = {
      format = "[[ $symbol ($version) ](fg:red)]($style)";
      version_format = "$raw";
      symbol = " ";
      disabled = false;
    };

    c = {
      format = "[[ $symbol ($version) ](fg:blue)]($style)";
      symbol = " ";
      version_format = "$raw";
      disabled = false;
    };

    zig = {
      format = "[[ $symbol ($version) ](fg:peach)]($style)";
      version_format = "$raw";
      disabled = false;
    };

    bun = {
      version_format = "$raw";
      format = "[[ $symbol ($version) ](fg:text)]($style)";
      disabled = false;
    };

    nodejs = {
      symbol = "";
      format = "[[ $symbol ($version) ](fg:green)]($style)";
    };

    rust = {
      symbol = "";
      format = "[[ $symbol ($version) ](fg:red)]($style)";
    };

    golang = {
      symbol = "";
      format = "[[ $symbol ($version) ](fg:teal)]($style)";
    };

    php = {
      symbol = "";
      format = "[[ $symbol ($version) ](fg:peach)]($style)";
    };

    time = {
      disabled = false;
      time_format = "%R";
      format = "[[   $time ](fg:subtext0)]($style)";
    };
  };

  gentlemanConfig = commonConfig // {
    palette = "gentleman";
    palettes = {
      catppuccin_mocha = catppuccinMocha;
      gentleman = gentlemanPalette;
    };
  };

  gentlemanCuteConfig = commonConfig // {
    palette = "gentleman-cute";
    palettes = {
      catppuccin_mocha = catppuccinMocha;
      "gentleman-cute" = gentlemanCutePalette;
    };
  };

  gentlemanProfile = toml.generate "starship-gentleman.toml" gentlemanConfig;
  gentlemanCuteProfile = toml.generate "starship-gentleman-cute.toml" gentlemanCuteConfig;
in
{
  home.file = {
    # Keep the existing top-level config for non-Fish consumers.
    ".config/starship.toml".source = gentlemanProfile;
    ".config/starship/gentleman.toml".source = gentlemanProfile;
    ".config/starship/gentleman-cute.toml".source = gentlemanCuteProfile;
  };
}
