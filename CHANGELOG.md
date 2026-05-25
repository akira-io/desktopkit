# Changelog

All notable changes to `onyx` are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.0] - 2026-05-25

### Added

- Add DeviceName for a human-friendly machine name.

## [1.1.0] - 2026-05-25

### Added

- Add clipboard package.
- Add notify package.
- Add keyring package.
- Add Hostname for best-effort OS host name.
- Add LoginPath and EnrichedEnviron for login-shell PATH recovery.
- Add GetOrCreate for a stable per-application machine identifier.
- Add IsDark for OS color scheme detection.
- Add Relaunch to start a fresh application instance.

### Documentation

- Explain ResolutionSource removal and Resolver simplification.
- Rewrite for adoption.
- Add Prior art section.
- Restructure docs/ to tens-block numbering.
- List new modules in index, README, and changelog.

### Fixed

- Use Windows.Forms clipboard on Windows.

## [1.0.2] - 2026-05-18

### Changed

- Drop fallback, lookup handles both names and paths.

## [1.0.1] - 2026-05-18

### Changed

- Rename Candidates builder to Resolver.

### Documentation

- Show before/after onyx comparison.

## [1.0.0] - 2026-05-17

### Fixed

- Swap EOL git-cliff-action for taiki-e/install-action.
- Replace unreleased block instead of prepending to avoid duplicate headers.
- Strip cliff header to avoid duplicating changelog top matter.

## [0.2.0] - 2026-05-17

### Added

- Add well-known bin path listers + WithCandidates bulk.

## [0.1.0] - 2026-05-17

### Added

- Desktopkit v0.1.0.


