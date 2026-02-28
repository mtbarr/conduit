# Conduit Discord Bot

A simple Discord bot for bug reports, feature requests, and listing GitHub issues. Supports i18n (English & Portuguese).

## Quick Start

1. **Clone the repo**
2. **Create a `.env` file** with:
   ```env
   DISCORD_TOKEN=your-discord-token
   GITHUB_TOKEN=your-github-token
   GITHUB_OWNER=your-github-username
   GITHUB_REPO=your-repo-name
   GUILD_ID=your-discord-guild-id
   APP_LANG=pt-BR   # or en
   ```
3. **Build and run:**
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
- Add more: create `locales/xx.json` and set `APP_LANG=xx`

## Testing

```powershell
go test ./...
```

## Notes
- Commands update instantly if `GUILD_ID` is set.
- All strings (modals, messages) are translatable.
- Cooldowns and ephemeral delete time are configurable.

---

**Ready to use!**
