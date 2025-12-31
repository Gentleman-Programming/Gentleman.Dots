class GentlemanDots < Formula
  desc "Interactive TUI installer for Gentleman.Dots development environment"
  homepage "https://github.com/Gentleman-Programming/Gentleman.Dots"
  version "1.0.0"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/Gentleman-Programming/Gentleman.Dots/releases/download/v#{version}/gentleman-dots-darwin-arm64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_DARWIN_ARM64"
    end
    on_intel do
      url "https://github.com/Gentleman-Programming/Gentleman.Dots/releases/download/v#{version}/gentleman-dots-darwin-amd64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_DARWIN_AMD64"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/Gentleman-Programming/Gentleman.Dots/releases/download/v#{version}/gentleman-dots-linux-arm64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_LINUX_ARM64"
    end
    on_intel do
      url "https://github.com/Gentleman-Programming/Gentleman.Dots/releases/download/v#{version}/gentleman-dots-linux-amd64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_LINUX_AMD64"
    end
  end

  def install
    bin.install "gentleman-dots"
  end

  test do
    assert_match "Gentleman.Dots Installer", shell_output("#{bin}/gentleman-dots --version")
  end
end
