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

- Install Enterprise Dashboard

```shell
./emqx5.py > emqx5.json
```

- Install OpenSource Dashboard 

```
./emqx5.py -c os > emqx5.json 
```




