# ProxyToGo

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/license/mit/)

## Description

ProxyToGo is a lightweight and customizable proxy server written in Go. It allows you to easily configure multiple listeners and corresponding backend addresses for proxying requests.

## Features

- Load balancing using a round-robin algorithm
- Makefile for building, running, and managing the application
- Support for multiple listeners and backend addresses
- Debug mode for detailed logging
- Easy configuration through YAML file

## Prerequisites

- Go 1.16 or higher

## Installation

1. Clone the repository:

   ```shell
   git clone https://github.com/whilezo/proxytogo.git
 
2. Change to the project directory:
   ```shell
   cd proxytogo

3. Install all dependencies:
   ```shell
   make deps

4. Build the proxy binary: 
   ```shell
   make build

## Configuration

The configuration of ProxyToGo is defined in the etc/config.yaml file. You can specify multiple listeners and their corresponding backend server addresses in the following format:
   ```yaml
   listeners:
     - listenerAddress: localhost:8000
       backendAddresses:
         - 127.0.0.1:9000
         - 127.0.0.1:9001
     - listenerAddress: localhost:8080
       backendAddresses:
         - 127.0.0.1:9090
         - 127.0.0.1:9091
   debug: true
   ```

- listenerAddress: The address on which the proxy server listens for client connections.
- backendAddresses: The addresses of the backend servers to which the client requests will be forwarded.
- debug: Enable or disable debug mode for logging (true or false).

## Usage
To build and run PoxyToGo, you can use the provided Makefile targets. Here are the available targets:

- `make build`: Build the proxy binary.
- `make run`: Run the proxy server.
- `make clean`: Clean build artifacts.
- `make deps`: Install project dependencies.
- `make test`: Run tests.
- `make docs`: Generate code documentation.
- `make help`: Display available targets and their descriptions.

2. Run the proxy:
   ```shell
   make run

## Contributing
Contributions to ProxyToGo are welcome! If you encounter any issues, have suggestions, or would like to contribute new features or improvements, please open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](https://opensource.org/license/mit/).