// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "py.h"

#include "types.h"

extern void runtime_entersyscall(void) __asm__("syscall.Entersyscall");
extern void runtime_exitsyscall(void) __asm__("syscall.Exitsyscall");
extern void runtime_LockOSThread(void) __asm__("runtime.LockOSThread");
extern void runtime_UnlockOSThread(void) __asm__("runtime.UnlockOSThread");
extern int runtime_lockedOSThread(void);

extern int thread_locked(void) {
    return runtime_lockedOSThread();
}

void simple_cgocall(void (*f)(void*), void *a) {
    int lock = !runtime_lockedOSThread();
    if (lock) runtime_LockOSThread();
    runtime_entersyscall();
    f(a);
    runtime_exitsyscall();
    if (lock) runtime_UnlockOSThread();
}

void simple_cgocallback(void (*f)(void*), void (*f2)(void*), void *a) {
    runtime_exitsyscall();
    f(a);
    runtime_entersyscall();
}

void (*cgocallback)(void (*)(void*), void (*)(void*), void*) = simple_cgocallback;

#include "callbacks.c"
#include "utils.c"

extern void incref(PyObject *o) {
    Py_INCREF(o);
}

extern void decref(PyObject *o) {
    Py_DECREF(o);
}

extern void xdecref(PyObject *o) {
    Py_XDECREF(o);
}

extern PyObject *get__Py_NoneStruct(void) {
    return &_Py_NoneStruct;
}

extern PyTypeObject *get_PyBaseObject_Type(void) {
    return &PyBaseObject_Type;
}

extern PyObject **get_PyExc_Exception(void) {
    return &PyExc_Exception;
}

extern void typeFree(PyTypeObject *type, PyObject *o) {
    type->tp_free(o);
}

extern PyObject *pyTrue(void) {
    return Py_True;
}

extern PyObject *pyFalse(void) {
    return Py_False;
}

extern struct _dfr { FILE *f; int err; } do_fopen(const char *path, const char *mode) {
    struct _dfr s;
    s.err = 0;
    s.f = fopen(path, mode);
    if (!s.f) s.err = errno;
    return s;
}

// --- class.go helper functions ---
//
extern PyTypeObject *newType(void) {
    return calloc(1, sizeof(PyTypeObject));
}
extern int typeReady(PyTypeObject *o) {
    if (o->tp_new == NULL && o->tp_base == NULL) {
        o->tp_new = PyType_GenericNew;
    }
    if (o->tp_flags & Py_TPFLAGS_HAVE_GC) {
        enableClassGc(o);
    }
    return PyType_Ready(o);
}
extern ClassContext *newContext(void) {
    // We don't use tp_methods, and it is read when calling PyType_Ready
    // - so we use it to hide a classContext struct.  The classContext
    // starts with a NULL pointer just in case, so it looks like an
    // empty methods list if Python does try to process it.
    return calloc(1, sizeof(ClassContext));
}
extern void storeContext(PyTypeObject *t, ClassContext *c) {
    t->tp_methods = (void *)c;
}
extern int setTypeAttr(PyTypeObject *tp, char *name, PyObject *o) {
    return PyDict_SetItemString(tp->tp_dict, name, o);
}
extern int doVisit(PyObject *o, void *v, void *a) {
    visitproc visit = v;
    return visit(o, a);
}
//
// ----------------------------------

// --- object.go helper functions ---
//
extern int exceptionCheck(PyObject *obj) {
    return PyExceptionClass_Check(obj);
}
//
// ----------------------------------

// --- memory.go helper functions ---
//
size_t __PyObject_VAR_SIZE(PyObject *obj, Py_ssize_t n) {
    return _PyObject_VAR_SIZE((PyTypeObject *)obj, n);
}
void _PyObject_INIT(PyObject *obj, PyObject *typ) {
    PyObject_INIT(obj, (PyTypeObject *)typ);
}
void _PyObject_INIT_VAR(PyObject *obj, PyObject *typ, Py_ssize_t n) {
    PyObject_INIT_VAR(obj, (PyTypeObject *)typ, n);
}
void __PyObject_GC_TRACK(PyObject *obj) {
    _PyObject_GC_TRACK(obj);
}
void setGcRefs(PyGC_Head *g, Py_ssize_t refs) {
    g->gc.gc_refs = refs;
}
//
// ----------------------------------

// --- python.go helper functions ---
//
int enterRecursive(char *w) { return Py_EnterRecursiveCall(w); }
void leaveRecursive(void) { Py_LeaveRecursiveCall(); }
//
// ----------------------------------

// --- tuple.go helper functions ---
//
size_t tupleItemSize() { return sizeof(PyObject *); }
//
// ----------------------------------

// --- type.go helper functions ---
//
PyObject *typeAlloc(PyObject *t, Py_ssize_t n) {
    return ((PyTypeObject *)t)->tp_alloc((PyTypeObject *)t, n);
}
int typeInit(PyObject *t, PyObject *o, PyObject *a, PyObject *k) {
    return ((PyTypeObject *)t)->tp_init(o, a, k);
}
//
// ----------------------------------
