class GentlemanDots < Formula
  desc "Interactive TUI installer for Gentleman.Dots development environment"
  homepage "https://github.com/Gentleman-Programming/Gentleman.Dots"
  version "2.7.7"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/Gentleman-Programming/Gentleman.Dots/releases/download/v#{version}/gentleman-installer-darwin-arm64"
      sha256 "d36aba07893727f7dc516e2ea8063d2ac86807a8a97efc753492960070d52204"
    end
    on_intel do
      url "https://github.com/Gentleman-Programming/Gentleman.Dots/releases/download/v#{version}/gentleman-installer-darwin-amd64"
      sha256 "612c9e25fbf53542100770e2f1e43440595048420939aecae985ad37a55ac9a8"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/Gentleman-Programming/Gentleman.Dots/releases/download/v#{version}/gentleman-installer-linux-arm64"
      sha256 "f96ac64498c4c4d9de54c5cb2df005367210c4632b5078b07ab2b3632336087a"
    end
    on_intel do
      url "https://github.com/Gentleman-Programming/Gentleman.Dots/releases/download/v#{version}/gentleman-installer-linux-amd64"
      sha256 "1a2b3b913487381f53721912ccf6c0daded39111193bdecf61da9859ec9ffc1e"
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
