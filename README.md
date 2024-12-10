# Gobatmon

Gobatmon is a simple battery monitoring utility for Linux.
It is used to keep an eye on the current charge level of the battery in your laptop and trigger a warning notification in case of
low charge.

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
