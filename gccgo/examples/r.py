#!/usr/bin/env python

import rgo
import threading

def two(i, count):
    if count == 0:
        print "= python: 0"
        return
    print "+ python: %d" % count
    rgo.one(i, two, count - 1)
    print "- python: %d" % count

def x():
    rgo.one(1, two, 10)

t = threading.Thread(target=x)

print "-- main --"
x()

print "-- thread --"
t.start()
t.join()
