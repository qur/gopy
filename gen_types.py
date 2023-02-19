#!/usr/bin/env python3

import os
import re
import subprocess
import sys

include_re = re.compile('^\s*#\s*include\s+"(?P<name>[^"]+)"\s*$')

check_re = re.compile(
    '(?:[^_A-Za-z0-9]|^)Py(?P<name>[A-Za-z0-9]+?)_Check(?:[^_A-Za-z0-9]|$)')
check_exact_re = re.compile(
    '(?:[^_A-Za-z0-9]|^)Py(?P<name>[A-Za-z0-9]+?)_CheckExact(?:[^_A-Za-z0-9]|$)')
type_re = re.compile(
    '(?:[^_A-Za-z0-9]|^)Py(?P<name>[A-Za-z0-9]+?)_Type(?:[^_A-Za-z0-9]|$)')
object_re = re.compile(
    '(?:[^_A-Za-z0-9]|^)Py(?P<name>[A-Za-z0-9]+?)Object(?:[^_A-Za-z0-9]|$)')


def get_includes(root: str, files: set) -> set:
    if len(files) == 0:
        return files
    new = set()
    for name in sorted(files):
        path = os.path.join(root, name)
        with open(path, 'r', encoding='utf-8') as input:
            for line in input:
                m = include_re.match(line)
                if m:
                    new.add(m.group('name'))
    files.update(get_includes(root, new))
    files.update(new)
    return files


def process(path, checks, exact_checks, types, objects):
    with open(path, 'r', encoding='utf-8') as input:
        for line in input:
            checks.update(check_re.findall(line))
            exact_checks.update(check_exact_re.findall(line))
            types.update(type_re.findall(line))
            objects.update(object_re.findall(line))


def write_h_header(f):
    print("// Code generated by gen_types.py. DO NOT EDIT.", file=f)
    print("", file=f)
    print("#ifndef _GO_PYTHON_TYPES_H_", file=f)
    print("#define _GO_PYTHON_TYPES_H_", file=f)
    print("", file=f)
    print("#include <ffi.h>", file=f)
    print("#include <python3.11/Python.h>", file=f)
    print("#include <python3.11/frameobject.h>", file=f)
    print("", file=f)


def write_h_footer(f):
    print("", file=f)
    print("#endif /* _GO_PYTHON_TYPES_H */", file=f)


def write_check_signatures(f, checks, exact_checks):
    print("// ===============================================================", file=f)
    print("", file=f)
    for check in sorted(checks):
        name = check[0].lower() + check[1:]
        print(
            f'extern int { name }Check(PyObject *o);', file=f)
    for check in sorted(exact_checks):
        name = check[0].lower() + check[1:]
        print(
            f'extern int { name }CheckE(PyObject *o);', file=f)
    print("", file=f)
    print("// ===============================================================", file=f)


def write_function_headers(f):
    print("// ===============================================================", file=f)
    print("", file=f)
    print('// --- types.go helper functions ---', file=f)
    print('//', file=f)
    print('extern PyTypeObject *getBasePyType(PyObject *o);', file=f)
    print('//', file=f)
    print('// ----------------------------------', file=f)
    print("", file=f)
    print("// ===============================================================", file=f)


def write_c_header(f):
    print('// Code generated by gen_types.py. DO NOT EDIT.', file=f)
    print('', file=f)
    print('#include "utils.h"', file=f)
    print('', file=f)
    print('#include "_cgo_export.h"', file=f)
    print('', file=f)


def write_check_functions(f, checks, exact_checks):
    print("// ===============================================================", file=f)
    print("", file=f)
    for check in sorted(checks):
        name = check[0].lower() + check[1:]
        print(
            f'int { name }Check(PyObject *o) {{ return Py{ check }_Check(o); }}', file=f)
    for check in sorted(exact_checks):
        name = check[0].lower() + check[1:]
        print(
            f'int { name }CheckE(PyObject *o) {{ return Py{ check }_CheckExact(o); }}', file=f)
    print("", file=f)
    print("// ===============================================================", file=f)


