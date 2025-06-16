-- Node.js configuration for Neovim
-- This ensures Neovim always uses Node.js from Nix, bypassing Volta and other PATH manipulations

local M = {}

-- Function to force Node.js from Nix for Neovim (bypassing Volta/PATH issues)
local function setup_nodejs()
    -- Force use of Nix-managed Node.js - hardcoded path to bypass PATH issues
    local nix_node_path = vim.fn.expand("~/.nix-profile/bin/node")

    -- Always try Nix first, don't rely on PATH detection
    if vim.fn.executable(nix_node_path) == 1 then
        local version_output = vim.fn.system(nix_node_path .. " --version 2>/dev/null")
        if vim.v.shell_error == 0 then
            local version = version_output:gsub("\n", ""):gsub("v", "")
            local major_version = tonumber(version:match("^(%d+)"))

            if major_version and major_version >= 14 then
                -- Force set the Node.js host program directly
                vim.g.node_host_prog = nix_node_path

                -- Set npm path from Nix too
                local nix_npm_path = vim.fn.expand("~/.nix-profile/bin/npm")
                if vim.fn.executable(nix_npm_path) == 1 then
                    vim.g.npm_host_prog = nix_npm_path
                end

                -- Override PATH to prioritize Nix Node.js directory
                local nix_bin_dir = vim.fn.expand("~/.nix-profile/bin")
                local current_path = vim.env.PATH or ""
                -- Remove any existing nix-profile entries and add at the beginning
                current_path = current_path:gsub(nix_bin_dir .. ":", "")
                current_path = current_path:gsub(":" .. nix_bin_dir, "")
                vim.env.PATH = nix_bin_dir .. ":" .. current_path

                if vim.g.debug_nodejs or vim.env.DEBUG_NODEJS then
                    print("✓ Node.js for Neovim: " .. nix_node_path .. " (v" .. version .. ")")
                    print("  Using Nix-managed Node.js (forced, bypassing Volta)")
                end

                return true, version
            end
        end
    end

    -- Fallback only if Nix Node.js is completely unavailable
    local fallback_paths = {
        "/opt/homebrew/bin/node",
        "/usr/local/bin/node",
        "/usr/bin/node",
    }

    for _, path in ipairs(fallback_paths) do
        if vim.fn.executable(path) == 1 then
            local version_output = vim.fn.system(path .. " --version 2>/dev/null")
            if vim.v.shell_error == 0 then
                local version = version_output:gsub("\n", ""):gsub("v", "")
                local major_version = tonumber(version:match("^(%d+)"))

                if major_version and major_version >= 18 then
                    vim.g.node_host_prog = path
                    local npm_path = path:gsub("/node$", "/npm")
                    if vim.fn.executable(npm_path) == 1 then
                        vim.g.npm_host_prog = npm_path
                    end

                    vim.notify(
                        "⚠️  Using fallback Node.js " .. version .. " from " .. path ..
                        "\nConsider installing Node.js via Nix for better isolation.",
                        vim.log.levels.WARN
                    )

                    return true, version
                end
            end
        end
    end

    vim.notify(
        "⚠️  Node.js 18+ not found! Some plugins may not work correctly.\nInstall Node.js via Nix: 'nix profile install nixpkgs#nodejs'",
        vim.log.levels.ERROR
    )
    return false, nil
end

-- Function to check if we're using a recent Node.js version
local function check_node_version()
    if not vim.g.node_host_prog then
        return false
    end

    local version_output = vim.fn.system(vim.g.node_host_prog .. " --version 2>/dev/null")
    if vim.v.shell_error ~= 0 then
        return false
    end

    local version = version_output:gsub("\n", ""):gsub("v", "")
    local major_version = tonumber(version:match("^(%d+)"))

    if major_version and major_version >= 22 then
        return true, version
    elseif major_version and major_version >= 18 then
        if vim.g.debug_nodejs then
            vim.notify(
                "ℹ️  Node.js v" .. version .. " works but v22+ is recommended for optimal performance.",
                vim.log.levels.INFO
            )
        end
        return true, version
    else
        vim.notify(
            "⚠️  Node.js version " .. version .. " is too old. Neovim requires v18+ (v22+ recommended).",
            vim.log.levels.WARN
        )
        return false, version
    end
end

-- Setup function to be called from init.lua
function M.setup(opts)
    opts = opts or {}

    -- Setup Node.js with forced Nix path
    local success, version = setup_nodejs()

    if success and not opts.silent then
        -- Check version and warn if outdated
        check_node_version()
    end

    return success
end

-- Function to get current Node.js info
function M.info()
    if not vim.g.node_host_prog then
        print("Node.js: Not configured")
        return
    end

    local version_output = vim.fn.system(vim.g.node_host_prog .. " --version 2>/dev/null")
    local version = "unknown"
    if vim.v.shell_error == 0 then
        version = version_output:gsub("\n", ""):gsub("v", "")
    end

    print("Node.js for Neovim: " .. vim.g.node_host_prog .. " (v" .. version .. ")")

    local nix_node = vim.fn.expand("~/.nix-profile/bin/node")
    if vim.g.node_host_prog == nix_node then
        print("Source: Nix-managed (isolated from project Node.js)")
    else
        print("Source: System/Fallback")
        if vim.g.node_host_prog:match("%.volta/") then
            print("⚠️  Warning: Using Volta Node.js - this may cause issues")
            print("   Consider using Nix: nix profile install nixpkgs#nodejs")
        end
    end

    if vim.g.npm_host_prog then
        print("npm: " .. vim.g.npm_host_prog)
    end

    -- Show PATH priority for debugging
    if vim.g.debug_nodejs then
        local path_entries = vim.split(vim.env.PATH, ":")
        print("PATH priority (first 5 entries):")
        for i = 1, math.min(5, #path_entries) do
            print("  " .. i .. ": " .. path_entries[i])
        end
    end
end

-- Command to manually refresh Node.js configuration
vim.api.nvim_create_user_command("NodeRefresh", function()
    M.setup({ silent = false })
    M.info()
end, { desc = "Refresh Node.js configuration" })

-- Command to show Node.js info
vim.api.nvim_create_user_command("NodeInfo", function()
    M.info()
end, { desc = "Show Node.js configuration info" })

-- Command to debug Node.js PATH issues
vim.api.nvim_create_user_command("NodeDebug", function()
    vim.g.debug_nodejs = true
    print("=== Node.js Debug Mode Enabled ===")
    M.setup({ silent = false })
    M.info()
    print("=== Current PATH ===")
    print(vim.env.PATH)
    vim.g.debug_nodejs = false
end, { desc = "Debug Node.js configuration and PATH" })

return M
