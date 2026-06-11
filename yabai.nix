{ lib, pkgs, ... }:
{
  home.activation.copyYabai = lib.hm.dag.entryAfter ["writeBoundary"] ''
    # ─── macOS Mission Control settings ───
    # Disable automatic space reordering based on most recent use
    # This is CRITICAL for yabai/sketchybar to work correctly with numbered spaces
    echo "Configuring Mission Control settings..."
    /usr/bin/defaults write com.apple.dock mru-spaces -bool false
    killall Dock 2>/dev/null || true

    echo "Copying Yabai configuration..."
    rm -rf "$HOME/.config/yabai"

    if [ ! -d "$HOME/.config/yabai" ]; then
      mkdir -p "$HOME/.config/yabai"
    fi

    cp -r ${toString ./yabai}/* "$HOME/.config/yabai/"
    chmod -R u+w "$HOME/.config/yabai"
    chmod +x "$HOME/.config/yabai/yabairc"
    chmod +x "$HOME/.config/yabai"/*.sh

    # ─── Stable-path binary install (macOS Accessibility persistence) ───
    # Nix store paths include a hash that changes on every package update,
    # which invalidates macOS Accessibility permission. Copy the binary to
    # a stable user-owned path and point the LaunchAgent there so the
    # permission survives home-manager bumps.
    YABAI_SRC="${lib.getExe pkgs.yabai}"
    YABAI_DST="$HOME/.local/bin/yabai"
    PLIST="$HOME/Library/LaunchAgents/com.koekeishiya.yabai.plist"
    UID_VAL=$(/usr/bin/id -u)

    mkdir -p "$HOME/.local/bin"
    mkdir -p "$HOME/Library/LaunchAgents"

    NEW_HASH=$(/usr/bin/shasum -a 256 "$YABAI_SRC" | /usr/bin/awk '{print $1}')
    OLD_HASH=""
    if [ -f "$YABAI_DST" ]; then
      OLD_HASH=$(/usr/bin/shasum -a 256 "$YABAI_DST" | /usr/bin/awk '{print $1}')
    fi

    if [ "$NEW_HASH" != "$OLD_HASH" ] || [ ! -f "$PLIST" ]; then
      echo "Installing yabai to $YABAI_DST"
      /bin/launchctl bootout gui/$UID_VAL/com.koekeishiya.yabai 2>/dev/null || true
      /bin/rm -f "$YABAI_DST"
      /bin/cp "$YABAI_SRC" "$YABAI_DST"
      /bin/chmod +x "$YABAI_DST"

      /bin/cat > "$PLIST" <<PLISTEOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.koekeishiya.yabai</string>
    <key>ProgramArguments</key>
    <array>
        <string>$YABAI_DST</string>
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
    <string>/tmp/yabai_$USER.out.log</string>
    <key>StandardErrorPath</key>
    <string>/tmp/yabai_$USER.err.log</string>
    <key>ProcessType</key>
    <string>Interactive</string>
    <key>Nice</key>
    <integer>-20</integer>
</dict>
</plist>
PLISTEOF

      /bin/launchctl bootstrap gui/$UID_VAL "$PLIST" 2>/dev/null || true
      echo "yabai LaunchAgent reloaded at $YABAI_DST"
      echo "⚠️  If yabai fails to start: grant Accessibility to $YABAI_DST in System Settings"
    else
      echo "yabai binary already up to date at $YABAI_DST"
    fi

    # Home Manager may only copy config files without changing the yabai binary.
    # Restart the custom LaunchAgent so changes in ~/.config/yabai/yabairc take effect.
    if [ -f "$PLIST" ]; then
      echo "Restarting yabai LaunchAgent to apply config changes..."
      /bin/launchctl kickstart -k gui/$UID_VAL/com.koekeishiya.yabai 2>/dev/null || true
    fi

    # ─── Yabai sudoers setup for scripting addition ───
    # Required for space switching (yabai --load-sa needs passwordless sudo).
    # Hash is computed from the stable-path copy and must match exactly.
    YABAI_HASH=$(/usr/bin/shasum -a 256 "$YABAI_DST" | /usr/bin/awk '{print $1}')
    EXPECTED_ENTRY="$USER ALL=(root) NOPASSWD: sha256:$YABAI_HASH $YABAI_DST --load-sa"

    NEEDS_UPDATE=false
    if [ ! -f /private/etc/sudoers.d/yabai ]; then
      NEEDS_UPDATE=true
    elif ! /usr/bin/grep -q "$YABAI_HASH" /private/etc/sudoers.d/yabai 2>/dev/null; then
      NEEDS_UPDATE=true
    fi

    if [ "$NEEDS_UPDATE" = true ]; then
      echo ""
      echo "══════════════════════════════════════════════════════════════"
      echo "  ⚠️  Yabai scripting addition needs sudoers update!"
      echo "  Run this command to enable space switching:"
      echo ""
      echo "  echo \"$EXPECTED_ENTRY\" | sudo tee /private/etc/sudoers.d/yabai"
      echo "  sudo $YABAI_DST --load-sa"
      echo "══════════════════════════════════════════════════════════════"
      echo ""
    else
      echo "Yabai sudoers entry is up to date."
    fi
  '';
}
