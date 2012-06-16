// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func rewriteForCgo(name string, sym Sym) Sym {
	parts := strings.Split(sym.declare, "__asm__")
	fmt.Printf("//cgo: %v\n", parts)
	sym.declare = fmt.Sprintf("%s __asm__ (\"cgo_%s\")", parts[0], name)
	return sym
}

func rewriteVar(name string, sym Sym) Sym {
	sym.declare = "= _get_" + name + "()"
	return sym
}

func main() {
	want, err := readLines(os.Args[2])
	if err != nil {
		log.Fatalln(err)
	}

	symbols, err := loadSym(os.Args[3])
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("package %s\n\n", os.Args[1])

	for _, name := range want {
		if name == "" || name[0] == '#' {
			continue
		}
		sym, ok := findSym(name, symbols)
		if !ok {
			log.Fatalf("Unknown symbol: %s\n", name)
		}
		fmt.Printf("//[%s] [%s] [%s]\n", sym.kind, name, sym.declare)
		switch sym.kind {
		case "func":
			sym = rewriteForCgo(name, sym)
		case "var":
			sym = rewriteVar(name, sym)
		}
		fmt.Printf("%s _%s %s\n", sym.kind, name, sym.declare)
	}
}
