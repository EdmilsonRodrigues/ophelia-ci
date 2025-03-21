from typing import Annotated

from fastapi import APIRouter, File, Form, Request, UploadFile, status
from fastapi.responses import HTMLResponse, RedirectResponse
from ophelia_ci_interface.routers.dependencies import Template, Authentication

router = APIRouter(tags=['Authentication'])


@router.get('/login', response_class=HTMLResponse)
def login_page(request: Request, template: Template):
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
    private_key: Annotated[UploadFile, File()],
    username: Annotated[str, Form()],
):
    token = authentication_service.authenticate(username=username, private_key=await private_key.read())
    response = RedirectResponse(url='/', status_code=status.HTTP_303_SEE_OTHER)
    response.set_cookie(key='session', value=token)
    return response


@router.get('/unique', response_class=HTMLResponse)
def unique_key_page(request: Request, template: Template):
    return template.TemplateResponse(
        'unique.html',
        {
            'request': request,
            'title': 'Ophelia CI - Unique Key',
        },
    )


@router.post('/unique', response_class=RedirectResponse)
def unique_key(unique_key: Annotated[str, Form()], authentication_service: Authentication):
    token = authentication_service.authenticate_with_unique_key(unique_key)
    response = RedirectResponse(url='/', status_code=status.HTTP_303_SEE_OTHER)
    response.set_cookie(key='session', value=token)
    return response
