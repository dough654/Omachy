# Integration Tests

Omachy's integration tests run a full install/uninstall cycle inside a disposable macOS VM using [Tart](https://tart.run/). This ensures the system returns to a clean state after uninstall — configs removed, processes killed, defaults restored.

## Prerequisites

```sh
brew install cirruslabs/cli/tart cirruslabs/cli/sshpass
```

- **Tart** — CLI tool for running macOS/Linux VMs on Apple Silicon via Apple's Virtualization.framework
- **sshpass** — non-interactive SSH password authentication (the base VM images use `admin`/`admin`)
- **Apple Silicon Mac** required (M1+)

## Setup

Create the base VM image. This pulls a macOS Tahoe image from Tart's registry, boots it, installs Homebrew, and saves it as `omachy-base`. Only needs to be run once.

```sh
make test-setup
```

This takes ~10-15 minutes the first time (downloads a ~15GB macOS image). The resulting `omachy-base` image is stored locally by Tart and reused for all future test runs.

## Running tests

```sh
make test-integration
```

Each run:

1. Clones a fresh VM from `omachy-base` (takes seconds, copy-on-write)
2. Boots it headless
3. Builds the omachy binary and copies it into the VM
4. Verifies pre-install state (no aerospace/sketchybar/borders running)
5. Runs `omachy install --force --skip-backup`
6. Verifies post-install state (configs deployed, defaults applied, zshrc updated)
7. Runs `omachy uninstall`
8. Verifies post-uninstall state (processes killed, configs removed, defaults restored)
9. Deletes the disposable VM clone

## Manual testing

You can open a VM with a full GUI window for manual testing. This is useful for visually verifying the tiling WM, menu bar, and other UI changes.

### 1. Create and launch the VM

```sh
# Clone a fresh VM from the base image
tart clone omachy-base omachy-manual

# Open it with a GUI window (this takes over the terminal)
tart run omachy-manual
```

A macOS desktop will appear in a native window. Leave this running.

### 2. Build and copy the binary

In a **separate terminal** on your host, build omachy and copy it into the VM:

```sh
cd ~/dev/omachy
go build -o /tmp/omachy-test .
sshpass -p admin scp -o StrictHostKeyChecking=no /tmp/omachy-test admin@$(tart ip omachy-manual):/tmp/omachy
```

### 3. Run omachy inside the VM GUI

In the VM's GUI window, open **Terminal.app** (Spotlight: Cmd+Space → "Terminal") and run:

```sh
eval "$(/opt/homebrew/bin/brew shellenv)"
chmod +x /tmp/omachy
/tmp/omachy install --force --skip-backup
```

When AeroSpace starts, a system dialog will appear asking for Accessibility permissions — grant access.

### 4. Verify

After install completes, verify directly in the VM GUI:
- AeroSpace is tiling windows
- SketchyBar is visible in the menu bar
- JankyBorders outlines the active window
- Dock is auto-hidden

### 5. Test uninstall

In the same Terminal window inside the VM:

```sh
/tmp/omachy uninstall
```

Verify:
- AeroSpace, SketchyBar, and JankyBorders are no longer running
- The Dock and menu bar are back to their default appearance

### 6. Clean up

Close the GUI window (or Ctrl+C the `tart run` terminal), then:

```sh
tart delete omachy-manual
```

## Managing VMs

```sh
# List all VMs
tart list

# Delete a VM
tart delete <vm-name>

# Rebuild the base image (if needed)
tart delete omachy-base
make test-setup
```

## Troubleshooting

- **"Base VM not found"** — Run `make test-setup` first.
- **VM won't boot** — Check that you're on Apple Silicon and have enough disk space (~30GB for the base image).
- **SSH connection refused** — The VM may still be booting. The test script retries for up to 5 minutes.
- **AeroSpace accessibility permissions** — In a VM, the accessibility permission dialog may not auto-resolve. Tests account for this by not asserting that AeroSpace is running post-install (it may be blocked by the permission prompt).
