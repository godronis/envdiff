# envdiff

Compare `.env` files across environments and highlight missing or mismatched keys with structured output.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git
cd envdiff
go build -o envdiff .
```

---

## Usage

```bash
envdiff --base .env --compare .env.production
```

**Example output:**

```
MISSING IN .env.production:
  - DATABASE_URL
  - REDIS_HOST

MISMATCHED VALUES:
  - APP_ENV: "development" → "production"
  - LOG_LEVEL: "debug" → "info"

OK: 12 keys match across both files.
```

### Flags

| Flag        | Description                          | Default  |
|-------------|--------------------------------------|----------|
| `--base`    | Path to the base `.env` file         | `.env`   |
| `--compare` | Path to the environment file to diff | required |
| `--format`  | Output format: `text`, `json`        | `text`   |
| `--strict`  | Exit with non-zero code on any diff  | `false`  |

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## License

[MIT](LICENSE)