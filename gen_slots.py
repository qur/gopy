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
    ("tp_repr", "reprfunc", "Repr", "() (Object, error)"),
    ("tp_hash", "hashfunc", "Hash", "() (int, error)"),
    ("tp_call", "ternaryfunc", "Call", "(*Tuple, *Dict) (Object, error)"),
    ("tp_str", "reprfunc", "Str", "() (Object, error)"),
    ("tp_getattro", "getattrofunc", "GetAttr", "(Object) (Object, error)"),
    ("tp_setattro", "setattrofunc", "SetAttr", "(Object, Object) error"),
    ("tp_richcompare", "richcmpfunc", "RichCompare", "(Object, Op) (Object, error)"),
    ("tp_iter", "getiterfunc", "Iter", "() (Iterator, error)"),
    ("tp_iternext", "iternextfunc", "Next", "() (Object, error)"),
    ("tp_descr_get", "descrgetfunc", "DescrGet",
     "(Object, Object) (Object, error)"),
    ("tp_descr_set", "descrsetfunc", "DescrSet", "(Object, Object) error"),
    ("tp_init", "initproc", "Init", "(*Tuple, *Dict) error"),

    # Async Slots
    ("am_await", "unaryfunc", "Await", "() (Object, error)"),
    ("am_aiter", "unaryfunc", "AsyncIter", "() (Object, error)"),
    ("am_anext", "unaryfunc", "AsyncNext", "() (Object, error)"),
    ("am_send", "sendfunc", "AsyncSend", "(Object) (Object, SendResult, error)"),

    # Number Slots
    ("nb_add", "binaryfunc", "Add", "(Object) (Object, error)"),
    ("nb_inplace_add", "binaryfunc", "InPlaceAdd", "(Object) (Object, error)"),
    ("nb_subtract", "binaryfunc", "Subtract", "(Object) (Object, error)"),
    ("nb_inplace_subtract", "binaryfunc",
     "InPlaceSubtract", "(Object) (Object, error)"),
    ("nb_multiply", "binaryfunc", "Multiply", "(Object) (Object, error)"),
    ("nb_inplace_multiply", "binaryfunc",
     "InPlaceMultiply", "(Object) (Object, error)"),
    ("nb_remainder", "binaryfunc", "Remainder", "(Object) (Object, error)"),
    ("nb_inplace_remainder", "binaryfunc",
     "InPlaceRemainder", "(Object) (Object, error)"),
    ("nb_divmod", "binaryfunc", "Divmod", "(Object) (Object, error)"),
    ("nb_power", "ternaryfunc", "Power", "(Object, Object) (Object, error)"),
    ("nb_inplace_power", "ternaryfunc", "InPlacePower",
     "(Object, Object) (Object, error)"),
    ("nb_negative", "unaryfunc", "Negative", "() (Object, error)"),
    ("nb_positive", "unaryfunc", "Positive", "() (Object, error)"),
    ("nb_absolute", "unaryfunc", "Absolute", "() (Object, error)"),
    ("nb_bool", "inquiry", "Bool", "() (bool, error)"),
    ("nb_invert", "unaryfunc", "Invert", "() (Object, error)"),
    ("nb_lshift", "binaryfunc", "LShift", "(Object) (Object, error)"),
    ("nb_inplace_lshift", "binaryfunc",
     "InPlaceLShift", "(Object) (Object, error)"),
    ("nb_rshift", "binaryfunc", "RShift", "(Object) (Object, error)"),
    ("nb_inplace_rshift", "binaryfunc",
     "InPlaceRShift", "(Object) (Object, error)"),
    ("nb_and", "binaryfunc", "And", "(Object) (Object, error)"),
    ("nb_inplace_and", "binaryfunc", "InPlaceAnd", "(Object) (Object, error)"),
    ("nb_xor", "binaryfunc", "Xor", "(Object) (Object, error)"),
    ("nb_inplace_xor", "binaryfunc", "PyInPlaceXor", "(Object) (Object, error)"),
    ("nb_or", "binaryfunc", "Or", "(Object) (Object, error)"),
    ("nb_inplace_or", "binaryfunc", "PyInPlaceOr", "(Object) (Object, error)"),
    ("nb_int", "unaryfunc", "Long", "() (*Long, error)"),
    ("nb_float", "unaryfunc", "Float", "() (*Float, error)"),
    ("nb_floor_divide", "binaryfunc", "FloorDivide", "(Object) (Object, error)"),
    ("nb_inplace_floor_divide", "binaryfunc",
     "InPlaceFloorDivide", "(Object) (Object, error)"),
    ("nb_true_divide", "binaryfunc", "TrueDivide", "(Object) (Object, error)"),
    ("nb_inplace_true_divide", "binaryfunc",
     "InPlaceTrueDivide", "(Object) (Object, error)"),
    ("nb_index", "unaryfunc", "Index", "() (*Long, error)"),
    ("nb_matrix_multiply", "binaryfunc",
     "MatrixMultiply", "(Object) (Object, error)"),
    ("nb_inplace_matrix_multiply", "binaryfunc",
     "InPlaceMatrixMultiply", "(Object) (Object, error)"),

    # Mapping Slots
    ("mp_length", "lenfunc", "PyMappingLen", "() int"),
    ("mp_subscript", "binaryfunc", "Subscript", "(Object) (Object, error)"),
    ("mp_ass_subscript", "objobjargproc",
     "AssSubscript", "(Object, Object) error"),

    # Sequence Slots
    ("sq_length", "lenfunc", "Len", "() int"),
    ("sq_concat", "binaryfunc", "Concat", "(Object) (Object, error)"),
    ("sq_repeat", "ssizeargfunc", "Repeat", "(Object, int) (Object, error)"),
    ("sq_item", "ssizeargfunc", "GetIndex", "(int) (Object, error)"),
    ("sq_ass_item", "ssizeobjargproc", "SetIndex", "(int, Object) error"),
    ("sq_contains", "objobjproc", "Contains", "(Object) (bool, error)"),
    ("sq_inplace_concat", "binaryfunc",
     "InPlaceConcat", "(Object) (Object, error)"),
    ("sq_inplace_repeat", "ssizeargfunc",
     "InPlaceRepeat", "(Object, int) (Object, error)"),

    # Buffer Slots
    ("bf_getbuffer", "getbufferproc", "GetBuffer", "(Object, int) error"),
    ("bf_releasebuffer", "releasebufferproc", "ReleaseBuffer", "(Object)"),
]

