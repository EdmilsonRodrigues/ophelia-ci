from typing import Annotated

from fastapi import APIRouter, Body, Request, status
from fastapi.responses import HTMLResponse, RedirectResponse
from ophelia_ci_interface.models.generals import Modal, ModalItem
from ophelia_ci_interface.models.health import HealthService
from ophelia_ci_interface.models.user import (
    CreateUserRequest,
    UpdateUserRequest,
    User,
)
from ophelia_ci_interface.routers.dependencies import Template

router = APIRouter(prefix='/users', tags=['User'])


users_modal = Modal(
    title='Create user',
    items=[
        ModalItem(
            id='user_username',
            label='Username',
            type='text',
            autocomplete='username',
        ),
        ModalItem(
            id='user_private_key',
            label='Private Key',
            type='textarea',
            autocomplete='off',
        ),
    ],
    submit='Add user',
    submit_id='user-create',
)

user_modal = Modal(
    title='Update user',
    items=[
        ModalItem(
            id='user_username',
            label='User name',
            type='text',
            autocomplete='username',
        ),
        ModalItem(
            id='user_private_key',
            label='Private Key',
            type='textarea',
            autocomplete='off',
        ),
    ],
    submit='Update user',
    submit_id='user-update',
)


@router.get('/', response_class=HTMLResponse)
def users(request: Request, template: Template):
    return template.TemplateResponse(
        'users.html',
        {
            'request': request,
            'title': 'Collaborators - Ophelia CI',
            'page_title': 'Collaborators',
            'modal': users_modal,
            'status': HealthService.get_status(),
            'repositories': User.get_all(metadata=metadata),
        },
    )


@router.post('/', response_class=HTMLResponse)
def create_user(request: Request, body: CreateUserRequest, template: Template):
    User.create(body.user_username, body.user_private_key, metadata=metadata)


@router.get('/{username}', response_class=HTMLResponse)
def repository(request: Request, username: str, template: Template):
    user = User.get_by_username(username, metadata=metadata)
    return template.TemplateResponse(
        'repository.html',
        {
            'request': request,
            'username': username,
            'status': HealthService.get_status(),
            'user': user,
            'id': user.id,
            'modal': user_modal,
        },
    )


@router.put('/{username}', status_code=204)
def update_repository(
    request: Request, body: UpdateUserRequest, template: Template
):
    User.update(
        body.id, body.user_username, body.user_private_key, metadata=metadata
    )


@router.delete('/{username}', response_class=RedirectResponse)
def delete_repository(
    request: Request, id: Annotated[str, Body(embed=True)], template: Template
):
    User.delete(id, metadata=metadata)
    return RedirectResponse(
        url=request.url_for('repositories'),
        status_code=status.HTTP_303_SEE_OTHER,
    )
