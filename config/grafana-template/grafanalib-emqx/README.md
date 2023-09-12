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

- Install EMQX Enterprise Dashboard

```shell
./emqx5.py -e ee > emqx5.json
```

- Install EMQX Community Edition Dashboard

```
./emqx5.py -e ce > emqx5.json
```
