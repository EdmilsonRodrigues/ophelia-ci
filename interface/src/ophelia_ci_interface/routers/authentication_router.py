from typing import Annotated

from fastapi import (
    APIRouter,
    File,
    Form,
    HTTPException,
    Request,
    UploadFile,
    status,
)
from fastapi.responses import HTMLResponse, RedirectResponse
from ophelia_ci_interface.models.generals import OpheliaException
from ophelia_ci_interface.routers.dependencies import Authentication, Template

router = APIRouter(tags=['Authentication'])


@router.get('/login', response_class=HTMLResponse)
def login_page(request: Request, template: Template):
    """
    Render the login page for the Ophelia CI interface.

    :return: An HTMLResponse containing the rendered login page.
    """
    return template.TemplateResponse(
        'login.html',
        {
            'request': request,
            'title': 'Ophelia CI - Login',
        },
    )


@router.post('/login', response_class=RedirectResponse)
async def login(
    authentication_service: Authentication,
    private_key: Annotated[
        UploadFile,
        File(
            title='Private Key',
            description="The user's private key file for authentication."
            ' Needs to be rsa256',
        ),
    ],
    username: Annotated[
        str,
        Form(
            title='Username',
            description='The username of the user attempting to log in.',
        ),
    ],
):
    """
    Authenticate a user with their username and private key, then redirect
    to the home page upon successful authentication.

    :param private_key: The user's private key file for authentication.
        Needs to be rsa256.
    :param username: The username of the user attempting to log in.
    :return: A RedirectResponse object that sets a session cookie and
        redirects to the home page.
    """

    try:
        token = authentication_service.authenticate(
            username=username,
            private_key=(await private_key.read()).decode('utf-8'),
        )
    except OpheliaException as e:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED, detail=str(e)
        ) from e
    response = RedirectResponse(url='/', status_code=status.HTTP_303_SEE_OTHER)
    response.set_cookie(key='session', value=token)
    return response


@router.get('/unique', response_class=HTMLResponse)
def unique_key_page(request: Request, template: Template):
    """
    Render the unique key login page for the Ophelia CI interface.

    :return: An HTMLResponse containing the rendered unique key login page.
    """
    return template.TemplateResponse(
        'unique.html',
        {
            'request': request,
            'title': 'Ophelia CI - Unique Key',
        },
    )


@router.post('/unique', response_class=RedirectResponse)
def unique_key(
    unique_key: Annotated[
        str,
        Form(
            title='Unique Key',
            description='The unique key generated when the server is started.',
        ),
    ],
    authentication_service: Authentication,
):
    """
    Authenticate a user with the server's unique key, then redirect to the
    home page upon successful authentication.

    :param unique_key: The unique key generated when the server is started.
    :return: A RedirectResponse object that sets a session cookie and redirects
        to the home page.
    """
    token = authentication_service.authenticate_with_unique_key(unique_key)
    response = RedirectResponse(url='/', status_code=status.HTTP_303_SEE_OTHER)
    response.set_cookie(key='session', value=token)
    return response
