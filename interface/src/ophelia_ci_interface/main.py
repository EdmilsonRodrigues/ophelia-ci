from pathlib import Path
from typing import Literal

from fastapi import FastAPI, Request
from fastapi.responses import HTMLResponse
from fastapi.staticfiles import StaticFiles

from ophelia_ci_interface.config import VERSION
from ophelia_ci_interface.routers.dependencies import Health, Template
from ophelia_ci_interface.routers.repository_router import (
    router as repository_router,
)
from ophelia_ci_interface.routers.user_router import router as user_router

app = FastAPI(version=VERSION)
app.mount(
    '/static',
    StaticFiles(directory=Path('resources', 'static')),
    name='static',
)

app.include_router(repository_router)
app.include_router(user_router)


@app.get('/health', tags=['Common'])
def root() -> dict[Literal['version'], str]:
    """
    Return the version of Ophelia CI Interface.

    Returns:
        dict[Literal['version'], str]: A dictionary containing a single key,
            'version', whose value is the version of Ophelia CI.
    """
    return {'version': VERSION}


@app.get('/', response_class=HTMLResponse, tags=['Common'])
def home(request: Request, template: Template, health_service: Health):
    """
    Return the homepage of Ophelia CI Interface.

    Returns:
        HTMLResponse: The rendered homepage.
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
