class GentlemanDots < Formula
  desc "Interactive TUI installer for Gentleman.Dots development environment"
  homepage "https://github.com/Gentleman-Programming/Gentleman.Dots"
  version "2.7.2"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/Gentleman-Programming/Gentleman.Dots/releases/download/v#{version}/gentleman-installer-darwin-arm64"
      sha256 "bfd550447bdb5a34c936ecc45878f213f766b4a237db838428c77b3eadfaedde"
    end
    on_intel do
      url "https://github.com/Gentleman-Programming/Gentleman.Dots/releases/download/v#{version}/gentleman-installer-darwin-amd64"
      sha256 "d408b2081737824ff041c2bf55cfed213acd73e10fa17444c6e37838fd315dfb"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/Gentleman-Programming/Gentleman.Dots/releases/download/v#{version}/gentleman-installer-linux-arm64"
      sha256 "abc83ddc5629ab5dc7110178137fa2d7b7b8ddcae16a889a82eb15ae1f9d54c3"
    end
    on_intel do
      url "https://github.com/Gentleman-Programming/Gentleman.Dots/releases/download/v#{version}/gentleman-installer-linux-amd64"
      sha256 "46200b61ec9949b65aefae91328596a3ea399a406782707a07163cd189ec3a34"
    end
  end

  def install
    if OS.mac? && Hardware::CPU.arm?
      bin.install "gentleman-installer-darwin-arm64" => "gentleman-dots"
    elsif OS.mac? && Hardware::CPU.intel?
      bin.install "gentleman-installer-darwin-amd64" => "gentleman-dots"
    elsif OS.linux? && Hardware::CPU.arm?
      bin.install "gentleman-installer-linux-arm64" => "gentleman-dots"
    elsif OS.linux? && Hardware::CPU.intel?
      bin.install "gentleman-installer-linux-amd64" => "gentleman-dots"
    end
  end

  test do
    system "#{bin}/gentleman-dots", "--help"
  end
end
