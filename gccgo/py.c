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

void simple_cgocallback(void (*f)(void*), void *a) {
    runtime_exitsyscall();
    f(a);
    runtime_entersyscall();
}

void (*cgocallback)(void (*)(void*), void*) = simple_cgocallback;

#include "callbacks.c"
#include "utils.c"

extern PyObject *get__Py_NoneStruct(void) {
    return &_Py_NoneStruct;
}

extern PyTypeObject *get_PyBaseObject_Type(void) {
    return &PyBaseObject_Type;
}

extern PyObject **get_PyExc_Exception(void) {
    return &PyExc_Exception;
}

extern struct _dfr { FILE *f; int err; } do_fopen(const char *path, const char *mode) {
    struct _dfr s;
    s.err = 0;
    s.f = fopen(path, mode);
    if (!s.f) s.err = errno;
    return s;
}
