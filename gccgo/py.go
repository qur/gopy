// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

import (
	"fmt"
	"unsafe"
)

type _PyTypeObject __typeobject

func _malloc(int) unsafe.Pointer __asm__ ("malloc")
func _free(unsafe.Pointer) __asm__ ("free")

func __fopen(*int8,*int8) (*_FILE,int) __asm__ ("do_fopen")

func _fopen(name, mode *int8) (*_FILE, error) {
	f, errno := __fopen(name, mode)
	if f == nil {
		return nil, fmt.Errorf("Failed to open %s: %d", name, errno)
	}
	return f, nil
}

func _GoString(data *int8) string {
	s := uintptr(unsafe.Pointer(data))
	b := make([]byte, 0, 1024)
	for i := uintptr(0); *(*byte)(unsafe.Pointer(s + i)) != 0; i++ {
		b = append(b, *(*byte)(unsafe.Pointer(s + i)))
	}
	return string(b)
}

func _CString(s string) *int8 {
	b := _malloc(len(s)+1)
	data := []byte(s)
	for i, c := range data {
		x := (*byte)(unsafe.Pointer(uintptr(b) + uintptr(i)))
		*x = c
	}
	x := (*byte)(unsafe.Pointer(uintptr(b) + uintptr(len(s))))
	*x = 0
	return (*int8)(b)
}
