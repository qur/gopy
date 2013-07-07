// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef _GO_PYTHON_PY_H_
#define _GO_PYTHON_PY_H_

#include <ffi.h>
#include <python2.7/Python.h>

#include "utils.h"
#include "ffi_type.h"

extern int thread_locked(void);

extern void incref(PyObject *o);
extern void decref(PyObject *o);
extern void xdecref(PyObject *o);
extern PyObject *get__Py_NoneStruct(void);
extern PyTypeObject *get_PyBaseObject_Type(void);
extern PyObject **get_PyExc_Exception(void);
extern void typeFree(PyTypeObject *type, PyObject *o);

extern void storeContext(PyTypeObject *t, ClassContext *c);
extern PyTypeObject *newType(void);
extern int typeReady(PyTypeObject *o);
extern ClassContext *newContext(void);
extern int setTypeAttr(PyTypeObject *tp, char *name, PyObject *o);
extern int doVisit(PyObject *o, void *v, void *a);

// --- memory.go helper functions ---
//
extern size_t __PyObject_VAR_SIZE(PyObject *obj, Py_ssize_t n);
extern void _PyObject_INIT(PyObject *obj, PyObject *typ);
extern void _PyObject_INIT_VAR(PyObject *obj, PyObject *typ, Py_ssize_t n);
extern void __PyObject_GC_TRACK(PyObject *obj);
extern void setGcRefs(PyGC_Head *g, Py_ssize_t refs);
//
// ----------------------------------

// --- python.go helper functions ---
//
extern int enterRecursive(char *w);
extern void leaveRecursive(void);
//
// ----------------------------------

// --- tuple.go helper functions ---
//
extern size_t tupleItemSize();
//
// ----------------------------------

// --- type.go helper functions ---
//
extern PyObject *typeAlloc(PyObject *t, Py_ssize_t n);
extern int typeInit(PyObject *t, PyObject *o, PyObject *a, PyObject *k);
//
// ----------------------------------

typedef struct _en _en;

#endif /* _GO_PYTHON_PY_H */
