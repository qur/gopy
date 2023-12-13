// Code generated by gen_types.py. DO NOT EDIT.

#include "utils.h"

#include "_cgo_export.h"

// ===============================================================

int aIterCheck(PyObject *o) { return PyAIter_Check(o); }
int anySetCheck(PyObject *o) { return PyAnySet_Check(o); }
int boolCheck(PyObject *o) { return PyBool_Check(o); }
int byteArrayCheck(PyObject *o) { return PyByteArray_Check(o); }
int bytesCheck(PyObject *o) { return PyBytes_Check(o); }
int cFunctionCheck(PyObject *o) { return PyCFunction_Check(o); }
int cMethodCheck(PyObject *o) { return PyCMethod_Check(o); }
int callIterCheck(PyObject *o) { return PyCallIter_Check(o); }
int callableCheck(PyObject *o) { return PyCallable_Check(o); }
int cellCheck(PyObject *o) { return PyCell_Check(o); }
int codeCheck(PyObject *o) { return PyCode_Check(o); }
int complexCheck(PyObject *o) { return PyComplex_Check(o); }
int dictCheck(PyObject *o) { return PyDict_Check(o); }
int dictItemsCheck(PyObject *o) { return PyDictItems_Check(o); }
int dictKeysCheck(PyObject *o) { return PyDictKeys_Check(o); }
int dictValuesCheck(PyObject *o) { return PyDictValues_Check(o); }
int dictViewSetCheck(PyObject *o) { return PyDictViewSet_Check(o); }
int exceptionClassCheck(PyObject *o) { return PyExceptionClass_Check(o); }
int exceptionInstanceCheck(PyObject *o) { return PyExceptionInstance_Check(o); }
int floatCheck(PyObject *o) { return PyFloat_Check(o); }
int frameCheck(PyObject *o) { return PyFrame_Check(o); }
int frozenSetCheck(PyObject *o) { return PyFrozenSet_Check(o); }
int functionCheck(PyObject *o) { return PyFunction_Check(o); }
int genCheck(PyObject *o) { return PyGen_Check(o); }
int indexCheck(PyObject *o) { return PyIndex_Check(o); }
int instanceMethodCheck(PyObject *o) { return PyInstanceMethod_Check(o); }
int iterCheck(PyObject *o) { return PyIter_Check(o); }
int listCheck(PyObject *o) { return PyList_Check(o); }
int longCheck(PyObject *o) { return PyLong_Check(o); }
int mappingCheck(PyObject *o) { return PyMapping_Check(o); }
int memoryViewCheck(PyObject *o) { return PyMemoryView_Check(o); }
int methodCheck(PyObject *o) { return PyMethod_Check(o); }
int moduleCheck(PyObject *o) { return PyModule_Check(o); }
int numberCheck(PyObject *o) { return PyNumber_Check(o); }
int oDictCheck(PyObject *o) { return PyODict_Check(o); }
int pickleBufferCheck(PyObject *o) { return PyPickleBuffer_Check(o); }
int rangeCheck(PyObject *o) { return PyRange_Check(o); }
int seqIterCheck(PyObject *o) { return PySeqIter_Check(o); }
int sequenceCheck(PyObject *o) { return PySequence_Check(o); }
int setCheck(PyObject *o) { return PySet_Check(o); }
int sliceCheck(PyObject *o) { return PySlice_Check(o); }
int traceBackCheck(PyObject *o) { return PyTraceBack_Check(o); }
int tupleCheck(PyObject *o) { return PyTuple_Check(o); }
int typeCheck(PyObject *o) { return PyType_Check(o); }
int unicodeCheck(PyObject *o) { return PyUnicode_Check(o); }
int weakrefCheck(PyObject *o) { return PyWeakref_Check(o); }
int anySetCheckE(PyObject *o) { return PyAnySet_CheckExact(o); }
int asyncGenCheckE(PyObject *o) { return PyAsyncGen_CheckExact(o); }
int asyncGenASendCheckE(PyObject *o) { return PyAsyncGenASend_CheckExact(o); }
int byteArrayCheckE(PyObject *o) { return PyByteArray_CheckExact(o); }
int bytesCheckE(PyObject *o) { return PyBytes_CheckExact(o); }
int cFunctionCheckE(PyObject *o) { return PyCFunction_CheckExact(o); }
int cMethodCheckE(PyObject *o) { return PyCMethod_CheckExact(o); }
int capsuleCheckE(PyObject *o) { return PyCapsule_CheckExact(o); }
int complexCheckE(PyObject *o) { return PyComplex_CheckExact(o); }
int contextCheckE(PyObject *o) { return PyContext_CheckExact(o); }
int contextTokenCheckE(PyObject *o) { return PyContextToken_CheckExact(o); }
int contextVarCheckE(PyObject *o) { return PyContextVar_CheckExact(o); }
int coroCheckE(PyObject *o) { return PyCoro_CheckExact(o); }
int dictCheckE(PyObject *o) { return PyDict_CheckExact(o); }
int floatCheckE(PyObject *o) { return PyFloat_CheckExact(o); }
int frozenSetCheckE(PyObject *o) { return PyFrozenSet_CheckExact(o); }
int genCheckE(PyObject *o) { return PyGen_CheckExact(o); }
int listCheckE(PyObject *o) { return PyList_CheckExact(o); }
int longCheckE(PyObject *o) { return PyLong_CheckExact(o); }
int moduleCheckE(PyObject *o) { return PyModule_CheckExact(o); }
int oDictCheckE(PyObject *o) { return PyODict_CheckExact(o); }
int setCheckE(PyObject *o) { return PySet_CheckExact(o); }
int tupleCheckE(PyObject *o) { return PyTuple_CheckExact(o); }
int typeCheckE(PyObject *o) { return PyType_CheckExact(o); }
int unicodeCheckE(PyObject *o) { return PyUnicode_CheckExact(o); }

