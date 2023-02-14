#!/usr/bin/env python3

import sys


methods = {
    "tp": (None, None),
    "am": ("async", "PyAsyncMethods"),
    "nb": ("number", "PyNumberMethods"),
    "mp": ("mapping", "PyMappingMethods"),
    "sq": ("sequence", "PySequenceMethods"),
    "bf": ("buffer", "PyBufferProcs"),
}

slots = [
    # Standard Slots
    ("tp_repr", "reprfunc", "PyRepr", "() string"),
    ("tp_hash", "hashfunc", "PyHash", "() (uint32, error)"),
    ("tp_call", "ternaryfunc", "PyCall", "(*Tuple, *Dict) (Object, error)"),
    ("tp_str", "reprfunc", "PyStr", "() string"),
    ("tp_getattro", "getattrofunc", "PyGetAttr", ""),
    ("tp_setattro", "setattrofunc", "PySetAttr", ""),
    ("tp_richcompare", "richcmpfunc", "PyRichCompare", ""),
    ("tp_iter", "getiterfunc", "PyIter", "() (Object, error)"),
    ("tp_iternext", "iternextfunc", "PyIterNext", "() (Object, error)"),
    ("tp_descr_get", "descrgetfunc", "", ""),
    ("tp_descr_set", "descrsetfunc", "", ""),
    ("tp_init", "initproc", "PyInit", "(*Tuple, *Dict) error"),

    # Async Slots
    ("am_await", "unaryfunc", "PyAwait", "() (Object, error)"),
    ("am_aiter", "unaryfunc", "PyAsyncIter", "() (Object, error)"),
    ("am_anext", "unaryfunc", "PyAsyncNext", "() (Object, error)"),
    ("am_send", "sendfunc", "", ""),

    # Number Slots

    # Mapping Slots
    ("mp_length", "lenfunc", "PyMappingLen", "() int"),
    ("mp_subscript", "binaryfunc", "", ""),
    ("mp_ass_subscript", "objobjargproc", "", ""),

    # Sequence Slots
    ("sq_length", "lenfunc", "PyLen", "() int"),
    ("sq_concat", "binaryfunc", "", ""),
    ("sq_repeat", "ssizeargfunc", "", ""),
    ("sq_item", "ssizeargfunc", "", ""),
    ("sq_ass_item", "ssizeobjargproc", "", ""),
    ("sq_contains", "objobjproc", "", ""),
    ("sq_inplace_concat", "binaryfunc", "", ""),
    ("sq_inplace_repeat", "ssizeargfunc", "", ""),

    # Buffer Slots
    ("bf_getbuffer", "getbufferproc", "", ""),
    ("bf_releasebuffer", "releasebufferproc", "", ""),
]

callbacks = {
    "() string": (
        "(obj unsafe.Pointer) unsafe.Pointer",
        [
            '	s := C.CString(co.%(fn)s())',
            '	defer C.free(unsafe.Pointer(s))',
            '',
            '	return unsafe.Pointer(C.PyUnicode_FromString(s))',
        ],
    ),
    "() int": (
        "(obj unsafe.Pointer) C.Py_ssize_t",
        [
            '	return C.Py_ssize_t(co.%(fn)s())',
        ],
    ),
    "() (Object, error)": (
        "(obj unsafe.Pointer) unsafe.Pointer",
        [
            '	ret, err := co.%(fn)s()',
            '	if err != nil {',
            '		raise(err)',
            '		return nil',
            '	}',
            '',
            '	return unsafe.Pointer(c(ret))',
        ],
    ),
    "(*Tuple, *Dict) error": (
        "(obj, args, kwds unsafe.Pointer) unsafe.Pointer",
        [
            '	a := newTuple((*C.PyObject)(args))',
            '	k := newDict((*C.PyObject)(kwds))',
            '	ret, err := co.%(fn)s(a, k)',
            '	if err != nil {',
            '		raise(err)',
            '		return nil',
            '	}',
            '',
            '	return unsafe.Pointer(c(ret))',
        ],
    ),
}


def write_h_header(f):
    print("#ifndef _GO_PYTHON_SLOTS_H_", file=f)
    print("#define _GO_PYTHON_SLOTS_H_", file=f)
    print("", file=f)
    print("#include <ffi.h>", file=f)
    print("#include <python3.11/Python.h>", file=f)
    print("#include <python3.11/frameobject.h>", file=f)
    print("", file=f)


def write_h_footer(f):
    print("", file=f)
    print("#endif /* _GO_PYTHON_SLOTS_H */", file=f)


def write_constants(f):
    print("// ===============================================================", file=f)
    print("", file=f)
    group = 0
    current_prefix = "tp"
    flag = 1
    for (slot, _, _, _) in slots:
        prefix = slot.split("_")[0]
        if prefix != current_prefix:
            print(
                f"#define CLASS_HAS_{current_prefix.upper()} 0x{group:016x}ull", file=f)
            print("", file=f)
            group = 0
            current_prefix = prefix
        print(f"#define CLASS_HAS_{slot.upper()} 0x{flag:016x}ull", file=f)
        group |= flag
        flag <<= 1
    print(
        f"#define CLASS_HAS_{current_prefix.upper()} 0x{group:016x}ull", file=f)
    print("", file=f)
    print("// ===============================================================", file=f)


