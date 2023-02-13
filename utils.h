// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef _GO_PYTHON_UTILS_H_
#define _GO_PYTHON_UTILS_H_

#include <ffi.h>
#include <python3.11/Python.h>
#include <python3.11/frameobject.h>

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
  PyMappingMethods mp_meth;
  PyNumberMethods nb_meth;
  PySequenceMethods sq_meth;
} ClassContext;

extern int doParseTuple(PyObject *args, char *fmt, void *values[], int c);
extern int doParseTupleKwds(PyObject *args, PyObject *kwds, char *fmt,
                            char *kwlist[], void *values[], int c);
extern PyObject *doBuildValue(char *fmt, ArgValue values[], int c);

extern PyObject *newMethod(char *name, void *func, int flags);
extern PyObject *newObjMember(int idx, char *doc);
extern PyObject *newNatMember(int idx, char *doc);
extern PyObject *newProperty(PyTypeObject *type, char *name, void *get,
                             void *set);

extern void enableClassGc(PyTypeObject *type);

extern void setClassContext(PyTypeObject *type, ClassContext *ctxt);

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
#define CLASS_HAS_TP_REPR 0x0000000000000001ull
#define CLASS_HAS_TP_HASH 0x0000000000000002ull
#define CLASS_HAS_TP_CALL 0x0000000000000004ull
#define CLASS_HAS_TP_STR 0x0000000000000008ull
#define CLASS_HAS_TP_GETATTRO 0x0000000000000010ull
#define CLASS_HAS_TP_SETATTRO 0x0000000000000020ull
#define CLASS_HAS_TP_RICHCOMPARE 0x0000000000000040ull
#define CLASS_HAS_TP_ITER 0x0000000000000080ull
#define CLASS_HAS_TP_ITERNEXT 0x0000000000000100ull
#define CLASS_HAS_TP_DESCR_GET 0x0000000000000200ull
#define CLASS_HAS_TP_DESCR_SET 0x0000000000000400ull
#define CLASS_HAS_TP_INIT 0x0000000000000800ull

#define CLASS_HAS_AM_AWAIT 0x0000000000001000ull
#define CLASS_HAS_AM_AITER 0x0000000000002000ull
#define CLASS_HAS_AM_ANEXT 0x0000000000004000ull
#define CLASS_HAS_AM_SEND 0x0000000000008000ull

#define CLASS_HAS_NB_ADD 0x0000000000010000ull
#define CLASS_HAS_NB_INPLACE_ADD 0x0000000000020000ull
#define CLASS_HAS_NB_SUBTRACT 0x0000000000040000ull
#define CLASS_HAS_NB_INPLACE_SUBTRACT 0x0000000000080000ull
#define CLASS_HAS_NB_MULTIPLY 0x0000000000100000ull
#define CLASS_HAS_NB_INPLACE_MULTIPLY 0x0000000000200000ull
#define CLASS_HAS_NB_REMAINDER 0x0000000000400000ull
#define CLASS_HAS_NB_INPLACE_REMAINDER 0x0000000000800000ull
#define CLASS_HAS_NB_DIVMOD 0x0000000001000000ull
#define CLASS_HAS_NB_POWER 0x0000000002000000ull
#define CLASS_HAS_NB_INPLACE_POWER 0x0000000004000000ull
#define CLASS_HAS_NB_NEGATIVE 0x0000000008000000ull
#define CLASS_HAS_NB_POSITIVE 0x0000000010000000ull
#define CLASS_HAS_NB_ABSOLUTE 0x0000000020000000ull
#define CLASS_HAS_NB_BOOL 0x0000000040000000ull
#define CLASS_HAS_NB_INVERT 0x0000000080000000ull
#define CLASS_HAS_NB_LSHIFT 0x0000000100000000ull
#define CLASS_HAS_NB_INPLACE_LSHIFT 0x0000000200000000ull
#define CLASS_HAS_NB_RSHIFT 0x0000000400000000ull
#define CLASS_HAS_NB_INPLACE_RSHIFT 0x0000000800000000ull
#define CLASS_HAS_NB_AND 0x0000001000000000ull
#define CLASS_HAS_NB_INPLACE_AND 0x0000002000000000ull
#define CLASS_HAS_NB_XOR 0x0000004000000000ull
#define CLASS_HAS_NB_INPLACE_XOR 0x0000008000000000ull
#define CLASS_HAS_NB_OR 0x0000010000000000ull
#define CLASS_HAS_NB_INPLACE_OR 0x0000020000000000ull
#define CLASS_HAS_NB_INT 0x0000040000000000ull
#define CLASS_HAS_NB_FLOAT 0x0000080000000000ull
#define CLASS_HAS_NB_FLOOR_DIVIDE 0x0000100000000000ull
#define CLASS_HAS_NB_INPLACE_FLOOR_DIVIDE 0x0000200000000000ull
#define CLASS_HAS_NB_TRUE_DIVIDE 0x0000400000000000ull
#define CLASS_HAS_NB_INPLACE_TRUE_DIVIDE 0x0000800000000000ull
#define CLASS_HAS_NB_INDEX 0x0001000000000000ull
#define CLASS_HAS_NB_MATRIX_MULTIPLY 0x0002000000000000ull
#define CLASS_HAS_NB_LENGTH 0x0004000000000000ull
#define CLASS_HAS_NB_SUBSCRIPT 0x0008000000000000ull
#define CLASS_HAS_NB_ASS_SUBSCRIPT 0x0010000000000000ull

#define CLASS_HAS_SQ_LENGTH 0x0020000000000000ull
#define CLASS_HAS_SQ_CONCAT 0x0040000000000000ull
#define CLASS_HAS_SQ_REPEAT 0x0080000000000000ull
#define CLASS_HAS_SQ_ITEM 0x0100000000000000ull
#define CLASS_HAS_SQ_ASS_ITEM 0x0200000000000000ull
#define CLASS_HAS_SQ_CONTAINS 0x0400000000000000ull
#define CLASS_HAS_SQ_INPLACE_CONCAT 0x0800000000000000ull
#define CLASS_HAS_SQ_INPLACE_REPEAT 0x1000000000000000ull

#define CLASS_HAS_BF_GET_BUFFER 0x2000000000000000ull
#define CLASS_HAS_BF_RELEASE_BUFFER 0x4000000000000000ull

extern PyTypeObject *newType(void);
extern int typeReady(PyTypeObject *o);
extern ClassContext *newContext(void);
extern void storeContext(PyTypeObject *t, ClassContext *c);
extern int setTypeAttr(PyTypeObject *tp, char *name, PyObject *o);
extern int doVisit(PyObject *o, void *v, void *a);
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
