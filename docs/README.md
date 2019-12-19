## Configula documentation

### Usage
Simple usage

```sh
configula <file-name>
```

Usage with Kubernetes
```sh
configula <file-name> | kubectl create -f -
```

#### Flags
* `--debug=true`: Just dump the generated code, don't execute it.
* `--python=</some/path/to/python>`: Specify the Python interpreter to use

#### Environment variables
* `${CONFIGULA_PYTHON}`: Overrides `--python` flag and sets the Python interpreter to use.

### Further Documentation
* [Building Configula](building.md)
* [Configula Language](language.md)
* [Configula Internals](internals.md)
