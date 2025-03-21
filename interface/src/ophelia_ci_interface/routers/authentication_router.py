from typing import Annotated

from fastapi import APIRouter, File, Form, Request, UploadFile
from fastapi.responses import HTMLResponse, RedirectResponse
from ophelia_ci_interface.routers.dependencies import Template

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
def login(
    private_key: Annotated[UploadFile, File()],
    username: Annotated[str, Form()],
):
    return RedirectResponse(url='/')


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
def unique_key():
    return RedirectResponse(url='/')
