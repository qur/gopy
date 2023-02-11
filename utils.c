// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "utils.h"
#include "_cgo_export.h"


PyMethodDef *newMethodDef(void) {
    return (PyMethodDef *)calloc(1, sizeof(PyMethodDef));
}

void set_call_noargs(PyCFunction *f) {
    *f = (PyCFunction)callWithoutArgs;
}

void set_call_args(PyCFunction *f) {
    *f = (PyCFunction)callWithArgs;
}

void set_call_keywords(PyCFunction *f) {
    *f = (PyCFunction)callWithKeywords;
}

int doParseTuple(PyObject *args, char *fmt, void *values[], int c) {
    ffi_cif cif;
    ffi_status status;
    ffi_sarg result;
    int i;

    ffi_type **arg_types = calloc(c + 2, sizeof(ffi_type *));
    void **arg_values = calloc(c + 2, sizeof(void *));

    for (i = 0; i < c + 2; i++) {
        arg_types[i] = &ffi_type_pointer;
    }

    status = ffi_prep_cif(&cif, FFI_DEFAULT_ABI, c+2,
                          &ffi_type_sint, arg_types);

    if (status != FFI_OK) {
        PyErr_Format(PyExc_TypeError, "ffi_prep_cif failed: %d", status);
        return -1;
    }

    arg_values[0] = &args;
    arg_values[1] = &fmt;
    for (i = 0; i < c; i++) {
        arg_values[i+2] = &values[i];
    }

    ffi_call(&cif, FFI_FN(PyArg_ParseTuple), &result, arg_values);

    free(arg_types);
    free(arg_values);

    return (int)result;
}

int doParseTupleKwds(PyObject *args, PyObject *kwds, char *fmt, char *kwlist[], void *values[], int c) {
    ffi_cif cif;
    ffi_status status;
    ffi_sarg result;
    int i;

    ffi_type **arg_types = calloc(c + 4, sizeof(ffi_type *));
    void **arg_values = calloc(c + 4, sizeof(void *));

    for (i = 0; i < c + 4; i++) {
        arg_types[i] = &ffi_type_pointer;
    }

    status = ffi_prep_cif(&cif, FFI_DEFAULT_ABI, c+4,
                          &ffi_type_sint, arg_types);

    if (status != FFI_OK) {
        PyErr_Format(PyExc_TypeError, "ffi_prep_cif failed: %d", status);
        return -1;
    }

    arg_values[0] = &args;
    arg_values[1] = &kwds;
    arg_values[2] = &fmt;
    arg_values[3] = &kwlist;
    for (i = 0; i < c; i++) {
        arg_values[i+4] = &values[i];
    }

    ffi_call(&cif, FFI_FN(PyArg_ParseTupleAndKeywords), &result, arg_values);

    free(arg_types);
    free(arg_values);

    return (int)result;
}

PyObject *doBuildValue(char *fmt, ArgValue values[], int c) {
    ffi_cif cif;
    ffi_status status;
    PyObject *result;
    int i;

    ffi_type **arg_types = calloc(c + 1, sizeof(ffi_type *));
    if (arg_types == NULL) {
        return PyErr_NoMemory();
    }

    void **arg_values = calloc(c + 1, sizeof(void *));
    if (arg_values == NULL) {
        return PyErr_NoMemory();
    }

    arg_types[0] = &ffi_type_pointer;
    for (i = 0; i < c; i++) {
        arg_types[i+1] = values[i].type;
    }

    status = ffi_prep_cif(&cif, FFI_DEFAULT_ABI, c+1,
                          &ffi_type_pointer, arg_types);

    if (status != FFI_OK) {
        PyErr_Format(PyExc_TypeError, "ffi_prep_cif failed: %d", status);
        return NULL;
    }

    arg_values[0] = &fmt;
    for (i = 0; i < c; i++) {
        arg_values[i+1] = values[i].value;
    }

    ffi_call(&cif, FFI_FN(Py_BuildValue), &result, arg_values);

    free(arg_types);
    free(arg_values);

    return result;
}

typedef struct {
    PyObject_HEAD
    void *func;
    PyMethodDef meth;
} PyGoMethod;

