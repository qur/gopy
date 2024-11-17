package main

import (
	"log"
	"os"

	"gopython.xyz/py/v3"
)

func main() {
	if err := py.PythonPreConfig(os.Args...).PreInitialize(); err != nil {
		log.Fatalf("Failed to preinit Python: %s", err)
	}

	// We don't want to call unlock, as Python will not be initialised after
	// RunMain returns, and that will cause Unlock to panic. For the same
	// reason, we don't defer any Decref calls here.
	if _, err := py.PythonConfig(os.Args...).InitAndLock(); err != nil {
		log.Fatalf("Failed to init Python: %s", err)
	}

	if ret := py.RunMain(); ret != 0 {
		os.Exit(ret)
	}
}
