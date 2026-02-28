# Conduit

A Discord bot that collects bug reports via a modal and creates GitHub issues.

## Build

```powershell
 go build -o conduit .
```

## Run

```powershell
 $env:DISCORD_TOKEN="your-token"
 $env:GITHUB_TOKEN="your-github-token"
 $env:GITHUB_OWNER="your-github-owner"
 $env:GITHUB_REPO="your-github-repo"
 $env:REPORTBUG_COOLDOWN_SECONDS="60"
 $env:EPHEMERAL_DELETE_SECONDS="3"
 $env:APP_LANG="en"
 $env:I18N_PATH="C:\path\to\locales"
 .\conduit.exe
```

## Test

```powershell
 go test ./...
```

## Docker

Build the image:

```powershell
 docker build -t conduit .
```

Run with an env file:

```powershell
 docker run --rm --env-file .env conduit
```

Run with explicit variables:

```powershell
 docker run --rm \
  -e DISCORD_TOKEN="your-token" \
  -e GITHUB_TOKEN="your-github-token" \
  -e GITHUB_OWNER="your-github-owner" \
  -e GITHUB_REPO="your-github-repo" \
  -e REPORTBUG_COOLDOWN_SECONDS="60" \
  -e EPHEMERAL_DELETE_SECONDS="3" \
  -e APP_LANG="en" \
  -e I18N_PATH="/config/locales" \
  -v ${PWD}/locales:/config/locales:ro \
  conduit
```

## Configuration

Required runtime environment variables:

- `DISCORD_TOKEN`
- `GITHUB_TOKEN`
- `GITHUB_OWNER`
- `GITHUB_REPO`

Optional:

- `REPORTBUG_COOLDOWN_SECONDS` (default: 60)
- `EPHEMERAL_DELETE_SECONDS` (default: 3, set to `0` to disable)
- `APP_LANG` (default: `en`)
- `I18N_PATH` (path to a directory with locale files, default: `locales`)

## i18n

The app uses `go-i18n` and loads locale files from `I18N_PATH` (defaults to `locales`). Each locale file is JSON and is named like `en.json`, `es.json`.

Example `locales/en.json`:

```json
[
  {"id": "command_name", "translation": "reportbug"},
  {"id": "modal_title", "translation": "Report a Bug"},
  {"id": "modal_title_label", "translation": "Title"},
  {"id": "modal_title_placeholder", "translation": "Short summary of the bug"},
  {"id": "modal_desc_label", "translation": "Description"},
  {"id": "modal_desc_placeholder", "translation": "Detailed description of the bug"},
  {"id": "cooldown_message", "translation": "Please wait %s before submitting another bug report."},
  {"id": "issue_failed", "translation": "Failed to create GitHub issue. Please try again later."},
  {"id": "issue_created_simple", "translation": "Bug report submitted successfully."}
]
```

To add a language:

1. Copy `locales/en.json` to `locales/<lang>.json`.
2. Translate the `translation` values.
3. Set `APP_LANG` to that language code (example: `es`).