// ===============================================================
// ===============================================================

PyTypeObject *getBasePyType(PyObject *o) {
  if (o == NULL) return NULL;

  if (PyBool_Check(o)) return &PyBool_Type;
  if (PyByteArray_Check(o)) return &PyByteArray_Type;
  if (PyBytes_Check(o)) return &PyBytes_Type;
  if (PyCFunction_Check(o)) return &PyCFunction_Type;
  if (PyCMethod_Check(o)) return &PyCMethod_Type;
  if (PyCallIter_Check(o)) return &PyCallIter_Type;
  if (PyCell_Check(o)) return &PyCell_Type;
  if (PyCode_Check(o)) return &PyCode_Type;
  if (PyComplex_Check(o)) return &PyComplex_Type;
  if (PyDict_Check(o)) return &PyDict_Type;
  if (PyDictItems_Check(o)) return &PyDictItems_Type;
  if (PyDictKeys_Check(o)) return &PyDictKeys_Type;
  if (PyDictValues_Check(o)) return &PyDictValues_Type;
  if (PyFloat_Check(o)) return &PyFloat_Type;
  if (PyFrame_Check(o)) return &PyFrame_Type;
  if (PyFrozenSet_Check(o)) return &PyFrozenSet_Type;
  if (PyFunction_Check(o)) return &PyFunction_Type;
  if (PyGen_Check(o)) return &PyGen_Type;
  if (PyInstanceMethod_Check(o)) return &PyInstanceMethod_Type;
  if (PyList_Check(o)) return &PyList_Type;
  if (PyLong_Check(o)) return &PyLong_Type;
  if (PyMemoryView_Check(o)) return &PyMemoryView_Type;
  if (PyMethod_Check(o)) return &PyMethod_Type;
  if (PyModule_Check(o)) return &PyModule_Type;
  if (PyODict_Check(o)) return &PyODict_Type;
  if (PyPickleBuffer_Check(o)) return &PyPickleBuffer_Type;
  if (PyRange_Check(o)) return &PyRange_Type;
  if (PySeqIter_Check(o)) return &PySeqIter_Type;
  if (PySet_Check(o)) return &PySet_Type;
  if (PySlice_Check(o)) return &PySlice_Type;
  if (PyTraceBack_Check(o)) return &PyTraceBack_Type;
  if (PyTuple_Check(o)) return &PyTuple_Type;
  if (PyType_Check(o)) return &PyType_Type;
  if (PyUnicode_Check(o)) return &PyUnicode_Type;

  return o->ob_type;
}

// ===============================================================
