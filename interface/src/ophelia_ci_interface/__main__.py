import uvicorn

from ophelia_ci_interface.config import Settings
from ophelia_ci_interface.main import app


def main():
    settings = Settings()
    optional_args = {
        'ssl_keyfile': settings.SSL_KEYFILE,
        'ssl_certfile': settings.SSL_CERTFILE,
        'workers': settings.WORKERS,
    }
    optional_args = {k: v for k, v in optional_args.items() if v is not None}
    uvicorn.run(
        app, loop='uvloop', host='0.0.0.0', port=settings.PORT, **optional_args
    )


if __name__ == '__main__':
    main()
