// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"os/exec"
)

func execCmd(cmd, fname string, args []string) error {
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

func findExecCmd(cmd string) (CmdFunc, error) {
	p, err := exec.LookPath(cmd)
	if e, ok := err.(*exec.Error); ok && e.Err == exec.ErrNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return func(args []string) error { return execCmd(cmd, p, args) }, nil
}
