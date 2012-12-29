goPy Shell
==========

This is a very simple shell that is intended as an example of embedding Python
in a Go program, extending it by implementing modules in Go - and calling into
Python code.

Basically this shell lets you run three types of command (in order of increasing
preference): external commands, Python commands and builtin commands.  The only
builtin commands currently are "echo" and "cd".

Commands can be implemented in Python by putting a function with the name of the
command in a file of the same name in the cmds directory (which must be in the
same directory as the binary).  This Python function will then be called with
the commands line arguments as parameters.  Some example commands are included.
Commands are loaded when they are run, so you can run a Python command, edit the
file, then re-run the command - and get the new behaviour (not best from a
performance point of view, but nicer for playing around).

The shell also makes a sh module available to Python code which exposes the
programs tokeniser, and command executer, to Python code.  These can be seen in
the token and run Python command examples.

You can also run a command in the background by ending the command with a '&'
character.  This is primarily so that the Python interaction code has to deal
with the fact that it may be called from a goroutine - and thus gives
justification to including locking in the example.
