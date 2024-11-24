package main

import (
	"log"
	"os"

	"gopython.xyz/py/v3"
	"gopython.xyz/py/v3/examples/example/pkg/example"
)

func main() {
	// We don't want to call unlock, as Python will not be initialised after
	// Main returns, and that will cause Unlock to panic. For the same reason,
	// we don't defer any Decref calls here.
	py.InitAndLockWithSignals()

	if err := example.Init(); err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	os.Exit(py.Main(os.Args))
}
