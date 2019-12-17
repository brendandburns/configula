import yaml

class YamlNode(yaml.YAMLObject):
    def __init__(self, value):
        self.value = value
    def __repr__(self):
        return "YamlNode(%s)" % self.value

def yaml_node_representer(dumper, data):
    return dumper.represent_str(data.value)

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
    elif (insinstance(value, str)):
        return dumper.represent_str(value)
    else:
        return dumper.represent_scalar("!~", value)
    # TODO: More here!

yaml.add_representer(YamlExpr, yaml_expr_representer)

class YamlVariable:
    def __init__(self, data):
        self.data = data

    def render(self):
        return yaml.dump(self.data)