static PyObject *method_get(PyGoMethod *self, PyObject *obj, PyObject *type) {
    PyObject *cap = PyCapsule_New(self->func, NULL, NULL);
    PyObject *o = PyTuple_Pack(2, obj, cap);
    PyObject *ret = PyCFunction_New(&self->meth, o);
    Py_DECREF(o);
    Py_DECREF(cap);
    return ret;
}

static PyTypeObject goMethodType = {
    PyObject_HEAD_INIT(NULL)
    0,                         /*ob_size*/
    "GoMethod",                /*tp_name*/
    sizeof(PyGoMethod),        /*tp_basicsize*/
    0,                         /*tp_itemsize*/
    0,                         /*tp_dealloc*/
    0,                         /*tp_print*/
    0,                         /*tp_getattr*/
    0,                         /*tp_setattr*/
    0,                         /*tp_compare*/
    0,                         /*tp_repr*/
    0,                         /*tp_as_number*/
    0,                         /*tp_as_sequence*/
    0,                         /*tp_as_mapping*/
    0,                         /*tp_hash */
    0,                         /*tp_call*/
    0,                         /*tp_str*/
    0,                         /*tp_getattro*/
    0,                         /*tp_setattro*/
    0,                         /*tp_as_buffer*/
    Py_TPFLAGS_DEFAULT,        /*tp_flags*/
    "GoMethod objects",        /* tp_doc */
    0,                         /* tp_traverse */
    0,                         /* tp_clear */
    0,                         /* tp_richcompare */
    0,                         /* tp_weaklistoffset */
    0,                         /* tp_iter */
    0,                         /* tp_iternext */
    0,                         /* tp_methods */
    0,                         /* tp_members */
    0,                         /* tp_getset */
    0,                         /* tp_base */
    0,                         /* tp_dict */
    (descrgetfunc)method_get,  /* tp_descr_get */
    0,                         /* tp_descr_set */
};
static int goMethodInit = 0;

PyObject *newMethod(char *name, void *func, int flags) {
    PyGoMethod *self;

    if (!goMethodInit) {
        goMethodType.tp_new = PyType_GenericNew;
        if (PyType_Ready(&goMethodType) < 0) return NULL;
        goMethodInit = 1;
    }

    self = (PyGoMethod *)goMethodType.tp_alloc(&goMethodType, 0);

    if (self != NULL) {
        self->func = func;
        self->meth.ml_name = name;
        switch (flags) {
            case METH_NOARGS:
                self->meth.ml_meth = (PyCFunction)goClassCallMethod;
                break;
            case METH_VARARGS:
                self->meth.ml_meth = (PyCFunction)goClassCallMethodArgs;
                break;
            case METH_VARARGS | METH_KEYWORDS:
                self->meth.ml_meth = (PyCFunction)goClassCallMethodKwds;
                break;
            default:
                fprintf(stderr, "Invalid method flags: %x\n", flags);
                return NULL;
        }
        self->meth.ml_flags = flags;
        self->meth.ml_doc = "";
    }

    return (PyObject *)self;
}

typedef struct {
    PyObject_HEAD
    char *doc;
    int field;
} PyGoObjMember;

static PyObject *objmemb_get(PyGoObjMember *self, PyObject *obj, PyObject *type) {
    return goClassObjGet(obj, self->field);
}

static int objmemb_set(PyGoObjMember *self, PyObject *obj, PyObject *value) {
    return goClassObjSet(obj, self->field, value);
}

static PyObject *objmemb_doc(PyGoObjMember *self, void *closure) {
    return PyString_FromString(self->doc);
}

static PyGetSetDef objmemb_getset[] = {
    {"__doc__", (getter)objmemb_doc},
    {NULL}
};

