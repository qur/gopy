import go

c = go.Chan(3)

print("put")
c.put("Hello")
c.put("there")
c.put("friend")

print("close")
c.close()
try:
    c.close()
except go.ChanClosedError as e:
    print(f"close failed: {e}")

print("put 2")
try:
    c.put("again")
except go.ChanClosedError as e:
    print(f"put failed: {e}")

print("loop")
for o in c:
    print(f"got: {o}")

print("get")
try:
    print(f"got: {c.get()}")
except go.ChanClosedError as e:
    print(f"get failed: {e}")

print("done")
