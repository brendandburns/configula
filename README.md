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

This is neither Python, nor YAML, it combines the syntax of both. Instead of being a templating language like Jinja or others, or a DSL like HCL, it combines the power of
a full programming language with the ease of a declarative syntax like YAML.

### Example

Consider the task of generating a Kubernetes namespace for a collection of users, where each
namespace is named after the user. Here is what this looks like in Configula:

```
# Simple example of creating 3 Kubernetes namespaces

# Our users in need of namespaces
users = ['jim', 'sally', 'sue']

# The namespaces objects from YAML
namespaces = map(lambda user: <
        apiVersion: v1
        kind: Namespace
        metadata:
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

*** Clone the git repo

*** Try out an example

*** Read the docs

### FAQ
*What about other Javascript? Java? Go? .NET? Other Language?*
We're open to other programming languages. What other languages would you want to see? Please file issues.

*What about JSON? HCL? Other Data Language?*
We're open to other data languages too! What other languages would you want to see? Please file issues.


