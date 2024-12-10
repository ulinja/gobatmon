# Gobatmon

Gobatmon is a simple battery level monitoring daemon for Linux systems.
It continuously keeps an eye on the current charge level of the battery in your laptop and triggers warning notifications in case of
low charge.

System requirements:

- a Linux system
- a battery
- `notify-send`

## Development

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
