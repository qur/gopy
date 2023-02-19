package main

import "text/template"

var code = template.Must(template.New("code").Parse(`// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// {{ .type }} represents objects of the {{ .type }}Type (or Py{{ .type }}Type
// in the Python API) type.{{ if .funcs.mp_subscript }}
//
// This type implements the Mapping protocol.
{{- end }}{{ if .funcs.sq_item }}
//
// This type implements the Sequence protocol.
{{- end }}{{ if or .funcs.nb_index .funcs.nb_int .funcs.nb_float (eq .type "Complex") }}
//
// This type implements the Number protocol.
{{- end }}
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
		clearErr();
		return 0
	}
	return int(ret)
}

{{ end }}

{{- if .funcs.mp_subscript -}}
func ({{ .name }} *{{ .type }}) AsMapping() *MappingMethods {
	return (*MappingMethods)(unsafe.Pointer({{ .name }}.Base()))
}

func ({{ .name }} *{{ .type }}) GetItemString(key string) (Object, error) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_GetItemString(c({{ .name }}), cKey)
	return obj2ObjErr(ret)
}

func ({{ .name }} *{{ .type }}) HasKey(key Object) bool {
	ret := C.PyMapping_HasKey(c({{ .name }}), c(key))
	clearErr();
	return ret > 0
}

func ({{ .name }} *{{ .type }}) HasKeyString(key string) bool {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_HasKeyString(c({{ .name }}), cKey)
	clearErr();
	return ret > 0
}

{{ end }}

{{- if or .funcs.mp_ass_subscript .funcs.sq_ass_item -}}
func ({{ .name }} *{{ .type }}) DelItem(key Object) error {
	ret := C.PyObject_DelItem(c({{ .name }}), c(key))
	return int2Err(ret)
}

{{ end }}

{{- if .funcs.mp_ass_subscript -}}
func ({{ .name }} *{{ .type }}) DelItemString(key string) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyObject_DelItemString(c({{ .name }}), cKey)
	return int2Err(ret)
}

func ({{ .name }} *{{ .type }}) SetItemString(key string, v Object) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_SetItemString(c({{ .name }}), cKey, c(v))
	return int2Err(ret)
}

{{ end }}

{{- if .funcs.sq_item -}}
func ({{ .name }} *{{ .type }}) AsSequence() *SequenceMethods {
	return (*SequenceMethods)(unsafe.Pointer({{ .name }}.Base()))
}

func ({{ .name }} *{{ .type }}) GetIndex(idx int) (Object, error) {
	ret := C.PySequence_GetItem(c({{ .name }}), C.Py_ssize_t(idx))
	return obj2ObjErr(ret)
}

func ({{ .name }} *{{ .type }}) Count(obj Object) (int, error) {
	ret := C.PySequence_Count(c({{ .name }}), c(obj))
	return ssize_t2IntErr(ret)
}

func ({{ .name }} *{{ .type }}) Index(obj Object) (int, error) {
	ret := C.PySequence_Index(c({{ .name }}), c(obj))
	return ssize_t2IntErr(ret)
}

{{ end }}

{{- if and .funcs.sq_item (ne .type "List") -}}
func ({{ .name }} *{{ .type }}) List() (*List, error) {
	ret := C.PySequence_List(c({{ .name }}))
	if ret == nil {
		return nil, exception()
	}
	return newList(ret), nil
}

{{ end }}

{{- if and .funcs.sq_item (ne .type "Tuple") -}}
func ({{ .name }} *{{ .type }}) Tuple() (*Tuple, error) {
	ret := C.PySequence_Tuple(c({{ .name }}))
	if ret == nil {
		return nil, exception()
	}
	return newTuple(ret), nil
}

{{ end }}

{{- if .funcs.sq_ass_item -}}
func ({{ .name }} *{{ .type }}) SetIndex(idx int, obj Object) error {
	ret := C.PySequence_SetItem(c({{ .name }}), C.Py_ssize_t(idx), c(obj))
	return int2Err(ret)
}

func ({{ .name }} *{{ .type }}) DelIndex(idx int) error {
	ret := C.PySequence_DelItem(c({{ .name }}), C.Py_ssize_t(idx))
	return int2Err(ret)
}

{{ end }}

{{- if .funcs.sq_concat -}}
func ({{ .name }} *{{ .type }}) Concat(obj Object) (Object, error) {
	ret := C.PySequence_Concat(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.sq_inplace_concat -}}
func ({{ .name }} *{{ .type }}) InPlaceConcat(obj Object) (Object, error) {
	ret := C.PySequence_InPlaceConcat(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.sq_repeat -}}
func ({{ .name }} *{{ .type }}) Repeat(count int) (Object, error) {
	ret := C.PySequence_Repeat(c({{ .name }}), C.Py_ssize_t(count))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.sq_inplace_repeat -}}
func ({{ .name }} *{{ .type }}) InPlaceRepeat(count int) (Object, error) {
	ret := C.PySequence_InPlaceRepeat(c({{ .name }}), C.Py_ssize_t(count))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.sq_contains -}}
func ({{ .name }} *{{ .type }}) Contains(obj Object) (bool, error) {
	ret := C.PySequence_Contains(c({{ .name }}), c(obj))
	return int2BoolErr(ret)
}

{{ end }}

{{- if and .funcs.sq_item .funcs.mp_subscript -}}
func ({{ .name }} *{{ .type }}) GetSlice(start, end int) (Object, error) {
	ret := C.PySequence_GetSlice(c({{ .name }}), C.Py_ssize_t(start), C.Py_ssize_t(end))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if and .funcs.sq_item .funcs.mp_ass_subscript -}}
func ({{ .name }} *{{ .type }}) SetSlice(start, end int, obj Object) error {
	ret := C.PySequence_SetSlice(c({{ .name }}), C.Py_ssize_t(start), C.Py_ssize_t(end), c(obj))
	return int2Err(ret)
}

func ({{ .name }} *{{ .type }}) DelSlice(start, end int) error {
	ret := C.PySequence_DelSlice(c({{ .name }}), C.Py_ssize_t(start), C.Py_ssize_t(end))
	return int2Err(ret)
}

{{ end }}

{{- if or .funcs.nb_index .funcs.nb_int .funcs.nb_float (eq .type "Complex") -}}
func ({{ .name }} *{{ .type }}) AsNumber() *NumberMethods {
	return (*NumberMethods)(unsafe.Pointer({{ .name }}.Base()))
}

{{ end }}

/*
set fields:
{{- range $name, $set := .funcs -}}
{{ if $set }}
  {{ $name }}
{{- end -}}
{{ end }}
*/

`))