def write_get_types(f, checks, types):
    print("// ===============================================================", file=f)
    print("", file=f)
    print('PyTypeObject *getBasePyType(PyObject *o) {', file=f)
    print('  if (o == NULL) return NULL;', file=f)
    print('', file=f)
    for type in sorted(types):
        if type in checks:
            print(
                f"  if (Py{ type }_Check(o)) return &Py{ type }_Type;", file=f)
    print('', file=f)
    print('  return o->ob_type;', file=f)
    print('}', file=f)
    print("", file=f)
    print("// ===============================================================", file=f)


def write_go_header(f):
    print('// Code generated by gen_types.py. DO NOT EDIT.', file=f)
    print('', file=f)
    print('// This file is automatically generated.  To regenerate:', file=f)
    print('//   go generate ./...', file=f)
    print('', file=f)
    print('//go:generate ./gen_types.py /usr/local/include/python3.11', file=f)
    print('', file=f)
    print('package py', file=f)
    print('', file=f)
    print('// #include "utils.h"', file=f)
    print('import "C"', file=f)
    print('', file=f)


def write_get_native_type(f, checks, types, objects):
    print("// ===============================================================", file=f)
    print("", file=f)
    print('func getNativeType(obj *C.PyObject) Object {', file=f)
    print('	switch C.getBasePyType(obj) {', file=f)
    for type in sorted(types):
        # Bool, FrozenSet and Module are manually generated, we want to include them
        if (type not in checks or type not in objects) and type not in ["Bool", "FrozenSet", "Module"]:
            continue
        print(f"	case &C.Py{ type }_Type:", file=f)
        print(f"		return new{ type }(obj)", file=f)
    print('	}', file=f)
    print('', file=f)
    print('	return nil;', file=f)
    print('}', file=f)
    print("", file=f)
    print("// ===============================================================", file=f)


def write_cmd_go_header(f):
    print('// Code generated by gen_types.py. DO NOT EDIT.', file=f)
    print('', file=f)
    print('package main', file=f)
    print('', file=f)
    print('// #cgo CFLAGS: -Werror', file=f)
    print('// #cgo LDFLAGS: -lpython3.11 -lm -lutil -ldl', file=f)
    print('// #cgo pkg-config: libffi', file=f)
    print('//', file=f)
    print('// #include <python3.11/Python.h>', file=f)
    print('import "C"', file=f)
    print('', file=f)


def write_type_table(f, checks, types, objects):
    print("// ===============================================================", file=f)
    print("", file=f)
    print('var types = map[string]any{', file=f)
    for type in sorted(types):
        if type not in checks or type not in objects:
            continue
        print(f'	"{type}": C.Py{type}_Type,', file=f)
    print('}', file=f)
    print("", file=f)
    print("// ===============================================================", file=f)


def main():
    root = sys.argv[1]

    files = get_includes(root, set(["Python.h", "frameobject.h"]))

    checks, exact_checks, types, objects = set(), set(), set(), set()
    for name in sorted(files):
        path = os.path.join(root, name)
        process(path, checks, exact_checks, types, objects)

    # PyGILState_Check is a special case, ignore it
    checks.remove("GILState")

    with open('types.h', 'w', encoding='utf-8') as output:
        write_h_header(output)
        write_check_signatures(output, checks, exact_checks)
        write_function_headers(output)
        write_h_footer(output)

    with open('types.c', 'w', encoding='utf-8') as output:
        write_c_header(output)
        write_check_functions(output, checks, exact_checks)
        write_get_types(output, checks, types)

    with open('types.go', 'w', encoding='utf-8') as output:
        write_go_header(output)
        write_get_native_type(output, checks, types, objects)

    with open('cmd/gen_types/types.go', 'w', encoding='utf-8') as output:
        write_cmd_go_header(output)
        write_type_table(output, checks, types, objects)

    subprocess.run(["go", "run", "qur.me/py/v3/cmd/gen_types"], check=True)


if __name__ == "__main__":
    main()
