# Gobatmon

Gobatmon is a very simple battery level monitoring daemon for Linux systems.

It keeps an eye on the current charge level of the battery in your laptop and triggers desktop notifications warning
you about low charge.

If charging with AC power or above 20% charge, you will see no notifications.
If running on battery power and below 20% / 10% charge, you will receive a low battery / very low battery notification
every 10 minutes / 5 minutes respectively. That's it.

**System requirements:**

- a Linux Desktop system
- a battery
- a running notification server (`swaync`, `dunst` etc.)

If you are using a desktop environment (Xfce/Gnome/Plasma etc.) you most likely don't need this software.

Gobatmon uses DBUS to dispatch desktop notifications.

## Usage

Simply run it by executing `gobatmon`.

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
    --poll-rate uint
        Poll rate for battery status in seconds (default 60)
    --version
        Show version information and exit
    --help
        Show help message and exit
```

## Installation

### Manual

Download the the precompiled binary (`gobatmon`) from the [latest release](https://github.com/ulinja/gobatmon/releases/latest).
Alternatively, install Go and build it yourself.

Save the binary and configure your window manager to start it on launch, by putting the following into your startup script:
```bash
/path/to/gobatmon &
```

The daemon will run continuously in the background.

Gobatmon is super-low on resources to conserve CPU cycles and thus its power requirement.

### NixOS

> :construction: Will be added in the future.

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
