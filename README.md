## Configula

Configula is a configuration generation language and processor. It's goal is to make the programattic
definition of declarative configuration easy and intuitive.

### Inspiration
Configula is inspired by the [JSX language](https://reactjs.org/docs/introducing-jsx.html) in [React](https://reactjs.org) that combines Javascript and HTML tags. Configula
defines a similar pattern for Python and YAML (for now).

Consider this declaration:
```
my_object = foo: bar
```

It's neither Python, nor YAML. It combines the syntax of both. Instead of being a templating language like Jinja or others, or a DSL like HCL, it combines the power of
a full programming language with the ease of a declarative syntax like YAML.

### Example

Let's generate a Kubernetes namespace for a collection of users. Each namespace get's the user's name.

Here is what this looks like in Configula:

```python
# Simple example of creating 3 Kubernetes namespaces

# Our users in need of namespaces
users = ['jim', 'sally', 'sue']

# The namespaces objects from YAML
namespaces = map(lambda user: <
        apiVersion: v1
        kind: Namespace
        metadata:
          # Use the !~ tag syntax to inline Python code
          name: !~ user
    >, users)

# Output
render(namespaces)
```

You run this through the `configula` interpreter:
```sh
configula namespaces.py
```

and the output of this looks like:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: bob
---
apiVersion: v1
kind: Namespace
metadata:
  name: sue
---
apiVersion: v1
kind: Namespace
metadata:
  name: sally
---
```

### Getting Started
*** Download Configula for your platform
* [darwin](https://github.com/brendandburns/configula/releases/download/0.0.1/darwin.tgz)
* [linux](https://github.com/brendandburns/configula/releases/download/0.0.1/linux.tgz)
* [windows](https://github.com/brendandburns/configula/releases/download/0.0.1/windows.zip)

*** Clone the git repo
```sh
git clone https://github.com/brendandburns/configula
cd configula
```

*** Try out an example

```sh
./configula examples/tiny.py
```

*** Read the docs

### FAQ
*What about other Javascript? Java? Go? .NET? Other Language?*
We're open to other programming languages. What other languages would you want to see? Please file issues.

*What about JSON? HCL? Other Data Language?*
We're open to other data languages too! What other languages would you want to see? Please file issues.


