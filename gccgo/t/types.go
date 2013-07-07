package main

import (
	"path/filepath"
	"fmt"
	"os"
	"strings"
)

func process(lines []string, base string) error {
	hOut, err := os.Create(base + ".h")
	if err != nil {
		return err
	}
	defer hOut.Close()

	cOut, err := os.Create(base + ".c")
	if err != nil {
		return err
	}
	defer cOut.Close()

	sOut, err := os.Create(base + ".sym")
	if err != nil {
		return err
	}
	defer sOut.Close()

	fmt.Fprintf(cOut, "#include \"py.h\"\n")
	fmt.Fprintf(cOut, "#include \"types.h\"\n")

	for _, name := range lines {
		if name == "" || name[0] == '#' {
			continue
		}
		flags := ""
		if strings.Contains(name, " ") {
			parts := strings.SplitN(name, " ", 2)
			name = parts[0]
			flags = strings.TrimSpace(parts[1])
		}
		lname := strings.ToLower(name[0:1]) + name[1:]
		fmt.Fprintf(sOut, "Py%s_Type\n", name)
		fmt.Fprintf(sOut, "%sCheck\n", lname)
		fmt.Fprintf(hOut, "extern int %sCheck(PyObject *o);\n", lname)
		if strings.Contains(flags, "E") {
			fmt.Fprintf(sOut, "%sCheckE\n", lname)
			fmt.Fprintf(hOut, "extern int %sCheckE(PyObject *o);\n", lname)
		}
		fmt.Fprintf(sOut, "get_Py%s_Type\n", name)
		fmt.Fprintf(hOut, "extern PyTypeObject *get_Py%s_Type(void);\n", name)
		fmt.Fprintf(cOut, "extern PyTypeObject *get_Py%s_Type(void) {\n", name)
		fmt.Fprintf(cOut, "    return &Py%s_Type;\n", name)
		fmt.Fprintf(cOut, "}\n\n")
	}

	return nil
}

func main() {
	base := strings.SplitN(filepath.Base(os.Args[1]), ".", 2)[0]
	lines, err := readLines(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load %s: %s", os.Args[1], err)
		os.Exit(1)
	}

	err = process(lines, base)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to process %s: %s", os.Args[1], err)
		os.Exit(1)
	}
}
