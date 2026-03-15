---
name: homebrew-release
description: >
  Release workflow for Gentleman-Programming homebrew-tap projects (GGA, Gentleman.Dots).
  Trigger: When user asks to release, bump version, update homebrew, or publish a new version.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "2.0"
---

## When to Use

- User asks to "release", "bump", "publish", or "update homebrew"
- User mentions a new version number
- User says "homebrew-tap" or "formula"
- After creating a git tag in GGA or Gentleman.Dots repos

## Supported Projects

| Project | Repo | Formula | Tag Format | Type |
|---------|------|---------|------------|------|
| GGA | `gentleman-guardian-angel` | `gga.rb` | `V{version}` (e.g., `V2.6.2`) | Tarball (builds from source) |
| Gentleman.Dots | `Gentleman.Dots` | `gentleman-dots.rb` | `v{version}` (e.g., `v2.5.1`) | Pre-built binaries |

---

## Gentleman.Dots Release Process (Pre-built Binaries)

### Step 1: Build Binaries

```bash
cd installer
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o gentleman-installer-darwin-amd64 ./cmd/gentleman-installer
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o gentleman-installer-darwin-arm64 ./cmd/gentleman-installer
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o gentleman-installer-linux-amd64 ./cmd/gentleman-installer
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o gentleman-installer-linux-arm64 ./cmd/gentleman-installer
```

### Step 2: Commit, Tag, and Push

```bash
git add -A
git commit -m "feat/fix: description"
git push origin main
git tag v{VERSION}
git push origin v{VERSION}
```

### Step 3: Create GitHub Release

```bash
gh release create v{VERSION} \
  installer/gentleman-installer-darwin-amd64 \
  installer/gentleman-installer-darwin-arm64 \
  installer/gentleman-installer-linux-amd64 \
  installer/gentleman-installer-linux-arm64 \
  --title "v{VERSION}" \
  --notes "## Changes
- {description of changes}"
```

### Step 4: Get SHA256 of Binaries

```bash
shasum -a 256 installer/gentleman-installer-darwin-arm64 installer/gentleman-installer-darwin-amd64 installer/gentleman-installer-linux-amd64 installer/gentleman-installer-linux-arm64
```

### Step 5: Update Formula

Update `homebrew-tap/Formula/gentleman-dots.rb`:

```ruby
class GentlemanDots < Formula
  desc "Interactive TUI installer for Gentleman.Dots development environment"
  homepage "https://github.com/Gentleman-Programming/Gentleman.Dots"
  version "{VERSION}"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/Gentleman-Programming/Gentleman.Dots/releases/download/v#{version}/gentleman-installer-darwin-arm64"
      sha256 "{SHA256_DARWIN_ARM64}"
    end
    on_intel do
      url "https://github.com/Gentleman-Programming/Gentleman.Dots/releases/download/v#{version}/gentleman-installer-darwin-amd64"
      sha256 "{SHA256_DARWIN_AMD64}"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/Gentleman-Programming/Gentleman.Dots/releases/download/v#{version}/gentleman-installer-linux-amd64"
      sha256 "{SHA256_LINUX_AMD64}"
    end
  end

  def install
    if OS.mac? && Hardware::CPU.arm?
      bin.install "gentleman-installer-darwin-arm64" => "gentleman-dots"
    elsif OS.mac? && Hardware::CPU.intel?
      bin.install "gentleman-installer-darwin-amd64" => "gentleman-dots"
    elsif OS.linux? && Hardware::CPU.intel?
      bin.install "gentleman-installer-linux-amd64" => "gentleman-dots"
    end
  end

  test do
    system "#{bin}/gentleman-dots", "--help"
  end
end
```

### Step 6: Commit to Both Repos

```bash
# In Gentleman.Dots repo
git add homebrew-tap/Formula/gentleman-dots.rb
git commit -m "chore(homebrew): bump version to v{VERSION}"
git push origin main

# In homebrew-tap repo
cd /tmp && rm -rf homebrew-tap
git clone git@github.com:Gentleman-Programming/homebrew-tap.git
cp {path-to}/Gentleman.Dots/homebrew-tap/Formula/gentleman-dots.rb /tmp/homebrew-tap/Formula/
cd /tmp/homebrew-tap
git add -A
git commit -m "chore: bump gentleman-dots to v{VERSION}"
git push origin main
```

---

## GGA Release Process (Tarball - Builds from Source)

### Step 1: Verify Tag Exists

```bash
git tag --list | tail -5
```

### Step 2: Get SHA256 of Tarball

```bash
curl -sL https://github.com/Gentleman-Programming/gentleman-guardian-angel/archive/refs/tags/V{VERSION}.tar.gz | shasum -a 256
```

### Step 3: Update Formula

Update `homebrew-tap/Formula/gga.rb`:

```ruby
url "https://github.com/Gentleman-Programming/gentleman-guardian-angel/archive/refs/tags/V{VERSION}.tar.gz"
sha256 "{NEW_SHA256}"
version "{VERSION}"
```

### Step 4: Commit and Push

```bash
cd ~/work/homebrew-tap
git add -A
git commit -m "chore: bump gga to V{VERSION}"
git push
```

---

## Quick Reference Commands

```bash
# Build Gentleman.Dots binaries
cd installer && GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o gentleman-installer-darwin-amd64 ./cmd/gentleman-installer && GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o gentleman-installer-darwin-arm64 ./cmd/gentleman-installer && GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o gentleman-installer-linux-amd64 ./cmd/gentleman-installer && GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o gentleman-installer-linux-arm64 ./cmd/gentleman-installer

# SHA256 for binaries
shasum -a 256 installer/gentleman-installer-*

# SHA256 for GGA tarball
curl -sL https://github.com/Gentleman-Programming/gentleman-guardian-angel/archive/refs/tags/V{VERSION}.tar.gz | shasum -a 256

# Create GitHub release with binaries
gh release create v{VERSION} installer/gentleman-installer-* --title "v{VERSION}" --notes "## Changes"
```

---

## Checklist

### Gentleman.Dots
- [ ] Binaries built for all platforms (darwin-amd64, darwin-arm64, linux-amd64, linux-arm64)
- [ ] Changes committed and pushed to main
- [ ] Tag created and pushed (v{VERSION})
- [ ] GitHub release created with binaries attached
- [ ] SHA256 computed for all binaries
- [ ] Formula updated in Gentleman.Dots repo
- [ ] Formula copied to homebrew-tap repo and pushed

### GGA
- [ ] Tag exists (V{VERSION})
- [ ] SHA256 computed from tarball
- [ ] Formula updated in homebrew-tap
- [ ] Committed and pushed
