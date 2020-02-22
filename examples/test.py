foobar = 'blah %d' % 1
c = 'a'
quoted = 'what\'ll happen here?'
newlines = 'print\r\nme!'

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
    boolList: [true, false, !~ True]
    nullVal: null
    nullVallExpr: !~ None
    stringWithQuotes: what'll happen here?
    exprWithQuotes: !~ quoted
    stringWithNewline: 'print\r\nme!'
    exprWithNewline: !~ newlines
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