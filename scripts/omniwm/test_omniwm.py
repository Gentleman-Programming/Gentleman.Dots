#!/usr/bin/env python3
"""Exercise the OmniWM Home Manager activation with fake macOS commands only."""

import json
import re
import shutil
import subprocess
import tempfile
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
ACTIVATION = ROOT / "scripts/omniwm/activate.sh"
LABELS = (
    "com.asmvik.yabai",
    "com.koekeishiya.yabai",
    "com.koekeishiya.skhd",
    "homebrew.mxcl.sketchybar",
    "com.felixkratz.sketchybar",
)


class ImportsTest(unittest.TestCase):
    def test_only_window_manager_import_is_replaced(self):
        flake = (ROOT / "flake.nix").read_text()
        imports = re.findall(r"^\s*(\./[\w-]+\.nix)\s*(?:#.*)?$", flake, re.M)
        self.assertEqual(imports, [
            "./nushell.nix", "./ghostty.nix", "./alacritty.nix", "./zed.nix",
            "./television.nix", "./wezterm.nix", "./kitty.nix", "./zellij.nix",
            "./tmux.nix", "./tmux-agents.nix", "./fish.nix", "./starship.nix",
            "./nvim.nix", "./zsh.nix", "./oil-scripts.nix", "./opencode.nix",
            "./claude.nix", "./engram.nix", "./herdr.nix", "./omniwm.nix",
            "./raycast.nix",
        ])
        self.assertTrue((ROOT / "omniwm.nix").is_file())


