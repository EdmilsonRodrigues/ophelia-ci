import uvicorn

from ophelia_ci_interface.main import app

if __name__ == '__main__':
    uvicorn.run(app, host='0.0.0.0', port=8000)
