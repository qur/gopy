// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "utils.h"
#include "gopy_types.h"
#include "_cgo_export.h"


PyMethodDef *newMethodDef() {
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
    PyVarObject_HEAD_INIT(NULL, 0)
    "GoMethod",                /*tp_name*/
    sizeof(PyGoMethod),        /*tp_basicsize*/
    0,                         /*tp_itemsize*/
    0,                         /*tp_dealloc*/
    0,                         /*tp_print*/
    0,                         /*tp_getattr*/
    0,                         /*tp_setattr*/
    0,                         /*tp_reserved*/
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
    return PyUnicode_FromString(self->doc);
}

static PyGetSetDef objmemb_getset[] = {
    {"__doc__", (getter)objmemb_doc},
    {NULL}
};

static PyTypeObject goObjMemberType = {
    PyVarObject_HEAD_INIT(NULL, 0)
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
    return PyUnicode_FromString(self->doc);
}

static PyGetSetDef natmemb_getset[] = {
    {"__doc__", (getter)natmemb_doc},
    {NULL}
};

static PyTypeObject goNatMemberType = {
    PyVarObject_HEAD_INIT(NULL, 0)
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
        if (ctxt->nb_remainder)    m->nb_remainder            = (binaryfunc)  goClassNumRemainder;
        if (ctxt->nb_divmod)       m->nb_divmod               = (binaryfunc)  goClassNumDivmod;
        if (ctxt->nb_power)        m->nb_power                = (ternaryfunc) goClassNumPower;
        if (ctxt->nb_negative)     m->nb_negative             = (unaryfunc)   goClassNumNegative;
        if (ctxt->nb_positive)     m->nb_positive             = (unaryfunc)   goClassNumPositive;
        if (ctxt->nb_absolute)     m->nb_absolute             = (unaryfunc)   goClassNumAbsolute;
        if (ctxt->nb_invert)       m->nb_invert               = (unaryfunc)   goClassNumInvert;
        if (ctxt->nb_lshift)       m->nb_lshift               = (binaryfunc)  goClassNumLshift;
        if (ctxt->nb_rshift)       m->nb_rshift               = (binaryfunc)  goClassNumRshift;
        if (ctxt->nb_and)          m->nb_and                  = (binaryfunc)  goClassNumAnd;
        if (ctxt->nb_xor)          m->nb_xor                  = (binaryfunc)  goClassNumXor;
        if (ctxt->nb_or)           m->nb_or                   = (binaryfunc)  goClassNumOr;
        if (ctxt->nb_int)          m->nb_int                  = (unaryfunc)   goClassNumInt;
        if (ctxt->nb_float)        m->nb_float                = (unaryfunc)   goClassNumFloat;
        if (ctxt->nb_ip_add)       m->nb_inplace_add          = (binaryfunc)  goClassNumInplaceAdd;
        if (ctxt->nb_ip_subtract)  m->nb_inplace_remainder    = (binaryfunc)  goClassNumInplaceSubtract;
        if (ctxt->nb_ip_multiply)  m->nb_inplace_multiply     = (binaryfunc)  goClassNumInplaceMultiply;
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

PyTypeObject *getBasePyType(int t) {
    switch (t) {
//    case GoPyNilType:         return &PyNone_Type;
    case GoPyBaseObject_Type: return &PyBaseObject_Type;
    case GoPyBool_Type:       return &PyBool_Type;
    case GoPyCFunction_Type:  return &PyCFunction_Type;
    case GoPyCode_Type:       return &PyCode_Type;
    case GoPyComplex_Type:    return &PyComplex_Type;
    case GoPyDict_Type:       return &PyDict_Type;
    case GoPyFloat_Type:      return &PyFloat_Type;
    case GoPyFrozenSet_Type:  return &PyFrozenSet_Type;
    case GoPyFunction_Type:   return &PyFunction_Type;
    case GoPyList_Type:       return &PyList_Type;
    case GoPyLong_Type:       return &PyLong_Type;
    case GoPyModule_Type:     return &PyModule_Type;
//    case GoPyObject_Type:     return &PyObject_Type;
    case GoPySet_Type:        return &PySet_Type;
    case GoPyTuple_Type:      return &PyTuple_Type;
    case GoPyType_Type:       return &PyType_Type;
    case GoPyUnicode_Type:    return &PyUnicode_Type;
//    case GoPyOtherType:       return &PyOther_Type;
    default:                  return NULL;
    }
}

int getBaseGoPyType(PyObject *o) {
    if (o == NULL) return GoPyNilType;

    if (PyTuple_Check(o))     return GoPyTuple_Type;
    if (PyDict_Check(o))      return GoPyDict_Type;
    if (PyList_Check(o))      return GoPyList_Type;
    if (PyUnicode_Check(o))   return GoPyUnicode_Type;
    if (PyBool_Check(o))      return GoPyBool_Type;
    if (PyLong_Check(o))      return GoPyLong_Type;
    if (PyFloat_Check(o))     return GoPyFloat_Type;
    if (PyModule_Check(o))    return GoPyModule_Type;
    if (PyType_Check(o))      return GoPyType_Type;
    if (PyCode_Check(o))      return GoPyCode_Type;
    if (PyCFunction_Check(o)) return GoPyCFunction_Type;
    if (PyComplex_Check(o))   return GoPyComplex_Type;
    if (PyFrozenSet_Check(o)) return GoPyFrozenSet_Type;
    if (PySet_Check(o))       return GoPySet_Type;
    if (PyFunction_Check(o))  return GoPyFunction_Type;
    return GoPyOtherType;
}

PyObject *compileFile(char *name) {
    return NULL;
    // PyObject *f = PyFile_FromString(name, "rU");
    // if (!f) return NULL;
    // struct _node *n = PyParser_SimpleParseFile(PyFile_AsFile(f), name,
    //                                            Py_file_input);
    // Py_DECREF(f);
    // if (!n) return NULL;
    // return (PyObject *)PyNode_Compile(n, name);
}

struct _en excName(PyObject *o) {
    struct _en en = { NULL, NULL };

    PyObject *m, *u;

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
    u = PyUnicode_AsUTF8String(m);
    if (u == NULL) {
        Py_DECREF(m);
        en.m = "<unknown>";
        return en;
    }

    en.m = PyBytes_AsString(u);
    if (en.m != NULL && !strcmp(en.m, "exceptions")) {
        en.m = NULL;
    }

    Py_DECREF(m);

    return en;
}
size_t var_size(PyObject *obj, Py_ssize_t n) {
    return _PyObject_VAR_SIZE((PyTypeObject *)obj, n);
}
void GoPyObject_INIT(PyObject *obj, PyObject *typ){
    PyObject_INIT(obj, (PyTypeObject *)typ);
}
void GoPyObject_INIT_VAR(PyObject *obj, PyObject *typ, Py_ssize_t n){
    PyObject_INIT_VAR(obj, (PyTypeObject *)typ, n);
}
void GoPyObject_GC_Track(PyObject *obj){
    PyObject_GC_Track(obj);
}
void setGcRefs(PyGC_Head *g, Py_ssize_t refs){
    g->gc.gc_refs = refs;
}
extern void stub(void* addr);
void cinit() {
    stub(&sigaltstack);
    stub(&signal);
    stub(&sigaction);
#ifdef __APPLE__
    extern void __sigaltstack();
    stub(&__sigaltstack);
#endif
}
