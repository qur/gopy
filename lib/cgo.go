package py

// #cgo CFLAGS: -Ic:/python33/include
// #cgo LDFLAGS: -Lc:/Python33/libs -Lc:/mingw/lib -ldl -lpython33
// #cgo pkg-config: libffi
import "C"