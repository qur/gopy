// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package py (github.com/qur/gopy/lib) provides access to the CPython C API.  This
package presents an idomatic Go interface to the CPython C API described at
http://docs.python.org/c-api/index.html

Instead of simply exposing the C API as-is, this package uses interfaces,
embedding, type assertions and methods to try and present the functionality of
the Python API in a manner that feels more natural in Go.

Embedding Python

Embedding Python is fully supported, with the ability to initialise the
interpreter, enable threading support, manipulate the GIL and call Python API
functions to manipulate Python objects.

In addition to providing the ability to use the API to call into Python calling
from Python back into Go is also supported.  New types can be implemented in Go
and exposed into Python.

In addition to the normal Python C API, a optional (i.e. it must be explictly
enabled) "go" package can be presented to the embedded Python.  This gives
access to some Go functionality - currently the only extra available is the Chan
class, which allows Go and Python code to communicate by exchanging Python
objects over a Go channel.

Python Extensions

Currently there is experimental support for creating CPython extensions in Go on
the ext branch.  This uses gccgo, as there is not currently any support for
creating shared libraries in the gc toolchain.

More information about this can be found at http://gopy.qur.me/extensions
*/
package py