static PyTypeObject goObjMemberType = {
    PyObject_HEAD_INIT(NULL)
    0,                         /*ob_size*/
    "GoObjMember",             /*tp_name*/
    sizeof(PyGoObjMember),     /*tp_basicsize*/
    0,                         /*tp_itemsize*/
    0,                         /*tp_dealloc*/
    0,                         /*tp_print*/
    0,                         /*tp_getattr*/
    0,                         /*tp_setattr*/
    0,                         /*tp_compare*/
    0,                         /*tp_repr*/
    0,                         /*tp_as_number*/
    0,                         /*tp_as_sequence*/
    0,                         /*tp_as_mapping*/
    0,                         /*tp_hash */
    0,                         /*tp_call*/
    0,                         /*tp_str*/
    0,                         /*tp_getattro*/
    0,                         /*tp_setattro*/
    0,                         /*tp_as_buffer*/
    Py_TPFLAGS_DEFAULT,        /*tp_flags*/
    "GoObjMember objects",     /* tp_doc */
    0,                         /* tp_traverse */
    0,                         /* tp_clear */
    0,                         /* tp_richcompare */
    0,                         /* tp_weaklistoffset */
    0,                         /* tp_iter */
    0,                         /* tp_iternext */
    0,                         /* tp_methods */
    0,                         /* tp_members */
    objmemb_getset,            /* tp_getset */
    0,                         /* tp_base */
    0,                         /* tp_dict */
    (descrgetfunc)objmemb_get, /* tp_descr_get */
    (descrsetfunc)objmemb_set, /* tp_descr_set */
};
static int goObjMemberInit = 0;

PyObject *newObjMember(int idx, char *doc) {
    PyGoObjMember *self;

    if (!goObjMemberInit) {
        goObjMemberType.tp_new = PyType_GenericNew;
        if (PyType_Ready(&goObjMemberType) < 0) return NULL;
        goObjMemberInit = 1;
    }

    self = (PyGoObjMember *)goObjMemberType.tp_alloc(&goObjMemberType, 0);

    if (self != NULL) {
        self->doc   = doc;
        self->field = idx;
    }

    return (PyObject *)self;
}

typedef struct {
    PyObject_HEAD
    char *doc;
    int field;
} PyGoNatMember;

static PyObject *natmemb_get(PyGoNatMember *self, PyObject *obj, PyObject *type) {
    return goClassNatGet(obj, self->field);
}

static int natmemb_set(PyGoNatMember *self, PyObject *obj, PyObject *value) {
    return goClassNatSet(obj, self->field, value);
}

static PyObject *natmemb_doc(PyGoNatMember *self, void *closure) {
    return PyString_FromString(self->doc);
}

static PyGetSetDef natmemb_getset[] = {
    {"__doc__", (getter)natmemb_doc},
    {NULL}
};

static PyTypeObject goNatMemberType = {
    PyObject_HEAD_INIT(NULL)
    0,                         /*ob_size*/
    "GoNatMember",             /*tp_name*/
    sizeof(PyGoNatMember),     /*tp_basicsize*/
    0,                         /*tp_itemsize*/
    0,                         /*tp_dealloc*/
    0,                         /*tp_print*/
    0,                         /*tp_getattr*/
    0,                         /*tp_setattr*/
    0,                         /*tp_compare*/
    0,                         /*tp_repr*/
    0,                         /*tp_as_number*/
    0,                         /*tp_as_sequence*/
    0,                         /*tp_as_mapping*/
    0,                         /*tp_hash */
    0,                         /*tp_call*/
    0,                         /*tp_str*/
    0,                         /*tp_getattro*/
    0,                         /*tp_setattro*/
    0,                         /*tp_as_buffer*/
    Py_TPFLAGS_DEFAULT,        /*tp_flags*/
    "GoNatMember objects",     /* tp_doc */
    0,                         /* tp_traverse */
    0,                         /* tp_clear */
    0,                         /* tp_richcompare */
    0,                         /* tp_weaklistoffset */
    0,                         /* tp_iter */
    0,                         /* tp_iternext */
    0,                         /* tp_methods */
    0,                         /* tp_members */
    natmemb_getset,            /* tp_getset */
    0,                         /* tp_base */
    0,                         /* tp_dict */
    (descrgetfunc)natmemb_get, /* tp_descr_get */
    (descrsetfunc)natmemb_set, /* tp_descr_set */
};
static int goNatMemberInit = 0;