callbacks = {
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
    "() (Iterator, error)": (
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
    "() (*Long, error)": (
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
    "() (*Float, error)": (
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
    "() (int, error)": (
        "(obj unsafe.Pointer) C.long",
        [
            '	ret, err := co.%(fn)s()',
            '	if err != nil {',
            '		raise(err)',
            '		return -1',
            '	} else if ret == -1 {',
            '		return -2',
            '	}',
            '',
            '	return C.long(ret)',
        ],
    ),
    "(int) (Object, error)": (
        "(obj unsafe.Pointer, arg1 C.Py_ssize_t) unsafe.Pointer",
        [
            '	ret, err := co.%(fn)s(int(arg1))',
            '	if err != nil {',
            '		raise(err)',
            '		return nil',
            '	}',
            '',
            '	return unsafe.Pointer(c(ret))',
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
    "(Object) (Object, SendResult, error)": (
        "(obj, arg, out unsafe.Pointer) C.PySendResult",
        [
            '	o := newObject((*C.PyObject)(arg))',
            '	result := (**C.PyObject)(out)',
            '	ret, res, err := co.%(fn)s(o)',
            '	if err != nil {',
            '		raise(err)',
            '		*result = nil',
            '		return C.PYGEN_ERROR',
            '	}',
            '',
            '	*result = c(ret)',
            '	return C.PySendResult(res)',
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
    "(Object, Op) (Object, error)": (
        "(obj, arg1 unsafe.Pointer, arg2 C.int) unsafe.Pointer",
        [
            '	o := newObject((*C.PyObject)(arg1))',
            '	ret, err := co.%(fn)s(o, Op(arg2))',
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
    "(int, Object) error": (
        "(obj unsafe.Pointer, arg1 C.Py_ssize_t, arg2 unsafe.Pointer) C.int",
        [
            '	a := newObject((*C.PyObject)(arg2))',
            '	if err := co.%(fn)s(int(arg1), a); err != nil {',
            '		raise(err)',
            '		return -1',
            '	}',
            '',
            '	return 0',
        ],
    ),
}


def check_slots():
    slots_ok = True

    # we can only have 64 slots with a single uint64_t
    if len(slots) > 64:
        print("TOO MANY SLOTS!", file=sys.stderr)
        print(
            f"Can only fit 64 slots in a uint64_t, we have {len(slots)}", file=sys.stderr)
        slots_ok = False

    # names should be unique
    slot_names = {}
    func_names = {}
    wanted_callbacks = {}
    for (slot, _, func, cb) in slots:
        if slot in slot_names:
            print(f"DUPLICATE SLOT: {slot}", file=sys.stderr)
            slots_ok = False
        if func in func_names:
            print(f"DUPLICATE FUNC: {func}", file=sys.stderr)
            slots_ok = False
        slot_names[slot] = True
        func_names[func] = True
        wanted_callbacks[cb] = True

    # all callbacks should be needed, and all wanted callbacks should be present
    for cb in callbacks:
        if not wanted_callbacks.get(cb, False):
            print(f"UNUSED CALLBACK: {cb}")
            slots_ok = False
    for cb in wanted_callbacks:
        if not cb in callbacks:
            print(f"MISSING CALLBACK: {cb}")
            slots_ok = False

    if not slots_ok:
        sys.exit(1)


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
    print("  // The concrete Methods structs for the protocols (PyTypeObjects only stores", file=f)
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
    print("  if (type->tp_flags & Py_TPFLAGS_HAVE_GC) {", file=f)
    print("    type->tp_traverse = (traverseproc)goClassTraverse;", file=f)
    print("    type->tp_clear = (inquiry)goClassClear;", file=f)
    print("  }", file=f)
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
    check_slots()

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
