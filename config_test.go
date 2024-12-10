package py_test

import (
	"errors"
	"log"
	"os"

	"gopython.xyz/py/v3"
)

func ExampleConfig_InitAndLock() {
	lock, err := py.PythonConfig(os.Args...).InitAndLock()
	if err != nil {
		var exit py.StatusExit
		if errors.As(err, &exit) && exit.ExitCode == 0 {
			// no error, but the init wants to exit anyway
			os.Exit(0)
		}
		log.Fatalf("Failed to init Python: %s", err)
	}
	defer lock.Finalize()
}
