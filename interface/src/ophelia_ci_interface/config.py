from pathlib import Path

from pydantic_settings import BaseSettings, SettingsConfigDict

VERSION = '0.0.1'
GITIGNORE_OPTIONS = ['None', 'python', 'go']


class Settings(BaseSettings):
    model_config = SettingsConfigDict(
        env_file='.env', env_prefix='OPHELIA_CI_', extra='ignore'
    )
    PORT: int = 8000
    GRPC_SERVER: str = 'localhost:50051'
    DEBUG: bool = False
    SSL_KEYFILE: Path | None = None
    SSL_CERTFILE: Path | None = None
    WORKERS: int | None = None
