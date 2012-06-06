package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func readLines(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b := bufio.NewReader(f)

	lines := make([]string, 0, 100)

	for {
		line, err := b.ReadString('\n')
		if err == nil && len(line) > 0 {
			line = line[:len(line)-1]
		}
		lines = append(lines, line)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}

	return lines, nil
}

type Sym struct {
	kind    string
	declare string
}

func getSym(line string) (string, Sym) {
	if strings.HasPrefix(line, "// ") && !strings.HasPrefix(line, "// undef ") {
		line = line[3:]
	}
	parts := strings.SplitN(line, " ", 3)
	return parts[1], Sym{parts[0], parts[2]}
}

func loadSym(filename string) (map[string]Sym, error) {
	lines, err := readLines(filename)
	if err != nil {
		return nil, err
	}

	symbols := make(map[string]Sym)

	for _, line := range lines {
		if line == "" {
			continue
		}
		name, entry := getSym(line)
		symbols[name] = entry
	}

	return symbols, err
}

func findSym(name string, symbols map[string]Sym) (Sym, bool) {
	sym, ok := symbols["_"+name]
	if ok {
		return sym, true
	}
	sym, ok = symbols[name]
	if ok && sym.kind == "unknowndefine" && !strings.Contains(sym.declare, " ") {
		sym, ok = symbols["_"+sym.declare]
		return sym, ok
	}
	return Sym{}, false
}
