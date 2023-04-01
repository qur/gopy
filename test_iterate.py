import example
print("iterate: ['a',1,2,3]")
example.iterate(['a',1,2,3])
print("iterate: 'hello'")
example.iterate('hello')
print("iterate: test.py")
with open('test.py', 'r', encoding='utf-8') as f:
    example.iterate(f)
