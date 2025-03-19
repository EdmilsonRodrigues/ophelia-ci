from functools import cache
from typing import Annotated

from fastapi import Depends
from fastapi.templating import Jinja2Templates
from ophelia_ci_interface.config import Settings, base_path
from ophelia_ci_interface.services.gRPC_service import AuthenticationService


@cache
def authentiction_service() -> AuthenticationService:
    return AuthenticationService(str(Settings().GRPC_SERVER))


@cache
def template_dependency() -> Jinja2Templates:
    return Jinja2Templates(directory=base_path / 'templates')


Template = Annotated[Jinja2Templates, Depends(template_dependency)]
Authentication = Annotated[
    AuthenticationService, Depends(authentiction_service)
]
