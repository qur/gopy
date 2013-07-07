====
gopy
====

--------------------------------------------------------------
Build a CPython extenstion module from Go package source files
--------------------------------------------------------------

:Author: julian@quantumfyre.co.uk
:Date: 2013-07-07
:Copyright: 2012-13 Julian Phillips
:Version: current
:Manual section: 1
:Manual group: 

Synopis
=======

gopy [-v] [-work] [-exe] <sourcefile>+

Description
===========

gopy compiles the given Go source files into a CPython extenstion module, or
into a CPython enabled program.

The Go package name is used as the module (or executable) name.  The output (.so
or executable) file is created in the current directory, intermediate files are
created in a temporary directory and deleted after the module is created.

If creating an executable, then the package must include a method called 'Main'.
This method takes and returns no arguments, i.e. it has the signature:

  func Main()

Options
=======

-h	Show help
-v      Verbose output.  Shows the commands that are being run.
-work	Print the path of the temporary directory used for intermediate files,
	and don't delete them after building.
-exe    Create an executable instead of an extension module.

Environment Variables
=====================

CFLAGS	Override the default compiler flags (-g)

GOPY_CC	Override the tool used to link the shared library (default: gcc)

GOPY_GO Override the tool used to build the Go package (default: gccgo)
