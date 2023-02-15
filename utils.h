// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef _GO_PYTHON_UTILS_H_
#define _GO_PYTHON_UTILS_H_

#include <ffi.h>
#include <python3.11/Python.h>
#include <python3.11/frameobject.h>

#include "class_slots.h"

typedef struct {
  ffi_type *type;
  void *value;
} ArgValue;

extern int doParseTuple(PyObject *args, char *fmt, void *values[], int c);
extern int doParseTupleKwds(PyObject *args, PyObject *kwds, char *fmt,
                            char *kwlist[], void *values[], int c);
extern PyObject *doBuildValue(char *fmt, ArgValue values[], int c);

extern void enableClassGc(PyTypeObject *type);

extern PyTypeObject *getBasePyType(PyObject *o);

struct _en {
  const char *m;
  const char *c;
};

extern struct _en excName(PyObject *o);

extern int PyCode_GetNumFree_(PyCodeObject *o);

extern void incref(PyObject *o);
extern void decref(PyObject *o);
extern void xincref(PyObject *o);
extern void xdecref(PyObject *o);

extern void *pyTrue(void);
extern void *pyFalse(void);

extern int boolCheck(PyObject *o);
extern int cFunctionCheck(PyObject *o);
extern int codeCheck(PyObject *o);
extern int complexCheck(PyObject *o);
extern int dictCheckE(PyObject *o);
extern int dictCheck(PyObject *o);
extern int exceptionCheck(PyObject *o);
extern int floatCheck(PyObject *o);
extern int frozenSetCheckE(PyObject *o);
extern int frozenSetCheck(PyObject *o);
extern int functionCheck(PyObject *o);
// extern int intCheck(PyObject *o);
extern int listCheckE(PyObject *o);
extern int listCheck(PyObject *o);
extern int longCheck(PyObject *o);
extern int moduleCheckE(PyObject *o);
extern int moduleCheck(PyObject *o);
extern int setCheck(PyObject *o);
// extern int stringCheck(PyObject *o);
extern int tupleCheckE(PyObject *o);
extern int tupleCheck(PyObject *o);
extern int typeCheckE(PyObject *o);
extern int typeCheck(PyObject *o);
extern int frameCheck(PyObject *o);
extern int methodCheck(PyObject *o);
// extern int fileCheck(PyObject *o);
extern int unicodeCheck(PyObject *o);

// --- abstract.go helper functions ---
//
void typeFree(PyTypeObject *type, PyObject *o);
//
// ----------------------------------

// --- cfunction.go helper functions ---
//
extern PyMethodDef *newMethodDef(void);
extern void set_call_noargs(PyCFunction *f);
extern void set_call_args(PyCFunction *f);
extern void set_call_keywords(PyCFunction *f);
//
// ----------------------------------

// --- class.go helper functions ---
//
extern PyTypeObject *newType(void);
extern int typeReady(PyTypeObject *o);
extern void storeContext(PyTypeObject *t, ClassContext *c);
extern int setTypeAttr(PyTypeObject *tp, char *name, PyObject *o);
extern int doVisit(PyObject *o, void *v, void *a);
extern PyObject *newMethod(char *name, PyObject *func, int flags);
extern PyObject *newObjMember(int idx, char *doc);
extern PyObject *newNatMember(int idx, char *doc);
extern PyObject *newProperty(PyTypeObject *type, char *name, void *get,
                             void *set);
//
// ----------------------------------

// --- memory.go helper functions ---
//
extern size_t __PyObject_VAR_SIZE(PyObject *obj, Py_ssize_t n);
extern void _PyObject_INIT(PyObject *obj, PyObject *typ);
extern void _PyObject_INIT_VAR(PyObject *obj, PyObject *typ, Py_ssize_t n);
// extern void __PyObject_GC_TRACK(PyObject *obj);
// extern void setGcRefs(PyGC_Head *g, Py_ssize_t refs);
//
//  ----------------------------------

// --- memory.go helper functions ---
//
extern PyModuleDef *newModuleDef();
//
//  ----------------------------------

// --- python.go helper functions ---
//
extern int enterRecursive(char *w);
extern void leaveRecursive(void);
//
// ----------------------------------

// --- type.go helper functions ---
//
extern size_t tupleItemSize(void);
//
// ----------------------------------

// --- type.go helper functions ---
//
extern PyObject *typeAlloc(PyObject *t, Py_ssize_t n);
extern int typeInit(PyObject *t, PyObject *o, PyObject *a, PyObject *k);
//
// ----------------------------------

// --- mapping.go helper functions ---
//
extern PyObject *PyMapping_Keys_(PyObject *o);
extern PyObject *PyMapping_Values_(PyObject *o);
extern PyObject *PyMapping_Items_(PyObject *o);
//
// ----------------------------------

#endif /* _GO_PYTHON_UTILS_H */
