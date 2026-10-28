# Configurations & Scripts
**Path:** `configs/` and `scripts/`

## Configs

### Purpose
The `configs` directory contains configuration templates and environment-specific settings for the application.

### Contains
- **Configuration files**: YAML (or JSON/TOML) files with shared defaults and per-environment overrides.
- **Templates**: Files with placeholders for environment-specific values.

### Principles & Rules
- **Defaults + Overrides**: Defaults live in one base file; each environment overrides only what it needs. Environment variables override files.
- **Loaded Only in Bootstrap**: Configuration is read in [bootstrap](bootstrap.md) and passed as plain values to adapters. The core never reads configuration files or environment variables.
- **Secret Management**: Sensitive information (e.g., API keys, database credentials) should never be stored in plain text. Use tools like [SOPS](https://github.com/getsops/sops) for encrypted secrets.

### Scaling Configs
Start with one config file. As the number of environments and components grows, split the files by **environment** and then by **component**.

```
configs/                    configs/
└── config.yaml      ──►    ├── base.yaml             # shared defaults
                            ├── local/
                            │   ├── app.yaml
                            │   └── database.yaml
                            ├── staging/
                            │   └── app.yaml
                            └── production/
                                ├── app.yaml
                                └── database.yaml
```

## Scripts

### Purpose
The `scripts` directory houses utility scripts for building, testing, deploying, and local development.

### Contains
- **Shell scripts**: Small, focused scripts for repeatable tasks (build, lint, test, deploy, seed, ...).

### Principles & Rules
- **Single Entry Point**: Expose scripts through a `Makefile` / `Taskfile` so people don't need to know the internal layout.
- **Focused and Documented**: Each script does one repeatable task and explains its usage.

### Scaling Scripts
Start with a few scripts at the top level. As they pile up, group them by **purpose**.

```
scripts/                    scripts/
├── build.sh         ──►    ├── build/
├── test.sh                 │   ├── docker.sh
└── deploy.sh               │   └── package.sh
                            ├── ci/
                            │   ├── lint.sh
                            │   └── test.sh
                            ├── deploy/
                            │   ├── staging.sh
                            │   └── production.sh
                            └── dev/
                                ├── seed.sh
                                └── setup.sh
```

See [Scale](../README.md#scale) for the general principles.
