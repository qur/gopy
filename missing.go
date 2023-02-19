package py

// #include "utils.h"
import "C"

// TODO(jp3): CallIter needs a proper implementation
func newCallIter(*C.PyObject) Object {
	return nil
}

// TODO(jp3): DictItems needs a proper implementation
func newDictItems(*C.PyObject) Object {
	return nil
}

// TODO(jp3): DictValues needs a proper implementation
func newDictValues(*C.PyObject) Object {
	return nil
}

// TODO(jp3): PickleBuffer needs a proper implementation
func newPickleBuffer(*C.PyObject) Object {
	return nil
}

// TODO(jp3): Range needs a proper implementation
func newRange(*C.PyObject) Object {
	return nil
}

// TODO(jp3): SeqIter needs a proper implementation
func newSeqIter(*C.PyObject) Object {
	return nil
}

// TODO(jp3): TraceBack needs a proper implementation
func newTraceBack(*C.PyObject) Object {
	return nil
}
