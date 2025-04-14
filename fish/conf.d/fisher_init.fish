# Fisher initialization
if status is-interactive
    # Set up fisher path
    set -g fisher_path $__fish_config_dir/fisher

    # Add Fisher functions to function path
    if test -d $fisher_path/functions
        set fish_function_path $fish_function_path[1] $fisher_path/functions $fish_function_path[2..-1]
    end

    # Add Fisher completions to completion path
    if test -d $fisher_path/completions
        set fish_complete_path $fish_complete_path[1] $fisher_path/completions $fish_complete_path[2..-1]
    end

    # Source Fisher configuration files
    if test -d $fisher_path/conf.d
        for file in $fisher_path/conf.d/*.fish
            if test -f $file
                source $file
            end
        end
    end
end
