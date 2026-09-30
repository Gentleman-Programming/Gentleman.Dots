# Dry-Run Safety Specification

## Purpose

Defines the required behavior of `--dry-run`: it MUST perform zero real-world mutations while still communicating what the installer would have done, per upstream issue #190. This is stricter than `--test`, which still performs package-manager, sudo, and network operations and is therefore not an acceptable substitute.

## Requirements

### Requirement: `--dry-run` MUST perform zero mutations of any kind

When invoked with `--dry-run`, the installer MUST perform no process execution with side effects, no package-manager operations, no network operations (including clone/fetch), no filesystem mutations (no writes, deletes, or backups), and no shell configuration changes. This MUST hold across every `Run*` entry point and at every mutating call site, including `CopyFile`, `CopyDir`, `CreateBackup`, and both `Gentleman.Dots` deletion sites.

#### Scenario: Dry run on a dirty home produces zero filesystem changes

- GIVEN a user's home directory with existing dotfiles and configuration
- AND the installer is invoked with `--dry-run`
- WHEN the full installer flow runs to completion
- THEN no files under the home directory MUST be created, modified, or deleted
- AND no backup directory MUST be created

#### Scenario: Dry run performs no clone

- GIVEN the installer is invoked with `--dry-run`
- WHEN the clone step would normally run
- THEN no `git clone` (or any network fetch) MUST occur
- AND no `Gentleman.Dots` directory MUST be created or removed

#### Scenario: Dry run performs no package-manager or sudo invocation

- GIVEN the installer is invoked with `--dry-run`
- WHEN the package install step would normally run
- THEN no `pacman`, `apt`, `dnf`, `apk`, `brew`, or `pkg` command MUST be executed
- AND no `sudo` invocation MUST occur

### Requirement: `--dry-run` MUST print the planned operations

A dry run MUST NOT be silent. For every operation it suppresses, the installer MUST print a description of the planned operation (e.g., which packages would be installed, which files would be backed up or copied, which directory would be cloned or removed) so the user can review the plan before running for real.

#### Scenario: Dry run prints the install plan

- GIVEN the installer is invoked with `--dry-run`
- WHEN the flow reaches the package install step
- THEN the installer MUST print the list of packages that would have been installed for the detected platform

#### Scenario: Dry run prints the backup and clone plan

- GIVEN the installer is invoked with `--dry-run`
- WHEN the flow reaches the backup and clone steps
- THEN the installer MUST print which paths would have been backed up
- AND MUST print that the repository would have been cloned (and its target path)

### Requirement: A single central dry-run gate MUST govern all mutating call sites

Dry-run detection MUST be centralized (e.g., `system.IsDryRun()`) rather than checked ad hoc, so every mutating call site consistently suppresses side effects and prints its planned operation. `--test` mode MUST NOT be treated as satisfying this requirement, since it still performs package-manager, sudo, and network operations.

#### Scenario: Every Run* function respects the same dry-run gate

- GIVEN `--dry-run` is active
- WHEN any of the four `Run*` functions is invoked
- THEN each MUST consult the same central dry-run gate
- AND each MUST suppress its mutations and print its planned operation consistently

#### Scenario: `--test` mode is not mistaken for dry-run

- GIVEN the installer is invoked with `--test` (not `--dry-run`)
- WHEN the package install step runs
- THEN package-manager, sudo, and network operations MUST still execute as `--test` does not suppress mutations
