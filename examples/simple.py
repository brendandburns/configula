foobars = 1
baz = 2

baz = (
  bar:
    blah: foo
    baz: bazl
    bar:
    - one
    - two
    - !~ foobars + baz
)
