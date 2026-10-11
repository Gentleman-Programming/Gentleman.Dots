# Fish listing colors

The bundled Fish configuration exports `LS_COLORS` and `LSCOLORS` to match its
neutral Gentleman palette: cyan directories, pink links, and green executables.
GNU `ls` uses exact RGB colors; BSD/macOS `ls` uses ANSI approximations controlled
by the terminal palette.

The `ls` alias uses native `ls -G` on macOS and `ls --color=auto` on Linux/Termux.
Arguments are forwarded unchanged; no separate `gls` installation is required.
Regular files retain their default color and permission-sensitive categories
remain distinguishable.

After installing the updated Fish configuration, open a new shell to load it.
This legacy configuration does not include the `gentleman-cute` profile or the
`fish-theme` selector provided by the Nix checkout.

Run the isolated configuration regression check with:

```sh
bash GentlemanFish/tests/ls-theme.sh
```
