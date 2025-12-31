# Neovim Keymaps Reference

Complete reference of all keybindings configured in Gentleman.Dots Neovim setup. The leader key is `<Space>`.

> **Tip:** Press `<leader>?` in Neovim to see context-aware keybindings, or `<leader>sk` to search all keymaps.

---

## Quick Reference

| Action | Keys |
|--------|------|
| Find files | `<leader><space>` |
| Grep in project | `<leader>/` |
| Switch buffer | `<leader>,` |
| File explorer | `<leader>e` |
| Go to definition | `gd` |
| Hover docs | `K` |
| Code actions | `<leader>ca` |
| Open Lazygit | `<leader>gg` |

---

## Harpoon (Quick File Navigation)

Mark and jump to files instantly - by ThePrimeagen

| Keys | Description | Mode |
|------|-------------|------|
| `<leader>H` | Add file to Harpoon list | n |
| `<leader>h` | Toggle Harpoon quick menu | n |
| `<leader>1` | Jump to Harpoon file 1 | n |
| `<leader>2` | Jump to Harpoon file 2 | n |
| `<leader>3` | Jump to Harpoon file 3 | n |
| `<leader>4` | Jump to Harpoon file 4 | n |
| `<leader>5` | Jump to Harpoon file 5 | n |

---

## Mini.files (File Explorer)

Edit filesystem like a buffer

| Keys | Description | Mode |
|------|-------------|------|
| `<leader>fm` | Open mini.files (current file dir) | n |
| `<leader>fM` | Open mini.files (cwd) | n |
| `g.` | Toggle hidden files | n |
| `gc` | Set current directory as cwd | n |
| `<C-w>s` | Open in horizontal split | n |
| `<C-w>v` | Open in vertical split | n |
| `<C-w>S` | Open in horizontal split and close | n |
| `<C-w>V` | Open in vertical split and close | n |

---

## Oil.nvim (Directory Editor)

Edit your filesystem like a buffer

