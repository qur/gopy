package main

import "text/template"

var code = template.Must(template.New("code").Parse(`// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"fmt"
	"unsafe"
)

type {{ .type }} struct {
	abstractObject
	o C.Py{{ .type }}Object
}

// {{ .type }}Type is the Type object that represents the {{ .type }} type.
var {{ .type }}Type = (*Type)(unsafe.Pointer(&C.Py{{ .type }}_Type))

func {{ .ltype }}Check(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.{{ .ltype }}Check(c(obj)) != 0
}

func new{{ .type }}(obj *C.PyObject) *{{ .type }} {
	return (*{{ .type }})(unsafe.Pointer(obj))
}

{{ if or .funcs.mp_length .funcs.sq_length -}}
func ({{ .name }} *{{ .type }}) Size() int {
	ret := C.PyObject_Size(c({{ .name }}))
	if ret < 0 {
		return 0
	}
	return int(ret)
}
{{- end }}

{{ if .funcs.mp_subscript -}}
func ({{ .name }} *{{ .type }}) GetItemString(key string) (Object, error) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_GetItemString(c({{ .name }}), cKey)
	return obj2ObjErr(ret)
}

func ({{ .name }} *{{ .type }}) HasKey(key Object) bool {
	ret := C.PyMapping_HasKey(c({{ .name }}), c(key))
	return ret > 0
}

func ({{ .name }} *{{ .type }}) HasKeyString(key string) bool {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_HasKeyString(c({{ .name }}), cKey)
	return ret > 0
}
{{- end }}

{{ if .funcs.mp_ass_subscript -}}
func ({{ .name }} *{{ .type }}) DelItem(key Object) error {
	ret := C.PyMapping_DelItem(c({{ .name }}), c(key))
	return int2Err(ret)
}

func ({{ .name }} *{{ .type }}) DelItemString(key string) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_DelItemString(c({{ .name }}), cKey)
	return int2Err(ret)
}

func ({{ .name }} *{{ .type }}) SetItemString(key string, v Object) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_SetItemString(c({{ .name }}), cKey, c(v))
	return int2Err(ret)
}
{{- end }}

/*
{{ range $name, $set := .funcs -}}
{{ if $set -}}
{{ $name }} = {{ $set }}
{{ end -}}
{{ end -}}
*/

`))
