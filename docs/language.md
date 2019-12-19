# Configula Language
Configula is a melding of YAML for declarative declarations and Python for programming.
There are three ways that this interaction works.

### One-line YAML

To declare a one-line YAML module you can simply place the YAML code in the same place
you would place an Python expression.

Some examples:
```python
# assignment
my_variable = some: yaml expression

# result
def my_yaml_fn():
    return some: other yaml

# constructor
def __init__(self):
    self.my_yaml = another: yaml here
```

### Multi-line YAML
You can also use multi-line YAML objects in your Python. For now this uses the `<` and `>` delimiters (we're investigating lifting this!).

Examples of multi-line YAML:

```python
return <
  my: yaml
  goes: here
  even:
    sub: objects
    and:
    - lists
    - are
    - ok
>
```

### Inline Python
Once you are in a YAML expression, you may want to substitute a Python expression for some YAML field
value. You can do that using the `!~` YAML tag. Everything on the line after that tag is interpretted
as a Python expression.

Example:
```yaml
my: yaml
includes: !~ 'a' + 'python' + 'expression'
```

A note on scoping. The scope for a Python expression is defined at the time the YAML expression is evaluated, *not* the time that it is defined. This means that:

```python
my_value = 1
my_yaml = val: !~ my_value

my_value = 2
my_yaml.render()
```

prints out `my_value: 2`.

### Rendering
By default YAML expressions aren't actually output anywhere. To output them, you need to call the `render` function. This is both a member function for the `YAMLExpression` class as well as a global
function in the `configula` namespace.

Examples:
```python
my_var = some: yaml
# output the YAML
my_var.render()

# output the same YAML again
render(my_var)

my_list = [
    a: b,
    c: d,
    e: f,
]

# output the list
render(my_list)
```