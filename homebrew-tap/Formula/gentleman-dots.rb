class GentlemanDots < Formula
  desc "Interactive TUI installer for Gentleman.Dots development environment"
  homepage "https://github.com/Gentleman-Programming/Gentleman.Dots"
  version "2.7.5"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/Gentleman-Programming/Gentleman.Dots/releases/download/v#{version}/gentleman-installer-darwin-arm64"
      sha256 "2d95b7b0907ab5eea4598d8cc8ead4cdeaa2b128d497af2f10195165898afcc1"
    end
    on_intel do
      url "https://github.com/Gentleman-Programming/Gentleman.Dots/releases/download/v#{version}/gentleman-installer-darwin-amd64"
      sha256 "52379443a12eff118d75bbd6a5ef2d72cdb137cd4e30b1ef333078d83536d4f9"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/Gentleman-Programming/Gentleman.Dots/releases/download/v#{version}/gentleman-installer-linux-arm64"
      sha256 "fd9170ecb8e9c074ea3e26ecb08641f1bac270fd7ffcdc9a5504ef6d185d2cc0"
    end
    on_intel do
      url "https://github.com/Gentleman-Programming/Gentleman.Dots/releases/download/v#{version}/gentleman-installer-linux-amd64"
      sha256 "ec93be2db8f7473dd261704a16361a4d7f1680a8000c5ff6e441171f7a5f2853"
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
