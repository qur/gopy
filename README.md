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

This codebase only supports a single version of Python at a time, each version
on a separate branch.

See the specific branches for more details:
* [Python 2.6](https://github.com/qur/gopy/tree/python2.6) (broken)
* [Python 2.6](https://github.com/qur/gopy/tree/python2.7) (broken)
* [Python 3.11](https://github.com/qur/gopy/tree/python3.11)

