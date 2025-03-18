import sys
from datetime import datetime
from functools import cache
from pathlib import Path

from fastapi import FastAPI, Request
from fastapi.responses import HTMLResponse
from fastapi.staticfiles import StaticFiles
from fastapi.templating import Jinja2Templates
from pydantic import UUID4, BaseModel

from ophelia_ci_interface.config import GITIGNORE_OPTIONS, VERSION, Settings
from ophelia_ci_interface.services.gRPC_service import RepositoryService

app = FastAPI(version=VERSION)

base_path = (
    Path(next(path for path in sys.path if 'app' in path))
    / 'ophelia_ci_interface'
)
app.mount(
    '/static', StaticFiles(directory=base_path / 'static'), name='static'
)

template = Jinja2Templates(directory=base_path / 'templates')


class ModalItem(BaseModel):
    id: str
    label: str
    type: str
    autocomplete: str = 'off'
    options: list[str] = []


class Modal(BaseModel):
    title: str
    items: list[ModalItem] = []
    submit: str


class Repository(BaseModel):
    model_config = {'arbitrary_types_allowed': True}

    id: UUID4
    name: str
    description: str
    last_updated: datetime

    @staticmethod
    @cache
    def get_service() -> RepositoryService:
        return RepositoryService(Settings().GRPC_SERVER)

    @property
    def truncated_description(self):
        if len(self.description) > 100:
            return self.description[:100] + '...'
        return self.description

    @property
    def clone_url(self):
        return f'git@{str(Settings().GRPC_SERVER).rstrip("/")}:{self.name}.git'

    @classmethod
    def get_status(cls):
        return cls.get_service().get_status()

    @classmethod
    def create(cls, name: str, description: str):
        response = cls.get_service().create_repository(name, description)
        return cls(
            id=response.id,
            name=response.name,
            description=response.description,
            last_updated=response.last_updated,
        )

    @classmethod
    def update(cls, id: str, name: str, description: str):
        response = cls.get_service().update_repository(id, name, description)
        return cls(
            id=response.id,
            name=response.name,
            description=response.description,
            last_updated=response.last_updated,
        )

    @classmethod
    def get(cls, id: str):
        response = cls.get_service().get_repository(id)
        return cls(
            id=response.id,
            name=response.name,
            description=response.description,
            last_updated=response.last_updated,
        )

    @classmethod
    def get_by_name(cls, name: str):
        response = cls.get_service().get_by_name(name)
        return cls(
            id=response.id,
            name=response.name,
            description=response.description,
            last_updated=response.last_updated,
        )

    @classmethod
    def delete(cls, id: str):
        cls.get_service().delete_repository(id)

    @classmethod
    def get_all(cls):
        response_list = cls.get_service().get_repositories()
        return [
            cls(
                id=response.id,
                name=response.name,
                description=response.description,
                last_updated=response.last_updated,
            )
            for response in response_list
        ]


@app.get('/')
async def root():
    return {'version': VERSION}


@app.get('/home', response_class=HTMLResponse)
async def home(request: Request):
    return template.TemplateResponse(
        'index.html',
        {
            'request': request,
            'title': 'Ophelia CI',
            'page_title': 'Welcome to the Ophelia CI',
            'status': Repository.get_status(),
        },
    )


@app.get('/repositories', response_class=HTMLResponse)
async def repositories(request: Request):
    return template.TemplateResponse(
        'repositories.html',
        {
            'request': request,
            'title': 'Repositories - Ophelia CI',
            'page_title': 'Your repositories',
            'modal': Modal(
                title='Create repository',
                items=[
                    ModalItem(
                        id='repository-name',
                        label='Repository name',
                        type='text',
                        autocomplete='off',
                    ),
                    ModalItem(
                        id='repository-description',
                        label='Repository description',
                        type='text',
                        autocomplete='off',
                    ),
                    ModalItem(
                        id='repository-gitignore',
                        label='Repository gitignore',
                        type='select',
                        options=GITIGNORE_OPTIONS,
                    ),
                ],
                submit='Add repository',
            ),
            'status': Repository.get_status(),
            'repositories': Repository.get_all(),
        },
    )


@app.get('/repositories/{repo_name}', response_class=HTMLResponse)
async def repository(request: Request, repo_name: str):
    return template.TemplateResponse(
        'repository.html',
        {
            'request': request,
            'repo_name': repo_name,
            'status': Repository.get_status(),
            'repository': Repository.get(repo_name),
        },
    )


@app.get('/users', response_class=HTMLResponse)
async def users(request: Request):
    return template.TemplateResponse(
        'users.html',
        {
            'request': request,
            'status': Repository.get_status(),
        },
    )


@app.get('/users/{user_name}', response_class=HTMLResponse)
async def user(request: Request, user_name: str):
    return template.TemplateResponse(
        'user.html',
        {
            'request': request,
            'user_name': user_name,
            'status': Repository.get_status(),
        },
    )
