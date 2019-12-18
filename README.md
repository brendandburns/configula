## Configula

Configula is a configuration generation language and processor. It's goal is to make the programattic
definition of declarative configuration easy and intuitive.

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