package main

// #cgo CFLAGS: -Werror
// #cgo LDFLAGS: -lpython3.11 -lm -lutil -ldl
// #cgo pkg-config: libffi
//
// #include <python3.11/Python.h>
import "C"

var types = map[string]any{
	"Bytes": C.PyBytes_Type,
}
