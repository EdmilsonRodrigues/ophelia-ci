from typing import Annotated

from fastapi import (
    APIRouter,
    Body,
    File,
    Form,
    Path,
    Request,
    UploadFile,
    status,
)
from fastapi.responses import HTMLResponse, RedirectResponse
from ophelia_ci_interface.models.generals import Modal, ModalItem
from ophelia_ci_interface.models.user import (
    User,
)
from ophelia_ci_interface.routers.dependencies import (
    Health,
    Metadata,
    Template,
    UserDependency,
)
from pydantic import UUID4

router = APIRouter(prefix='/users', tags=['User'])


users_modal = Modal(
    title='Create user',
    action='/users/',
    method='POST',
    items=[
        ModalItem(
            id='user_username',
            label='Username',
            type='text',
            autocomplete='username',
        ),
        ModalItem(
            id='user_public_key',
            label='Public Key',
            type='file',
            autocomplete='off',
        ),
    ],
    submit='Add user',
)

user_modal = Modal(
    title='Update user',
    action='/users/{username}/',
    method='PUT',
    items=[
        ModalItem(
            id='user_username',
            label='User name',
            type='text',
            autocomplete='username',
        ),
        ModalItem(
            id='user_public_key',
            label='Public Key',
            type='file',
            autocomplete='off',
        ),
    ],
    submit='Update user',
)


@router.get('/', response_class=HTMLResponse)
def users(
    request: Request,
    user_service: UserDependency,
    template: Template,
    health_service: Health,
    metadata: Metadata,
):
    """
    Show all existing users in the database.

    :return: the HTML response
    """
    return template.TemplateResponse(
        'users.html',
        {
            'request': request,
            'title': 'Collaborators - Ophelia CI',
            'page_title': 'Collaborators',
            'modal': users_modal,
            'status': health_service.get_status(),
            'users': User.get_all(user_service, metadata=metadata),
        },
    )


@router.post('/', status_code=204)
async def create_user(
    user_service: UserDependency,
    user_username: Annotated[
        str, Form(title='Username', description='The username of the user.')
    ],
    user_public_key: Annotated[
        UploadFile,
        File(title='Private Key', description='The public key of the user.'),
    ],
    template: Template,
    metadata: Metadata,
):
    """
    Create a new user in the database with the provided username and
    public key.

    :param user_username: The username of the user to be created.
    :param user_public_key: The file containing the public key of the user.

    :return: The response to the request.
    """
    User.create(
        user_service,
        user_username,
        (await user_public_key.read()).decode('utf-8'),
        metadata=metadata,
    )
    return RedirectResponse(url='/users/', status_code=status.HTTP_201_CREATED)


@router.get('/{username}', response_class=HTMLResponse)
def repository(
    request: Request,
    user_service: UserDependency,
    template: Template,
    health_service: Health,
    metadata: Metadata,
    username: Annotated[
        str, Path(title='Username', description='The username of the user')
    ],
):
    """
    Show a user in the database.

    :param username: The username of the user.

    :return: The HTML response.
    """
    user = User.get_by_username(user_service, username, metadata=metadata)
    return template.TemplateResponse(
        'user.html',
        {
            'request': request,
            'username': username,
            'status': health_service.get_status(),
            'user': user,
            'id': user.id,
            'modal': user_modal.format_action(username=username),
        },
    )


@router.put('/{username}', status_code=204)
async def update_repository(
    user_service: UserDependency,
    id: Annotated[UUID4, Form(title='ID', description='The ID of the user')],
    user_username: Annotated[
        str, Form(title='Username', description='The username of the user')
    ],
    user_public_key: Annotated[
        UploadFile,
        File(title='Private Key', description='The public key of the user.'),
    ],
    template: Template,
    metadata: Metadata,
):
    """
    Update an existing user in the database.

    :param body: the request containing the user ID, username and public key

    :return: None
    """
    User.update(
        user_service,
        str(id),
        user_username,
        (await user_public_key.read()).decode('utf-8'),
        metadata=metadata,
    )


@router.delete('/{username}', response_class=RedirectResponse)
def delete_repository(
    request: Request,
    user_service: UserDependency,
    id: Annotated[UUID4, Body(embed=True)],
    template: Template,
    metadata: Metadata,
):
    """
    Delete an existing user from the database and redirect to the users
    page.

    :param id: The ID of the user to be deleted.

    :return: A RedirectResponse object that redirects to the users page.
    """
    User.delete(user_service, str(id), metadata=metadata)
    return RedirectResponse(
        url=request.url_for('users'),
        status_code=status.HTTP_303_SEE_OTHER,
    )
