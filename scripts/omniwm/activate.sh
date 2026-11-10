#!/bin/bash
# With --check-platform, perform only the pre-write Home Manager guard.
# The optional app path argument permits isolated tests of /Applications.
set -u

unsupported() {
  printf 'OmniWM requires Apple Silicon and macOS 26 or later.\n' >&2
  exit 1
}

if [ "$(uname -s)" != Darwin ] || [ "$(uname -m)" != arm64 ]; then
  unsupported
fi
version="$(sw_vers -productVersion 2>/dev/null)" || unsupported
major="${version%%.*}"
case "$version" in ''|*[!0-9.]*|.*|*..*|*.) unsupported;; esac
case "$major" in ''|*[!0-9]*) unsupported;; esac
if ! [ "$major" -ge 26 ] 2>/dev/null; then
  unsupported
fi

if [ "${1:-}" = --check-platform ]; then
  exit 0
fi

user_app="$HOME/Applications/OmniWM.app"
system_app="${1:-/Applications/OmniWM.app}"
if [ -e "$user_app" ] || [ -L "$user_app" ] || [ -e "$system_app" ] || [ -L "$system_app" ]; then
  printf 'OmniWM: existing app preserved.\n'
else
  if brew_cmd="$(command -v brew)"; then
    :
  elif [ -x /opt/homebrew/bin/brew ]; then
    brew_cmd=/opt/homebrew/bin/brew
  elif [ -x /usr/local/bin/brew ]; then
    brew_cmd=/usr/local/bin/brew
  else
    printf 'OmniWM: Homebrew not found; legacy services unchanged.\n' >&2
    exit 1
  fi
  printf 'OmniWM: installing Homebrew cask.\n'
  if ! "$brew_cmd" install --cask omniwm; then
    printf 'OmniWM: Homebrew installation failed; legacy services unchanged.\n' >&2
    exit 1
  fi
  if [ ! -e "$user_app" ] && [ ! -L "$user_app" ] && [ ! -e "$system_app" ] && [ ! -L "$system_app" ]; then
    printf 'OmniWM: installed app not found; legacy services unchanged.\n' >&2
    exit 1
  fi
fi

uid="$(id -u 2>/dev/null)" || uid=""
case "$uid" in ''|*[!0-9]*)
  printf 'OmniWM: current-user GUI launchd domain unavailable.\n' >&2
  exit 1;;
esac
domain="gui/$uid"
if ! launchctl print "$domain" >/dev/null 2>&1; then
  printf 'OmniWM: current-user GUI launchd domain unavailable.\n' >&2
  exit 1
fi

for label in com.asmvik.yabai com.koekeishiya.yabai com.koekeishiya.skhd \
             homebrew.mxcl.sketchybar com.felixkratz.sketchybar; do
  target="$domain/$label"
  if ! launchctl disable "$target" >/dev/null 2>&1; then
    printf 'OmniWM: could not disable %s.\n' "$label" >&2
    exit 1
  fi
  if launchctl print "$target" >/dev/null 2>&1; then
    if ! launchctl bootout "$target" >/dev/null 2>&1; then
      printf 'OmniWM: could not unload %s.\n' "$label" >&2
      exit 1
    fi
  else
    # launchctl uses 113 for a missing service; other errors are not absence.
    status=$?
    if [ "$status" -ne 113 ]; then
      printf 'OmniWM: could not inspect %s.\n' "$label" >&2
      exit 1
    fi
  fi
done
printf 'OmniWM: legacy startup disabled.\n'
