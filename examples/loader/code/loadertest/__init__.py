class Wibble:
    def do(self, *args, **kwargs):
        print(f"PYTHON DO: {args} {kwargs}")

print("HELLO")

def foobar(t):
    w = t()
    w.do("hello", "there", a=12)
