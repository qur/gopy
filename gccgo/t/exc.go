package main

import (
	"fmt"
	"os"
	"strings"
)

func filter(lines []string) []string {
	var exceptions []string

	for _, line := range lines {
		if strings.HasPrefix(line, "var _PyExc_") && strings.HasSuffix(line, " *_PyObject") {
			exceptions = append(exceptions, line[11:len(line)-11])
		}
	}

	return exceptions
}

func writeGo(exceptions []string, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "package py\n\n")
	fmt.Fprintf(f, "type _exceptions_ struct {\n")
	for _, e := range exceptions {
		fmt.Fprintf(f, "\t%s *_PyObject\n", e)
	}
	fmt.Fprintf(f, "}\n\n")
	fmt.Fprintf(f, "func _get_exceptions () *_exceptions_ __asm__(\"_gopy_get_exceptions\")\n")

	fmt.Fprintf(f, "var _exc = _get_exceptions()\n\n")

	fmt.Fprintf(f, "var (\n")
	for _, e := range exceptions {
		fmt.Fprintf(f, "\t%s = newException(_exc.%s)\n", e, e)
	}
	fmt.Fprintf(f, ")\n")

	return nil
}

func writeC(exceptions []string, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "#include \"py.h\"\n\n")
	fmt.Fprintf(f, "static struct __exc {\n")
	for _, e := range exceptions {
		fmt.Fprintf(f, "    PyObject *%s;\n", e)
	}
	fmt.Fprintf(f, "} _exceptions;\n\n")
	fmt.Fprintf(f, "extern struct __exc *_gopy_get_exceptions(void) {\n")
	for _, e := range exceptions {
		fmt.Fprintf(f, "    _exceptions.%s = PyExc_%s;\n", e, e)
	}
	fmt.Fprintf(f, "    return &_exceptions;\n")
	fmt.Fprintf(f, "}\n")

	return nil
}

func main() {
	lines, err := readLines(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load %s: %s", os.Args[1], err)
		os.Exit(1)
	}

	exceptions := filter(lines)

	if strings.HasSuffix(os.Args[2], ".go") {
		err := writeGo(exceptions, os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write %s: %s", os.Args[2], err)
			os.Exit(1)
		}
	} else {
		err := writeC(exceptions, os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write %s: %s", os.Args[2], err)
			os.Exit(1)
		}
	}
}
