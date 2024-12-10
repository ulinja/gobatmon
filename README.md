# Gobatmon

Gobatmon is a very simple battery level monitoring daemon for Linux systems.

It continuously keeps an eye on the current charge level of the battery in your laptop and triggers warning notifications in case of
low charge.

If charging with AC power or above 20% charge, you will see no notifications.
If running on battery power and below 20% / 10% charge, you will receive a low battery / very low battery notification
every 10 minutes / 5 minutes respectively. That's it.

**System requirements:**

- a Linux system
- a battery
- a running notification server (`swaync`, `dunst` etc.)
- the `notify-send` program on your path (should come preinstalled with your notification server)

If you are using a desktop environment like Gnome or KDE Plasma, you likely won't need this software.

## Installation

### Manual

Download the the precompiled binary (`gobatmon`) from the [latest release page](https://github.com/ulinja/gobatmon/releases/latest).
Alternatively, install Go and build it yourself.

Configure your window manager to start it on launch, or put the following into your startup script:
```bash
/path/to/gobatmon &
```

The daemon will run continuously in the background.

Gobatmon is designed to run on as few resources as possible to conserve CPU cycles and thus its power requirement.

### NixOS

> :construction: Will be added in the future.

## Configuration

> :construction: Will be added in the future.
>
> You can edit the source file and compile a custom binary to configure `gobatmon` yourself.

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
