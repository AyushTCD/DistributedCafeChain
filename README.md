# CS7NS6 Exercise 3

Group A's codebase for CS7NS6 Exercise 3.

## Setup

This codebase uses a combination of Go, for the API server and the worker, and Python, for the load and fault testing.

### Go

The code requires Go >= 1.18.

If you want to use interactive debugging in VSCode:

* Ensure you have the official Go extension installed
* Ensure you have the Delve debugger installed with `go install github.com/go-delve/delve/cmd/dlv@latest`
* Add a `.env` file in `loyalty-service` with the `MYSQL_URI` environment variable

A `launch.json` is included in the repo, so all you need to do is hit F5 to start debugging once you're set up

### Python

This repository is using [pipenv](https://pipenv.pypa.io/en/latest/) for dependency management.

* Install pipenv: `pip install --user pipenv`
* Install dependencies: `pipenv install` (in repository root)
* Activate environment: `pipenv shell`

*If using VSCode, you can select an interpreter using the Command Palette, or by clicking the Python version number in the gutter while a Python file is open*

## Load Testing

For load testing the application, we are using the [Locust](https://locust.io/) framework.

To run a test, ensure you have activated the Python environment using the instructions above, then run `locust -f <test_file.py>` using one of the test files in `load_testing/`.
