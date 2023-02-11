#!/usr/bin/env python3

# Copyright 2011 Julian Phillips.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import os
import sys
import re
import subprocess
import threading


header = """// Code generated by gen_exc.py. DO NOT EDIT.

// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file is automatically generated.  To regenerate:
//   go generate ./...

//go:generate ./gen_exc.py exc.go

package py

// #include "utils.h"
import "C"
"""

exc_re = re.compile('^extern .* PyExc_(?P<name>\w+);$')


def get_ffi_flags():
    cmd = ['pkg-config', '--cflags', 'libffi']
    p = subprocess.Popen(cmd, stdout=subprocess.PIPE)
    out, err = p.communicate()
    args = out.decode('utf-8').strip()
    if len(args) == 0:
        return []
    return args.split(' ')


def process(inp, outp):
    exceptions = []
    maxlen = 0

    for rawline in inp:
        line = rawline.strip().decode('utf-8')
        excm = exc_re.match(line)
        if excm is None:
            continue
        exc = excm.group('name')
        exceptions.append(exc)
        if len(exc) > maxlen:
            maxlen = len(exc)

    print(header, file=outp)

    print('var (', file=outp)
    for exception in exceptions:
        print('\t%s%s = newException(C.PyExc_%s)' %
              (exception, ' ' * (maxlen - len(exception)), exception),
              file=outp)
    print(')', file=outp)


def main():
    with open("_cgo_export.h", 'w', encoding='utf-8') as output:
        output.write("// stub file")
    with open(sys.argv[1], 'w', encoding='utf-8') as output:
        cmd = ["gcc", "-E", "-o", "-", "utils.c"] + get_ffi_flags()
        p = subprocess.Popen(cmd, stdout=subprocess.PIPE)
        t = threading.Thread(target=process, args=(p.stdout, output))
        t.start()
        p.wait()
        t.join()
    os.remove("_cgo_export.h")


if __name__ == "__main__":
    main()
