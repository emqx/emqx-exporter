# grafanalib-emqx

## Setup

Usage of Python 3 is required. It can be installed [on Python.org](https://www.python.org/downloads/)

- Install Pipenv

```shell
pip install --user pipenv
```

- Install all packages

```shell
pipenv install
```

- Activate the virtual environment

```shell
pipenv shell
```

## Generate Dashboard JSON

- Generate EMQX 5.x Enterprise Dashboard

```shell
make
```

- Generate EMQX 5.x Community Edition Dashboard

```shell
make EDITION_ARG=ce VERSION_ARG=5
```

- Generate EMQX 4.x Enterprise Dashboard

```shell
make EDITION_ARG=ee VERSION_ARG=4
```

- Generate EMQX 4.x Community Edition Dashboard

```shell
make EDITION_ARG=ce VERSION_ARG=4
```

## Remove Dashboard JSON

```shell
make clean
```
