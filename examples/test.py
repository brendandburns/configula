foobar = 'blah %d' % 1
c = 'a'
specialChars = 'what\'ll\r\nhappen\r\nhere\\?'

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
    specialCharsString: "what'll\r\nhappen\r\nhere\\?"
    specialCharsExpr: !~ specialChars
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