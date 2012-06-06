// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"debug/elf"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

func read_go_export(filename string) ([]byte, error) {
	f, err := elf.Open(filename)
	if err != nil {
		return []byte{}, err
	}
	defer f.Close()

	s := f.Section(".go_export")
	if s == nil {
		return []byte{}, fmt.Errorf(".go_export section missing!\n")
	}

	return s.Data()
}

type imp struct {
	n string
	f string
	p int
}

type impList []imp

func (i impList) Len() int {
	return len(i)
}

func (i impList) Less(a, b int) bool {
	return i[a].p < i[b].p
}

func (i impList) Swap(a, b int) {
	i[a], i[b] = i[b], i[a]
}

func get_imports(filename string) (string, []imp, error) {
	data, err := read_go_export(filename)
	if err != nil {
		return "", nil, err
	}

	var name string
	var imports []imp

	for _, b := range bytes.Split(data, []byte{'\n'}) {
		if bytes.HasPrefix(b, []byte("package ")) {
			name = string(b[8 : len(b)-1])
			continue
		}
		if !bytes.HasPrefix(b, []byte("init ")) {
			continue
		}
		bits := bytes.Split(b[5:len(b)-1], []byte{' '})
		for i := 0; i < len(bits); i += 3 {
			p, err := strconv.Atoi(string(bits[i+2]))
			if err != nil {
				return "", nil, err
			}
			I := imp{
				string(bits[i]),
				string(bits[i+1]),
				p,
			}
			imports = append(imports, I)
		}
	}

	sort.Sort(impList(imports))

	return name, imports, nil
}

func process_imports(fname, name string, imports []imp) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "#include <python2.7/Python.h>\n")
	fmt.Fprintf(f, "\n")
	fmt.Fprintf(f, "extern void _init_go(void *funcs[]);\n")
	for _, i := range imports {
		fmt.Fprintf(f, "extern void __%s_init(void) __asm__ (\"%s\");\n", i.n, i.f)
	}
	fmt.Fprintf(f, "\n")
	fmt.Fprintf(f, "static void *funcs[] = {\n")
	for _, i := range imports {
		fmt.Fprintf(f, "    __%s_init,\n", i.n)
	}
	fmt.Fprintf(f, "    NULL\n")
	fmt.Fprintf(f, "};\n")
	fmt.Fprintf(f, "\n")
	fmt.Fprintf(f, "PyMODINIT_FUNC init%s(void)\n", name)
	fmt.Fprintf(f, "{\n")
	fmt.Fprintf(f, "    _init_go(funcs);\n")
	fmt.Fprintf(f, "}\n")
	return nil
}

var GCCGO string
var GCC string
var CFLAGS []string

func build(out string, in []string) error {
	args := []string{"-c", "-fPIC", "-o", out}
	args = append(args, "-I/usr/lib/gccgo", "-I/usr/local/lib/gccgo")
	args = append(args, in...)
	args = append(args, CFLAGS...)
	cmd := exec.Command(GCCGO, args...)
	var errOut bytes.Buffer
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		fmt.Println(&errOut)
	}
	return err
}

func link(out string, in ...string) error {
	args := []string{"-shared", "-fPIC", "-o", out}
	args = append(args, in...)
	args = append(args, CFLAGS...)
	args = append(args, "-L/usr/lib/gccgo", "-L/usr/local/lib/gccgo")
	args = append(args, "-lgo", "-lgopy", "-lgopy.ext", "-static-libgcc")
	cmd := exec.Command(GCC, args...)
	var errOut bytes.Buffer
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		fmt.Println(&errOut)
	}
	return err
}

var work = flag.Bool("work", false, "print and keep work directory")

func main() {
	flag.Parse()

	tDir, err := ioutil.TempDir("", "gopy_")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if *work {
		fmt.Println("WORK =", tDir)
	} else {
		defer os.RemoveAll(tDir)
	}

	cFlags, found := syscall.Getenv("CFLAGS")
	if !found {
		cFlags = "-g"
	}
	CFLAGS = strings.Split(cFlags, " ")

	GCC, found = syscall.Getenv("GOPY_CC")
	if !found {
		GCC = "gcc-4.7"
	}

	GCCGO, found = syscall.Getenv("GOPY_GO")
	if !found {
		GCCGO = "gccgo"
	}

	objFile := filepath.Join(tDir, "module.o")
	cFile := filepath.Join(tDir, "module.c")

	err = build(objFile, flag.Args())
	if err != nil {
		fmt.Println("build failed:", err)
		os.Exit(1)
	}

	name, imp, err := get_imports(objFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = process_imports(cFile, name, imp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = link(name+".so", cFile, objFile)
	if err != nil {
		fmt.Println("link failed:", err)
		os.Exit(1)
	}
}
