package main

import "text/template"

var code = template.Must(template.New("code").Parse(`// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

{{ if .settings.Doc -}}
{{ .settings.Doc -}}
{{ else -}}
// {{ .type }} represents objects of the {{ .type }}Type (or Py{{ .type }}_Type
// in the Python API) type.{{ if .funcs.mp_subscript }}
//
// This type implements the Mapping protocol.
{{- end }}{{ if .funcs.sq_item }}
//
// This type implements the Sequence protocol.
{{- end }}{{ if or .funcs.nb_index .funcs.nb_int .funcs.nb_float (eq .type "Complex") }}
//
// This type implements the Number protocol.
{{- end }}{{ if .funcs.tp_iternext }}
//
// This type implements the Iterator protocol.
{{- end }}{{ if .funcs.am_anext }}
//
// This type implements the AsyncIterator protocol.
{{- end }}
{{- end }}
type {{ .type }} struct {
	o C.{{ .ctype }}
}

var _ Object = (*{{ .type }})(nil)

{{ if .settings.Type -}}
// {{ .type }}Type is the Type object that represents the {{ .type }} type.
var {{ .type }}Type = newType(&C.Py{{ .type }}_Type)
{{- end }}

{{ if .settings.Check -}}
func {{ .ltype }}Check(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.{{ .ltype }}Check(c(obj)) != 0
}

// As{{ .type }} casts the given obj to a {{ .type }} (i.e. the underlying
// Python Object is the same, just the type is changed). If the value cannot be
// cast to a {{ .type }}, then nil is returned.
//
// Return value: Borrowed Reference.
func As{{ .type }}(obj Object) *{{ .type }} {
	if obj == nil {
		return nil
	}
	o := c(obj)
	if C.{{ .ltype }}Check(o) == 0 {
		return nil
	}
	return (*{{ .type }})(unsafe.Pointer(o))
}

{{ end }}

{{- if and .settings.New (ne .type "Type") -}}
func new{{ .type }}(obj *C.PyObject) *{{ .type }} {
	return (*{{ .type }})(unsafe.Pointer(obj))
}

{{ end }}

{{- if eq .type "Type" -}}
func new{{ .type }}(obj *C.PyTypeObject) *{{ .type }} {
	return (*{{ .type }})(unsafe.Pointer(obj))
}
{{- end }}

func ({{ .name }} *{{ .type }}) c() *C.{{ .ctype }} {
	return (*C.{{ .ctype }})(unsafe.Pointer({{ .name }}))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func ({{ .name }} *{{ .type }}) Base() *BaseObject {
{{- if eq .type "BaseObject" }}
	return {{ .name }}
{{- else }}
	return (*BaseObject)(unsafe.Pointer({{ .name }}))
{{- end }}
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func ({{ .name }} *{{ .type }}) Type() *Type {
	return newType(c({{ .name }}).ob_type)
}

{{ if .settings.Immortal -}}
// Decref is a nop, since {{ .type }} values are immortal.
func ({{ .name }} *{{ .type }}) Decref() {}

// Incref is a nop, since {{ .type }} values are immortal.
func ({{ .name }} *{{ .type }}) Incref() {}
{{- else }}
// Decref decrements {{ .name }}'s reference count, {{ .name }} may not be nil.
func ({{ .name }} *{{ .type }}) Decref() {
	obj := (*C.PyObject)(unsafe.Pointer({{ .name }}))
	refcnt := (*int)(unsafe.Pointer(&obj.anon0[0]))
	*refcnt--
	if *refcnt == 0 {
		C._Py_Dealloc(obj)
	}
}

// Incref increments {{ .name }}'s reference count, {{ .name }} may not be nil.
func ({{ .name }} *{{ .type }}) Incref() {
	refcnt := (*int)(unsafe.Pointer(&(*C.PyObject)(unsafe.Pointer({{ .name }})).anon0[0]))
	*refcnt++
}
{{- end }}

{{ if .funcs.tp_repr -}}
// Repr returns a String representation of "{{ .name }}". This is equivalent to the
// Python "repr({{ .name }})".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) Repr() (*Unicode, error) {
	ret := C.PyObject_Repr(c({{ .name }}))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret).(*Unicode), nil
}

{{ end }}

{{- if .funcs.tp_hash -}}
// Hash computes and returns the hash value of {{ .name }}. The equivalent
// Python is "hash({{ .name }})".
func ({{ .name }} *{{ .type }}) Hash() (int, error) {
	ret := C.PyObject_Hash(c({{ .name }}))
	if ret == -1 {
		return 0, exception()
	}
	return int(ret), nil
}

{{ end }}

{{- if .funcs.tp_call -}}
// Call calls {{ .name }} with the given args and kwds. kwds may be nil, args may not
// (an empty Tuple must be used if no arguments are wanted). Returns the result
// of the call, or an Error on failure.  This is equivalent to
// "{{ .name }}(*args, **kwds)" in Python.
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) Call(args *Tuple, kwds *Dict) (Object, error) {
	ret := C.PyObject_Call(c({{ .name }}), c(args), c(kwds))
	return obj2ObjErr(ret)
}

// CallGo calls {{ .name }} with the given args and kwds, either may be nil. Returns the
// result of the call, or an Error on failure.  This is equivalent to
// "{{ .name }}(*args, **kwds)" in Python.
//
// The values are converted to Objects using NewValue. A TypeError will be
// returned if a value cannot be converted.
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) CallGo(args []any, kwds map[string]any) (Object, error) {
	obj1, err := NewTupleFromValues(args...)
	if err != nil {
		return nil, err
	}
	defer obj1.Decref()
	obj2, err := NewDictFromValuesString(kwds)
	if err != nil {
		return nil, err
	}
	defer obj2.Decref()
	ret := C.PyObject_Call(c({{ .name }}), c(obj1), c(obj2))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.tp_str -}}
// Str returns a String representation of "{{ .name }}". This is equivalent to the
// Python "str({{ .name }})".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) Str() (*Unicode, error) {
	ret := C.PyObject_Str(c({{ .name }}))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret).(*Unicode), nil
}

{{ end }}

{{- if .funcs.tp_getattro -}}
// HasAttr returns true if "{{ .name }}" has the attribute "name".  This is equivalent
// to the Python "hasattr({{ .name }}, name)".
func ({{ .name }} *{{ .type }}) HasAttr(name Object) bool {
	ret := C.PyObject_HasAttr(c({{ .name }}), c(name))
	return ret == 1
}

// GetAttr returns the attribute of "{{ .name }}" with the name "name".  This is
// equivalent to the Python "{{ .name }}.name".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) GetAttr(name Object) (Object, error) {
	ret := C.PyObject_GetAttr(c({{ .name }}), c(name))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.tp_setattro -}}
// SetAttr sets the attribute of "{{ .name }}" with the name "name" to "value".  This is
// equivalent to the Python "{{ .name }}.name = value".
func ({{ .name }} *{{ .type }}) SetAttr(name, value Object) error {
	ret := C.PyObject_SetAttr(c({{ .name }}), c(name), c(value))
	return int2Err(ret)
}

// DelAttr deletes the attribute with the name "name" from "{{ .name }}".  This is
// equivalent to the Python "del {{ .name }}.name".
func ({{ .name }} *{{ .type }}) DelAttr(name, value Object) error {
	ret := C.PyObject_SetAttr(c({{ .name }}), c(name), nil)
	return int2Err(ret)
}

{{ end }}

{{- if .funcs.tp_richcompare -}}
// RichCompare compares "{{ .name }}" with "obj" using the specified operation (LE, GE
// etc.), and returns the result.  The equivalent Python is "{{ .name }} op obj", where
// op is the corresponding Python operator for op.
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) RichCompare(obj Object, op Op) (Object, error) {
	ret := C.PyObject_RichCompare(c({{ .name }}), c(obj), C.int(op))
	return obj2ObjErr(ret)
}

// RichCompare compares "{{ .name }}" with "obj" using the specified operation (LE, GE
// etc.), and returns true or false.  The equivalent Python is "{{ .name }} op obj",
// where op is the corresponding Python operator for op.
func ({{ .name }} *{{ .type }}) RichCompareBool(obj Object, op Op) (bool, error) {
	ret := C.PyObject_RichCompareBool(c({{ .name }}), c(obj), C.int(op))
	return int2BoolErr(ret)
}

{{ end }}

{{- if .funcs.tp_iter -}}
// Iter returns an Iterator that will iterate over the members of {{ .name }}.
func ({{ .name }} *{{ .type }}) Iter() (Iterator, error) {
	ret := C.PyObject_GetIter(c({{ .name }}))
	if ret == nil {
		return nil, exception()
	}
	return newIterator(ret)
}

{{ end }}


{{- if .funcs.tp_iternext -}}
// AsIteratorMethods returns a IteratorMethods instance that refers to the same
// underlying Python object as {{ .name }}.
//
// This method also means that {{ .type }} implements the Iterator interface.
func ({{ .name }} *{{ .type }}) AsIteratorMethods() *IteratorMethods {
	return (*IteratorMethods)(unsafe.Pointer({{ .name }}.Base()))
}

// Next is the Iterator method, it returns the next item from the Object being
// iterated. When the end is reached then both the Object and the error will be
// nil.
func ({{ .name }} *{{ .type }}) Next() (Object, error) {
	ret := C.PyIter_Next(c({{ .name }}))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.am_anext -}}
// AsAsyncIterator returns a AyncIteratorMethods instance that refers to the
// same underlying Python object as {{ .name }}.
//
// This method also means that {{ .type }} implements the AsyncIterator interface.
func ({{ .name }} *{{ .type }}) AsAsyncIterator() *AsyncIteratorMethods {
	return (*AsyncIteratorMethods)(unsafe.Pointer({{ .name }}.Base()))
}

{{ end }}

{{- if or .funcs.mp_length .funcs.sq_length -}}
// Size returns the size of {{ .name }}. The equivalent Python is "len({{ .name }})".
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
// AsMappingMethods returns a MappingMethods instance that refers to the same
// underlying Python object as {{ .name }}.
//
// This method also means that {{ .type }} implements the MappingProtocol interface.
func ({{ .name }} *{{ .type }}) AsMappingMethods() *MappingMethods {
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
// AsSequenceMethods returns a SequenceMethods instance that refers to the same
// underlying Python object as {{ .name }}.
//
// This method also means that {{ .type }} implements the SequenceProtocol
// interface.
func ({{ .name }} *{{ .type }}) AsSequenceMethods() *SequenceMethods {
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
// AsNumberMethods returns a NumberMethods instance that refers to the same
// underlying Python object as {{ .name }}.
//
// This method also means that {{ .type }} implements the NumberProtocol
// interface.
func ({{ .name }} *{{ .type }}) AsNumberMethods() *NumberMethods {
	return (*NumberMethods)(unsafe.Pointer({{ .name }}.Base()))
}

{{ end }}

{{- if .funcs.nb_add -}}
// Add returns the result of adding {{ .name }} and obj. The equivalent Python is
// "{{ .name }} + obj".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) Add(obj Object) (Object, error) {
	ret := C.PyNumber_Add(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_subtract -}}
// Subtract returns the result of subtracting obj from {{ .name }}. The equivalent Python
// is "{{ .name }} - obj".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) Subtract(obj Object) (Object, error) {
	ret := C.PyNumber_Subtract(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_multiply -}}
// Multiply returns the result of multiplying {{ .name }} by obj. The equivalent Python
// is "{{ .name }} * obj".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) Multiply(obj Object) (Object, error) {
	ret := C.PyNumber_Multiply(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_matrix_multiply -}}
func ({{ .name }} *{{ .type }}) MatrixMultiply(obj Object) (Object, error) {
	ret := C.PyNumber_MatrixMultiply(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_floor_divide -}}
// FloorDivide returns the floor of dividing {{ .name }} by obj. The equivalent Python is
// "{{ .name }} // obj".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) FloorDivide(obj Object) (Object, error) {
	ret := C.PyNumber_FloorDivide(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_true_divide -}}
// TrueDivide returns the approximate result of dividing {{ .name }} by obj. The result is
// approximate due to the limited representational accuracy of binary floating
// point numbers. The equivalent Python is "{{ .name }} / obj".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) TrueDivide(obj Object) (Object, error) {
	ret := C.PyNumber_TrueDivide(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_remainder -}}
// Remainder returns the remainder of dividing {{ .name }} by obj. The equivalent Python
// is "{{ .name }} % obj".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) Remainder(obj Object) (Object, error) {
	ret := C.PyNumber_Remainder(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_divmod -}}
// Divmod returns the result of the Python "divmod({{ .name }}, obj)".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) Divmod(obj Object) (Object, error) {
	ret := C.PyNumber_Divmod(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_power -}}
// Power returns the result of the Python "pow({{ .name }}, obj1, obj2)", where
// obj2 is optional.
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) Power(obj1, obj2 Object) (Object, error) {
	if obj2 == nil {
		None.Incref()
		obj2 = None
	}
	ret := C.PyNumber_Power(c({{ .name }}), c(obj1), c(obj2))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_negative -}}
// Negative returns the negation of {{ .name }}. The equivalent Python is "-{{ .name }}".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) Negative() (Object, error) {
	ret := C.PyNumber_Negative(c({{ .name }}))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_positive -}}
// Positive returns the positive of {{ .name }}. The equivalent Python is "+{{ .name }}".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) Positive() (Object, error) {
	ret := C.PyNumber_Positive(c({{ .name }}))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_absolute -}}
// Absolute returns the absolute value of {{ .name }}. The equivalent Python is "abs({{ .name }})".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) Absolute() (Object, error) {
	ret := C.PyNumber_Absolute(c({{ .name }}))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_invert -}}
// Invert returns the bitwise negation of {{ .name }}. The equivalent Python is "-{{ .name }}".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) Invert() (Object, error) {
	ret := C.PyNumber_Invert(c({{ .name }}))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_lshift -}}
// LShift returns the result of left shifting {{ .name }} by obj. The equivalent Python
// is "{{ .name }} << obj".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) LShift(obj Object) (Object, error) {
	ret := C.PyNumber_Lshift(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_rshift -}}
// RShift returns the result of right shifting {{ .name }} by obj. The equivalent Python
// is "{{ .name }} << obj".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) RShift(obj Object) (Object, error) {
	ret := C.PyNumber_Rshift(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_and -}}
// And returns the bitwise and of {{ .name }} and obj. The equivalent Python is
// "{{ .name }} & obj".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) And(obj Object) (Object, error) {
	ret := C.PyNumber_And(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_xor -}}
// Xor returns the bitwise xor of {{ .name }} and obj. The equivalent Python is
// "{{ .name }} ^ obj".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) Xor(obj Object) (Object, error) {
	ret := C.PyNumber_Xor(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_or -}}
// Or returns the bitwise or of {{ .name }} and obj. The equivalent Python is
// "{{ .name }} | obj".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) Or(obj Object) (Object, error) {
	ret := C.PyNumber_Or(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_inplace_add -}}
// InPlaceAdd returns the result of adding {{ .name }} and obj. This is done in place.
// The equivalent Python is "{{ .name }} += obj".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) InPlaceAdd(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceAdd(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_inplace_subtract -}}
// InPlaceSubtract returns the result of subtracting obj from {{ .name }}. This is done
// in place. The equivalent Python is "{{ .name }} -= obj".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) InPlaceSubtract(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceSubtract(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_inplace_multiply -}}
// InPlaceMultiply returns the result of multiplying {{ .name }} by obj. This is done in
// place. The equivalent Python is "{{ .name }} *= obj".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) InPlaceMultiply(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceMultiply(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_inplace_matrix_multiply -}}
func ({{ .name }} *{{ .type }}) InPlaceMatrixMultiply(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceMatrixMultiply(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_floor_inplace_divide -}}
// InPlaceFloorDivide returns the floor of dividing {{ .name }} by obj. This is done in
// place. The equivalent Python is "{{ .name }} //= obj".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) InPlaceFloorDivide(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceFloorDivide(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_inplace_true_divide -}}
// InPlaceTrueDivide returns the approximate result of dividing {{ .name }} by obj. This
// is done in place. The result is approximate due to the limited
// representational accuracy of binary floating point numbers. The equivalent
// Python is "{{ .name }} /= obj".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) InPlaceTrueDivide(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceTrueDivide(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_inplace_remainder -}}
// InPlaceRemainder returns the remainder of {{ .name }} divided by obj. This is done in
// place. The equivalent Python is "{{ .name }} %= obj".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) InPlaceRemainder(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceRemainder(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_inplace_power -}}
// InPlacePower returns the result of the Python "pow({{ .name }}, obj1, obj2)". This is
// done in place. If obj2 is None, then the Python "{{ .name }} **= obj" is also
// equivalent, if obj2 is not None, there is no equivalent in Python.
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) InPlacePower(obj1, obj2 Object) (Object, error) {
	ret := C.PyNumber_InPlacePower(c({{ .name }}), c(obj1), c(obj2))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_inplace_lshift -}}
// InPlaceLShift returns the result of left shifting {{ .name }} by obj. This is done in
// place. The equivalent Python is "{{ .name }} <<= obj".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) InPlaceLShift(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceLshift(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_inplace_rshift -}}
// InPlaceRShift returns the result of right shifting {{ .name }} by obj. This is done in
// place. The equivalent Python is "{{ .name }} >>= obj".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) InPlaceRShift(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceRshift(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_inplace_and -}}
// InPlaceAnd returns the bitwise and of {{ .name }} and obj. This is done in place. The
// equivalent Python is "{{ .name }} &= obj".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) InPlaceAnd(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceAnd(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_inplace_xor -}}
// InPlaceXor returns the bitwise xor of {{ .name }} and obj. This is done in place. The
// equivalent Python is "{{ .name }} ^= obj".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) InPlaceXor(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceXor(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if .funcs.nb_inplace_or -}}
// InPlaceOr returns the bitwise or of {{ .name }} and obj. This is done in place. The
// equivalent Python is "{{ .name }} |= obj".
//
// Return value: New Reference.
func ({{ .name }} *{{ .type }}) InPlaceOr(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceOr(c({{ .name }}), c(obj))
	return obj2ObjErr(ret)
}

{{ end }}

{{- if and .funcs.nb_int (ne .type "Long") -}}
func ({{ .name }} *{{ .type }}) Long() (*Long, error) {
	ret := C.PyNumber_Long(c({{ .name }}))
	return newLong(ret), exception()
}

{{ end }}

{{- if and .funcs.nb_float (ne .type "Float") -}}
func ({{ .name }} *{{ .type }}) Float() (*Float, error) {
	ret := C.PyNumber_Float(c({{ .name }}))
	return newFloat(ret), exception()
}

{{ end }}

{{- if .funcs.nb_index -}}
func ({{ .name }} *{{ .type }}) Index() (*Long, error) {
	ret := C.PyNumber_Index(c({{ .name }}))
	return newLong(ret), exception()
}

func ({{ .name }} *{{ .type }}) ToBase(base int) (*Long, error) {
	ret := C.PyNumber_ToBase(c({{ .name }}), C.int(base))
	return newLong(ret), exception()
}

func ({{ .name }} *{{ .type }}) AsInt(exc *ExceptionClass) (int, error) {
	ret := C.PyNumber_AsSsize_t(c({{ .name }}), c(exc))
	return int(ret), exception()
}

{{ end }}

{{- if .funcs.bf_getbuffer -}}
// AsBufferMethods returns a BufferMethods instance that refers to the same
// underlying Python object as {{ .name }}.
//
// This method also means that {{ .type }} implements the BufferProtocol
// interface.
func ({{ .name }} *{{ .type }}) AsBufferMethods() *BufferMethods {
	return (*BufferMethods)(unsafe.Pointer({{ .name }}.Base()))
}

func({{ .name }} *{{ .type }}) GetBuffer(flags BufferFlags) (*Buffer, error) {
	buf := newBuffer()
	ret := C.PyObject_GetBuffer(c({{ .name }}), buf.c(), C.int(flags))
	if ret < 0 {
		return nil, exception()
	}
	return buf, nil
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
