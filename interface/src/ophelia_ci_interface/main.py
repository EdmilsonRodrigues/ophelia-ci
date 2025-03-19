from fastapi import FastAPI, Request
from fastapi.responses import HTMLResponse, RedirectResponse
from fastapi.staticfiles import StaticFiles

from ophelia_ci_interface.config import VERSION, base_path
from ophelia_ci_interface.models.health import (
    HealthService,
)
from ophelia_ci_interface.routers.dependencies import Template

app = FastAPI(version=VERSION)
app.mount(
    '/static', StaticFiles(directory=base_path / 'static'), name='static'
)


@app.get('/login')
def login():
    return RedirectResponse(url='/')


@app.get('/health')
def root():
    return {'version': VERSION}


@app.get('/', response_class=HTMLResponse)
def home(request: Request, template: Template):
    return template.TemplateResponse(
        'index.html',
        {
            'request': request,
            'title': 'Ophelia CI',
            'page_title': 'Welcome to the Ophelia CI',
            'status': HealthService.get_status(),
        },
    )