| Keys | Description | Mode |
|------|-------------|------|
| `-` | Open Oil (parent dir) | n |
| `<leader>E` | Open Oil (floating) | n |
| `<leader>-` | Open Oil in current file's directory | n |
| `g?` | Show help | n |
| `<CR>` | Select file/directory | n |
| `<C-s>` | Open in vertical split | n |
| `<C-v>` | Open in horizontal split | n |
| `<C-t>` | Open in new tab | n |
| `<C-p>` | Preview file | n |
| `<C-c>` | Close Oil | n |
| `<C-r>` | Refresh | n |
| `_` | Open cwd | n |
| `` ` `` | cd to current directory | n |
| `~` | :tcd to current directory | n |
| `gs` | Change sort | n |
| `gx` | Open external | n |
| `g.` | Toggle hidden files | n |
| `g\` | Toggle trash | n |
| `q` | Close Oil | n |

---

## Mini.surround (Text Manipulation)

Add, delete, replace surroundings (brackets, quotes, etc)

| Keys | Description | Mode |
|------|-------------|------|
| `sa` | Add surrounding (e.g., `saiw"` = surround word with ") | n,v |
| `sd` | Delete surrounding (e.g., `sd"` = delete surrounding ") | n |
| `sr` | Replace surrounding (e.g., `sr"'` = replace " with ') | n |
| `sf` | Find surrounding (to the right) | n |
| `sF` | Find surrounding (to the left) | n |
| `sh` | Highlight surrounding | n |
| `sn` | Update number of search lines | n |

---

## Snacks Picker (Fuzzy Finding)

Find files, grep text, search everything

| Keys | Description | Mode |
|------|-------------|------|
| `<leader><space>` | Find Files (Root Dir) | n |
| `<leader>,` | Switch Buffer | n |
| `<leader>/` | Grep (Root Dir) | n |
| `<leader>:` | Command History | n |
| `<leader>fb` | Buffers | n |
| `<leader>fB` | Buffers (all, including hidden) | n |
| `<leader>fc` | Find Config File | n |
| `<leader>ff` | Find Files (Root Dir) | n |
| `<leader>fF` | Find Files (cwd) | n |
| `<leader>fg` | Find Files (git-files) | n |
| `<leader>fr` | Recent Files | n |
| `<leader>fR` | Recent Files (cwd) | n |
| `<leader>fp` | Projects | n |
| `<leader>n` | Notification History | n |

---

## Search Commands

Search registers, marks, diagnostics, and more

| Keys | Description | Mode |
|------|-------------|------|
| `<leader>s"` | Registers | n |
| `<leader>s/` | Search History | n |
| `<leader>sa` | Auto Commands | n |
| `<leader>sb` | Buffer Lines | n |
| `<leader>sB` | Grep Open Buffers | n |
| `<leader>sc` | Command History | n |
| `<leader>sC` | Commands | n |
| `<leader>sd` | Document Diagnostics | n |
| `<leader>sD` | Workspace Diagnostics | n |
| `<leader>sg` | Grep (Root Dir) | n |
| `<leader>sG` | Grep (cwd) | n |
| `<leader>sh` | Help Pages | n |
| `<leader>sH` | Highlights | n |
| `<leader>si` | Icons | n |
| `<leader>sj` | Jumps | n |
| `<leader>sk` | Keymaps | n |
| `<leader>sl` | Location List | n |
| `<leader>sm` | Marks | n |
| `<leader>sM` | Man Pages | n |
| `<leader>sp` | Search Plugin Spec | n |
| `<leader>sq` | Quickfix List | n |
| `<leader>sR` | Resume Last Search | n |
| `<leader>ss` | LSP Symbols (Document) | n |
| `<leader>sS` | LSP Symbols (Workspace) | n |
| `<leader>su` | Undotree | n |
| `<leader>sw` | Word under cursor (Root Dir) | n,v |
| `<leader>sW` | Word under cursor (cwd) | n,v |

---

## LSP & Code Navigation

Go-to definitions, references, and previews

| Keys | Description | Mode |
|------|-------------|------|
| `gd` | Go to Definition | n |
| `gr` | Go to References | n |
| `gI` | Go to Implementation | n |
| `gy` | Go to Type Definition | n |
| `gD` | Go to Declaration | n |
| `K` | Hover Documentation | n |
| `gK` | Signature Help | n |
| `<leader>ca` | Code Actions | n,v |
| `<leader>cc` | Run Codelens | n |
| `<leader>cC` | Refresh Codelens | n |
| `<leader>cr` | Rename Symbol | n |
| `<leader>cR` | Rename File | n |
| `<leader>cs` | Symbols Outline | n |

### Preview Windows (goto-preview)

| Keys | Description | Mode |
|------|-------------|------|
| `gpd` | Preview Definition (popup) | n |
| `gpD` | Preview Declaration (popup) | n |
| `gpi` | Preview Implementation (popup) | n |
| `gpy` | Preview Type Definition (popup) | n |
| `gpr` | Preview References (popup) | n |
| `gP` | Close all preview windows | n |

---

## Git Integration

Git operations from within Neovim

| Keys | Description | Mode |
|------|-------------|------|
| `<leader>gb` | Git Blame (show who changed line) | n |
| `<leader>go` | Open file/folder in git repo | n |
| `<leader>gB` | Git Browse | n |
| `<leader>gc` | Git Commits | n |
| `<leader>gd` | Git Diff (hunks) | n |
| `<leader>gf` | Git File History (current) | n |
| `<leader>gg` | Open Lazygit | n |
| `<leader>gl` | Git Log | n |
| `<leader>gL` | Git Log (line) | n |
| `<leader>gs` | Git Status | n |
| `<leader>gS` | Git Stash | n |

### Hunk Navigation

| Keys | Description | Mode |
|------|-------------|------|
| `]h` | Next Hunk | n |
| `[h` | Previous Hunk | n |
| `<leader>ghp` | Preview Hunk | n |
| `<leader>ghs` | Stage Hunk | n |
| `<leader>ghr` | Reset Hunk | n |
| `<leader>ghS` | Stage Buffer | n |
| `<leader>ghR` | Reset Buffer | n |

---

## AI Assistants (Copilot + OpenCode)

AI-powered coding assistance

### Copilot (Insert Mode)

| Keys | Description | Mode |
|------|-------------|------|
| `<Tab>` | Accept Copilot suggestion | i |
| `<C-]>` | Dismiss Copilot suggestion | i |
| `<M-]>` | Next Copilot suggestion | i |
| `<M-[>` | Previous Copilot suggestion | i |

### OpenCode

| Keys | Description | Mode |
|------|-------------|------|
| `<leader>aa` | Toggle OpenCode | n |
| `<leader>as` | OpenCode select | n,x |
| `<leader>ai` | OpenCode ask | n,x |
| `<leader>aI` | OpenCode ask with context | n,x |
| `<leader>ab` | OpenCode ask about buffer | n,x |
| `<leader>ap` | OpenCode prompt | n,x |
| `<leader>ape` | OpenCode explain | n,x |
| `<leader>apf` | OpenCode fix | n,x |
| `<leader>apd` | OpenCode diagnose | n,x |
| `<leader>apr` | OpenCode review | n,x |
| `<leader>apt` | OpenCode test | n,x |
| `<leader>apo` | OpenCode optimize | n,x |

---

## Obsidian (Notes)

Note-taking and knowledge management

| Keys | Description | Mode |
|------|-------------|------|
| `<leader>oc` | Obsidian Check Checkbox | n |
| `<leader>ot` | Insert Obsidian Template | n |
| `<leader>oo` | Open in Obsidian App | n |
| `<leader>ob` | Show Obsidian Backlinks | n |
| `<leader>ol` | Show Obsidian Links | n |
| `<leader>on` | Create New Note | n |
| `<leader>os` | Search Obsidian | n |
| `<leader>oq` | Quick Switch | n |

---

## Tmux Navigation

Seamless navigation between Neovim and Tmux panes

| Keys | Description | Mode |
|------|-------------|------|
| `<C-h>` | Navigate to the left pane | n |
| `<C-j>` | Navigate to the bottom pane | n |
| `<C-k>` | Navigate to the top pane | n |
| `<C-l>` | Navigate to the right pane | n |
| `<C-\>` | Navigate to last active pane | n |
| `<C-Space>` | Navigate to next pane | n |

---

## Custom Keymaps

Gentleman.Dots custom keybindings

| Keys | Description | Mode |
|------|-------------|------|
| `<C-b>` | Delete to end of word (insert) | i |
| `<C-c>` | Escape from any mode | i,n,v |
| `<leader>uk` | Toggle Screenkey | n |
| `-` | Open Oil (parent directory) | n |
| `<leader>bq` | Delete other buffers but current | n |
| `<C-s>` | Save file | n |
| `<leader>sg` | Grep Selected Text | v |
| `<leader>sG` | Grep Selected Text (Root Dir) | v |
| `<leader>md` | Delete all marks | n |

---

## Buffers & Tabs

Buffer and tab management

### Buffers

| Keys | Description | Mode |
|------|-------------|------|
| `<leader>bb` | Switch to Other Buffer | n |
| `<leader>bd` | Delete Buffer | n |
| `<leader>bD` | Delete Buffer and Window | n |
| `<leader>bo` | Delete Other Buffers | n |
| `<leader>bp` | Toggle Buffer Pin | n |
| `<leader>bP` | Delete Non-Pinned Buffers | n |
| `[b` | Previous Buffer | n |
| `]b` | Next Buffer | n |

### Tabs

| Keys | Description | Mode |
|------|-------------|------|
| `<leader><tab>l` | Last Tab | n |
| `<leader><tab>f` | First Tab | n |
| `<leader><tab><tab>` | New Tab | n |
| `<leader><tab>d` | Close Tab | n |
| `<leader><tab>]` | Next Tab | n |
| `<leader><tab>[` | Previous Tab | n |

---

## Windows & Splits

Window navigation and management

| Keys | Description | Mode |
|------|-------------|------|
| `<leader>w` | Windows menu (which-key) | n |
| `<leader>wd` | Delete Window | n |
| `<leader>wm` | Maximize Window | n |
| `<leader>-` | Split Below | n |
| `<leader>\|` | Split Right | n |
| `<C-Up>` | Increase Height | n |
| `<C-Down>` | Decrease Height | n |
| `<C-Left>` | Decrease Width | n |
| `<C-Right>` | Increase Width | n |

> **Note:** Window navigation (`<C-h/j/k/l>`) is handled by Tmux Navigation for seamless pane switching.

---

## Debugging (DAP)

Debug Adapter Protocol integration

| Keys | Description | Mode |
|------|-------------|------|
| `<leader>db` | Toggle Breakpoint | n |
| `<leader>dB` | Breakpoint Condition | n |
| `<leader>dc` | Continue | n |
| `<leader>da` | Run with Args | n |
| `<leader>dC` | Run to Cursor | n |
| `<leader>dg` | Go to Line (no execute) | n |
| `<leader>di` | Step Into | n |
| `<leader>dj` | Down | n |
| `<leader>dk` | Up | n |
| `<leader>dl` | Run Last | n |
| `<leader>do` | Step Out | n |
| `<leader>dO` | Step Over | n |
| `<leader>dp` | Pause | n |
| `<leader>dr` | Toggle REPL | n |
| `<leader>ds` | Session | n |
| `<leader>dt` | Terminate | n |
| `<leader>du` | Toggle DAP UI | n |
| `<leader>dw` | Widgets | n |

---

## UI Toggles

Toggle UI elements and settings

| Keys | Description | Mode |
|------|-------------|------|
| `<leader>ub` | Toggle Background (dark/light) | n |
| `<leader>uC` | Colorscheme Picker | n |
| `<leader>ud` | Toggle Diagnostics | n |
| `<leader>uf` | Toggle Auto Format (global) | n |
| `<leader>uF` | Toggle Auto Format (buffer) | n |
| `<leader>ug` | Toggle Indent Guides | n |
| `<leader>uh` | Toggle Inlay Hints | n |
| `<leader>uI` | Inspect Pos (Treesitter) | n |
| `<leader>ul` | Toggle Line Numbers | n |
| `<leader>uL` | Toggle Relative Numbers | n |
| `<leader>un` | Dismiss Notifications | n |
| `<leader>us` | Toggle Spelling | n |
| `<leader>uT` | Toggle Treesitter Highlight | n |
| `<leader>uw` | Toggle Word Wrap | n |
| `<leader>uz` | Toggle Zen Mode | n |
| `<leader>uZ` | Toggle Zoom | n |

---

## Diagnostics & Quickfix

Navigate errors, warnings, and quickfix

| Keys | Description | Mode |
|------|-------------|------|
| `<leader>xx` | Document Diagnostics (Trouble) | n |
| `<leader>xX` | Workspace Diagnostics (Trouble) | n |
| `<leader>xL` | Location List (Trouble) | n |
| `<leader>xQ` | Quickfix List (Trouble) | n |
| `[d` | Previous Diagnostic | n |
| `]d` | Next Diagnostic | n |
| `[e` | Previous Error | n |
| `]e` | Next Error | n |
| `[w` | Previous Warning | n |
| `]w` | Next Warning | n |
| `[q` | Previous Quickfix | n |
| `]q` | Next Quickfix | n |

---

## Session & Quit

Session management and exiting

| Keys | Description | Mode |
|------|-------------|------|
| `<leader>qq` | Quit All | n |
| `<leader>qs` | Restore Session | n |
| `<leader>qS` | Select Session | n |
| `<leader>ql` | Restore Last Session | n |
| `<leader>qd` | Don't Save Current Session | n |

---

## General & Essential

Essential keybindings every Neovim user needs

| Keys | Description | Mode |
|------|-------------|------|
| `<leader>?` | Show Buffer Local Keymaps | n |
| `<leader>e` | Toggle Explorer (neo-tree) | n |
| `<leader>E` | Explorer (neo-tree cwd) | n |
| `<leader>l` | Lazy (plugin manager) | n |
| `<leader>L` | LazyVim Changelog | n |
| `<Esc>` | Escape and Clear hlsearch | n,i |
| `<leader>ur` | Redraw / Clear hlsearch | n |
| `gx` | Open with system app | n |
| `j` | Down (respects wrapped lines) | n,x |
| `k` | Up (respects wrapped lines) | n,x |

---

## Mode Legend

| Mode | Description |
|------|-------------|
| `n` | Normal mode |
| `i` | Insert mode |
| `v` | Visual mode |
| `x` | Visual block mode |
| `n,v` | Normal and Visual modes |

---

## Learning More

- Press `<leader>?` to see context-aware keybindings in which-key
- Press `<leader>sk` to fuzzy search all keymaps
- Run `:map` to see all current mappings
- Check the [LazyVim documentation](https://lazyvim.org) for more details