class ActivationTest(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory(prefix="omniwm test ; ")
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        self.home = self.root / "home $HOME; with spaces"
        self.system_app = self.root / "system Applications;" / "OmniWM.app"
        self.bin = self.root / "mock bin"
        self.state = self.root / "service state"
        self.log = self.root / "commands.log"
        for path in (self.home, self.bin, self.state / "loaded", self.state / "disabled",
                     self.state / "stopped"):
            path.mkdir(parents=True)
        self.log.write_text("")
        self.settings = self.home / ".config/omniwm/settings.toml"
        self.settings.parent.mkdir(parents=True)
        self.settings.write_text("# GUI-owned sentinel\n")
        self.settings.chmod(0o640)
        self.legacy = [self.home / ".config/yabai/yabairc",
                       self.home / ".config/skhd/skhdrc",
                       self.home / ".config/sketchybar/sketchybarrc",
                       self.home / ".config/omniwm/settings.toml.backup"]
        for path in self.legacy:
            path.parent.mkdir(parents=True, exist_ok=True)
            path.write_text("legacy sentinel\n")
        self.spaces = self.home / ".config/omniwm/nativeSpaces"
        self.spaces.write_text("spaces sentinel\n")
        (self.state / "loaded/com.example.unrelated").touch()

        self.mock("uname", r'''printf 'uname %s\n' "$*" >> "$MOCK_LOG"
case "$1" in -s) printf '%s\n' "$MOCK_OS";; -m) printf '%s\n' "$MOCK_ARCH";; *) exit 90;; esac''')
        self.mock("sw_vers", r'''printf 'sw_vers %s\n' "$*" >> "$MOCK_LOG"
[ "$1" = -productVersion ] || exit 90
printf '%s\n' "$MOCK_VERSION"''')
        self.mock("id", r'''printf 'id %s\n' "$*" >> "$MOCK_LOG"
[ "$1" = -u ] || exit 90
printf '501\n' ''')
        self.mock("brew", r'''printf 'brew %s\n' "$*" >> "$MOCK_LOG"
[ "$*" = 'install --cask omniwm' ] || exit 90
if [ "$MOCK_BREW" = fail ]; then
  printf 'mock download failed\n' >&2
  exit 42
fi
if [ "$MOCK_BREW" = success ]; then mkdir -p "$MOCK_SYSTEM_APP"; fi''')
        self.mock("launchctl", r'''printf 'launchctl %s\n' "$*" >> "$MOCK_LOG"
if [ "$1" = print ] && [ "$2" = gui/501 ]; then
  [ "$MOCK_DOMAIN_FAIL" != 1 ]
  exit $?
fi
case "$2" in gui/501/*) label="${2#gui/501/}";; *) exit 90;; esac
case "$1" in
  disable)
    [ "$label" != "$MOCK_FAIL_DISABLE" ] || exit 42
    touch "$MOCK_STATE/disabled/$label";;
  print)
    [ "$label" != "$MOCK_FAIL_PRINT" ] || exit 42
    if [ -e "$MOCK_STATE/loaded/$label" ] && [ ! -e "$MOCK_STATE/stopped/$label" ]; then
      exit 0
    fi
    exit 113;;
  bootout)
    [ "$label" != "$MOCK_FAIL_BOOTOUT" ] || exit 42
    touch "$MOCK_STATE/stopped/$label";;
  *) exit 90;;
esac''')
        self.env = {
            "HOME": str(self.home),
            "PATH": str(self.bin) + ":/usr/bin:/bin",
            "MOCK_LOG": str(self.log),
            "MOCK_STATE": str(self.state),
            "MOCK_SYSTEM_APP": str(self.system_app),
            "MOCK_OS": "Darwin",
            "MOCK_ARCH": "arm64",
            "MOCK_VERSION": "26.0",
            "MOCK_BREW": "success",
            "MOCK_DOMAIN_FAIL": "0",
            "MOCK_FAIL_DISABLE": "",
            "MOCK_FAIL_PRINT": "",
            "MOCK_FAIL_BOOTOUT": "",
        }
        # Execute an exact copy from a path containing spaces and shell metacharacters.
        self.script = self.root / "repo ; $not-a-command" / "activate.sh"
        self.script.parent.mkdir()
        shutil.copyfile(ACTIVATION, self.script)

    def mock(self, name, body):
        path = self.bin / name
        path.write_text("#!/bin/bash\nset -u\n" + body + "\n")
        path.chmod(0o755)

    def activate(self, **overrides):
        env = dict(self.env, **overrides)
        return subprocess.run(["/bin/bash", str(self.script), str(self.system_app)],
                              env=env, text=True, capture_output=True, check=False)

    def preflight(self, **overrides):
        env = dict(self.env, **overrides)
        return subprocess.run(["/bin/bash", str(self.script), "--check-platform"],
                              env=env, text=True, capture_output=True, check=False)

    def evaluated_activation(self):
        nix = shutil.which("nix-instantiate")
        if not nix:
            self.skipTest("nix-instantiate is unavailable")
        expression = (
            'let dag = { entryBefore = before: data: { inherit before data; }; '
            'entryAfter = after: data: { inherit after data; }; }; '
            'in (import ./omniwm.nix { lib.hm.dag = dag; pkgs.bash = ""; }).home.activation'
        )
        result = subprocess.run([nix, "--eval", "--strict", "--json", "-E", expression],
                                cwd=ROOT, env=self.env, text=True, capture_output=True, check=False)
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(result.stderr, "")
        return json.loads(result.stdout)

    def assert_result(self, result, status, stdout, stderr, commands):
        self.assertEqual(result.returncode, status)
        self.assertEqual(result.stdout, stdout)
        self.assertEqual(result.stderr, stderr)
        self.assertEqual(self.log.read_text().splitlines(), commands)
        self.assertEqual(self.settings.read_text(), "# GUI-owned sentinel\n")
        self.assertEqual(self.settings.stat().st_mode & 0o777, 0o640)
        self.assertFalse(self.settings.is_symlink())
        for path in self.legacy:
            self.assertEqual(path.read_text(), "legacy sentinel\n")
        self.assertEqual(self.spaces.read_text(), "spaces sentinel\n")
        self.assertTrue((self.state / "loaded/com.example.unrelated").exists())
        self.assertFalse((self.state / "disabled/com.example.unrelated").exists())
        self.assertFalse((self.state / "stopped/com.example.unrelated").exists())

    def assert_services_unchanged(self):
        self.assertEqual(list((self.state / "disabled").iterdir()), [])
        self.assertEqual(list((self.state / "stopped").iterdir()), [])

    @staticmethod
    def platform_commands():
        return ["uname -s", "uname -m", "sw_vers -productVersion"]

    @staticmethod
    def service_commands(loaded=()):
        commands = ["id -u", "launchctl print gui/501"]
        for label in LABELS:
            commands.extend([f"launchctl disable gui/501/{label}",
                             f"launchctl print gui/501/{label}"])
            if label in loaded:
                commands.append(f"launchctl bootout gui/501/{label}")
        return commands

    def test_existing_user_app_and_configuration_are_preserved_and_repeat_is_safe(self):
        (self.home / "Applications/OmniWM.app").mkdir(parents=True)
        (self.state / "loaded" / LABELS[0]).touch()
        commands = self.platform_commands() + self.service_commands([LABELS[0]])
        self.assert_result(self.activate(), 0,
                           "OmniWM: existing app preserved.\nOmniWM: legacy startup disabled.\n",
                           "", commands)
        self.assertEqual({p.name for p in (self.state / "disabled").iterdir()}, set(LABELS))
        self.assertEqual({p.name for p in (self.state / "stopped").iterdir()}, {LABELS[0]})
        self.log.write_text("")
        self.assert_result(self.activate(), 0,
                           "OmniWM: existing app preserved.\nOmniWM: legacy startup disabled.\n",
                           "", self.platform_commands() + self.service_commands())

    def test_all_five_loaded_labels_are_unloaded_on_a_later_supported_version(self):
        (self.home / "Applications/OmniWM.app").mkdir(parents=True)
        for label in LABELS:
            (self.state / "loaded" / label).touch()
        self.assert_result(self.activate(MOCK_VERSION="27.0"), 0,
                           "OmniWM: existing app preserved.\nOmniWM: legacy startup disabled.\n",
                           "", self.platform_commands() + self.service_commands(LABELS))
        self.assertEqual({p.name for p in (self.state / "disabled").iterdir()}, set(LABELS))
        self.assertEqual({p.name for p in (self.state / "stopped").iterdir()}, set(LABELS))

    def test_existing_system_app_is_preserved_without_brew(self):
        self.system_app.mkdir(parents=True)
        self.assert_result(self.activate(), 0,
                           "OmniWM: existing app preserved.\nOmniWM: legacy startup disabled.\n",
                           "", self.platform_commands() + self.service_commands())
        self.assertTrue(self.system_app.is_dir())

    def test_missing_app_installs_cask_before_disabling_services(self):
        self.assert_result(self.activate(), 0,
                           "OmniWM: installing Homebrew cask.\nOmniWM: legacy startup disabled.\n",
                           "", self.platform_commands() + ["brew install --cask omniwm"]
                           + self.service_commands())
        self.assertTrue(self.system_app.is_dir())

    def test_home_manager_activation_wiring_and_dry_run(self):
        activation = self.evaluated_activation()
        self.assertEqual(activation["checkOmniWM"]["before"], ["writeBoundary"])
        self.assertEqual(activation["setupOmniWM"]["after"], ["linkGeneration"])
        guard = subprocess.run(["/bin/bash", "-e", "-c", activation["checkOmniWM"]["data"]],
                               env=dict(self.env, MOCK_ARCH="x86_64"),
                               text=True, capture_output=True, check=False)
        self.assert_result(guard, 1, "",
                           "OmniWM requires Apple Silicon and macOS 26 or later.\n",
                           ["uname -s", "uname -m"])
        self.assert_services_unchanged()
        self.log.write_text("")
        (self.home / "Applications/OmniWM.app").mkdir(parents=True)
        script = activation["setupOmniWM"]["data"]
        dry_run = subprocess.run(["/bin/bash", "-e", "-c", script],
                                 env=dict(self.env, DRY_RUN_CMD="echo"),
                                 text=True, capture_output=True, check=False)
        self.assert_result(dry_run, 0, "", "", [])
        self.assert_services_unchanged()
        actual = subprocess.run(["/bin/bash", "-e", "-c", script],
                                env=self.env, text=True, capture_output=True, check=False)
        self.assert_result(actual, 0,
                           "OmniWM: existing app preserved.\nOmniWM: legacy startup disabled.\n",
                           "", self.platform_commands() + self.service_commands())

    def test_preflight_rejects_unsupported_platform_before_any_changes(self):
        self.assert_result(self.preflight(), 0, "", "", self.platform_commands())
        self.assert_services_unchanged()
        self.log.write_text("")
        self.assert_result(self.preflight(MOCK_ARCH="x86_64"), 1, "",
                           "OmniWM requires Apple Silicon and macOS 26 or later.\n",
                           ["uname -s", "uname -m"])
        self.assert_services_unchanged()

    def test_unsupported_platforms_leave_configuration_and_services_unchanged(self):
        cases = [
            ({"MOCK_OS": "Linux"}, ["uname -s"]),
            ({"MOCK_ARCH": "x86_64"}, ["uname -s", "uname -m"]),
            ({"MOCK_VERSION": "25.7.1"}, self.platform_commands()),
            ({"MOCK_VERSION": "unknown"}, self.platform_commands()),
            ({"MOCK_VERSION": "999999999999999999999999999999"}, self.platform_commands()),
        ]
        for overrides, commands in cases:
            with self.subTest(overrides=overrides):
                self.log.write_text("")
                self.assert_result(self.activate(**overrides), 1, "",
                                   "OmniWM requires Apple Silicon and macOS 26 or later.\n",
                                   commands)
                self.assert_services_unchanged()

    def test_failed_download_leaves_configuration_and_services_unchanged(self):
        self.assert_result(self.activate(MOCK_BREW="fail"), 1,
                           "OmniWM: installing Homebrew cask.\n",
                           "mock download failed\nOmniWM: Homebrew installation failed; legacy services unchanged.\n",
                           self.platform_commands() + ["brew install --cask omniwm"])
        self.assert_services_unchanged()
        self.assertFalse(self.system_app.exists())

    def test_successful_brew_without_app_does_not_disable_services(self):
        self.assert_result(self.activate(MOCK_BREW="no-app"), 1,
                           "OmniWM: installing Homebrew cask.\n",
                           "OmniWM: installed app not found; legacy services unchanged.\n",
                           self.platform_commands() + ["brew install --cask omniwm"])
        self.assert_services_unchanged()

    def test_unavailable_gui_domain_does_not_disable_services(self):
        (self.home / "Applications/OmniWM.app").mkdir(parents=True)
        self.assert_result(self.activate(MOCK_DOMAIN_FAIL="1"), 1,
                           "OmniWM: existing app preserved.\n",
                           "OmniWM: current-user GUI launchd domain unavailable.\n",
                           self.platform_commands() + ["id -u", "launchctl print gui/501"])
        self.assert_services_unchanged()

    def test_inspection_errors_are_not_mistaken_for_absent_services(self):
        (self.home / "Applications/OmniWM.app").mkdir(parents=True)
        label = LABELS[0]
        self.assert_result(self.activate(MOCK_FAIL_PRINT=label), 1,
                           "OmniWM: existing app preserved.\n",
                           f"OmniWM: could not inspect {label}.\n",
                           self.platform_commands() + ["id -u", "launchctl print gui/501",
                           f"launchctl disable gui/501/{label}",
                           f"launchctl print gui/501/{label}"])
        self.assertEqual({p.name for p in (self.state / "disabled").iterdir()}, {label})
        self.assertEqual(list((self.state / "stopped").iterdir()), [])

    def test_disable_and_bootout_failures_are_reported(self):
        (self.home / "Applications/OmniWM.app").mkdir(parents=True)
        prefix = self.platform_commands() + ["id -u", "launchctl print gui/501"]
        label = LABELS[0]
        self.assert_result(self.activate(MOCK_FAIL_DISABLE=label), 1,
                           "OmniWM: existing app preserved.\n",
                           f"OmniWM: could not disable {label}.\n",
                           prefix + [f"launchctl disable gui/501/{label}"])
        self.assert_services_unchanged()
        self.log.write_text("")
        (self.state / "loaded" / label).touch()
        self.assert_result(self.activate(MOCK_FAIL_BOOTOUT=label), 1,
                           "OmniWM: existing app preserved.\n",
                           f"OmniWM: could not unload {label}.\n",
                           prefix + [f"launchctl disable gui/501/{label}",
                                     f"launchctl print gui/501/{label}",
                                     f"launchctl bootout gui/501/{label}"])
        self.assertEqual({p.name for p in (self.state / "disabled").iterdir()}, {label})
        self.assertEqual(list((self.state / "stopped").iterdir()), [])


if __name__ == "__main__":
    unittest.main(verbosity=2)
