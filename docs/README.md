# Documentation

## Building

To view the documentation locally create a virtual environment and install [requirements](./requirements.txt)

```bash
$ virtualenv ~/docs && source ~/docs/bin/activate

$ pip install -r requirements.txt

$ mkdocs serve

Serving on http://127.0.0.1:8000
```

## Tooling

| Tool              | Documentation                       | Sources                           |
|-------------------|-------------------------------------|-----------------------------------|
| mkdocs            | [documentation][mkdocs]             | [Sources][mkdocs-src]             |
| mkdocs-material   | [documentation][mkdocs-material]    | [Sources][mkdocs-material-src]    |
| pymdown-extensions| [documentation][pymdown-extensions] | [Sources][pymdown-extensions-src] |


[mkdocs]: https://www.mkdocs.org "Mkdocs"
[mkdocs-src]: https://github.com/mkdocs/mkdocs "Mkdocs - Sources"

[mkdocs-material]: https://squidfunk.github.io/mkdocs-material/ "Material for MkDocs"
[mkdocs-material-src]: https://github.com/squidfunk/mkdocs-material "Material for MkDocs - Sources"

[pymdown-extensions]: https://facelessuser.github.io/pymdown-extensions/ "PyMdown Extensions"
[pymdown-extensions-src]: https://github.com/facelessuser/pymdown-extensions "PyMdown Extensions - Sources"
