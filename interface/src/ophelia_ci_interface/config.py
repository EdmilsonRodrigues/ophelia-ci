from pydantic_settings import BaseSettings, SettingsConfigDict

VERSION = '0.0.1'
GITIGNORE_OPTIONS = ['None', 'python', 'go']


class Settings(BaseSettings):
    model_config = SettingsConfigDict(
        env_file='.env', env_prefix='OPHELIA_CI_', extra='ignore'
    )
    GRPC_SERVER: str = 'localhost:50051'
    DEBUG: bool = False
