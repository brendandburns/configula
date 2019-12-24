# Simple example of creating 3 Kubernetes namespaces

# Our users in need of namespaces
users = ['jim', 'sally', 'sue']

# The namespaces objects from YAML
ns = \
  apiVersion: v1
  kind: Namespace
  metadata:
    name: !~ user

for user in users:
  ns.render()
