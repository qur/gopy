// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef _GO_PYTHON_UTILS_H_
#define _GO_PYTHON_UTILS_H_

#include <ffi.h>
#include <python2.7/Python.h>

#include "ffi_type.h"

enum arg_type {
    at_pointer = 0,
    at_uint,
    at_uint8,
    at_uint16,
    at_uint32,
    at_uint64,
    at_int,
    at_int8,
    at_int16,
    at_int32,
    at_int64,
    at_float,
    at_double
};

typedef struct {
    ffi_type *type;
    void *value;
} ArgValue;

typedef struct {
    // Protective NULL pointer
    void *zero;

    // Standard Methods
    void *call;
    void *compare;
    void *dealloc;
    void *getattr;
    void *getattro;
    void *hash;
    void *init;
    void *iter;
    void *iternext;
    void *repr;
    void *richcmp;
    void *setattr;
    void *setattro;
    void *str;

    // Mapping Protocol
    int has_mp;
    void *mp_len;
    void *mp_get;
    void *mp_set;

    // Number Protocol
    int has_nb;
    void *nb_add;
    void *nb_subtract;
    void *nb_multiply;
    void *nb_divide;
    void *nb_remainder;
    void *nb_divmod;
    void *nb_power;
    void *nb_negative;
    void *nb_positive;
    void *nb_absolute;
    void *nb_nonzero;
    void *nb_invert;
    void *nb_lshift;
    void *nb_rshift;
    void *nb_and;
    void *nb_xor;
    void *nb_or;
    void *nb_int;
    void *nb_long;
    void *nb_float;
    void *nb_oct;
    void *nb_hex;
    void *nb_ip_add;
    void *nb_ip_subtract;
    void *nb_ip_multiply;
    void *nb_ip_divide;
    void *nb_ip_remainder;
    void *nb_ip_power;
    void *nb_ip_lshift;
    void *nb_ip_rshift;
    void *nb_ip_and;
    void *nb_ip_xor;
    void *nb_ip_or;
    void *nb_floordiv;
    void *nb_truediv;
    void *nb_ip_floordiv;
    void *nb_ip_truediv;
    void *nb_index;

    // Sequence Protocol
    int has_sq;
    void *sq_length;
    void *sq_concat;
    void *sq_repeat;
    void *sq_get;
    void *sq_set;
    void *sq_contains;
    void *sq_ip_concat;
    void *sq_ip_repeat;

    // The concrete Methods structs for the protocols (PyTypeObjects only store
    // a pointer).
    PyMappingMethods  mp_meth;
    PyNumberMethods   nb_meth;
    PySequenceMethods sq_meth;
} ClassContext;

extern PyMethodDef *newMethodDef();
extern void set_call_noargs(PyCFunction *f);
extern void set_call_args(PyCFunction *f);
extern void set_call_keywords(PyCFunction *f);

extern void incref(PyObject *o);
extern void decref(PyObject *o);
extern void xdecref(PyObject *o);
extern PyObject *get__Py_NoneStruct(void);
extern PyTypeObject *get_PyBaseObject_Type(void);
extern PyObject **get_PyExc_Exception(void);
extern void typeFree(PyTypeObject *type, PyObject *o);

extern int doParseTuple(PyObject *args, char *fmt, void *values[], int c);
extern int doParseTupleKwds(PyObject *args, PyObject *kwds, char *fmt, char *kwlist[], void *values[], int c);
extern PyObject *doBuildValue(char *fmt, ArgValue values[], int c);

extern PyObject *pyTrue(void);
extern PyObject *pyFalse(void);

extern PyObject *newMethod(char *name, void *func, int flags);
extern PyObject *newObjMember(int idx, char *doc);
extern PyObject *newNatMember(int idx, char *doc);
extern PyObject *newProperty(PyTypeObject *type, char *name, void *get, void *set);

extern void enableClassGc(PyTypeObject *type);

extern void setClassContext(PyTypeObject *type, ClassContext *ctxt);

extern PyTypeObject *getBasePyType(PyObject *o);

extern PyObject *compileFile(char *name);

extern void storeContext(PyTypeObject *t, ClassContext *c);
extern PyTypeObject *newType(void);
extern int typeReady(PyTypeObject *o);
extern ClassContext *newContext(void);
extern int setTypeAttr(PyTypeObject *tp, char *name, PyObject *o);
extern int doVisit(PyObject *o, void *v, void *a);

typedef struct _en {
    char *m;
    char *c;
} _en;

extern struct _en excName(PyObject *o);

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

#endif /* _GO_PYTHON_UTILS_H */
