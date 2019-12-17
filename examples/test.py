foobar = "blah" + 1

baz = (
  foo: "bar bar"
  baz: blah
  bar: 1
  bob:
    foo: bar
    blah:
      blah: blahber
    someList:
    - a
    - b
    - c
    otherList: [1, 2, 3]
    anotherList: [
      a, b, c, ! c + 'd' + fn(baz)
    ]
)

foo: ! a + b

(
  alist:
  - one
  !! 1 + 4
  - two
  - three
)

