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