PyObject *newNatMember(int idx, char *doc) {
    PyGoNatMember *self;

    if (!goNatMemberInit) {
        goNatMemberType.tp_new = PyType_GenericNew;
        if (PyType_Ready(&goNatMemberType) < 0) return NULL;
        goNatMemberInit = 1;
    }

    self = (PyGoNatMember *)goNatMemberType.tp_alloc(&goNatMemberType, 0);

    if (self != NULL) {
        self->doc   = doc;
        self->field = idx;
    }

    return (PyObject *)self;
}

PyObject *newProperty(PyTypeObject *type, char *name, void *get, void *set) {
    PyGetSetDef *gsp = calloc(1, sizeof(PyGetSetDef));

    if (gsp == NULL) return NULL;

    gsp->name    = name;
    gsp->get     = NULL;
    gsp->set     = NULL;
    gsp->doc     = "";
    gsp->closure = PyTuple_Pack(2, Py_None, Py_None);

    if (get != NULL) {
        gsp->get = (getter)goClassGetProp;
        PyTuple_SetItem(gsp->closure, 0, PyCapsule_New(get, NULL, NULL));
    }

    if (set != NULL) {
        gsp->set = (setter)goClassSetProp;
        PyTuple_SetItem(gsp->closure, 1, PyCapsule_New(set, NULL, NULL));
    }

    return PyDescr_NewGetSet(type, gsp);
}

void enableClassGc(PyTypeObject *type) {
    type->tp_traverse = (traverseproc) goClassTraverse;
    type->tp_clear    = (inquiry)      goClassClear;
}

void overrideGenericAlloc(PyTypeObject *type) {
    if (type->tp_alloc == PyType_GenericAlloc) {
        type->tp_alloc   = (allocfunc)  goGenericAlloc;
        type->tp_free    = (freefunc)   goGenericFree;
    }
}

