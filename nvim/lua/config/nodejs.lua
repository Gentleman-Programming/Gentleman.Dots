-- Node.js configuration for Neovim
-- This ensures Neovim always uses the latest stable Node.js version

local M = {}

-- Function to force Node.js 22 for Neovim (independent of project Node version)
local function setup_nodejs()
    -- Force use of Node.js 22 from Nix for Neovim plugins and LSPs
    -- This ensures Neovim always works regardless of project Node.js version
    local forced_node_path = vim.fn.expand("~/.nix-profile/bin/node")

    local selected_node = nil
    local selected_version = nil

    -- First, try to use the Nix-managed Node.js 22
    if vim.fn.executable(forced_node_path) == 1 then
        local version_output = vim.fn.system(forced_node_path .. " --version 2>/dev/null")
        if vim.v.shell_error == 0 then
            local version = version_output:gsub("\n", ""):gsub("v", "")
            local major_version = tonumber(version:match("^(%d+)"))

            -- Verify it's actually Node.js 18+ (preferably 22+)
            if major_version and major_version >= 18 then
                selected_node = forced_node_path
                selected_version = version
            end
        end
    end

    -- Fallback to other paths only if Nix Node.js is not available or too old
    if not selected_node then
        local fallback_paths = {
            "/opt/homebrew/bin/node",
            "/usr/local/bin/node",
            vim.fn.exepath("node"),
        }

        for _, path in ipairs(fallback_paths) do
            if vim.fn.executable(path) == 1 then
                local version_output = vim.fn.system(path .. " --version 2>/dev/null")
                if vim.v.shell_error == 0 then
                    local version = version_output:gsub("\n", ""):gsub("v", "")
                    local major_version = tonumber(version:match("^(%d+)"))

                    -- Only accept Node.js 18+ as fallback
                    if major_version and major_version >= 18 then
                        selected_node = path
                        selected_version = version
                        break
                    end
                end
            end
        end
    end

    if selected_node then
        -- Set the Node.js host program for Neovim
        vim.g.node_host_prog = selected_node

        -- Also set npm path if available
        local npm_path = selected_node:gsub("/node$", "/npm")
        if vim.fn.executable(npm_path) == 1 then
            vim.g.npm_host_prog = npm_path
        end

        -- Print version info (only in debug mode or when explicitly requested)
        if vim.g.debug_nodejs or vim.env.DEBUG_NODEJS then
            print("✓ Node.js for Neovim: " .. selected_node .. " (v" .. (selected_version or "unknown") .. ")")
            if selected_node == forced_node_path then
                print("  Using Nix-managed Node.js (isolated from project)")
            else
                print("  Using fallback Node.js - consider installing Node.js 22 via Nix")
            end
        end

        return true
    else
        vim.notify(
        "⚠️  Node.js 18+ not found! Some plugins may not work correctly.\nConsider installing Node.js 22 via Nix.",
            vim.log.levels.WARN)
        return false
    end
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

    -- Setup Node.js
    local success = setup_nodejs()

    if success and not opts.silent then
        -- Check version and warn if outdated
        check_node_version()
    end

    -- Set environment variables for LSPs and other tools
    if vim.g.node_host_prog then
        local node_dir = vim.fn.fnamemodify(vim.g.node_host_prog, ":h")

        -- Ensure Node.js bin directory is in PATH
        local current_path = vim.env.PATH or ""
        if not current_path:find(node_dir, 1, true) then
            vim.env.PATH = node_dir .. ":" .. current_path
        end
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
    end

    if vim.g.npm_host_prog then
        print("npm: " .. vim.g.npm_host_prog)
    end
end

-- Command to manually refresh Node.js configuration
vim.api.nvim_create_user_command("NodeRefresh", function()
    M.setup({ silent = false })
end, { desc = "Refresh Node.js configuration" })

-- Command to show Node.js info
vim.api.nvim_create_user_command("NodeInfo", function()
    M.info()
end, { desc = "Show Node.js configuration info" })

return M
