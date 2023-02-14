#!/usr/bin/env python3

import sys


methods = {
    "tp": (None, None),
    "am": ("async", "PyAsyncMethods"),
    "nb": ("number", "PyNumberMethods"),
    "mp": ("mapping", "PyMappingMethods"),
    "sq": ("sequence", "PySequenceMethods"),
    "bf": ("buffer", "PyBufferMethods"),
}

slots = [
    # Standard Slots
    ("tp_repr", "reprfunc", ""),
    ("tp_hash", "hashfunc", ""),
    ("tp_call", "terneryfunc", ""),
    ("tp_str", "reprfunc", ""),
    ("tp_getattro", "getattrofunc", ""),
    ("tp_setattro", "setattrofunc", ""),
    ("tp_richcompare", "richcmpfunc", ""),
    ("tp_iter", "getiterfunc", ""),
    ("tp_iternext", "iternextfunc", ""),
    ("tp_descr_get", "descrgetfunc", ""),
    ("tp_descr_set", "descrsetfunc", ""),
    ("tp_init", "initproc", ""),

    # Async Slots
    ("am_await", "unaryfunc", ""),
    ("am_aiter", "unaryfunc", ""),
    ("am_anext", "unaryfunc", ""),
    ("am_send", "sendfunc", ""),

    # Number Slots

    # Mapping Slots
    ("mp_length", "lenfunc", ""),
    ("mp_subscript", "binaryfunc", ""),
    ("mp_ass_subscript", "objobjargproc", ""),

    # Sequence Slots
    ("sq_length", "lenfunc", ""),
    ("sq_concat", "binaryfunc", ""),
    ("sq_repeat", "ssizeargfunc", ""),
    ("sq_item", "ssizeargfunc", ""),
    ("sq_ass_item", "ssizeobjargproc", ""),
    ("sq_contains", "objobjproc", ""),
    ("sq_inplace_concat", "binaryfunc", ""),
    ("sq_inplace_repeat", "ssizeargfunc", ""),

    # Buffer Slots
    ("bf_getbuffer", "", ""),
    ("bf_releasebuffer", "", ""),
]


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
    for (slot, _, _) in slots:
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
    for (slot, pySig, _) in slots:
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
            f"    if (slotFlags & CLASS_HAS_{slot.upper()}) m->{slot} = ({pySig})goClassSlot_{slot}", file=f)
    print("  }", file=f)
    print("", file=f)
    print("  return ctxt;", file=f)
    print("}", file=f)
    print("", file=f)
    print("// ===============================================================", file=f)


def main():
    with open("slots.c", "w", encoding='utf-8') as output:
        write_setSlots(output)

    with open("slots.h", "w", encoding='utf-8') as output:
        write_h_header(output)
        write_constants(output)
        write_ClassContext(output)
        write_extern_setSlots(output)
        write_h_footer(output)


if __name__ == "__main__":
    main()
