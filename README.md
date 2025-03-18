# Ophelia CI

[![License](https://img.shields.io/badge/License-GPLv3-blue.svg)](LICENSE) [![Docker Pulls](https://img.shields.io/docker/pulls/edmilsonrodrigues/ophelia-server)](https://hub.docker.com/r/edmilsonrodrigues/ophelia-server)

Ophelia CI is an open-source Git server that allows teams to host and manage their own Git repositories. It provides a comprehensive solution with a server, a CLI client, and a graphical web interface.

## Features

* **Self-Hosted Git Server:** Host your own Git server on your infrastructure.
* **Multi-Platform Support:** Available as .deb, .snap, .rock packages, and Docker images for the server and CLI. The web interface is available as .deb, .snap, .rock, .AppImage, .flatpak and Docker.
* **Command-Line Interface (CLI):** Manage repositories efficiently from the command line.
* **Web Interface:** An intuitive web interface for easy repository management.
* **Open Source:** Released under the GPLv3 License.
* **Coming Soon: Integrated CI/CD:** Plan to add a built-in CI tool with YAML workflow support.

## Components

* **Ophelia Server:** A Linux daemon written in Go, responsible for managing Git repositories.
* **Ophelia CLI:** A command-line client written in Go for interacting with the Ophelia Server.
* **Ophelia Web:** A graphical web interface built with Python and FastAPI for user-friendly repository management.

## Installation

### Docker

* **Server:**

    ```bash
    docker pull edmilsonrodrigues/ophelia-server
    docker run -d -p 50051:50051 edmilsonrodrigues/ophelia-server
    ```

* **Web Interface:**

    ```bash
    docker pull edmilsonrodrigues/ophelia-web
    docker run -d -p 8000:8000 edmilsonrodrigues/ophelia-web
    ```

### Package Managers

* **.deb, .snap, .rock, .AppImage, .flatpak:**

    Download the appropriate package from the [releases page](YOUR_RELEASES_PAGE_URL) and install it using your system's package manager.

    Example for .deb:

    ```bash
    sudo dpkg -i ophelia-server.deb
    ```

    (Replace `ophelia-server.deb` with the actual file name.)

    Example for snap:

    ```bash
    sudo snap install ophelia-server.snap
    ```

    (Replace `ophelia-server.snap` with the actual file name.)

    Follow similar procedures for other package types.

## Usage

* **Server:**

    The server runs as a daemon. Refer to the server's documentation for configuration options.

* **CLI:**

    Use the `ophelia-cli` command followed by the desired subcommands. Refer to the CLI's documentation for command details.

* **Web Interface:**

    Access the web interface in your browser at `http://localhost:8000` (or the port you configured).

## Roadmap

* **Continuous Integration (CI):** Implement a CI tool with YAML workflow support for automated testing and deployment.
* **Enhanced Web Interface:** Add more features and improve the user experience of the web interface.
* **Improved Security:** Strengthen the security of the server and web interface.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

This project is licensed under the GPLv3 License.

## Contact

For questions and support, please open an issue or contact [edmilon.monteiro.rodrigues@gmail.com](mailto:edmilon.monteiro.rodrigues@gmail.com).

[YOUR_RELEASES_PAGE_URL]: (Replace with the URL to your releases page)
[edmilon.monteiro.rodrigues@gmail.com](mailto:edmilon.monteiro.rodrigues@gmail.com)
