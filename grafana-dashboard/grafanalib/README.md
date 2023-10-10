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

- Install EMQX 5.x Enterprise Dashboard

```shell
./emqx.py -e ee -v 5 > emqx5-ee.json
```

- Install EMQX 5.x Community Edition Dashboard

```
./emqx.py -e ce -v 5 > emqx5-ce.json
```

- Install EMQX 4.x Enterprise Dashboard

```shell
./emqx.py -e ee -v 4 > emqx4-ee.json
```

- Install EMQX 4.x Community Edition Dashboard

```
./emqx.py -e ce -v 4 > emqx4-ce.json
```
