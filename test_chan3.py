import go
import time


def item(o):
    print(f"GOT: {o}")

def closed():
    print("CLOSED")

c = go.Chan(3)

c.monitor(item, closed=closed)

print("put")
c.put("Hello")
c.put("there")
c.put("friend")

print("sleep 3")
time.sleep(3)

print("close")
c.close()
try:
    c.close()
except go.ChanClosedError as e:
    print(f"close failed: {e}")

print("sleep 3")
time.sleep(3)

print("put 2")
try:
    c.put("again")
except go.ChanClosedError as e:
    print(f"put failed: {e}")

print("done")
