// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strings"
)

type CmdFunc func(args []string) error

func findCmdFunc(cmd string) (fn CmdFunc, err error) {
	fn = builtins[cmd]

	if fn == nil && err == nil {
		fn, err = findPyCmd(cmd)
	}

	if fn == nil && err == nil {
		fn, err = findExecCmd(cmd)
	}

	return
}

//////////////////////////////////////////////////////////////////////////////
// builtin functions follow ...
//
var builtins = map[string]CmdFunc{
	"echo": echo,
	"cd":   chdir,
}
//////////////////////////////////////////////////////////////////////////////

func echo(args []string) error {
	fmt.Printf("%s\n", strings.Join(args, " "))
	return nil
}

func chdir(args []string) error {
	if len(args) > 0 {
		return os.Chdir(args[0])
	}

	home, err := os.Getenverror("HOME")
	if err != nil {
		return fmt.Errorf("HOME not set")
	}

	return os.Chdir(home)
}
