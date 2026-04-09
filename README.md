# dataforseo-cli

DataForSEO CLI — SEO data from the command line.

## Installation

### Download Binary

Download the latest release from [GitHub Releases](https://github.com/Robben-Media/dataforseo-cli/releases).

### Build from Source

```bash
git clone https://github.com/Robben-Media/dataforseo-cli.git
cd dataforseo-cli
go build ./cmd/dataforseo
```

## Configuration

dataforseo-cli requires your DataForSEO API credentials (base64-encoded `login:password`). Credentials are stored securely in your system keyring.

**Store credentials:**

```bash
dataforseo-cli auth set-key
```

The CLI will prompt for your base64-encoded auth string interactively. You can also pipe it:

```bash
echo "base64-encoded-login:password" | dataforseo-cli auth set-key
```

**Environment variable override:**

```bash
export DATAFORSEO_AUTH="base64-encoded-login:password"
```

**Check status:**

```bash
dataforseo-cli auth status
```

**Remove credentials:**

```bash
dataforseo-cli auth remove
```

## Commands

### auth

Manage API credentials.

| Command | Description |
|---------|-------------|
| `auth set-key` | Store base64-encoded auth credentials in keyring |
| `auth status` | Show authentication status |
| `auth remove` | Remove stored credentials |

### keywords

| Command | Description |
|---------|-------------|
| `keywords search-volume --keywords <kw>` | Get search volume for keywords |
| `keywords difficulty --keywords <kw>` | Get keyword difficulty scores |

**Flags:** `--keywords` (required, comma-separated), `--location` (default 2840), `--language` (default en)

### labs

DataForSEO Labs — keyword research tools.

| Command | Description |
|---------|-------------|
| `labs related --keyword <kw>` | Get related keywords |
| `labs suggestions --keyword <kw>` | Get keyword suggestions |

**Flags:** `--keyword` (required), `--location` (default 2840), `--language` (default en), `--limit` (default 10)

### serp

SERP analysis.

| Command | Description |
|---------|-------------|
| `serp google --keyword <kw>` | Get Google organic SERP results |
| `serp paa --keyword <kw>` | Get People Also Ask results |

**Flags (google):** `--keyword` (required), `--location` (default 2840), `--language` (default en), `--device` (desktop or mobile, default desktop)

**Flags (paa):** `--keyword` (required), `--location` (default 2840), `--language` (default en)

## Global Flags

| Flag | Description |
|------|-------------|
| `--json` | Output JSON to stdout (best for scripting) |
| `--plain` | Output stable, parseable text to stdout (TSV; no colors) |
| `--verbose` | Enable verbose logging |
| `--force` | Skip confirmations for destructive commands |
| `--no-input` | Never prompt; fail instead (useful for CI) |
| `--color` | Color output: auto, always, or never (default auto) |

## License

MIT
