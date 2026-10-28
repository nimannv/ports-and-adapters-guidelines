# Configurations & Scripts
**Path:** `configs/` and `scripts/`

## Purpose
The `configs` directory contains configuration templates and environment-specific settings for the application. The `scripts` directory houses utility and deployment scripts for different purposes.

## Contains
- **`configs/`**: YAML (or JSON/TOML) configuration files. Start with a single `config.yaml`. As environments and components grow, move to a `base.yaml` with shared defaults plus one directory per environment (`local/`, `staging/`, `production/`), split by component (`app.yaml`, `database.yaml`). See [configs/ scaling](../README.md#configs).
- **`scripts/`**: Shell scripts for building, testing, deployment, and local development. Start flat (`build.sh`, `test.sh`, `deploy.sh`); group them by purpose (`build/`, `ci/`, `deploy/`, `dev/`) as they pile up. See [scripts/ scaling](../README.md#scripts).

## Principles & Rules
- **Defaults + Overrides**: Defaults live in one base file; each environment overrides only what it needs. Environment variables override files.
- **Loaded Only in Bootstrap**: Configuration is read in [bootstrap](bootstrap.md) and passed as plain values to adapters. The core never reads configuration files or environment variables.
- **Secret Management**: Sensitive information (e.g., API keys, database credentials) should never be stored in plain text. Use tools like [SOPS](https://github.com/getsops/sops) for encrypted secrets.
- **Config Templates**: Maintain configuration templates with placeholders for environment-specific values.
- **Single Entry Point for Scripts**: Expose scripts through a `Makefile` / `Taskfile` so people don't need to know the internal layout.
- **Utility Scripts**: Scripts in `scripts/` should be well-documented and focused on specific, repeatable tasks.
