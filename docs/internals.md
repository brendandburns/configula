# Configula Internals

## Operation
Internally, Configula extracts YAML expressions from the Python code and replaces the YAML expressions
with an instance of the `YamlNode` Python class. That class holds the parse tree for the YAML expression.
Within that parse tree, inline Python expressions are represented as `YamlExpr` classes.

The generated Python code is then sent to the `python` interpreter. The actual executable can be controlled using either the `${CONFIGULA_PYTHON}` environment variable or the `--python` flag on
the `configula` tool.

If you want to see the generated Python code instead of executing it inline, you can pass the `--debug=true` flag which will simply print the generated code instead of trying to execute it.

## Components
The operation of configula is split into three abstract interfaces:
* Parser - Splits the input into a collection of Python and YAML expressions.
* Processor - Processes each section to transform the YAML expression into Python code.
* Generator - Stiches all of the sections back into a single Python file that can be executed.

There are currently a single implementation of each of these classes, but they are easily replaceable
with alternate (or better) implementations.
