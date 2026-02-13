# Release

Releases are automated via GitHub Actions.

## Tagging

Pushing a tag like `v0.1.0` triggers `.github/workflows/release.yml` which:
- builds release archives for supported platforms
- optionally builds a Linux AppImage for the desktop app
- generates a Homebrew formula artifact
- publishes all artifacts as GitHub release assets

## Scripts

- `scripts/release/build-artifacts.sh`: builds tarballs/checksums
- `scripts/release/build-fogapp-appimage.sh`: AppImage packaging for `fogapp`
- `scripts/release/generate-homebrew-formula.sh`: formula generator

## CI

`.github/workflows/ci.yml` runs:
- `go test ./...` on macOS + Linux
- `go build` for all CLI binaries
- `go build -tags desktop ./cmd/fogapp` on Linux (with GTK/WebKit deps installed)