def write_ClassContext(f):
    print("// ===============================================================", file=f)
    print("", file=f)
    print("typedef struct {", file=f)
    print("  // Protective NULL pointer", file=f)
    print("  void *zero;", file=f)
    print("", file=f)
    print("  // The concrete Methods structs for the protocols (PyTypeObjects only store", file=f)
    print("  // a pointer).", file=f)
    for prefix, (method_set, method_type) in methods.items():
        if method_set is not None:
            print(f"  {method_type} {prefix}_meth;", file=f)
    print("} ClassContext;", file=f)
    print("", file=f)
    print("// ===============================================================", file=f)


def write_extern_setSlots(f):
    print("// ===============================================================", file=f)
    print("", file=f)
    print(
        "extern ClassContext *setSlots(PyTypeObject *type, uint64_t slotFlags);", file=f)
    print("", file=f)
    print("// ===============================================================", file=f)


def write_c_header(f):
    print('#include "utils.h"', file=f)
    print('', file=f)
    print('#include "_cgo_export.h"', file=f)
    print('', file=f)


def write_setSlots(f):
    print("// ===============================================================", file=f)
    print("", file=f)
    print(
        "ClassContext *setSlots(PyTypeObject *type, uint64_t slotFlags) {", file=f)
    print("  ClassContext *ctxt = calloc(1, sizeof(ClassContext));", file=f)
    print("  ctxt->zero = NULL;", file=f)
    print("", file=f)
    print("  type->tp_new = (newfunc)goClassNew;", file=f)
    print("  type->tp_dealloc = (destructor)goClassDealloc;", file=f)
    print("", file=f)
    print("  {", file=f)
    print("    PyTypeObject *m = type;", file=f)
    current_prefix = "tp"
    for (slot, pySig, _, _) in slots:
        prefix = slot.split("_")[0]
        if prefix != current_prefix:
            print("  }", file=f)
            print("", file=f)
            print(f"  if (slotFlags & CLASS_HAS_{prefix.upper()}) {{", file=f)
            (method_set, method_type) = methods.get(prefix, (None, None))
            if method_set is not None:
                print(f"    {method_type} *m = &ctxt->{prefix}_meth;", file=f)
                print(f"    type->tp_as_{method_set} = m;", file=f)
            current_prefix = prefix
        print(
            f"    if (slotFlags & CLASS_HAS_{slot.upper()}) m->{slot} = ({pySig})goClassSlot_{slot};", file=f)
    print("  }", file=f)
    print("", file=f)
    print("  return ctxt;", file=f)
    print("}", file=f)
    print("", file=f)
    print("// ===============================================================", file=f)


def write_go_header(f):
    print('package py', file=f)
    print('', file=f)
    print('// #include "utils.h"', file=f)
    print('import "C"', file=f)
    print('', file=f)
    print('import (', file=f)
    print('	"reflect"', file=f)
    print(')', file=f)
    print('', file=f)


def write_interfaces(f):
    print("// ===============================================================", file=f)
    print("", file=f)
    for (slot, _, fn, goSig) in slots:
        print(f"type {slot} interface {{", file=f)
        print(f"	{fn}{goSig}", file=f)
        print("}", file=f)
    print("", file=f)
    print("// ===============================================================", file=f)


def write_callbacks(f):
    print("// ===============================================================", file=f)
    print("", file=f)
    for (slot, _, fn, goSig) in slots:
        (cbSig, body) = callbacks.get(goSig, (None, []))
        if cbSig is None:
            continue
        print(f'//export goClassSlot_{slot}', file=f)
        print(
            f'func goClassSlot_{slot}{cbSig} {{', file=f)
        print(f'	co := newObject((*C.PyObject)(obj)).({slot})', file=f)
        print('', file=f)
        for line in body:
            print(line % {'fn': fn}, file=f)
        print('}', file=f)
    print("", file=f)
    print("// ===============================================================", file=f)


def write_slotMap(f):
    print("// ===============================================================", file=f)
    print("", file=f)
    print("var slotMap = map[C.uint64_t]reflect.Type{", file=f)
    for (slot, _, _, _) in slots:
        print(
            f"	C.CLASS_HAS_{slot.upper()}: reflect.TypeOf((*{slot})(nil)).Elem(),", file=f)
    print("}", file=f)
    print("", file=f)
    print("// ===============================================================", file=f)


def main():
    with open("slots.c", "w", encoding='utf-8') as output:
        write_c_header(output)
        write_setSlots(output)

    with open("slots.h", "w", encoding='utf-8') as output:
        write_h_header(output)
        write_constants(output)
        write_ClassContext(output)
        write_extern_setSlots(output)
        write_h_footer(output)

    with open("slots.go", "w", encoding='utf-8') as output:
        write_go_header(output)
        write_interfaces(output)
        write_callbacks(output)
        write_slotMap(output)

    output = sys.stdout


if __name__ == "__main__":
    main()
