# Docker Testing

Test the Gentleman.Dots installer in an isolated Ubuntu environment without affecting your system.

## Quick Start

```bash
# Build the image
docker build -f Dockerfile.test -t gentleman-test .

# Run the installer
docker run -it --rm gentleman-test
```

## User Credentials

- **Username**: `testuser`
- **Password**: `test`

The user has passwordless sudo configured, so you won't need the password for most operations.

## Update to Latest Version

If you've already built the image and want to test a newer version:

```bash
# Option 1: Rebuild from scratch (recommended)
docker build -f Dockerfile.test -t gentleman-test --no-cache .
docker run -it --rm gentleman-test

# Option 2: Update inside running container
docker run -it --rm gentleman-test bash
# Then inside the container:
cd /app/installer
git pull origin main
go build -o /usr/local/bin/gentleman.dots ./cmd/gentleman-installer
gentleman.dots
```

## Testing Specific Versions

```bash
# Checkout a specific tag before building
git checkout v2.4.2
docker build -f Dockerfile.test -t gentleman-test .
docker run -it --rm gentleman-test
```

## Interactive Shell (for debugging)

```bash
# Start bash instead of the installer
docker run -it --rm gentleman-test bash

# Then run the installer manually
gentleman.dots
```

## Clean Up

```bash
# Remove the image when done
docker rmi gentleman-test

# Remove all unused Docker resources
docker system prune
```

## Notes

- The container runs Ubuntu 22.04
- All changes are lost when the container exits (`--rm` flag)
- Your host system is completely isolated - nothing is modified
- The installer builds from the local source code at build time
