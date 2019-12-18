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
