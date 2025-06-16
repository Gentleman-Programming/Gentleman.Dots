-- Ultra minimal Neovim configuration - ONLY Oil.nvim
-- No plugins, no extras, just Oil for file management

-- Configure Node.js (force Node.js 22 from Nix for Neovim, independent of project version)
local function setup_nodejs()
    -- Force use of Nix-managed Node.js 22 for Neovim (isolated from project Node.js)
    local forced_node_path = vim.fn.expand("~/.nix-profile/bin/node")

    -- First priority: Nix-managed Node.js 22
    if vim.fn.executable(forced_node_path) == 1 then
        local version_output = vim.fn.system(forced_node_path .. " --version 2>/dev/null")
        if vim.v.shell_error == 0 then
            local version = version_output:gsub("\n", ""):gsub("v", "")
            local major_version = tonumber(version:match("^(%d+)"))

            -- Verify it's Node.js 18+ (preferably 22+)
            if major_version and major_version >= 18 then
                vim.g.node_host_prog = forced_node_path
                local node_dir = vim.fn.fnamemodify(forced_node_path, ":h")
                vim.env.PATH = node_dir .. ":" .. (vim.env.PATH or "")
                return
            end
        end
    end

    -- Fallback to other Node.js versions only if Nix version not available
    local fallback_paths = {
        "/opt/homebrew/bin/node",
        "/usr/local/bin/node",
    }

    for _, path in ipairs(fallback_paths) do
        if vim.fn.executable(path) == 1 then
            local version_output = vim.fn.system(path .. " --version 2>/dev/null")
            if vim.v.shell_error == 0 then
                local version = version_output:gsub("\n", ""):gsub("v", "")
                local major_version = tonumber(version:match("^(%d+)"))

                -- Only accept Node.js 18+ as fallback
                if major_version and major_version >= 18 then
                    vim.g.node_host_prog = path
                    local node_dir = vim.fn.fnamemodify(path, ":h")
                    vim.env.PATH = node_dir .. ":" .. (vim.env.PATH or "")
                    break
                end
            end
        end
    end
end

setup_nodejs()

-- Disable all default plugins to speed up startup
vim.g.loaded_gzip = 1
vim.g.loaded_zip = 1
vim.g.loaded_zipPlugin = 1
vim.g.loaded_tar = 1
vim.g.loaded_tarPlugin = 1
vim.g.loaded_getscript = 1
vim.g.loaded_getscriptPlugin = 1
vim.g.loaded_vimball = 1
vim.g.loaded_vimballPlugin = 1
vim.g.loaded_2html_plugin = 1
vim.g.loaded_logiPat = 1
vim.g.loaded_rrhelper = 1
vim.g.loaded_netrw = 1
vim.g.loaded_netrwPlugin = 1
vim.g.loaded_netrwSettings = 1
vim.g.loaded_netrwFileHandlers = 1
vim.g.loaded_matchit = 1
vim.g.loaded_matchparen = 1
vim.g.loaded_sql_completion = 1
vim.g.loaded_syntax_completion = 1
vim.g.loaded_xmlformat = 1

-- Minimal settings
vim.opt.termguicolors = true
vim.opt.number = false
vim.opt.relativenumber = false
vim.opt.signcolumn = "no"
vim.opt.laststatus = 0
vim.opt.cmdheight = 1
vim.opt.showmode = false
vim.opt.ruler = false
vim.opt.showcmd = false

-- Leader key
vim.g.mapleader = " "

-- Bootstrap lazy.nvim
local lazypath = vim.fn.stdpath("data") .. "/lazy/lazy.nvim"
if not vim.loop.fs_stat(lazypath) then
    vim.fn.system({
        "git", "clone", "--filter=blob:none",
        "https://github.com/folke/lazy.nvim.git",
        "--branch=stable", lazypath,
    })
end
vim.opt.rtp:prepend(lazypath)

-- Setup only Oil.nvim and theme
require("lazy").setup({
    {
        "Gentleman-Programming/gentleman-kanagawa-blur",
        name = "gentleman-kanagawa-blur",
        priority = 1000,
        config = function()
            vim.cmd.colorscheme("gentleman-kanagawa-blur")
        end,
    },
    {
        "stevearc/oil.nvim",
        dependencies = { "nvim-tree/nvim-web-devicons" },
        config = function()
            require("oil").setup({
                default_file_explorer = true,
                skip_confirm_for_simple_edits = true,

                keymaps = {
                    ["<CR>"] = {
                        function()
                            local oil = require("oil")
                            local entry = oil.get_cursor_entry()
                            local dir = oil.get_current_dir()

                            -- Only intercept files, let Oil handle directories normally
                            if entry and entry.type == "file" and dir and vim.g.oil_open_in_zed then
                                local file_path = dir .. entry.name
                                vim.fn.jobstart({ "zed", file_path }, { detach = true })
                                vim.cmd("qa!")
                            else
                                -- Use normal Oil behavior for directories and non-Zed contexts
                                require("oil.actions").select.callback()
                            end
                        end,
                    },
                    ["<C-s>"] = { "actions.select", opts = { vertical = true } },
                    ["<C-h>"] = { "actions.select", opts = { horizontal = true } },
                    ["<C-t>"] = { "actions.select", opts = { tab = true } },
                    ["<C-p>"] = "actions.preview",
                    ["<C-c>"] = "actions.close",
                    ["<C-r>"] = "actions.refresh",
                    ["-"] = "actions.parent",
                    ["_"] = "actions.open_cwd",
                    ["`"] = "actions.cd",
                    ["gs"] = "actions.change_sort",
                    ["gx"] = "actions.open_external",
                    ["g."] = "actions.toggle_hidden",
                    ["q"] = "actions.close",
                    ["<Esc>"] = "actions.close",
                },

                use_default_keymaps = false,

                view_options = {
                    show_hidden = true,
                    is_always_hidden = function(name)
                        return name == ".." or name == ".git"
                    end,
                    natural_order = true,
                    sort = {
                        { "type", "asc" },
                        { "name", "asc" },
                    },
                },

                float = {
                    padding = 2,
                    max_width = 90,
                    max_height = 25,
                    border = "rounded",
                },
            })

            -- Global mappings
            vim.keymap.set("n", "-", "<CMD>Oil<CR>")
            vim.keymap.set("n", "<leader>e", "<CMD>Oil<CR>")
        end,
    },
    {
        "nvim-tree/nvim-web-devicons",
        config = function()
            require("nvim-web-devicons").setup({
                default = true,
            })
        end,
    },
}, {
    defaults = { lazy = false },
    install = { missing = true },
    checker = { enabled = false },
    change_detection = { enabled = false },
    performance = {
        rtp = {
            disabled_plugins = {
                "gzip", "matchit", "matchparen", "netrwPlugin",
                "tarPlugin", "tohtml", "tutor", "zipPlugin",
            },
        },
    },
})

-- Quick quit mapping
vim.keymap.set("n", "q", "<CMD>qa!<CR>")

-- Auto-open Oil if started with a directory
vim.api.nvim_create_autocmd("VimEnter", {
    callback = function()
        local args = vim.fn.argv()
        if #args == 1 and vim.fn.isdirectory(args[1]) == 1 then
            vim.defer_fn(function()
                require("oil").open(args[1])
            end, 10)
        elseif #args == 0 then
            vim.defer_fn(function()
                require("oil").open(".")
            end, 10)
        end
    end,
})
