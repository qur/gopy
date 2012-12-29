// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef _GO_PYTHON_UTILS_H_
#define _GO_PYTHON_UTILS_H_

#include <ffi.h>
#include <python2.7/Python.h>

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

extern int doParseTuple(PyObject *args, char *fmt, void *values[], int c);
extern int doParseTupleKwds(PyObject *args, PyObject *kwds, char *fmt, char *kwlist[], void *values[], int c);
extern PyObject *doBuildValue(char *fmt, ArgValue values[], int c);

extern PyObject *newMethod(char *name, void *func, int flags);
extern PyObject *newObjMember(int idx, char *doc);
extern PyObject *newNatMember(int idx, char *doc);
extern PyObject *newProperty(PyTypeObject *type, char *name, void *get, void *set);

extern void enableClassGc(PyTypeObject *type);

extern void overrideGenericAlloc(PyTypeObject *type);
extern void setClassContext(PyTypeObject *type, ClassContext *ctxt);

extern PyTypeObject *getBasePyType(PyObject *o);

extern PyObject *compileFile(char *name);

struct _en {
    char *m;
    char *c;
};

extern struct _en excName(PyObject *o);

#endif /* _GO_PYTHON_UTILS_H */
