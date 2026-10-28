# Configurations & Scripts
**Path:** `configs/` and `scripts/`

## Purpose
The `configs` directory contains environment-specific configuration templates and settings for the application. The `scripts` directory houses utility and deployment scripts for different purposes.

## Contains
- **`configs/`**: Environment-specific YAML or JSON configuration files (e.g., `config.yaml`, `config.prod.yaml`).
- **`scripts/`**: Shell scripts for automation, data migration, and deployment.

## Principles & Rules
- **Environment Specificity**: Configuration files must be segregated by environment (e.g., development, testing, staging, production).
- **Secret Management**: Sensitive information (e.g., API keys, database credentials) should never be stored in plain text. Use tools like [SOPS](https://github.com/getsops/sops) for encrypted secrets.
- **Config Templates**: Maintain configuration templates with placeholders for environment-specific values.
- **Utility Scripts**: Scripts in `scripts/` should be well-documented and focused on specific, repeatable tasks.
