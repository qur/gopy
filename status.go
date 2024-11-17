package py

// #include "utils.h"
import "C"

import (
	"fmt"
)

type StatusExit struct {
	ExitCode int
}

func (s StatusExit) Error() string {
	return fmt.Sprintf("Python exited: %d", s.ExitCode)
}

func status2Err(status C.PyStatus) error {
	switch status._type {
	case C._PyStatus_TYPE_OK:
		return nil
	case C._PyStatus_TYPE_EXIT:
		return StatusExit{ExitCode: int(status.exitcode)}
	default:
		// TODO(jp3): implement
		panic("not implemented")
	}
}
