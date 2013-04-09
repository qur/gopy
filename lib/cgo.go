package py

// #cgo CFLAGS: -I/usr/local/Cellar/python3/3.3.0/Frameworks/Python.framework/Versions/3.3/include/python3.3m -I/usr/local/Cellar/python3/3.3.0/Frameworks/Python.framework/Versions/3.3/include/python3.3m
// #cgo LDFLAGS: -L/usr/local/Cellar/python3/3.3.0/Frameworks/Python.framework/Versions/3.3/lib/python3.3/config-3.3m -ldl -framework CoreFoundation -lpython3.3
// #cgo pkg-config: libffi
import "C"
