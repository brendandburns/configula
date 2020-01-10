## Configula

Configula is a configuration generation language and processor. It's goal is to make the programmatic
definition of declarative configuration easy and intuitive.

### Inspiration
Configula is inspired by the [JSX language](https://reactjs.org/docs/introducing-jsx.html) in [React](https://reactjs.org) that combines Javascript and HTML tags. Configula
defines a similar pattern for Python and YAML (for now).

Consider this declaration:
```
my_object = foo: bar
```

It's neither Python, nor YAML. It combines the syntax of both. Instead of being a templating language like [Jinja](https://www.palletsprojects.com/p/jinja/) or others, or a Domain-Specific Language (DSL) like [HCL](https://www.terraform.io/docs/configuration/syntax.html), [ARM templates](https://docs.microsoft.com/en-us/azure/azure-resource-manager/templates/template-syntax) or [jsonnet](https://jsonnet.org/), it combines the power of 
a full programming language with the ease of a declarative syntax like YAML.

### Programmatic Example

Let's generate a Kubernetes namespace for a collection of users. Each namespace gets the user's name.

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

An alternate version that uses a for loop looks like this:

```python
# Simple example of creating 3 Kubernetes namespaces
users = [ 'bob', 'sue', 'sally']

# This Python object is defined by inline YAML
namespace = <
  apiVersion: v1
  kind: Namespace
  metadata:
    # The !~ syntax enables python substitution, in this case a variable named `userName`
    name: !~ user
>

for user in users:
  namespace.render()
  print("---")
```

In either case you you run this through the `configula` interpreter:
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

### Templating Example
The previous example is programming centric. If your configuration is mostly data and not much code, sometimes it makes
sense to have a data-centric template. Such templates look more like YAML than Python.

Here is a simple example:

```yaml
first = "brendan"
last = "burns"

apiVersion: v1
kind: Namespace
metadata:
  name: !~ first + ' ' + last
```

When you run this with `configula namespace.yml` it will output:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: brendan burns
```


### Getting Started
Download Configula for your platform
* [darwin](https://github.com/brendandburns/configula/releases/download/0.0.3/darwin.tgz)
* [linux](https://github.com/brendandburns/configula/releases/download/0.0.3/linux.tgz)
* [windows](https://github.com/brendandburns/configula/releases/download/0.0.3/windows.zip)

Clone the git repo
```sh
git clone https://github.com/brendandburns/configula
cd configula
```

Try out an example

```sh
./configula examples/tiny.py
```

Read the [docs](https://github.com/brendandburns/configula/tree/master/docs)

### Integrating with Kubernetes
Configula can be used as a `kubectl` plugin.

To install the plugin run this command:

```
# assumes that configula is in /usr/local/bin/ change paths as needed.
ln -s /usr/local/bin/kubectl-configula /usr/local/bin/configula
```

Once that's done you can run:

```sh
kubectl configula create -f examples/namespaces.py
```

### FAQ
*What about other Javascript? Java? Go? .NET? Other Language?*

We're open to other programming languages. What other languages would you want to see? Please file issues.

*What about JSON? HCL? Other Data Language?*

We're open to other data languages too! What other languages would you want to see? Please file issues.


