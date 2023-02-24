/*
Package py (gopython.xyz/py/v3) provides access to the CPython C API.  This
package presents an idiomatic Go interface to the CPython C API described at
http://docs.python.org/3/c-api/index.html

Instead of simply exposing the C API as-is, this package uses interfaces,
embedding, type assertions and methods to try and present the functionality of
the Python API in a manner that feels more natural in Go.

# Embedding Python

Embedding Python is fully supported, with the ability to initialise the
interpreter, enable threading support, manipulate the GIL and call Python API
functions to manipulate Python objects.

In addition to providing the ability to use the API to call into Python, calling
from Python back into Go is also supported.  New types can be implemented in Go
and exposed into Python.

In addition to the normal Python C API, an optional (i.e. it must be explicitly
enabled) "go" package can be presented to the embedded Python.  This gives
access to some Go functionality - currently the only extra available is the Chan
class, which allows Go and Python code to communicate by exchanging Python
objects over a Go channel.

# Python Extensions

Building extensions is supported using the c-shared buildmode. Some boilerplate
code is needed to meet the Python C API - this can be generated using the
gen_extension command in this repository.
*/
package py
