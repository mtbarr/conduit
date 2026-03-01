# Conduit Discord Bot

A Discord bot for bug reports, feature requests, and listing GitHub issues. Supports i18n (English & Portuguese).

## Quick Start

1. Clone the repository.
2. Create a `.env` file with:
   ```env
   DISCORD_TOKEN=your-discord-token
   GITHUB_TOKEN=your-github-token
   GITHUB_OWNER=your-github-username
   GITHUB_REPO=your-repo-name
   GUILD_ID=your-discord-guild-id
   APP_LANG=pt-BR   # or en
   ```
3. Build and run:
   ```powershell
   go build -o conduit .
   .\conduit.exe
   ```

## Usage

- `/reportarbug` — Report a bug (modal)
- `/requestfeature` — Request a feature (modal)
- `/issues` — List 10 latest open GitHub issues

## Environment Variables

- `DISCORD_TOKEN` — Discord bot token (required)
- `GITHUB_TOKEN` — GitHub token (required)
- `GITHUB_OWNER` — GitHub username (required)
- `GITHUB_REPO` — GitHub repo name (required)
- `GUILD_ID` — Discord server ID (for instant command updates)
- `APP_LANG` — Language (`pt-BR` or `en`, default: pt-BR)
- `REPORTBUG_COOLDOWN_SECONDS` — Bug report cooldown (default: 60)
- `REQUESTFEATURE_COOLDOWN_SECONDS` — Feature request cooldown (default: 60)
- `EPHEMERAL_DELETE_SECONDS` — Ephemeral message delete time (default: 3)
- `I18N_PATH` — Path to locale files (default: `locales`)

## Docker

Build:
```powershell
docker build -t conduit .
```
Run:
```powershell
docker run --rm --env-file .env conduit
```

## i18n (Languages)

- Default: Portuguese (pt-BR)
- English supported
- To add more: create `locales/xx.json` and set `APP_LANG=xx`

## Testing

Run tests locally (host only):
```powershell
go test ./...
```

### Cross-platform builds with Zig (`build.zig`)

Requires Zig installed and available in your PATH.

- Build for all targets:
  ```bash
  zig build
  ```
- Filter by OS:
  ```bash
  zig build -Dgoos=linux
  zig build -Dgoos=windows
  ```
- Filter by architecture:
  ```bash
  zig build -Dgoarch=amd64
  zig build -Dgoarch=arm64
  ```
- Filter by OS and architecture:
  ```bash
  zig build -Dgoos=linux -Dgoarch=arm64
  zig build -Dgoos=windows -Dgoarch=amd64
  ```
- Build only for a specific target:
  ```bash
  zig build build-linux-amd64
  zig build build-linux-arm64
  zig build build-windows-amd64
  zig build build-windows-arm64
  ```
- Run tests (host target only):
  ```bash
  zig build test-linux-amd64
  zig build test-windows-amd64
  ```

Binaries are placed in `dist/<os>-<arch>/conduit` (or `conduit.exe` on Windows).

## Notes
- Commands update instantly if `GUILD_ID` is set.
- All strings (modals, messages) are translatable.
- Cooldowns and ephemeral delete time are configurable.

---

Ready to use.
