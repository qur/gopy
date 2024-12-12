package py

// #include "utils.h"
import "C"

import (
	"errors"
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
		if status._func != nil {
			f := C.GoString(status._func)
			msg := C.GoString(status.err_msg)

			return fmt.Errorf("%s: %s", f, msg)
		}

		msg := C.GoString(status.err_msg)

		return errors.New(msg)
	}
}
