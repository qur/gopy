// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"github.com/qur/gopy/lib"
)

var mydir string

func setupPython() error {
	lock := py.InitAndLock()
	defer lock.Unlock()

	err := setupShModule()
	if err != nil {
		return err
	}

	dir, _ := path.Split(os.Args[0])
	mydir, err = filepath.Abs(dir)
	if err != nil {
		return err
	}

	return nil
}

func callPyCmd(fn py.Object, args []string) error {
	lock := py.NewLock()
	defer lock.Unlock()

	pyArgs, err := py.NewTuple(int64(len(args)))
	if err != nil {
		return nil
	}
	defer pyArgs.Decref()

	for i, arg := range args {
		s, err := py.NewString(arg)
		if err != nil {
			return err
		}
		err = pyArgs.SetItem(int64(i), s)
		if err != nil {
			s.Decref()
			return err
		}
	}

	ret, err := fn.Base().CallObject(pyArgs)
	if err != nil {
		return err
	}
	ret.Decref()

	return nil
}

func tryPyFile(cmd, filename string) (CmdFunc, error) {
	lock := py.NewLock()
	defer lock.Unlock()

	code, err := py.CompileFile(filename)
	if err != nil {
		return nil, err
	}

	mod, err := py.ExecCodeModule("sh.__cmds__."+cmd, code)
	if err != nil {
		return nil, err
	}

	fn, err := mod.Base().GetAttrString(cmd)
	if err != nil {
		return nil, err
	}

	return func(args []string) error { return callPyCmd(fn, args) }, nil
}

func findPyCmd(cmd string) (CmdFunc, error) {
	// Python commands can't have a / in them ...
	if strings.Contains(cmd, "/") {
		return nil, nil
	}

	name := path.Join(mydir, "cmds", cmd+".py")

	_, err := os.Stat(name)
	if err != nil {
		perr := err.(*os.PathError)
		if perr.Err == os.ENOENT {
			return nil, nil
		}
		return nil, err
	}

	return tryPyFile(cmd, name)
}
