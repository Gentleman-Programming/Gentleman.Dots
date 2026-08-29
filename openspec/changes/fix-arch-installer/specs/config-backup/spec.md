# Config Backup Specification

## Purpose

Defines backup behavior for `~/.config` and other backed-up directories, requiring tolerance of non-regular files (Unix sockets, FIFOs, device files) so a single unreadable special file cannot abort the backup of every configuration.

## Requirements

### Requirement: Backup MUST tolerate non-regular files without aborting

`CopyDir` (and the `CreateBackup` flow that calls it) MUST NOT treat a non-regular file (Unix domain socket, named pipe/FIFO, character or block device, or any file `os.FileInfo.Mode().IsRegular()` reports false for) as a fatal error for the whole backup. The installer MUST skip that entry, emit a warning naming the skipped path, and continue copying the remaining regular files and directories.

#### Scenario: Backup directory contains a live Unix socket

- GIVEN a config directory being backed up contains a live Unix domain socket (created via `net.Listen("unix", ...)`) alongside regular files
- WHEN `CreateBackup` runs over that directory
- THEN the backup MUST complete successfully
- AND the socket path MUST be skipped (not copied, not causing failure)
- AND a warning identifying the skipped socket path MUST be emitted
- AND every regular file and subdirectory in the same tree MUST be copied

#### Scenario: Backup directory contains only regular files (nominal path)

- GIVEN a config directory containing only regular files and directories
- WHEN `CreateBackup` runs over that directory
- THEN the backup MUST complete successfully
- AND no warnings about skipped files MUST be emitted

#### Scenario: Multiple non-regular files in the same tree

- GIVEN a config directory containing a socket, a FIFO, and regular files
- WHEN `CreateBackup` runs over that directory
- THEN both the socket and the FIFO MUST be skipped with individual warnings
- AND the backup MUST still complete successfully for the regular files

### Requirement: Skipped files during backup MUST be surfaced to the user

Warnings about skipped non-regular files MUST NOT be silent or logged only at a level the user cannot see. They MUST be visible in the installer's output and MUST be listed in the TUI install summary shown at the end of the run.

#### Scenario: Skipped-file warning reaches the TUI summary

- GIVEN a backup run skipped one or more non-regular files
- WHEN the installer reaches the final TUI summary screen
- THEN the summary MUST list every skipped path from that run
