# Define a YAML object where the 'foo' field has the value of evaluating 1 + 2 (e.g. 3)
my_obj = foo: !~ 1 + 2

my_obj.render()