# Docker Testing

Test the Gentleman.Dots installer in an isolated Ubuntu environment without affecting your system.

## Table of Contents

- [Quick Start](#quick-start)
- [Container Environment](#container-environment)
- [Common Tasks](#common-tasks)
  - [Update to Latest Version](#update-to-latest-version)
  - [Test Specific Version](#test-specific-version)
  - [Interactive Debugging](#interactive-debugging)
- [Clean Up](#clean-up)

## Quick Start

```bash
# Build the image
docker build -f Dockerfile.test -t gentleman-test .

# Run the installer
docker run -it --rm gentleman-test
```

## Container Environment

| Setting | Value |
|---------|-------|
| Base Image | Ubuntu 22.04 |
| Username | `testuser` |
| Password | `test` |
| Sudo | Passwordless (NOPASSWD) |
| Installer Path | `/usr/local/bin/gentleman.dots` |

> **Note**: The user has passwordless sudo, so you won't need the password for most operations.

## Common Tasks

### Update to Latest Version

If you've already built the image and want to test a newer version:

**Option 1: Rebuild from scratch (recommended)**

```bash
docker build -f Dockerfile.test -t gentleman-test --no-cache .
docker run -it --rm gentleman-test
```

**Option 2: Update inside running container**

```bash
docker run -it --rm gentleman-test bash
```

Then inside the container:

```bash
cd /app/installer
git pull origin main
go build -o /usr/local/bin/gentleman.dots ./cmd/gentleman-installer
gentleman.dots
```

### Test Specific Version

```bash
# Checkout a specific tag before building
git checkout v2.4.2
docker build -f Dockerfile.test -t gentleman-test .
docker run -it --rm gentleman-test
```

### Interactive Debugging

```bash
# Start bash instead of the installer
docker run -it --rm gentleman-test bash

# Then run the installer manually
gentleman.dots
```

## Clean Up

```bash
# Remove the image
docker rmi gentleman-test

# Remove all unused Docker resources
docker system prune
```

## Notes

- All changes are lost when the container exits (`--rm` flag)
- Your host system is completely isolated - nothing is modified
- The installer builds from local source at image build time
