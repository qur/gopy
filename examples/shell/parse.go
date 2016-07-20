// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var stdin = bufio.NewReader(os.Stdin)

func readLine() (string, error) {
	data, isPrefix, err := stdin.ReadLine()
	if err != nil {
		return "", err
	}

	buffer := make([]byte, len(data))
	copy(buffer, data)

	for isPrefix {
		buffer = append(buffer, data...)
		data, isPrefix, err = stdin.ReadLine()
		if err != nil {
			return "", nil
		}
	}

	return string(buffer), nil
}

func glob(dest []string, s, raw string) []string {
	chunk, e := filepath.Glob(s)

	if e != nil || len(chunk) == 0 || (len(chunk) == 1 && chunk[0] == s) {
		return append(dest, raw)
	}

	wantDot := strings.HasPrefix(filepath.Base(s), ".")

	for _, exp := range chunk {
		if strings.HasPrefix(filepath.Base(exp), ".") == wantDot {
			dest = append(dest, exp)
		}
	}

	return dest
}

func tokenise(line string) []string {
	globBuf := make([]byte, 0, 8192)
	buffer := make([]byte, 0, 8192)
	parts := make([]string, 0)

	inQuote := false
	for i := 0; i < len(line); {
		switch c := line[i]; c {
		case '\\':
			globBuf = append(globBuf, '\\')
			i++
			if i < len(line) {
				globBuf = append(globBuf, line[i])
				buffer = append(buffer, line[i])
			}
		case '"':
			inQuote = !inQuote
		case ' ':
			if inQuote {
				globBuf = append(globBuf, '\\')
				globBuf = append(globBuf, c)
				buffer = append(buffer, c)
			} else {
				parts = glob(parts, string(globBuf), string(buffer))
				globBuf = globBuf[0:0]
				buffer = buffer[0:0]
			}
		default:
			if inQuote {
				globBuf = append(globBuf, '\\')
				globBuf = append(globBuf, c)
				buffer = append(buffer, c)
			} else {
				globBuf = append(globBuf, c)
				buffer = append(buffer, c)
			}
		}
		i++
	}
	if len(buffer) > 0 {
		parts = glob(parts, string(globBuf), string(buffer))
	}

	return parts
}

func getCommand() (cmd string, args []string, bg bool, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	fmt.Printf("%s> ", cwd)

	line, err := readLine()
	if err != nil {
		return
	}

	if line == "" {
		return
	}

	if len(line) > 1 && line[len(line)-1] == '&' {
		bg = true
		line = line[:len(line)-1]
	}

	parts := tokenise(line)

	cmd = parts[0]
	args = parts[1:]

	return
}
