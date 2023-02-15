package py

import (
	"sync"
	"unsafe"
)

type ClassObject interface {
	Object
	setBase(base *BaseObject)
}

type ClassBaseObject struct {
	base *BaseObject
}

var _ ClassObject = (*ClassBaseObject)(nil)

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

func (c *ClassBaseObject) IsTrue() bool {
	return c.base.IsTrue()
}

func (c *ClassBaseObject) Not() bool {
	return c.base.Not()
}

func (c *ClassBaseObject) Free() {
	c.base.Free()
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
