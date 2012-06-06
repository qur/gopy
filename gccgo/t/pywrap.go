// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"debug/elf"
	"fmt"
	"os"
	"strconv"
	"sort"
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

func get_imports(filename string) ([]imp, error) {
	data, err := read_go_export(filename)
	if err != nil {
		return []imp{}, err
	}

	var imports []imp

	for _, b := range bytes.Split(data, []byte{'\n'}) {
		if !bytes.HasPrefix(b, []byte("init ")) {
			continue
		}
		bits := bytes.Split(b[7:len(b)-1], []byte{' '})
		for i := 0; i < len(bits); i += 3 {
			p, err := strconv.Atoi(string(bits[i+2]))
			if err != nil {
				return []imp{}, err
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

	return imports, nil
}

func process_imports(name string, imports []imp) {
	fmt.Printf("#include <stdlib.h>\n")
	fmt.Printf("\n")
	for _, i := range imports {
		fmt.Printf("extern void __%s_init(void) __asm__ (\"%s\");\n", i.n, i.f)
	}
	fmt.Printf("\n")
	fmt.Printf("void *py_init_funcs[] = {\n")
	for _, i := range imports {
		fmt.Printf("    __%s_init,\n", i.n)
	}
	fmt.Printf("    NULL\n")
	fmt.Printf("};\n")
}

func main() {
	i, err := get_imports(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	process_imports(os.Args[1][:len(os.Args[1])-2], i)
}
