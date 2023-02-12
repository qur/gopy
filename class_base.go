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
	AbstractObject
	base *BaseObject
}

var _ ClassObject = (*ClassBaseObject)(nil)

func (c *ClassBaseObject) Base() *BaseObject {
	return c.base
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

	c := classObjMap[pyObj]

	if c == nil {
		return
	}

	delete(classObjMap, pyObj)
	c.setBase(nil)
}
