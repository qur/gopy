// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"exec"
	"os"
)

func execCmd(cmd, fname string, args []string) os.Error {
	attr := &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}

	argv := make([]string, 0, len(args)+1)
	argv = append(argv, cmd)
	argv = append(argv, args...)

	p, err := os.StartProcess(fname, argv, attr)
	if err != nil {
		return err
	}

	_, err = p.Wait(0)
	if err != nil {
		return err
	}

	return p.Release()
}

func findExecCmd(cmd string) (CmdFunc, os.Error) {
	p, err := exec.LookPath(cmd)
	if e, ok := err.(*exec.Error); ok && e.Error == exec.ErrNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return func(args []string) os.Error { return execCmd(cmd, p, args) }, nil
}
