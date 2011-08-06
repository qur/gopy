// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "utils.h"
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

PyObject *newMethod(char *name, void *func) {
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
        self->meth.ml_meth = (PyCFunction)callClassMethod;
        self->meth.ml_flags = METH_VARARGS | METH_KEYWORDS;
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

PyObject *newProperty(PyTypeObject *type, char *name, void *get, void *set) {
    PyGetSetDef *gsp = calloc(1, sizeof(PyGetSetDef));

    if (gsp == NULL) return NULL;

    gsp->name    = name;
    gsp->get     = NULL;
    gsp->set     = NULL;
    gsp->doc     = "";
    gsp->closure = PyTuple_New(2);

    if (get != NULL) {
        gsp->get = (getter)getClassProperty;
    }

    if (set != NULL) {
        gsp->set = (setter)setClassProperty;
    }

    PyTuple_SetItem(gsp->closure, 0, PyCapsule_New(get, NULL, NULL));
    PyTuple_SetItem(gsp->closure, 1, PyCapsule_New(set, NULL, NULL));

    return PyDescr_NewGetSet(type, gsp);
}

void setClassContext(PyTypeObject *type, ClassContext *ctxt) {
    ctxt->zero = NULL;

    type->tp_new     = (newfunc)    newGoClass;
    type->tp_dealloc = (destructor) deallocGoClass;

    if (ctxt->init)    type->tp_init    = (initproc)    initGoClass;
    if (ctxt->repr)    type->tp_repr    = (reprfunc)    reprGoClass;
    if (ctxt->str)     type->tp_str     = (reprfunc)    strGoClass;
    if (ctxt->call)    type->tp_call    = (ternaryfunc) callGoClass;
    if (ctxt->compare) type->tp_compare = (cmpfunc)     compareGoClass;

    if (ctxt->mp_len || ctxt->mp_get || ctxt->mp_set) {
        PyMappingMethods *m = &ctxt->mp_meth;
        type->tp_as_mapping = m;
        if (ctxt->mp_len) m->mp_length        = (lenfunc)       mapLenGoClass;
        if (ctxt->mp_get) m->mp_subscript     = (binaryfunc)    mapGetGoClass;
        if (ctxt->mp_set) m->mp_ass_subscript = (objobjargproc) mapSetGoClass;
    }
}

PyTypeObject *getBasePyType(PyObject *o) {
    if (o == NULL) return NULL;
    if (PyTuple_Check(o)) return &PyTuple_Type;
    if (PyDict_Check(o)) return &PyDict_Type;
    if (PyList_Check(o)) return &PyList_Type;
    if (PyString_Check(o)) return &PyString_Type;
    if (PyBool_Check(o)) return &PyBool_Type;
    if (PyLong_Check(o)) return &PyLong_Type;
    if (PyInt_Check(o)) return &PyInt_Type;
    if (PyFloat_Check(o)) return &PyFloat_Type;
    if (PyModule_Check(o)) return &PyModule_Type;
    if (PyType_Check(o)) return &PyType_Type;
    if (PyCode_Check(o)) return &PyCode_Type;
    if (PyCFunction_Check(o)) return &PyCFunction_Type;
    return o->ob_type;
}
