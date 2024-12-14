# Gobatmon

[![Go](https://img.shields.io/badge/Go-%2300ADD8.svg?&logo=go&logoColor=white)](#)
[![GitHub Release](https://img.shields.io/github/v/release/ulinja/gobatmon?logo=GitHub&label=Version&color=green)](https://github.com/ulinja/gobatmon/releases/latest)
[![AUR Version](https://img.shields.io/aur/version/gobatmon?logo=Arch%20Linux&label=AUR)](https://aur.archlinux.org/packages/gobatmon)

Gobatmon is a simple battery level monitoring daemon for Linux systems.

It keeps an eye on the current charge level of the battery in your laptop and triggers desktop notifications to warn you
if your battery is low.

Gobatmon is super-low on resources to conserve CPU cycles and thus battery life.

**System requirements:**

- a Linux Desktop system
- a battery
- a running notification server (`swaync`, `dunst` etc.)

> If you are using a desktop environment (Xfce/Gnome/Plasma etc.) you most likely don't need this software.

> Gobatmon uses DBUS to dispatch desktop notifications.

Gobatmon's behavior is fully configurable using commandline arguments.

## Usage

Simply run it by executing `gobatmon`. Gobatmon will run continuously.

The following options can be configured:

```
gobatmon [OPTIONS]

Options:
    --normal-warning-threshold uint
        Threshold percentage below which a normal low battery warning is triggered (default 20)
    --critical-warning-threshold uint
        Threshold percentage below which a critical low battery warning is triggered (default 10)
    --normal-warning-reminder-timeout uint
        Timeout in seconds after which a normal low battery warning is repeated (default 600)
    --critical-warning-reminder-timeout uint
        Timeout in seconds after which a critical low battery warning is repeated (default 300)
    --disable-icons (default false)
        Do not show icons in warning notifications
    --normal-warning-icon-name string
        Name of the icon to use for normal low battery warning notifications (default "battery-low")
    --critical-warning-icon-name string
        Name of the icon to use for critical low battery warning notifications (default "battery-caution")
    --poll-rate uint
        Poll rate for checking battery status in seconds (default 60)
    --version
        Show version information and exit
    --help
        Show help message and exit
```

While charging the battery or above the normal warning threshold, gobatmon will not display any notifications, and will
poll the battery status to watch for changes.

When running on battery power and below the normal/critical warning thresholds (20%/10% by default), gobatmon will check
the battery status and notify you with a reminder every 10 minutes / 5 minutes (by default) respectively.

That's it.

## Installation

### Arch Linux

Gobatmon is available in the AUR as [gobatmon](https://aur.archlinux.org/packages/gobatmon).

### NixOS

> :construction: Will be added in the future.

### Manual

Download the the precompiled binary (`gobatmon`) from the [latest release](https://github.com/ulinja/gobatmon/releases/latest).
Alternatively, install Go and build it yourself.

Save the binary and configure your window manager to start it on launch, by putting the following into your startup script:
```bash
/path/to/gobatmon &
```

## Development

Build requirements:

- `go`

To build locally, run:

```bash
make build
```

This will create the `gobatmon` executable.

The clean up built files, run:

```bash
make clean
```

## Contributing

Please make sure to properly format your source code using `gofmt` before committing to `main`.

A [pre-commit script](/.pre-commit) is provided, you can activate it in your local repository with the following command:
```bash
ln -sr .pre-commit .git/hooks/pre-commit
```

## Roadmap

- add notification sounds
