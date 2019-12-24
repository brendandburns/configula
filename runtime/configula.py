import yaml

def maybe_render():
    if YamlVariable.last:
        if hasattr(type(YamlVariable.last), 'render'):
            YamlVariable.last.render()

def render(obj):
    if hasattr(type(obj), '__iter__'):
        for o in obj:
            o.render()
            print("---")
    else:
        o.render()
    return None

class YamlNode(yaml.YAMLObject):
    def __init__(self, value):
        self.value = value
    def __repr__(self):
        return "YamlNode(%s)" % self.value

def yaml_node_representer(dumper, data):
    value = data.value
    if (isinstance(value, int)):
        return dumper.represent_int(value)
    elif (isinstance(value, str)):
        return dumper.represent_str(value)

    return dumper.represent_scalar("!~", data.value)

yaml.add_representer(YamlNode, yaml_node_representer)

class YamlExpr(yaml.YAMLObject):
    def __init__(self, expr):
        self.expr = expr
    def __repr__(self):
        return "YamlExpr(%s)" % self.expr

def yaml_expr_representer(dumper, data):
    value = data.expr()
    if (isinstance(value, int)):
        return dumper.represent_int(value)
    elif (isinstance(value, str)):
        return dumper.represent_str(value)
    else:
        return dumper.represent_scalar("!~", value)
    # TODO: More here!

yaml.add_representer(YamlExpr, yaml_expr_representer)

class YamlVariable:
    # This is used to auto-magically print un-rendered bits
    last = None
    def __init__(self, data):
        self.data = data
        YamlVariable.last = self

    def render(self):
        print(yaml.dump(self.data))
        YamlVariable.last = None
