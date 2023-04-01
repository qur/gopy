import go

c = go.Chan(3)

print("put")
c.put("Hello")
c.put("there")
c.put("friend")

print("close")
c.close()

print("loop")
for o in c:
    print(f"got: {o}")

print("done")
