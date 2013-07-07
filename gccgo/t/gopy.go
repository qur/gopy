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

func get_imports(filename string) (string, string, []imp, error) {
	data, err := read_go_export(filename)
	if err != nil {
		return "", "", nil, err
	}

	var name, prefix string
	var imports []imp
	has_main := false

	for _, b := range bytes.Split(data, []byte{'\n'}) {
		if bytes.HasPrefix(b, []byte("package ")) {
			name = string(b[8 : len(b)-1])
			continue
		}
		if bytes.HasPrefix(b, []byte("prefix ")) {
			prefix = string(b[7 : len(b)-1])
			continue
		}
		if bytes.HasPrefix(b, []byte("pkgpath ")) {
			path := string(b[8 : len(b)-1])
			if strings.Contains(path, ".") {
				parts := strings.Split(path, ".")
				prefix = strings.Join(parts[:len(parts)-1], ".")
				name = parts[len(parts)-1]
			}
			continue
		}
		if bytes.Equal(b, []byte("func Main ();")) {
			has_main = true
		}
		if !bytes.HasPrefix(b, []byte("init ")) {
			continue
		}
		bits := bytes.Split(b[5:len(b)-1], []byte{' '})
		for i := 0; i < len(bits); i += 3 {
			p, err := strconv.Atoi(string(bits[i+2]))
			if err != nil {
				return "", "", nil, err
			}
			I := imp{
				string(bits[i]),
				string(bits[i+1]),
				p,
			}
			imports = append(imports, I)
		}
	}

	if *exe && !has_main {
		err := fmt.Errorf("Package %s does not declare 'func Main()'", name)
		return "", "", nil, err
	}

	sort.Sort(impList(imports))

	return name, prefix, imports, nil
}

func process_exe_imports(fname, name, prefix string, imports []imp) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "#include <stdlib.h>\n")
	fmt.Fprintf(f, "\n")
	fmt.Fprintf(f, "extern int _init_go_main(int argc, char *argv[],")
	fmt.Fprintf(f, " void *funcs[]);\n")
	fmt.Fprintf(f, "\n")
	for _, i := range imports {
		fmt.Fprintf(f, "extern void __%s_init(void) __asm__ (\"%s\");\n", i.n, i.f)
	}
	fmt.Fprintf(f, "extern void __main(void) __asm__ (\"%s.%s.Main\");\n",
		prefix, name);
	fmt.Fprintf(f, "\n")
	fmt.Fprintf(f, "static void *funcs[] = {\n")
	for _, i := range imports {
		fmt.Fprintf(f, "    __%s_init,\n", i.n)
	}
	fmt.Fprintf(f, "    __main,\n")
	fmt.Fprintf(f, "    NULL\n")
	fmt.Fprintf(f, "};\n")
	fmt.Fprintf(f, "\n")
	fmt.Fprintf(f, "int main(int argc, char *argv[])\n")
	fmt.Fprintf(f, "{\n")
	fmt.Fprintf(f, "    return _init_go_main(argc, argv, funcs);\n")
	fmt.Fprintf(f, "}\n")
	return nil
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

var (
	GCCGO string
	GCC string
	CFLAGS []string
)

func build(out string, in []string) error {
	args := []string{"-c", "-fPIC", "-o", out}
	args = append(args, "-I/usr/lib/gccgo", "-I/usr/local/lib/gccgo")
	args = append(args, CFLAGS...)
	args = append(args, in...)
	if *verbose {
		fmt.Printf("build: %s %v\n", GCCGO, args)
	}
	cmd := exec.Command(GCCGO, args...)
	var errOut bytes.Buffer
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		fmt.Println(&errOut)
	}
	return err
}

func link(name string, in ...string) error {
	args := []string{"-fPIC"}
	if *exe {
		args = append(args, "-o", name)
	} else {
		args = append(args, "-shared", "-o", name+".so")
	}
	args = append(args, CFLAGS...)
	args = append(args, in...)
	args = append(args, "-L/usr/lib/gccgo", "-L/usr/local/lib/gccgo")
	args = append(args, "-lgo", "-lgopy", "-lgopy.ext", "-lpython2.7")
	args = append(args, "-lpthread", "-static-libgcc")
	if *verbose {
		fmt.Printf("link: %s %v\n", GCC, args)
	}
	cmd := exec.Command(GCC, args...)
	var errOut bytes.Buffer
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		fmt.Println(&errOut)
	}
	return err
}

var (
	work = flag.Bool("work", false, "print and keep work directory")
	verbose = flag.Bool("v", false, "verbose output")
	exe = flag.Bool("exe", false, "create an executable instead of a package")
)

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
		GCC = "gcc"
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

	name, prefix, imp, err := get_imports(objFile)
	if err != nil {
		fmt.Println("get imports failed:", err)
		os.Exit(1)
	}

	if *exe {
		err = process_exe_imports(cFile, name, prefix, imp)
	} else {
		err = process_imports(cFile, name, imp)
	}
	if err != nil {
		fmt.Println("process imports failed:", err)
		os.Exit(1)
	}

	err = link(name, cFile, objFile)
	if err != nil {
		fmt.Println("link failed:", err)
		os.Exit(1)
	}
}
