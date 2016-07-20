// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"os"
)

func runCmd(command string, args []string, bg bool) error {
	cmdFunc, err := findCmdFunc(command)
	if err != nil {
		return err
	}

	if cmdFunc == nil {
		return fmt.Errorf("Unknown Command: %s", command)
	}

	if bg {
		go func() {
			err := cmdFunc(args)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
			}
		}()
		return nil
	}

	return cmdFunc(args)
}

func main() {
	err := setupPython()
	if err != nil {
		fmt.Printf("Setup Error: %s\n", err)
	}

	for {
		command, args, bg, err := getCommand()
		if err == io.EOF {
			fmt.Print("\n")
			os.Exit(0)
		} else if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		if command == "" {
			continue
		}

		err = runCmd(command, args, bg)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
	}
}