void setClassContext(PyTypeObject *type, ClassContext *ctxt) {
    ctxt->zero = NULL;

    type->tp_new     = (newfunc)    goClassNew;
    type->tp_alloc   = (allocfunc)  goGenericAlloc;
    type->tp_dealloc = (destructor) goClassDealloc;
    type->tp_free    = (freefunc)   goGenericFree;

    if (ctxt->call)     type->tp_call        = (ternaryfunc)  goClassCall;
    if (ctxt->compare)  type->tp_compare     = (cmpfunc)      goClassCompare;
    if (ctxt->getattr)  type->tp_getattr     = (getattrfunc)  goClassGetAttr;
    if (ctxt->getattro) type->tp_getattro    = (getattrofunc) goClassGetAttrObj;
    if (ctxt->hash)     type->tp_hash        = (hashfunc)     goClassHash;
    if (ctxt->init)     type->tp_init        = (initproc)     goClassInit;
    if (ctxt->iter)     type->tp_iter        = (getiterfunc)  goClassIter;
    if (ctxt->iternext) type->tp_iternext    = (iternextfunc) goClassIterNext;
    if (ctxt->repr)     type->tp_repr        = (reprfunc)     goClassRepr;
    if (ctxt->richcmp)  type->tp_richcompare = (richcmpfunc)  goClassRichCmp;
    if (ctxt->setattr)  type->tp_setattr     = (setattrfunc)  goClassSetAttr;
    if (ctxt->setattro) type->tp_setattro    = (setattrofunc) goClassSetAttrObj;
    if (ctxt->str)      type->tp_str         = (reprfunc)     goClassStr;

    if (ctxt->has_mp) {
        PyMappingMethods *m = &ctxt->mp_meth;
        type->tp_as_mapping = m;
        if (ctxt->mp_len) m->mp_length        = (lenfunc)       goClassMapLen;
        if (ctxt->mp_get) m->mp_subscript     = (binaryfunc)    goClassMapGet;
        if (ctxt->mp_set) m->mp_ass_subscript = (objobjargproc) goClassMapSet;
    }

    if (ctxt->has_nb) {
        PyNumberMethods *m = &ctxt->nb_meth;
        type->tp_as_number = m;
        if (ctxt->nb_add)          m->nb_add                  = (binaryfunc)  goClassNumAdd;
        if (ctxt->nb_subtract)     m->nb_subtract             = (binaryfunc)  goClassNumSubtract;
        if (ctxt->nb_multiply)     m->nb_multiply             = (binaryfunc)  goClassNumMultiply;
        if (ctxt->nb_divide)       m->nb_divide               = (binaryfunc)  goClassNumDivide;
        if (ctxt->nb_remainder)    m->nb_remainder            = (binaryfunc)  goClassNumRemainder;
        if (ctxt->nb_divmod)       m->nb_divmod               = (binaryfunc)  goClassNumDivmod;
        if (ctxt->nb_power)        m->nb_power                = (ternaryfunc) goClassNumPower;
        if (ctxt->nb_negative)     m->nb_negative             = (unaryfunc)   goClassNumNegative;
        if (ctxt->nb_positive)     m->nb_positive             = (unaryfunc)   goClassNumPositive;
        if (ctxt->nb_absolute)     m->nb_absolute             = (unaryfunc)   goClassNumAbsolute;
        if (ctxt->nb_nonzero)      m->nb_nonzero              = (inquiry)     goClassNumNonzero;
        if (ctxt->nb_invert)       m->nb_invert               = (unaryfunc)   goClassNumInvert;
        if (ctxt->nb_lshift)       m->nb_lshift               = (binaryfunc)  goClassNumLshift;
        if (ctxt->nb_rshift)       m->nb_rshift               = (binaryfunc)  goClassNumRshift;
        if (ctxt->nb_and)          m->nb_and                  = (binaryfunc)  goClassNumAnd;
        if (ctxt->nb_xor)          m->nb_xor                  = (binaryfunc)  goClassNumXor;
        if (ctxt->nb_or)           m->nb_or                   = (binaryfunc)  goClassNumOr;
        if (ctxt->nb_int)          m->nb_int                  = (unaryfunc)   goClassNumInt;
        if (ctxt->nb_long)         m->nb_long                 = (unaryfunc)   goClassNumLong;
        if (ctxt->nb_float)        m->nb_float                = (unaryfunc)   goClassNumFloat;
        if (ctxt->nb_oct)          m->nb_oct                  = (unaryfunc)   goClassNumOct;
        if (ctxt->nb_hex)          m->nb_hex                  = (unaryfunc)   goClassNumHex;
        if (ctxt->nb_ip_add)       m->nb_inplace_add          = (binaryfunc)  goClassNumInplaceAdd;
        if (ctxt->nb_ip_subtract)  m->nb_inplace_remainder    = (binaryfunc)  goClassNumInplaceSubtract;
        if (ctxt->nb_ip_multiply)  m->nb_inplace_multiply     = (binaryfunc)  goClassNumInplaceMultiply;
        if (ctxt->nb_ip_divide)    m->nb_inplace_divide       = (binaryfunc)  goClassNumInplaceDivide;
        if (ctxt->nb_ip_remainder) m->nb_inplace_remainder    = (binaryfunc)  goClassNumInplaceRemainder;
        if (ctxt->nb_ip_power)     m->nb_inplace_power        = (ternaryfunc) goClassNumInplacePower;
        if (ctxt->nb_ip_lshift)    m->nb_inplace_lshift       = (binaryfunc)  goClassNumInplaceLshift;
        if (ctxt->nb_ip_rshift)    m->nb_inplace_rshift       = (binaryfunc)  goClassNumInplaceRshift;
        if (ctxt->nb_ip_and)       m->nb_inplace_and          = (binaryfunc)  goClassNumInplaceAnd;
        if (ctxt->nb_ip_xor)       m->nb_inplace_xor          = (binaryfunc)  goClassNumInplaceXor;
        if (ctxt->nb_ip_or)        m->nb_inplace_or           = (binaryfunc)  goClassNumInplaceOr;
        if (ctxt->nb_floordiv)     m->nb_floor_divide         = (binaryfunc)  goClassNumFloorDivide;
        if (ctxt->nb_truediv)      m->nb_true_divide          = (binaryfunc)  goClassNumTrueDivide;
        if (ctxt->nb_ip_floordiv)  m->nb_inplace_floor_divide = (binaryfunc)  goClassNumInplaceFloorDivide;
        if (ctxt->nb_ip_truediv)   m->nb_inplace_true_divide  = (binaryfunc)  goClassNumInplaceTrueDivide;
        if (ctxt->nb_index)        m->nb_index                = (unaryfunc)   goClassNumIndex;
    }

    if (ctxt->has_sq) {
        PySequenceMethods *m = &ctxt->sq_meth;
        type->tp_as_sequence = m;
        if (ctxt->sq_length)    m->sq_length         = (lenfunc)         goClassSeqLength;
        if (ctxt->sq_concat)    m->sq_concat         = (binaryfunc)      goClassSeqConcat;
        if (ctxt->sq_repeat)    m->sq_repeat         = (ssizeargfunc)    goClassSeqRepeat;
        if (ctxt->sq_get)       m->sq_item           = (ssizeargfunc)    goClassSeqGetItem;
        if (ctxt->sq_set)       m->sq_ass_item       = (ssizeobjargproc) goClassSeqSetItem;
        if (ctxt->sq_contains)  m->sq_contains       = (objobjproc)      goClassSeqContains;
        if (ctxt->sq_ip_concat) m->sq_inplace_concat = (binaryfunc)      goClassSeqIpConcat;
        if (ctxt->sq_ip_repeat) m->sq_inplace_repeat = (ssizeargfunc)    goClassSeqIpRepeat;
    }
}

