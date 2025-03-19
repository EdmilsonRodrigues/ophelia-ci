from typing import Annotated

from fastapi import APIRouter, Body, Request, status
from fastapi.responses import HTMLResponse, RedirectResponse
from ophelia_ci_interface.config import GITIGNORE_OPTIONS
from ophelia_ci_interface.models.generals import Modal, ModalItem
from ophelia_ci_interface.models.health import HealthService
from ophelia_ci_interface.models.repository import (
    CreateRepositoryRequest,
    Repository,
    UpdateRepositoryRequest,
)
from ophelia_ci_interface.routers.dependencies import Template

router = APIRouter(prefix='/repositories', tags=['Repository'])

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


@router.get('/', response_class=HTMLResponse)
def repositories(request: Request, template: Template):
    return template.TemplateResponse(
        'repositories.html',
        {
            'request': request,
            'title': 'Repositories - Ophelia CI',
            'page_title': 'Your repositories',
            'modal': repositories_modal,
            'status': HealthService.get_status(),
            'repositories': Repository.get_all(metadata=metadata),
        },
    )


@router.post('/', response_class=HTMLResponse)
def create_repository(
    request: Request, body: CreateRepositoryRequest, template: Template
):
    Repository.create(
        body.name, body.description, body.gitignore, metadata=metadata
    )


@router.get('/{repo_name}', response_class=HTMLResponse)
def repository(request: Request, repo_name: str, template: Template):
    repository = Repository.get_by_name(repo_name, metadata=metadata)
    return template.TemplateResponse(
        'repository.html',
        {
            'request': request,
            'repo_name': repo_name,
            'status': HealthService.get_status(),
            'repository': repository,
            'id': repository.id,
            'modal': repository_modal,
        },
    )


@router.put('/{repo_name}', status_code=204)
def update_repository(
    request: Request, body: UpdateRepositoryRequest, template: Template
):
    Repository.update(body.id, body.name, body.description, metadata=metadata)


@router.delete('/{repo_name}', response_class=RedirectResponse)
def delete_repository(
    request: Request, id: Annotated[str, Body(embed=True)], template: Template
):
    Repository.delete(id, metadata=metadata)
    return RedirectResponse(
        url=request.url_for('repositories'),
        status_code=status.HTTP_303_SEE_OTHER,
    )
