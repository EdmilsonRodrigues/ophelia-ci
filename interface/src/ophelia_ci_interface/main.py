from typing import Literal

from fastapi import FastAPI, Request
from fastapi.responses import HTMLResponse, RedirectResponse
from fastapi.staticfiles import StaticFiles

from ophelia_ci_interface.config import VERSION, base_path
from ophelia_ci_interface.routers.authentication_router import (
    router as authentication_router,
)
from ophelia_ci_interface.routers.dependencies import Health, Template
from ophelia_ci_interface.routers.repository_router import (
    router as repository_router,
)
from ophelia_ci_interface.routers.user_router import router as user_router

app = FastAPI(version=VERSION)
app.mount(
    '/static',
    StaticFiles(directory=base_path / 'resources' / 'static'),
    name='static',
)

app.include_router(authentication_router)
app.include_router(repository_router)
app.include_router(user_router)


@app.middleware('http')
async def redirect_401(request: Request, call_next):
    """
    Redirect all 401 responses to the login page.

    :param request: The current request.
    :param call_next: The next middleware to call.
    :return: The response from the next middleware if the status code
    is not 401, otherwise a RedirectResponse to the login page.
    """
    response = await call_next(request)
    if response.status_code == 401:
        return RedirectResponse(url='/login')
    return response


@app.get('/health', tags=['Common'])
def root() -> dict[Literal['version'], str]:
    """
    Return the version of Ophelia CI Interface.

    :return: A dictionary containing the version of Ophelia CI Interface.
    """
    return {'version': VERSION}


@app.get('/', response_class=HTMLResponse, tags=['Common'])
def home(request: Request, template: Template, health_service: Health):
    """
    Return the homepage of Ophelia CI Interface.

    :return: An HTMLResponse containing the rendered homepage.
    """
    return template.TemplateResponse(
        'index.html',
        {
            'request': request,
            'title': 'Ophelia CI',
            'page_title': 'Welcome to the Ophelia CI',
            'status': health_service.get_status(),
        },
    )
