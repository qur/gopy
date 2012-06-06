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

type CSym struct {
	ret string
	args []string
	aStct string
}

func typeConv(g string) string {
	ptr := ""
	for g[0] == '*' {
		ptr += "*"
		g = g[1:]
	}
	if g[0] == '_' {
		return g[1:] + ptr
	}
	switch g {
	case "int":
		return g + ptr
	case "byte":
		return "void" + ptr
	case "int8":
		return "char" + ptr
	case "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64":
		return g + "_t" + ptr
	case "float64":
		return "double" + ptr
	default:
		log.Fatalf("unknown go type: %s", g)
	}
	panic("unreachable")
}

func convertToC(name string, sym Sym) CSym {
	cSym := CSym{"void", nil, ""}
	parts := strings.Split(sym.declare, " __asm__")
	parts = strings.Split(parts[0], ") ")
	fmt.Printf("//cgo: %v\n", parts)
	if len(parts) > 1 {
		cSym.ret = typeConv(parts[1])
	}
	if len(parts[0]) > 2 {
		p := parts[0][1:]
		if p[len(p)-1] == ')' {
			p = p[:len(p)-1]
		}
		for _, gt := range strings.Split(p, ", ") {
			cSym.args = append(cSym.args, typeConv(gt))
		}
	}
	cSym.aStct = "struct { "
	if cSym.ret != "void" {
		cSym.aStct += fmt.Sprintf("%s ret; ", cSym.ret)
	}
	for i, t := range cSym.args {
		cSym.aStct += fmt.Sprintf("%s a%d; ", t, i+1)
	}
	cSym.aStct += "}"
	return cSym
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

	fmt.Printf("#include \"py.h\"\n")
	fmt.Printf("#include \"types.h\"\n")
	fmt.Printf("\n")
	fmt.Printf("extern void simple_cgocall(void (*)(void*), void*);\n")
	fmt.Printf("\n")

	for _, name := range want {
		if name == "" || name[0] == '#' {
			continue
		}
		sym, ok := findSym(name, symbols)
		if !ok {
			log.Fatalf("Unknown symbol: %s\n", name)
		}
		if sym.kind != "func" {
                        continue
		}
		cSym := convertToC(name, sym)
		fmt.Printf("//[%s] [%s] [%s]\n", sym.kind, name, sym.declare)
		fmt.Printf("static void _do_%s(void *_a) {\n", name)
		fmt.Printf("    %s *a = _a;\n", cSym.aStct)
		fmt.Printf("    ")
		if cSym.ret != "void" {
			fmt.Printf("a->ret = ")
		}
		fmt.Printf("%s(", name)
		for i := range cSym.args {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("a->a%d", i+1)
		}
		fmt.Printf(");\n")
		fmt.Printf("}\n")
		fmt.Printf("extern %s cgo_%s(", cSym.ret, name)
		if len(cSym.args) == 0 {
			fmt.Printf("void")
		} else {
			for i, t := range cSym.args {
				if i > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%s a%d", t, i+1)
			}
		}
		fmt.Printf(") {\n")
		fmt.Printf("    %s a;\n", cSym.aStct)
		for i := range cSym.args {
			fmt.Printf("    a.a%d = a%d;\n", i+1, i+1)
		}
        fmt.Printf("    simple_cgocall(_do_%s, &a);\n", name)
		if cSym.ret != "void" {
			fmt.Printf("    return a.ret;\n")
		}
		fmt.Printf("}\n\n")
	}
}
