#!/usr/bin/env python

import simple
import threading
import time

def x():
    print "x is running"
    time.sleep(1)
    simple.example("this is a thread")

print "start"
simple.example("this is the main thread")

t = threading.Thread(target=x)

print "start thread"
t.start()

print "wait for thread"
t.join()
