====
gopy
====

--------------------------------------------------------------
Build a CPython extenstion module from Go package source files
--------------------------------------------------------------

:Author: julian@quantumfyre.co.uk
:Date: 2012-05-27
:Copyright: 2012 Julian Phillips
:Version: current
:Manual section: 1
:Manual group: 

Synopis
=======

gopy [-work] <sourcefile>+

Description
===========

gopy compiles the given Go source files into a CPython extenstion module.

The Go package name is used as the module name.  The .so file is created in the
current directory, intermediate files are created in a temporary directory and
deleted after the module is created.

Options
=======

-h	Show help
-work	Print the path of the temporary directory used for intermediate files,
	and don't delete them after building.

Environment Variables
=====================

CFLAGS	Override the default compiler flags (-g)

GOPY_CC	Override the tool used to link the shared library (default: gcc-4.7)

GOPY_GO Override the tool used to build the Go package (default: gccgo)
