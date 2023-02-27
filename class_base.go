package py

// #include "utils.h"
import "C"

import (
	"reflect"
	"sync"
	"unsafe"
)

type ClassObject interface {
	Object
	getCBO() *ClassBaseObject
	setBase(base *BaseObject)
}

// ClassBaseObject should be embedded into structs that want to implement
// ClassObject to be the instance type of a Class.
type ClassBaseObject struct {
	base *BaseObject
}

var _ ClassObject = (*ClassBaseObject)(nil)

var cboType = reflect.TypeOf((*ClassBaseObject)(nil)).Elem()

func (c *ClassBaseObject) Base() *BaseObject {
	return c.base
}

func (c *ClassBaseObject) Type() *Type {
	return c.base.Type()
}

func (c *ClassBaseObject) Decref() {
	c.base.Decref()
}

func (c *ClassBaseObject) Incref() {
	c.base.Incref()
}

func (c *ClassBaseObject) raw() *C.PyObject {
	return (*C.PyObject)(unsafe.Pointer(c.base))
}

func (c *ClassBaseObject) getCBO() *ClassBaseObject {
	return c
}

func (c *ClassBaseObject) setBase(base *BaseObject) {
	c.base = base
}

var (
	classObjLock sync.Mutex
	classObjMap  = map[unsafe.Pointer]ClassObject{}
)

func registerClassObject(pyObj unsafe.Pointer, goObj ClassObject) {
	classObjLock.Lock()
	defer classObjLock.Unlock()

	classObjMap[pyObj] = goObj
}

func getClassObject(pyObj unsafe.Pointer) ClassObject {
	classObjLock.Lock()
	defer classObjLock.Unlock()

	return classObjMap[pyObj]
}

func clearClassObject(pyObj unsafe.Pointer) {
	classObjLock.Lock()
	defer classObjLock.Unlock()

	goObj := classObjMap[pyObj]

	if goObj == nil {
		return
	}

	delete(classObjMap, pyObj)
	goObj.setBase(nil)
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
