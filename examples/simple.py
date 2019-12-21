foobars = 1
baz = 2

bazr = <
  bar:
    blah: foo
    baz: bazl
    bar:
    - one
    - two
    - !~ foobars + baz
>

bazr.render()

baz = 5
bazr.render()