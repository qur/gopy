# Copyright 2011 Julian Phillips.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

BINDIR ?= ~/misc/gcc/bin
CFLAGS ?= -g -fsplit-stack

GCC = $(BINDIR)/gcc
GCCGO = $(BINDIR)/gccgo
ECHO = /bin/echo -e

ifndef V
.SILENT:
endif

all: libgopy.so libgopy.ext.so wrap gopy

GOFILES=\
	abstract.go \
	arg.go \
	base.go \
	bool.go \
	cfunction.go \
	chan.go \
	class.go \
	class_mapping.go \
	class_number.go \
	class_sequence.go \
	class_standard.go \
	code.go \
	complex.go \
	dict.go \
	doc.go \
	err.go \
	eval.go \
	exc.go \
	exceptions.go \
	float.go \
	function.go \
	gomod.go \
	int.go \
	list.go \
	lock.go \
	long.go \
	memory.go \
	module.go \
	number.go \
	object.go \
	python.go \
	run.go \
	set.go \
	string.go \
	tuple.go \
	type.go \

FFI=$(shell pkg-config --cflags --libs libffi)

install: libgopy.so libgopy.ext.so gopy gopy.1.gz
	@$(ECHO) "  INSTALL"
	install -d $(DESTDIR)/usr/lib/gccgo $(DESTDIR)/usr/bin \
		$(DESTDIR)/etc/ld.so.conf.d $(DESTDIR)/usr/share/man/man1
	install -t $(DESTDIR)/usr/lib/gccgo/ libgopy.so libgopy.ext.so
	install -t $(DESTDIR)/usr/bin/ gopy
	install -t $(DESTDIR)/etc/ld.so.conf.d/ gccgo.conf
	install -t $(DESTDIR)/usr/share/man/man1/ gopy.1.gz

%.go:: ../lib/%.go t/decgo typedefs.map gen-py.go
	@$(ECHO) "  DE-CGO\t$@"
	./t/decgo $< $@ || (rm $@; exit 1)

typedefs.map: gen-py.go t/typedefs
	@$(ECHO) "  TYPEDEFS\t$@"
	./t/typedefs < $< > $@ || rm $@

gopy.1.gz: gopy.rst
	@$(ECHO) "  RST2MAN\t$@"
	rst2man $< | gzip > $@

py.o: $(GOFILES) py.go bindings.go exc.go
	@$(ECHO) "  GCCGO\t\t$@"
	$(GCCGO) -c $(CFLAGS) -fPIC -o $@ $^

pybase.c: py.o pywrap
	@$(ECHO) "  PYWRAP\t$@"
	./pywrap $< > $@

libgopy.ext.so: base.c pybase.c
	@$(ECHO) "  GCC\t\t$@"
	$(GCC) -shared $(CFLAGS) -fPIC -o $@ $^ -lgo -static-libgcc

libgopy.so: py.c ffi_type.c cgo.c types.c exc.c py.o
	@$(ECHO) "  GCC\t\t$@"
	$(GCC) -shared $(CFLAGS) -fPIC -o $@ $^ -lgo $(FFI) -static-libgcc

callbacks.c: t/callbacks callbacks.dat
	@$(ECHO) "  CALLBACKS\t$@"
	./t/callbacks > $@

bindings.go: sym.in gen-py.go sym
	@$(ECHO) "  SYM\t\t$@"
	./sym py $^ > $@

sym.in: symbols.in types.sym
	@$(ECHO) "  CAT\t\t$@"
	cat $^ > $@

cgo.c: sym.in gen-py.go cgo
	@$(ECHO) "  CGO\t\t$@"
	./cgo py $^ > $@

gen-py.go: py.c callbacks.c types.h
	@$(ECHO) "  DUMP-GO-SPEC\t$@"
	$(GCC) -S -o /dev/null -fdump-go-spec=$@ $< || (rm -f $@; false)

types.c types.h types.sym: types.in types
	@$(ECHO) "  TYPES\t\t$@"
	./types $<

exc.c exc.go: exc gen-py.go
	@$(ECHO) "  EXC\t\t$@"
	./$^ $@

gopy: t/gopy.go
	@$(ECHO) "  GCCGO\t\t$@"
	$(GCCGO) -O2 -s -o $@ $^ -static-libgo -static-libgcc

wrap: t/wrap.go
	@$(ECHO) "  GCCGO\t\t$@"
	$(GCCGO) -O2 -s -o $@ $^ -static-libgo -static-libgcc

pywrap: t/pywrap.go
	@$(ECHO) "  GCCGO\t\t$@"
	$(GCCGO) -O2 -s -o $@ $^ -static-libgo -static-libgcc

cgo: t/cgo.go t/utils.go
	@$(ECHO) "  GCCGO\t\t$@"
	$(GCCGO) -O2 -s -o $@ $^ -static-libgo -static-libgcc

sym: t/sym.go t/utils.go
	@$(ECHO) "  GCCGO\t\t$@"
	$(GCCGO) -O2 -s -o $@ $^ -static-libgo -static-libgcc

types: t/types.go t/utils.go
	@$(ECHO) "  GCCGO\t\t$@"
	$(GCCGO) -O2 -s -o $@ $^ -static-libgo -static-libgcc

exc: t/exc.go t/utils.go
	@$(ECHO) "  GCCGO\t\t$@"
	$(GCCGO) -O2 -s -o $@ $^ -static-libgo -static-libgcc

clean:
	@$(ECHO) "  CLEAN"
	rm -f *.o *.so *.[68] wrap sym bindings.go gen-py.go types.c types.sym types sym.in exc exc.c exc.go types.h pywrap pybase.c cgo cgo.c gopy callbacks.c gopy.1.gz $(GOFILES) typedefs.map

.PHONY: all clean
.INTERMEDIATE: bindings.go gen-py.go py.o sym.in types.c types.h types.sym exc.go exc.c exc types sym pywrap pybase.c cgo cgo.c callbacks.c typedefs.map
