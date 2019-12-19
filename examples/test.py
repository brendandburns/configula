foobar = 'blah %d' % 1
c = 'a'

def fn(val):
  return val + '!!'

baz = <
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
      a, b, c, !~ c + 'd' + fn('bazzer')
    ]
>

foo: !~ a + b

<
  alist:
  - one
  - !~ 1 + 4
  - two
  - three
>

baz.render()