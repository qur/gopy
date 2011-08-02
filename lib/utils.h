// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef _GO_PYTHON_UTILS_H_
#define _GO_PYTHON_UTILS_H_

#include <ffi/ffi.h>
#include "Python.h"

typedef struct {
    ffi_type *type;
    void *value;
} ArgValue;

typedef struct {
    void *zero;
    void *new;
    void *init;
    void *repr;
    void *str;
    void *call;
    void *compare;
    void *mp_len;
    void *mp_get;
    void *mp_set;
    PyMappingMethods mp_meth;
} ClassContext;

extern PyMethodDef *newMethodDef();
extern void set_call_args(PyCFunction *f);
extern void set_call_keywords(PyCFunction *f);

extern int doParseTuple(PyObject *args, char *fmt, void *values[], int c);
extern PyObject *doBuildValue(char *fmt, ArgValue values[], int c);

extern PyObject *newMethod(char *name, void *func);
extern PyObject *newProperty(PyTypeObject *type, char *name, void *get, void *set);

extern void setClassContext(PyTypeObject *type, ClassContext *ctxt);

extern PyTypeObject *getBasePyType(PyObject *o);

#endif /* _GO_PYTHON_UTILS_H */
