import sys
from pathlib import Path
from typing import Annotated

from fastapi import Body, FastAPI, Request, status
from fastapi.responses import HTMLResponse, RedirectResponse
from fastapi.staticfiles import StaticFiles
from fastapi.templating import Jinja2Templates
from pydantic import BaseModel

from ophelia_ci_interface.config import GITIGNORE_OPTIONS, VERSION, Settings
from ophelia_ci_interface.models.repository import (
    CreateRepositoryRequest,
    Repository,
    UpdateRepositoryRequest,
)

app = FastAPI(version=VERSION)
base_path = (
    Path('ophelia_ci_interface')
    if Settings().DEBUG
    else Path(next(path for path in sys.path if 'app' in path))
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
    submit_id: str


repositories_modal = Modal(
    title='Create repository',
    items=[
        ModalItem(
            id='repository_name',
            label='Repository name',
            type='text',
            autocomplete='off',
        ),
        ModalItem(
            id='repository_description',
            label='Repository description',
            type='text',
            autocomplete='off',
        ),
        ModalItem(
            id='repository_gitignore',
            label='Repository gitignore',
            type='select',
            options=GITIGNORE_OPTIONS,
        ),
    ],
    submit='Add repository',
    submit_id='repository-create',
)

repository_modal = Modal(
    title='Update repository',
    items=[
        ModalItem(
            id='repository_name',
            label='Repository name',
            type='text',
            autocomplete='off',
        ),
        ModalItem(
            id='repository_description',
            label='Repository description',
            type='text',
            autocomplete='off',
        ),
    ],
    submit='Update repository',
    submit_id='repository-update',
)


@app.get('/health')
def root():
    return {'version': VERSION}


@app.get('/', response_class=HTMLResponse)
def home(request: Request):
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
def repositories(request: Request):
    return template.TemplateResponse(
        'repositories.html',
        {
            'request': request,
            'title': 'Repositories - Ophelia CI',
            'page_title': 'Your repositories',
            'modal': repositories_modal,
            'status': Repository.get_status(),
            'repositories': Repository.get_all(),
        },
    )


@app.post('/repositories', response_class=HTMLResponse)
def create_repository(request: Request, body: CreateRepositoryRequest):
    Repository.create(body.name, body.description, body.gitignore)

    return repositories(request)


@app.get('/repositories/{repo_name}', response_class=HTMLResponse)
def repository(request: Request, repo_name: str):
    repository = Repository.get_by_name(repo_name)
    return template.TemplateResponse(
        'repository.html',
        {
            'request': request,
            'repo_name': repo_name,
            'status': Repository.get_status(),
            'repository': repository,
            'id': repository.id,
            'modal': repository_modal,
        },
    )


@app.put('/repositories/{repo_name}', status_code=204)
def update_repository(request: Request, body: UpdateRepositoryRequest):
    Repository.update(body.id, body.name, body.description)


@app.delete('/repositories/{repo_name}', response_class=RedirectResponse)
def delete_repository(request: Request, id: Annotated[str, Body(embed=True)]):
    Repository.delete(id)
    return RedirectResponse(
        url=request.url_for('repositories'),
        status_code=status.HTTP_303_SEE_OTHER,
    )


@app.get('/users', response_class=HTMLResponse)
def users(request: Request):
    return template.TemplateResponse(
        'users.html',
        {
            'request': request,
            'status': Repository.get_status(),
        },
    )


@app.get('/users/{user_name}', response_class=HTMLResponse)
def user(request: Request, user_name: str):
    return template.TemplateResponse(
        'user.html',
        {
            'request': request,
            'user_name': user_name,
            'status': Repository.get_status(),
        },
    )