PyTypeObject *getBasePyType(PyObject *o) {
    if (o == NULL) return NULL;

    if (PyTuple_Check(o))     return &PyTuple_Type;
    if (PyDict_Check(o))      return &PyDict_Type;
    if (PyList_Check(o))      return &PyList_Type;
    if (PyString_Check(o))    return &PyString_Type;
    if (PyBool_Check(o))      return &PyBool_Type;
    if (PyLong_Check(o))      return &PyLong_Type;
    if (PyInt_Check(o))       return &PyInt_Type;
    if (PyFloat_Check(o))     return &PyFloat_Type;
    if (PyModule_Check(o))    return &PyModule_Type;
    if (PyType_Check(o))      return &PyType_Type;
    if (PyCode_Check(o))      return &PyCode_Type;
    if (PyCFunction_Check(o)) return &PyCFunction_Type;
    if (PyComplex_Check(o))   return &PyComplex_Type;
    if (PyFrozenSet_Check(o)) return &PyFrozenSet_Type;
    if (PySet_Check(o))       return &PySet_Type;
    if (PyFunction_Check(o))  return &PyFunction_Type;
    if (PyFrame_Check(o))     return &PyFrame_Type;
    if (PyMethod_Check(o))    return &PyMethod_Type;
    if (PyFile_Check(o))      return &PyFile_Type;
    if (PyUnicode_Check(o))   return &PyUnicode_Type;

    return o->ob_type;
}

PyObject *compileFile(char *name) {
    PyObject *f = PyFile_FromString(name, "rU");
    if (!f) return NULL;
    struct _node *n = PyParser_SimpleParseFile(PyFile_AsFile(f), name,
                                               Py_file_input);
    Py_DECREF(f);
    if (!n) return NULL;
    return (PyObject *)PyNode_Compile(n, name);
}

struct _en excName(PyObject *o) {
    struct _en en = { NULL, NULL };

    PyObject *m;

    if (!PyExceptionClass_Check(o)) {
        return en;
    }

    en.c = PyExceptionClass_Name(o);
    if (en.c != NULL) {
        char *d = strrchr(en.c, '.');
        if (d != NULL) en.c = d + 1;
    } else {
        en.c = "<unknown>";
    }

    m = PyObject_GetAttrString(o, "__module__");
    if (m == NULL) {
        en.m = "<unknown>";
        return en;
    }

    en.m = PyString_AsString(m);
    if (en.m != NULL && !strcmp(en.m, "exceptions")) {
        en.m = NULL;
    }

    Py_DECREF(m);

    return en;
}

int PyCode_GetNumFree_(PyCodeObject *o) { return PyCode_GetNumFree(o); }

