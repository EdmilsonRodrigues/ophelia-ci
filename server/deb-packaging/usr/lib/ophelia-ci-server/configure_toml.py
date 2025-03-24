#!/usr/bin/python3
import sys
from pathlib import Path

import toml


def configure_toml(
    config_path,
    server_port,
    db_path,
    expiration_time,
    secret,
    ssl_cert_file,
    ssl_key_file,
):
    config_path = Path(config_path)
    try:
        with config_path.open('r', encoding='utf-8') as f:
            config = toml.load(f)

        config['server'] |= {
            'port': int(server_port),
            'db_path': db_path,
            'expiration_time': int(expiration_time),
        }
        if secret:
            config['server']['secret'] = secret

        elif 'secret' in config['server']:
            del config['server']['secret']

        if ssl_cert_file and ssl_key_file:
            config.setdefault('ssl', {}).update({
                'cert_file': ssl_cert_file,
                'key_file': ssl_key_file,
            })

        elif 'ssl' in config:
            del config['ssl']

        with config_path.open('w', encoding='utf-8') as f:
            toml.dump(config, f)

    except Exception as e:
        print(f'Error configuring TOML: {e}', file=sys.stderr)
        sys.exit(1)


if __name__ == '__main__':
    if len(sys.argv) != 7:
        print(
            'Usage: configure_toml.py <config_path> <server_port> <db_path> '
            '<expiration_time> <secret> <ssl_cert_file> <ssl_key_file>',
            file=sys.stderr,
        )
        sys.exit(1)

    configure_toml(*sys.argv[1:])
