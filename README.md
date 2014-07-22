goPy
====

What is this?
-------------

Idiomatic Python bindings for Google Go

What does that mean?
--------------------

Instead of simply duplicating the Python C API in Go, this is an attempt to
provide the Python API in a form that fits with writing normal Go code.

This means that `PyFoo_XXX(foo, ...)` is `foo.XXX(...)`, and `PyFoo_Check(obj)`
is `foo, ok := obj.(*Foo)` to give a couple of examples.

What versions of Python are supported?
--------------------------------------

Currently, only Python 2.7 is supported.  Though I do hope to add support for
other versions in the future.

Note: Python 2.6 support is currently available on the python2.6 branch.

Docker image
-------------

To provide a low entry point a docker image ([docker-gopy](https://github.com/ChristianKniep/docker-gopy)) is available.
The image ([qnib/gopy](https://registry.hub.docker.com/u/qnib/gopy/)) is part of the official docker registry.

To run the project's examples, one has to install docker and run the following. If it's not available locally it is
going to be pulled from the 'internet'.

```
$ docker run --rm -ti qnib/gopy /bin/bash
root@10515b87b862:/# su -
root@10515b87b862:~# gopy simple.go 
root@10515b87b862:~# python simple.py 
simple.example: [hello map[123:true]]
root@10515b87b862:~# gopy parallel.go 
root@10515b87b862:~# python parallel.py
x = [1, 2, 3, 4]
```
