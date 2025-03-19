import sys
from pathlib import Path

from pydantic_settings import BaseSettings, SettingsConfigDict

VERSION = '0.0.1'
GITIGNORE_OPTIONS = ['None', 'python', 'go']


class Settings(BaseSettings):
    model_config = SettingsConfigDict(
        env_file='.env', env_prefix='OPHELIA_CI_', extra='ignore'
    )
    GRPC_SERVER: str = 'localhost:50051'
    DEBUG: bool = False


base_path = (
    Path('ophelia_ci_interface')
    if Settings().DEBUG
    else Path(next(path for path in sys.path if 'app' in path))
    / 'ophelia_ci_interface'
)
