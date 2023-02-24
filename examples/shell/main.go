package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func runCmd(command string, args []string, bg bool) error {
	cmdFunc, err := findCmdFunc(command)
	if err != nil {
		return err
	}

	if cmdFunc == nil {
		return fmt.Errorf("unknown Command: %s", command)
	}

	if bg {
		go func() {
			if err := cmdFunc(args); err != nil {
				fmt.Printf("Error: %s\n", err)
			}
		}()
		return nil
	}

	return cmdFunc(args)
}

func main() {

	if err := setupPython(); err != nil {
		log.Fatalf("Setup Error: %s\n", err)
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

		if err := runCmd(command, args, bg); err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
	}
}
