[![Go Reference](https://pkg.go.dev/badge/gopython.xyz/py/v3.svg)](https://pkg.go.dev/gopython.xyz/py/v3)
[![Go Report Card](https://goreportcard.com/badge/gopython.xyz/py/v3)](https://goreportcard.com/report/gopython.xyz/py/v3)

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

This codebase only supports a single version of Python at a time, currently this
is version 3.11 (note that this only applies to Major.Minor, not the patch
version). For a different version of Python a different branch is needed.

Installation and Building
-------------------------

This module can be installed as normal using go get:

```
go get gopython.xyz/py/v3
```

However, you will need Python and libffi libraries installed (along with their
dependencies) to build. To build extensions, Python will need to be built with
shared libraries.

Python Extensions
-----------------

Building extensions is supported using the c-shared buildmode. Some boilerplate
code is needed to meet the Python C API - this can be generated using the
gen_extension command in this repository.
