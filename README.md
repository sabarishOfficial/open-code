# devopen

A small Go CLI that lets you pick a project folder and an editor via [fzf](https://github.com/junegunn/fzf), then opens the project in your chosen editor.

## Requirements

- **Go** 1.25+
- **fzf** (optional but recommended) — for the fuzzy picker UI. Without it, a simple numbered list fallback is used.
- One or more editors on your `PATH`: `cursor`, `code`, `vim`, `nvim`, `nano`

## Build

```bash
go build -o scriptname .
```

## Usage

```bash
./scriptname                    # Scan default dir (~/Desktop/GitHub)
./scriptname ~/projects         # Scan a specific directory
DEVOPEN_DIR=~/repos ./devopen   # Use env var for scan directory
```

1. The tool lists non-hidden folders in the scan directory.
2. You pick a **project** (fzf or by number).
3. You pick an **editor** (cursor, code, vim, nvim, nano).
4. The project is opened in that editor.

## Configuration

| Source | Description |
|--------|-------------|
| `~/Desktop/GitHub` | Default directory if nothing else is set |
| `DEVOPEN_DIR` | Environment variable; overrides default (supports `~`) |
| First argument | Command-line path; overrides env and default (supports `~`) |

## Editors

Supported editors (must be on your `PATH`):

- `cursor` — Cursor
- `code` — VS Code
- `vim` — Vim
- `nvim` — Neovim
- `nano` — Nano

