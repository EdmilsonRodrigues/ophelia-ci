from typing import Annotated

from fastapi import APIRouter, Body, Form, Path, Request, status
from fastapi.responses import HTMLResponse, RedirectResponse
from ophelia_ci_interface.config import GITIGNORE_OPTIONS
from ophelia_ci_interface.models.generals import Modal, ModalItem
from ophelia_ci_interface.models.repository import (
    CreateRepositoryRequest,
    Repository,
    UpdateRepositoryRequest,
)
from ophelia_ci_interface.routers.dependencies import (
    Health,
    Metadata,
    RepositoryDependency,
    SettingsDependency,
    Template,
)
from pydantic import UUID4

router = APIRouter(prefix='/repositories', tags=['Repository'])

repositories_modal = Modal(
    title='Create repository',
    action='/repositories/',
    method='POST',
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
)

repository_modal = Modal(
    title='Update repository',
    action='/repositories/{repo_name}',
    method='PUT',
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
)


@router.get('/', response_class=HTMLResponse)
def repositories(
    request: Request,
    template: Template,
    health_service: Health,
    metadata: Metadata,
    repository_service: RepositoryDependency,
):
    return template.TemplateResponse(
        'repositories.html',
        {
            'request': request,
            'title': 'Repositories - Ophelia CI',
            'page_title': 'Your repositories',
            'modal': repositories_modal,
            'status': health_service.get_status(),
            'repositories': Repository.get_all(
                repository_service, metadata=metadata
            ),
        },
    )


@router.post('/', response_class=HTMLResponse)
def create_repository(
    repository_service: RepositoryDependency,
    req: Annotated[
        CreateRepositoryRequest,
        Form(title='Repository data', description='The repository data'),
    ],
    template: Template,
    health_service: Health,
    metadata: Metadata,
):
    """
    Create a new repository in the database.

    :param req: the request data

    :return: the HTML response
    """
    Repository.create(
        repository_service,
        req.repository_name,
        req.repository_description,
        req.repository_gitignore,
        metadata=metadata,
    )
    return RedirectResponse(
        url='/repositories/', status_code=status.HTTP_201_CREATED
    )


@router.get('/{repo_name}', response_class=HTMLResponse)
def repository(
    request: Request,
    settings: SettingsDependency,
    repository_service: RepositoryDependency,
    template: Template,
    health_service: Health,
    metadata: Metadata,
    repo_name: Annotated[
        str,
        Path(
            title='Repository Name', description='The name of the repository'
        ),
    ],
):
    """
    Get a repository by its name.

    :param repo_name: the name of the repository

    :return: the HTML response
    """
    repository = Repository.get_by_name(
        settings, repository_service, repo_name, metadata=metadata
    )
    return template.TemplateResponse(
        'repository.html',
        {
            'request': request,
            'repo_name': repo_name,
            'status': health_service.get_status(),
            'repository': repository,
            'id': repository.id,
            'modal': repository_modal.format_action(repo_name=repo_name),
        },
    )


@router.put('/{repo_name}', status_code=204)
def update_repository(
    repository_service: RepositoryDependency,
    req: Annotated[
        UpdateRepositoryRequest,
        Form(title='Repository data', description='The repository data'),
    ],
    template: Template,
    health_service: Health,
    metadata: Metadata,
) -> None:
    """
    Update an existing repository in the database.

    :param req: the request data containing the repository update details

    :return: None
    """
    Repository.update(
        repository_service,
        str(req.id),
        req.repository_name,
        req.repository_description,
        metadata=metadata,
    )


@router.delete('/{repo_name}', response_class=RedirectResponse)
def delete_repository(
    request: Request,
    repository_service: RepositoryDependency,
    id: Annotated[
        UUID4,
        Body(
            title='Repository ID', description='The repository ID', embed=True
        ),
    ],
    template: Template,
    health_service: Health,
    metadata: Metadata,
):
    """
    Delete an existing repository from the database.

    :param id: the ID of the repository

    :return: a RedirectResponse object that redirects to the repositories page
    """
    Repository.delete(repository_service, str(id), metadata=metadata)
    return RedirectResponse(
        url=request.url_for('repositories'),
        status_code=status.HTTP_303_SEE_OTHER,
    )
