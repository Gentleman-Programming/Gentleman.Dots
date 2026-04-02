class GentlemanDots < Formula
  desc "Interactive TUI installer for Gentleman.Dots development environment"
  homepage "https://github.com/Gentleman-Programming/Gentleman.Dots"
  version "2.9.11"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/Gentleman-Programming/Gentleman.Dots/releases/download/v#{version}/gentleman-installer-darwin-arm64"
      sha256 "21ede39afc3980a6ca77e4fc6724535575aeb8e34eca84730fad1a20e3dfd09d"
    end
    on_intel do
      url "https://github.com/Gentleman-Programming/Gentleman.Dots/releases/download/v#{version}/gentleman-installer-darwin-amd64"
      sha256 "ef0330f6d22f76bb1457f5cfc26dca00a8b8b0a57b53517f6662da5a5973886f"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/Gentleman-Programming/Gentleman.Dots/releases/download/v#{version}/gentleman-installer-linux-arm64"
      sha256 "be5255e84c263833e490a32a20e8ea2a8041407240b2266f430effd516e7cdde"
    end
    on_intel do
      url "https://github.com/Gentleman-Programming/Gentleman.Dots/releases/download/v#{version}/gentleman-installer-linux-amd64"
      sha256 "fd1d8bec6ded94ee6c294efdd4612915626b71926f8c8e78981f17dd1c93f7d0"
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
