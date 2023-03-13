gen_extension
=============

What is this?
-------------

A command to generate boiler-plate code for a Python extension module.

How to Use
----------

When building a Python extension module, you have to have an exported C function
in your module with the name expected by Python. This command will generate a
file containing such a function as a thin wrapper around a `func() (*py.Module,
error)` function.

To run the command use a go generate line something like this:
```go
//go:generate go run gopython.xyz/py/v3/cmd/gen_extension name file.go
```
Where `file.go` is the name of the Go file to create, and `name` is the name of
the extension (both the compiled library and the function should have the same
name, so `name` should match the compiled `name.so`).

By default, the package of the generated code will be `main`, and the wrapped
function will be assumed to be called `initExtension`. However both of these can
be changed. The full usage is (keeping the default values):
```go
//go:generate go run gopython.xyz/py/v3/cmd/gen_extension -package main -func initExtension name file.go
```
