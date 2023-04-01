#!/usr/bin/env python3.11

import extension


o = extension.ExampleClass()
print(f"wibble: {o.wibble}")
o.bar()
print(f"dir: {dir(o)}")
print(f"foo: {o.foo}")
o.foo = 12
print(f"foo: {o.foo}")
o.bar()
print(f"x: {o.x}")
o.x = [1,2,3]
print(f"x: {o.x}")
o.bar()
print(f"y: {o.y}")
o.y = (1,2,3)
print(f"x: {o.y}")
o.bar()
