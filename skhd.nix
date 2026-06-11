{ lib, pkgs, ... }:
{
  home.activation.copySkhd = lib.hm.dag.entryAfter ["writeBoundary"] ''
    echo "Copying Skhd configuration..."
    rm -rf "$HOME/.config/skhd"

    if [ ! -d "$HOME/.config/skhd" ]; then
      mkdir -p "$HOME/.config/skhd"
    fi

    cp -r ${toString ./skhd}/* "$HOME/.config/skhd/"
    chmod -R u+w "$HOME/.config/skhd"

    # ─── Stable-path binary install (macOS Accessibility persistence) ───
    # Nix store paths include a hash that changes on every package update,
    # which invalidates macOS Accessibility permission. Copy the binary to
    # a stable user-owned path and point the LaunchAgent there so the
    # permission survives home-manager bumps.
    SKHD_SRC="${lib.getExe pkgs.skhd}"
    SKHD_DST="$HOME/.local/bin/skhd"
    PLIST="$HOME/Library/LaunchAgents/com.koekeishiya.skhd.plist"
    UID_VAL=$(/usr/bin/id -u)

    mkdir -p "$HOME/.local/bin"
    mkdir -p "$HOME/Library/LaunchAgents"

    NEW_HASH=$(/usr/bin/shasum -a 256 "$SKHD_SRC" | /usr/bin/awk '{print $1}')
    OLD_HASH=""
    if [ -f "$SKHD_DST" ]; then
      OLD_HASH=$(/usr/bin/shasum -a 256 "$SKHD_DST" | /usr/bin/awk '{print $1}')
    fi

    if [ "$NEW_HASH" != "$OLD_HASH" ] || [ ! -f "$PLIST" ]; then
      echo "Installing skhd to $SKHD_DST"
      /bin/launchctl bootout gui/$UID_VAL/com.koekeishiya.skhd 2>/dev/null || true
      /bin/rm -f "$SKHD_DST"
      /bin/cp "$SKHD_SRC" "$SKHD_DST"
      /bin/chmod +x "$SKHD_DST"

      /bin/cat > "$PLIST" <<PLISTEOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.koekeishiya.skhd</string>
    <key>ProgramArguments</key>
    <array>
        <string>$SKHD_DST</string>
    </array>
    <key>EnvironmentVariables</key>
    <dict>
        <key>PATH</key>
        <string>$HOME/.local/bin:$HOME/.local/state/nix/profiles/home-manager/home-path/bin:$HOME/.nix-profile/bin:/nix/var/nix/profiles/default/bin:/opt/homebrew/bin:/opt/homebrew/sbin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin</string>
    </dict>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <dict>
        <key>SuccessfulExit</key>
        <false/>
        <key>Crashed</key>
        <true/>
    </dict>
    <key>StandardOutPath</key>
    <string>/tmp/skhd_$USER.out.log</string>
    <key>StandardErrorPath</key>
    <string>/tmp/skhd_$USER.err.log</string>
    <key>ProcessType</key>
    <string>Interactive</string>
    <key>Nice</key>
    <integer>-20</integer>
</dict>
</plist>
PLISTEOF

      /bin/launchctl bootstrap gui/$UID_VAL "$PLIST" 2>/dev/null || true
      echo "skhd LaunchAgent reloaded at $SKHD_DST"
      echo "⚠️  If skhd fails to start: grant Accessibility to $SKHD_DST in System Settings"
    else
      echo "skhd binary already up to date at $SKHD_DST"
    fi
  '';
}
