import example

class t(example.ExampleClass):
    def __new__(cls, *args, **kwds):
        print(f"NEW: {args} {kwds}")
        return super().__new__(cls, *args, **kwds)

    def __init__(self, *args, **kwds):
        print(f"INIT: {args} {kwds}")
        self.kwds = kwds
        return super().__init__(self, *args, **kwds)

    def foobar(self):
        self.wibble += 4

o = t("arg1", "arg2", kwd1="kwd1", kwd2="kwd2")
o.x = "hello"
print(o)
o.bar()
print(o.kwds)
o.foobar()
print(o)
print(o > example.ExampleClass())
del o
