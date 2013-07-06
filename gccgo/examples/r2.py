#!/usr/bin/env python

import rgo
import threading

data  = threading.local()
data.x = 1

def two(i, count):
    if count == 0:
        print "%d = 0 py %d" % (i, data.x)
        return
    print "%d + %d py %d" % (i, count, data.x)
    rgo.one(i, two, count - 1)
    print "%d - %d py %d" % (i, count, data.x)

def x():
    data.x = 2
    rgo.one(2, two, 1)
    rgo.one(2, two, 2)
    rgo.one(2, two, 3)

def x2():
    data.x = 3
    rgo.one(3, two, 1)
    rgo.one(3, two, 2)
    rgo.one(3, two, 3)

t = threading.Thread(target=x)
t2 = threading.Thread(target=x2)

t.start()
t2.start()

rgo.one(1, two, 1)
rgo.one(1, two, 2)
rgo.one(1, two, 3)

t.join()
t2.join()
