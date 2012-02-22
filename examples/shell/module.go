// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "github.com/qur/gopy/lib"

var shMod *py.Module

func pyTokenise(args *py.Tuple) (py.Object, error) {
	s := ""
	err := py.ParseTuple(args, "s", &s)
	if err != nil {
		return nil, err
	}

	parts := tokenise(s)

	t, err := py.NewTuple(int64(len(parts)))
	if err != nil {
		return nil, err
	}

	for i, arg := range parts {
		pyS, err := py.NewString(arg)
		if err != nil {
			t.Decref()
			return nil, err
		}
		t.SetItem(int64(i), pyS)
	}

	return t, nil
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
		return nil, py.NewTypeErrorString("args must be tuple or list")
	}

	cArgs := make([]string, len(argS))
	for i, obj := range argS {
		str, err := obj.Base().Str()
		if err != nil {
			return nil, err
		}
		pyS := str.(*py.String)
		cArgs[i] = pyS.String()
	}

	if bg == nil {
		bg = py.False
	}

	err = runCmd(command, cArgs, bg.IsTrue())
	if err != nil {
		return nil, err
	}

	py.None.Incref()
	return py.None, nil
}

func setupShModule() error {
	methods := []py.Method{
		{"tokenise", pyTokenise, "tokenise the given string"},
		{"run", pyRun, "run the given command"},
	}

	shMod, err := py.InitModule("sh", methods)
	if err != nil {
		return err
	}

	mod, err := py.InitModule("sh.__cmds__", []py.Method{})
	if err != nil {
		return err
	}

	err = shMod.AddObject("__cmds__", mod)
	if err != nil {
		return err
	}

	return nil
}
