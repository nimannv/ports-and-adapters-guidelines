import os
from pathlib import Path
from typing import Optional
from pydantic import BaseModel
import yaml


class Config(BaseModel):
    db_driver: str = "sqlite"
    database_url: str = "sqlite:///tasks.db"
    api_host: str = "0.0.0.0"
    api_port: int = 8000
    log_level: str = "INFO"


def load_config(config_path: Optional[str] = None) -> Config:
    data = {}
    target_path = None

    if config_path:
        target_path = Path(config_path)
    else:
        # Default fallback to configs/base.yaml
        default_file = Path(__file__).parent / "base.yaml"
        if default_file.exists():
            target_path = default_file

    if target_path and target_path.exists():
        with open(target_path, "r", encoding="utf-8") as f:
            loaded = yaml.safe_load(f)
            if isinstance(loaded, dict):
                data.update(loaded)

    # Environment variables override config file values
    if "DB_DRIVER" in os.environ:
        data["db_driver"] = os.environ["DB_DRIVER"]
    if "DATABASE_URL" in os.environ:
        data["database_url"] = os.environ["DATABASE_URL"]
    if "API_HOST" in os.environ:
        data["api_host"] = os.environ["API_HOST"]
    if "API_PORT" in os.environ:
        try:
            data["api_port"] = int(os.environ["API_PORT"])
        except ValueError:
            pass
    if "LOG_LEVEL" in os.environ:
        data["log_level"] = os.environ["LOG_LEVEL"]

    return Config(**data)

