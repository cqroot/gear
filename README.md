# gear

Lightweight utilities around the Go toolchain.

[![test](https://github.com/cqroot/gear/actions/workflows/test.yml/badge.svg)](https://github.com/cqroot/gear/actions/workflows/test.yml)

## Requirements

- Go 1.27.1 or newer

## Installation

Install the latest release:

```sh
go install github.com/cqroot/gear@latest
```

Or build from source:

```sh
git clone https://github.com/cqroot/gear
cd gear
go build -o gear .
```

## Usage

Run `gear` with no arguments to list the available commands:

```sh
gear
```

### remod

Rebuild `go.mod` and `go.sum` in place:

```sh
gear remod
```

`remod` reconstructs the module files for the target directory:

1. Resolve the target directory (`--dir`/`-C`, default the working directory).
2. Read the existing `go.mod` to capture the declared module path.
3. Delete the current `go.mod` and `go.sum` files.
4. Re-run `go mod init <module>` followed by `go mod tidy`.

The module path is preserved so existing import paths keep working. To run
against another directory, pass `-C`:

```sh
gear remod -C path/to/module
```

> **Note:** `remod` deletes and regenerates `go.mod` / `go.sum`. If a step
> fails, the original files are restored automatically.

## License

This project is licensed under the GNU General Public License v3.0 — see
[`LICENSE`](LICENSE) for the full text.
