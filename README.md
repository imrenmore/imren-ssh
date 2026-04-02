# imren-ssh

An interactive SSH portfolio — connect from any terminal and explore projects, links, and contact info through a live TUI.

```
ssh ssh.imren.online
```

## What it is

A [Bubbletea](https://github.com/charmbracelet/bubbletea) TUI served over SSH via [Wish](https://github.com/charmbracelet/wish), deployed on Fly.io. Visitors connect anonymously — no auth required — and get a navigable, typewriter-animated interface with an ASCII portrait, project listings, and contact details.

## Features

- ASCII art portrait that scales to terminal width
- Typewriter animation on bio text
- Responsive layout — adapts to small, medium, and large terminals
- Keyboard navigation: `tab` / `h` / `l` to move between sections, `1`/`2`/`3` to jump, `q` to quit
- Rate limiting (10 connections/min per IP, max 20 concurrent sessions)
- Command execution blocked — shell access is not possible

## Stack

| Layer | Tool |
|---|---|
| TUI framework | [Bubbletea](https://github.com/charmbracelet/bubbletea) |
| SSH server | [Wish](https://github.com/charmbracelet/wish) |
| Styling | [Lipgloss](https://github.com/charmbracelet/lipgloss) |
| Hosting | [Fly.io](https://fly.io) |
| Language | Go 1.22 |

## Running locally

```bash
# Generate a host key
ssh-keygen -t ed25519 -f host_key -N ""

# Run
go run .
```

Then connect with:

```bash
ssh localhost -p 2222
```

## Deploying

```bash
flyctl deploy
```

The host key is persisted on a Fly volume mounted at `/data/host_key` so it survives redeploys.

## Security notes

- `host_key` / `host_key.pub` are gitignored — never committed
- All incoming exec commands are blocked at the middleware layer
- Log inputs are sanitized (non-printable chars stripped, truncated to 64 chars)
