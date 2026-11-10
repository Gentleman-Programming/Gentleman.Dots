"""Independent S2-S5 preservation and dual-installation regressions."""
import hashlib
import unittest
from test_omniwm import ActivationTest


class IndependentTest(ActivationTest):
    def snapshot(self):
        result = {}
        for base in (self.home, self.state, self.system_app):
            paths = [base] + (sorted(base.rglob("*")) if base.is_dir() else [])
            for path in paths:
                if path.is_symlink():
                    value = ("link", str(path.readlink()))
                elif path.is_file():
                    value = ("file", path.stat().st_mode & 0o777,
                             hashlib.sha256(path.read_bytes()).hexdigest())
                elif path.is_dir():
                    value = ("dir",)
                else:
                    value = ("absent",)
                result[str(path)] = value
        return result

    def test_rejection_hashes(self):
        cases = [
            ({"MOCK_OS": "Linux"}, "", ["uname -s"]),
            ({"MOCK_ARCH": "x86_64"}, "", ["uname -s", "uname -m"]),
            ({"MOCK_VERSION": "25.7.1"}, "", self.platform_commands()),
            ({"MOCK_VERSION": "26..1"}, "", self.platform_commands()),
            ({"MOCK_BREW": "fail"}, "failed", self.platform_commands()
             + ["brew install --cask omniwm"]),
            ({"MOCK_BREW": "no-app"}, "missing", self.platform_commands()
             + ["brew install --cask omniwm"]),
        ]
        for overrides, failure, commands in cases:
            with self.subTest(overrides=overrides):
                self.log.write_text("")
                before = self.snapshot()
                stdout = "OmniWM: installing Homebrew cask.\n" if failure else ""
                stderr = {
                    "": "OmniWM requires Apple Silicon and macOS 26 or later.\n",
                    "failed": "mock download failed\nOmniWM: Homebrew installation failed; legacy services unchanged.\n",
                    "missing": "OmniWM: installed app not found; legacy services unchanged.\n",
                }[failure]
                self.assert_result(self.activate(**overrides), 1, stdout, stderr, commands)
                self.assertEqual(self.snapshot(), before)

    def test_dual_app_preservation(self):
        user = self.home / "Applications/OmniWM.app"
        for app, data in ((user, b"user install\x00\n"),
                          (self.system_app, b"system install\x00\n")):
            app.mkdir(parents=True)
            (app / "marker").write_bytes(data)
        self.assert_result(self.activate(MOCK_BREW="fail", MOCK_VERSION="26.0.1"), 0,
                           "OmniWM: existing app preserved.\nOmniWM: legacy startup disabled.\n",
                           "", self.platform_commands() + self.service_commands())
        self.assertEqual((user / "marker").read_bytes(), b"user install\x00\n")
        self.assertEqual((self.system_app / "marker").read_bytes(), b"system install\x00\n")


if __name__ == "__main__":
    suite = unittest.TestSuite(IndependentTest(name) for name in
                              ("test_rejection_hashes", "test_dual_app_preservation"))
    result = unittest.TextTestRunner(verbosity=2).run(suite)
    raise SystemExit(not result.wasSuccessful())
