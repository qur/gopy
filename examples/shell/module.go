package main

import (
	"gopython.xyz/py/v3"
)

var shMod *py.Module

func pyTokenise(args *py.Tuple) (py.Object, error) {
	s := ""
	err := py.ParseTuple(args, "s", &s)
	if err != nil {
		return nil, err
	}

	parts := tokenise(s)

	pyParts := make([]py.Object, len(parts))
	for i, arg := range parts {
		pyS, err := py.NewUnicode(arg)
		if err != nil {
			return nil, err
		}
		pyParts[i] = pyS
	}

	return py.PackTuple(pyParts...)
}

func pyRun(args *py.Tuple, kw *py.Dict) (py.Object, error) {
	command := ""
	var bg py.Object
	var argO py.Object

	kwlist := []string{"cmd", "args", "bg"}

	err := py.ParseTupleAndKeywords(args, kw, "s|OO", kwlist, &command, &argO, &bg)
	if err != nil {
		return nil, err
	}

	var argS []py.Object

	switch a := argO.(type) {
	case *py.Tuple:
		argS = a.Slice()
	case *py.List:
		argS = a.Slice()
	default:
		return nil, py.NewError(py.TypeError, "args must be tuple or list")
	}

	cArgs := make([]string, len(argS))
	for i, obj := range argS {
		str, err := obj.Base().Str()
		if err != nil {
			return nil, err
		}
		cArgs[i] = str.String()
	}

	if bg == nil {
		bg = py.False
	}

	if err = runCmd(command, cArgs, bg.Base().IsTrue()); err != nil {
		return nil, err
	}

	return py.None, nil
}

var shModule = py.ModuleDef{
	Name: "sh",
	Methods: []py.GoMethod{
		{"tokenise", pyTokenise, "tokenise the given string"},
		{"run", pyRun, "run the given command"},
	},
	Package: true,
}

var cmdsModule = py.ModuleDef{
	Name:    "sh.__cmds__",
	Package: true,
}

func setupShModule() error {
	shMod, err := py.CreateModule(&shModule)
	if err != nil {
		return err
	}
	defer shMod.Decref()

	if err := shMod.Register(); err != nil {
		return err
	}

	mod, err := py.CreateModule(&cmdsModule)
	if err != nil {
		return err
	}
	defer mod.Decref()

	if err := mod.Register(); err != nil {
		return err
	}

	if err := shMod.AddObjectRef("__cmds__", mod); err != nil {
		return err
	}

	return nil
}
