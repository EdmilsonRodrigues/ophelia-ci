import logging
from pathlib import Path
from typing import Literal

from pydantic_settings import BaseSettings, SettingsConfigDict

VERSION = '0.7.0'
GITIGNORE_OPTIONS = ['None', 'python', 'go']


class Settings(BaseSettings):
    model_config = SettingsConfigDict(
        env_file='.env', env_prefix='OPHELIA_CI_', extra='ignore'
    )
    PORT: int = 8000
    GRPC_SERVER: str = 'localhost:50051'
    DEBUG: bool = False
    LOG_LEVEL: Literal['DEBUG', 'INFO', 'WARNING', 'ERROR'] = 'INFO'
    SSL_KEYFILE: Path | None = None
    SSL_CERTFILE: Path | None = None
    WORKERS: int | None = None


base_path = (
    Path('ophelia_ci_interface')
    if Settings().DEBUG
    else Path('/usr/lib/ophelia-ci-interface/app/ophelia_ci_interface')
)

logging.basicConfig(
    level=logging.getLevelName(Settings().LOG_LEVEL),
    format='%(asctime)s - %(levelname)s - %(message)s',
)
