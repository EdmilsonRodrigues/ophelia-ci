from functools import cache
from pathlib import Path
from typing import Annotated

from fastapi import Cookie, Depends, HTTPException
from fastapi.templating import Jinja2Templates
from ophelia_ci_interface.config import Settings
from ophelia_ci_interface.models.generals import MetadataTuple
from ophelia_ci_interface.services.gRPC_service import (
    AuthenticationService,
    HealthService,
    RepositoryService,
    UserService,
)


@cache
def settings_dependency() -> Settings:
    return Settings()


SettingsDependency = Annotated[Settings, Depends(settings_dependency)]


def authentication_service(
    settings: SettingsDependency,
) -> AuthenticationService:
    return AuthenticationService(settings.GRPC_SERVER)


def health_service(settings: SettingsDependency) -> HealthService:
    return HealthService(settings.GRPC_SERVER)


def repository_service(settings: SettingsDependency) -> RepositoryService:
    return RepositoryService(settings.GRPC_SERVER)


@cache
def template_dependency() -> Jinja2Templates:
    return Jinja2Templates(directory=Path('resources', 'templates'))


def user_service(settings: SettingsDependency) -> UserService:
    return UserService(settings.GRPC_SERVER)


def get_metadata(
    session: Annotated[
        str | None,
        Cookie(title='Session', description='The token of the user'),
    ] = None,
) -> MetadataTuple:
    if session is None:
        raise HTTPException(status_code=401, detail='Unauthorized')
    return (('authorization', f'Bearer {session}'),)


Authentication = Annotated[
    AuthenticationService, Depends(authentication_service)
]
Health = Annotated[HealthService, Depends(health_service)]
RepositoryDependency = Annotated[
    RepositoryService, Depends(repository_service)
]
Template = Annotated[Jinja2Templates, Depends(template_dependency)]
UserDependency = Annotated[UserService, Depends(user_service)]
Metadata = Annotated[MetadataTuple, Depends(get_metadata)]