void incref(PyObject *o) { Py_INCREF(o); }
void decref(PyObject *o) { Py_DECREF(o); }
void xincref(PyObject *o) { Py_INCREF(o); }
void xdecref(PyObject *o) { Py_XDECREF(o); }

void *pyTrue(void) { return Py_True; }
void *pyFalse(void) { return Py_False; }

int boolCheck(PyObject *o) { return PyBool_Check(o); }
int cFunctionCheck(PyObject *o) { return PyCFunction_Check(o); }
int codeCheck(PyObject *o) { return PyCode_Check(o); }
int complexCheck(PyObject *o) { return PyComplex_Check(o); }
int dictCheckE(PyObject *o) { return PyDict_CheckExact(o); }
int dictCheck(PyObject *o) { return PyDict_Check(o); }
int exceptionCheck(PyObject *o) { return PyExceptionClass_Check(o); }
int floatCheck(PyObject *o) { return PyFloat_Check(o); }
int frozenSetCheckE(PyObject *o) { return PyFrozenSet_CheckExact(o); }
int frozenSetCheck(PyObject *o) { return PyFrozenSet_Check(o); }
int functionCheck(PyObject *o) { return PyFunction_Check(o); }
int intCheck(PyObject *o) { return PyInt_Check(o); }
int listCheckE(PyObject *o) { return PyList_CheckExact(o); }
int listCheck(PyObject *o) { return PyList_Check(o); }
int longCheck(PyObject *o) { return PyLong_Check(o); }
int moduleCheckE(PyObject *o) { return PyModule_CheckExact(o); }
int moduleCheck(PyObject *o) { return PyModule_Check(o); }
int setCheck(PyObject *o) { return PySet_Check(o); }
int stringCheck(PyObject *o) { return PyString_Check(o); }
int tupleCheckE(PyObject *o) { return PyTuple_CheckExact(o); }
int tupleCheck(PyObject *o) { return PyTuple_Check(o); }
int typeCheckE(PyObject *o) { return PyType_CheckExact(o); }
int typeCheck(PyObject *o) { return PyType_Check(o); }
int frameCheck(PyObject *o) { return PyFrame_Check(o); }
int methodCheck(PyObject *o) { return PyMethod_Check(o); }
int fileCheck(PyObject *o) { return PyFile_Check(o); }
int unicodeCheck(PyObject *o) { return PyUnicode_Check(o); }

// --- abstract.go helper functions ---
//
void typeFree(PyTypeObject *type, PyObject *o) { type->tp_free(o); }
//
// ----------------------------------

// --- class.go helper functions ---
//
PyTypeObject *newType(void) {
    return calloc(1, sizeof(PyTypeObject));
}
int typeReady(PyTypeObject *o) {
    if (o->tp_new == NULL && o->tp_base == NULL) {
        o->tp_new = PyType_GenericNew;
    }
    if (o->tp_flags & Py_TPFLAGS_HAVE_GC) {
        enableClassGc(o);
    }
    return PyType_Ready(o);
}
ClassContext *newContext(void) {
    // We don't use tp_methods, and it is read when calling PyType_Ready
    // - so we use it to hide a classContext struct.  The classContext
    // starts with a NULL pointer just in case, so it looks like an
    // empty methods list if Python does try to process it.
    return calloc(1, sizeof(ClassContext));
}
void storeContext(PyTypeObject *t, ClassContext *c) {
    t->tp_methods = (void *)c;
}
int setTypeAttr(PyTypeObject *tp, char *name, PyObject *o) {
    return PyDict_SetItemString(tp->tp_dict, name, o);
}
int doVisit(PyObject *o, void *v, void *a) {
    visitproc visit = v;
    return visit(o, a);
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
size_t tupleItemSize(void) { return sizeof(PyObject *); }
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

// --- mapping.go helper functions ---
//
PyObject *PyMapping_Keys_(PyObject *o) {
    return PyMapping_Keys(o);
}
PyObject *PyMapping_Values_(PyObject *o) {
    return PyMapping_Values(o);
}
PyObject *PyMapping_Items_(PyObject *o) {
    return PyMapping_Items(o);
}
//
// ----------------------------------
