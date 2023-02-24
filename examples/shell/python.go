package main

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopython.xyz/py/v3"
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

	pyParts := make([]py.Object, len(args))
	for i, arg := range args {
		pyS, err := py.NewUnicode(arg)
		if err != nil {
			return err
		}
		defer pyS.Decref()
		pyParts[i] = pyS
	}

	pyArgs, err := py.PackTuple(pyParts...)
	if err != nil {
		return nil
	}
	defer pyArgs.Decref()

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

	if _, err := os.Stat(name); errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return tryPyFile(cmd, name)
}
