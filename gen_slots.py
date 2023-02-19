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
    ("tp_getattro", "getattrofunc", "PyGetAttr", "(Object) (Object, error)"),
    ("tp_setattro", "setattrofunc", "PySetAttr", "(Object, Object) error"),
    #    TODO(jp3): ("tp_richcompare", "richcmpfunc", "PyRichCompare", ""),
    ("tp_iter", "getiterfunc", "PyIter", "() (Object, error)"),
    ("tp_iternext", "iternextfunc", "PyIterNext", "() (Object, error)"),
    ("tp_descr_get", "descrgetfunc", "PyDescrGet",
     "(Object, Object) (Object, error)"),
    ("tp_descr_set", "descrsetfunc", "PyDescrSet", "(Object, Object) error"),
    ("tp_init", "initproc", "PyInit", "(*Tuple, *Dict) error"),

    # Async Slots
    ("am_await", "unaryfunc", "PyAwait", "() (Object, error)"),
    ("am_aiter", "unaryfunc", "PyAsyncIter", "() (Object, error)"),
    ("am_anext", "unaryfunc", "PyAsyncNext", "() (Object, error)"),
    #    TODO(jp3): ("am_send", "sendfunc", "", ""),

    # Number Slots
    ("nb_add", "binaryfunc", "PyAdd", "(Object) (Object, error)"),
    ("nb_inplace_add", "binaryfunc", "PyInplaceAdd", "(Object) (Object, error)"),
    ("nb_subtract", "binaryfunc", "PySubtract", "(Object) (Object, error)"),
    ("nb_inplace_subtract", "binaryfunc",
     "PyInplaceSubtract", "(Object) (Object, error)"),
    ("nb_multiply", "binaryfunc", "PyMultiply", "(Object) (Object, error)"),
    ("nb_inplace_multiply", "binaryfunc",
     "PyInplaceMultiply", "(Object) (Object, error)"),
    ("nb_remainder", "binaryfunc", "PyRemainder", "(Object) (Object, error)"),
    ("nb_inplace_remainder", "binaryfunc",
     "PyInplaceRemainder", "(Object) (Object, error)"),
    ("nb_divmod", "binaryfunc", "PyDivmod", "(Object) (Object, error)"),
    ("nb_power", "ternaryfunc", "PyPower", "(Object, Object) (Object, error)"),
    ("nb_inplace_power", "ternaryfunc", "PyInplacePower",
     "(Object, Object) (Object, error)"),
    ("nb_negative", "unaryfunc", "PyNegative", "() (Object, error)"),
    ("nb_positive", "unaryfunc", "PyPositive", "() (Object, error)"),
    ("nb_absolute", "unaryfunc", "PyAbsolute", "() (Object, error)"),
    ("nb_bool", "inquiry", "PyBool", "() (bool, error)"),
    ("nb_invert", "unaryfunc", "PyInvert", "() (Object, error)"),
    ("nb_lshift", "binaryfunc", "PyLShift", "(Object) (Object, error)"),
    ("nb_inplace_lshift", "binaryfunc",
     "PyInplaceLShift", "(Object) (Object, error)"),
    ("nb_rshift", "binaryfunc", "PyRShift", "(Object) (Object, error)"),
    ("nb_inplace_rshift", "binaryfunc",
     "PyInplaceRShift", "(Object) (Object, error)"),
    ("nb_and", "binaryfunc", "PyAnd", "(Object) (Object, error)"),
    ("nb_inplace_and", "binaryfunc", "PyInplaceAnd", "(Object) (Object, error)"),
    ("nb_xor", "binaryfunc", "PyXor", "(Object) (Object, error)"),
    ("nb_inplace_xor", "binaryfunc", "PyInplaceXor", "(Object) (Object, error)"),
    ("nb_or", "binaryfunc", "PyOr", "(Object) (Object, error)"),
    ("nb_inplace_or", "binaryfunc", "PyInplaceOr", "(Object) (Object, error)"),
    ("nb_int", "unaryfunc", "PyInt", "() (Object, error)"),
    ("nb_float", "unaryfunc", "PyFloat", "() (Object, error)"),
    ("nb_floor_divide", "binaryfunc", "PyFloorDivide", "(Object) (Object, error)"),
    ("nb_inplace_floor_divide", "binaryfunc",
     "PyInplaceFloorDivide", "(Object) (Object, error)"),
    ("nb_true_divide", "binaryfunc", "PyTrueDivide", "(Object) (Object, error)"),
    ("nb_inplace_true_divide", "binaryfunc",
     "PyInplaceTrueDivide", "(Object) (Object, error)"),
    ("nb_index", "unaryfunc", "PyIndex", "() (Object, error)"),
    ("nb_matrix_multiply", "binaryfunc",
     "PyMatrixMultiply", "(Object) (Object, error)"),
    ("nb_inplace_matrix_multiply", "binaryfunc",
     "PyInplaceMatrixMultiply", "(Object) (Object, error)"),

    # Mapping Slots
    ("mp_length", "lenfunc", "PyMappingLen", "() int"),
    ("mp_subscript", "binaryfunc", "PySubscript", "(Object) (Object, error)"),
    ("mp_ass_subscript", "objobjargproc",
     "PyAssSubscript", "(Object, Object) error"),

    # Sequence Slots
    ("sq_length", "lenfunc", "PyLen", "() int"),
    ("sq_concat", "binaryfunc", "PyConcat", "(Object) (Object, error)"),
    ("sq_repeat", "ssizeargfunc", "PyRepeat", "(Object, int) (Object, error)"),
    ("sq_item", "ssizeargfunc", "PyItem", "(Object, int) (Object, error)"),
    ("sq_ass_item", "ssizeobjargproc", "PyAssItem", "(Object, int, Object) error"),
    ("sq_contains", "objobjproc", "PyContains", "(Object) (bool, error)"),
    ("sq_inplace_concat", "binaryfunc",
     "PyInplaceConcat", "(Object) (Object, error)"),
    ("sq_inplace_repeat", "ssizeargfunc",
     "PyInplaceRepeat", "(Object, int) (Object, error)"),

    # Buffer Slots
    ("bf_getbuffer", "getbufferproc", "PyGetBuffer", "(Object, int) error"),
    ("bf_releasebuffer", "releasebufferproc", "PyReleaseBuffer", "(Object)"),
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
    "() (bool, error)": (
        "(obj unsafe.Pointer) int",
        [
            '	ret, err := co.%(fn)s()',
            '	if err != nil {',
            '		raise(err)',
            '		return -1',
            '	}',
            '',
            '	if ret {',
            '		return 1',
            '	}',
            '',
            '	return 0',
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
    "() (uint32, error)": (
        "(obj unsafe.Pointer) C.long",
        [
            '	ret, err := co.%(fn)s()',
            '	if err != nil {',
            '		raise(err)',
            '		return -1',
            '	} else if C.long(ret) == -1 {',
            '		return -2',
            '	}',
            '',
            '	return C.long(ret)',
        ],
    ),
    "(Object)": (
        "(obj, arg unsafe.Pointer)",
        [
            '	o := newObject((*C.PyObject)(arg))',
            '	co.%(fn)s(o)',
        ],
    ),
    "(Object) (Object, error)": (
        "(obj, arg unsafe.Pointer) unsafe.Pointer",
        [
            '	o := newObject((*C.PyObject)(arg))',
            '	ret, err := co.%(fn)s(o)',
            '	if err != nil {',
            '		raise(err)',
            '		return nil',
            '	}',
            '',
            '	return unsafe.Pointer(c(ret))',
        ],
    ),
    "(Object) (bool, error)": (
        "(obj, arg unsafe.Pointer) int",
        [
            '	o := newObject((*C.PyObject)(arg))',
            '	ret, err := co.%(fn)s(o)',
            '	if err != nil {',
            '		raise(err)',
            '		return -1',
            '	}',
            '',
            '	if ret {',
            '		return 1',
            '	}',
            '',
            '	return 0',
        ],
    ),
    "(Object, int) error": (
        "(obj, arg1 unsafe.Pointer, arg2 C.int) int",
        [
            '	o := newObject((*C.PyObject)(arg1))',
            '	if err := co.%(fn)s(o, int(arg2)); err != nil {',
            '		raise(err)',
            '		return -1',
            '	}',
            '',
            '	return 0',
        ],
    ),
    "(Object, int) (Object, error)": (
        "(obj, arg1 unsafe.Pointer, arg2 C.Py_ssize_t) unsafe.Pointer",
        [
            '	o := newObject((*C.PyObject)(arg1))',
            '	ret, err := co.%(fn)s(o, int(arg2))',
            '	if err != nil {',
            '		raise(err)',
            '		return nil',
            '	}',
            '',
            '	return unsafe.Pointer(c(ret))',
        ],
    ),
    "(Object, Object) (Object, error)": (
        "(obj, arg1, arg2 unsafe.Pointer) unsafe.Pointer",
        [
            '	o1 := newObject((*C.PyObject)(arg1))',
            '	o2 := newObject((*C.PyObject)(arg2))',
            '	ret, err := co.%(fn)s(o1, o2)',
            '	if err != nil {',
            '		raise(err)',
            '		return nil',
            '	}',
            '',
            '	return unsafe.Pointer(c(ret))',
        ],
    ),
    "(Object, Object) error": (
        "(obj, arg1, arg2 unsafe.Pointer) int",
        [
            '	o1 := newObject((*C.PyObject)(arg1))',
            '	o2 := newObject((*C.PyObject)(arg2))',
            '	if err := co.%(fn)s(o1, o2); err != nil {',
            '		raise(err)',
            '		return -1',
            '	}',
            '',
            '	return 0',
        ],
    ),
    "(*Tuple, *Dict) error": (
        "(obj, args, kwds unsafe.Pointer) int",
        [
            '	a := newTuple((*C.PyObject)(args))',
            '	k := newDict((*C.PyObject)(kwds))',
            '	if err := co.%(fn)s(a, k); err != nil {',
            '		raise(err)',
            '		return -1',
            '	}',
            '',
            '	return 0',
        ],
    ),
    "(*Tuple, *Dict) (Object, error)": (
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
    "(Object, int, Object) error": (
        "(obj, arg1 unsafe.Pointer, arg2 C.Py_ssize_t, arg3 unsafe.Pointer) C.int",
        [
            '	o1 := newObject((*C.PyObject)(arg1))',
            '	o3 := newObject((*C.PyObject)(arg1))',
            '	if err := co.%(fn)s(o1, int(arg2), o3); err != nil {',
            '		raise(err)',
            '		return -1',
            '	}',
            '',
            '	return 0',
        ],
    ),
}


def write_h_header(f):
    print("// Code generated by gen_slots.py. DO NOT EDIT.", file=f)
    print("", file=f)
    print("#ifndef _GO_PYTHON_CLASS_SLOTS_H_", file=f)
    print("#define _GO_PYTHON_CLASS_SLOTS_H_", file=f)
    print("", file=f)
    print("#include <ffi.h>", file=f)
    print("#include <python3.11/Python.h>", file=f)
    print("#include <python3.11/frameobject.h>", file=f)
    print("", file=f)


def write_h_footer(f):
    print("", file=f)
    print("#endif /* _GO_PYTHON_CLASS_SLOTS_H */", file=f)


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
    print('// Code generated by gen_slots.py. DO NOT EDIT.', file=f)
    print('', file=f)
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
    print('// Code generated by gen_slots.py. DO NOT EDIT.', file=f)
    print('', file=f)
    print('// This file is automatically generated.  To regenerate:', file=f)
    print('//   go generate ./...', file=f)
    print('', file=f)
    print('package py', file=f)
    print('', file=f)
    print('// #include "utils.h"', file=f)
    print('import "C"', file=f)
    print('', file=f)
    print('import (', file=f)
    print('	"reflect"', file=f)
    print('	"unsafe"', file=f)
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
            print(f"WARNING: unknown callback signature: {goSig}")
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
    if len(slots) > 64:
        print("TOO MANY SLOTS!", file=sys.stderr)
        print(
            f"Can only fit 64 slots in a uint64_t, we have {len(slots)}", file=sys.stderr)
        sys.exit(1)

    with open("class_slots.c", "w", encoding='utf-8') as output:
        write_c_header(output)
        write_setSlots(output)

    with open("class_slots.h", "w", encoding='utf-8') as output:
        write_h_header(output)
        write_constants(output)
        write_ClassContext(output)
        write_extern_setSlots(output)
        write_h_footer(output)

    with open("class_slots.go", "w", encoding='utf-8') as output:
        write_go_header(output)
        write_interfaces(output)
        write_callbacks(output)
        write_slotMap(output)

    output = sys.stdout


if __name__ == "__main__":
    main()
