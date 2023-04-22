package py

// #include "utils.h"
import "C"

import (
	"reflect"
	"sync"
	"unsafe"
)

// ClassObject is the interface that has to be implemented by types that are the
// Go instance type of a Class.
//
// This interface should be implemented by embedding a ClassBaseObject in a
// struct.
type ClassObject interface {
	Object
	getCBO() *ClassBaseObject
	setBase(base *BaseObject, class *Class)
}

// ClassBaseObject should be embedded into structs that want to implement
// ClassObject to be the instance type of a Class.
type ClassBaseObject struct {
	base  *BaseObject
	class *Class
	flags uint64
}

const (
	classBaseClear uint64 = 1 << iota
	classBaseTypeRefCountOdd
)

var _ ClassObject = (*ClassBaseObject)(nil)

var cboType = reflect.TypeOf((*ClassBaseObject)(nil)).Elem()

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (c *ClassBaseObject) Base() *BaseObject {
	return c.base
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (c *ClassBaseObject) Type() *Type {
	return c.base.Type()
}

// Decref decrements c's reference count, c may not be nil.
func (c *ClassBaseObject) Decref() {
	c.base.Decref()
}

// Incref increments c's reference count, c may not be nil.
func (c *ClassBaseObject) Incref() {
	c.base.Incref()
}

func (c *ClassBaseObject) raw() *C.PyObject {
	return (*C.PyObject)(unsafe.Pointer(c.base))
}

func (c *ClassBaseObject) getCBO() *ClassBaseObject {
	return c
}

func (c *ClassBaseObject) setBase(base *BaseObject, class *Class) {
	c.base = base
	c.class = class
}

// Super returns a Super object for the instance.
//
// This is equivalent to `super(type(c), c)` in Python.
//
// Return value: New Reference.
func (c *ClassBaseObject) Super() (*Super, error) {
	return c.class.Super(c)
}

// Class returns the *Class for this object.
func (c *ClassBaseObject) Class() *Class {
	return c.class
}

var (
	classObjLock sync.Mutex
	classObjMap  = map[*C.PyObject]map[*C.PyTypeObject]ClassObject{}
)

func registerClassObject(pyObj *C.PyObject, pyType *C.PyTypeObject, goObj ClassObject) {
	classObjLock.Lock()
	defer classObjLock.Unlock()

	typeMap, found := classObjMap[pyObj]
	if !found {
		typeMap = make(map[*C.PyTypeObject]ClassObject)
		classObjMap[pyObj] = typeMap
	}

	typeMap[pyType] = goObj
}

func getClassObject(pyObj *C.PyObject) ClassObject {
	if pyObj == nil {
		return nil
	}
	return getClassObjectByType(pyObj, pyObj.ob_type)
}

func getClassObjectByType(pyObj *C.PyObject, pyType *C.PyTypeObject) ClassObject {
	classObjLock.Lock()
	defer classObjLock.Unlock()

	return classObjMap[pyObj][pyType]
}

func clearClassObject(pyObj *C.PyObject) {
	classObjLock.Lock()
	defer classObjLock.Unlock()

	typeMap := classObjMap[pyObj]

	if typeMap == nil {
		return
	}

	delete(classObjMap, pyObj)
	for _, goObj := range typeMap {
		goObj.setBase(nil, nil)
	}
}

type ClassIteratorProtocol struct {
	base *ClassBaseObject
}

var cipType = reflect.TypeOf((*ClassIteratorProtocol)(nil)).Elem()

func (c *ClassIteratorProtocol) setCBO(b *ClassBaseObject) {
	c.base = b
}

func (c *ClassIteratorProtocol) AsIteratorMethods() *IteratorMethods {
	return (*IteratorMethods)(unsafe.Pointer(c.base.Base()))
}

type ClassNumberProtocol struct {
	base *ClassBaseObject
}

var cnpType = reflect.TypeOf((*ClassNumberProtocol)(nil)).Elem()

func (c *ClassNumberProtocol) setCBO(b *ClassBaseObject) {
	c.base = b
}

func (c *ClassNumberProtocol) AsNumberMethods() *NumberMethods {
	return (*NumberMethods)(unsafe.Pointer(c.base.Base()))
}

type ClassMappingProtocol struct {
	base *ClassBaseObject
}

var cmpType = reflect.TypeOf((*ClassMappingProtocol)(nil)).Elem()

func (c *ClassMappingProtocol) setCBO(b *ClassBaseObject) {
	c.base = b
}

func (c *ClassMappingProtocol) AsMappingMethods() *MappingMethods {
	return (*MappingMethods)(unsafe.Pointer(c.base.Base()))
}

type ClassSequenceProtocol struct {
	base *ClassBaseObject
}

var cspType = reflect.TypeOf((*ClassSequenceProtocol)(nil)).Elem()

func (c *ClassSequenceProtocol) setCBO(b *ClassBaseObject) {
	c.base = b
}

func (c *ClassSequenceProtocol) AsSequenceMethods() *SequenceMethods {
	return (*SequenceMethods)(unsafe.Pointer(c.base.Base()))
}

type classProtocol interface {
	setCBO(b *ClassBaseObject)
}
