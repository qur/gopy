VERSION=$(shell git describe --match 'ext/v*' --dirty | cut -c6-)

gofiles: $(GOFILES) utils.c utils.h

%.go:: ../lib/%.go t/decgo typedefs.map gen-py.go
	@$(ECHO) "  DE-CGO\t$@"
	./t/decgo $< $@ || (rm $@; exit 1)

utils.c: ../lib/utils.c utils.h
	@$(ECHO) "  COPY\t\t$@"
	cp $< $@

utils.h: ../lib/utils.h
	@$(ECHO) "  COPY\t\t$@"
	cp $< $@

VERSION:
	@$(ECHO) "  VERSION"
	echo $(VERSION) > $@

LICENSE: ../LICENSE
	@$(ECHO) "  LICENSE"
	cp $< $@

gopy-$(VERSION).tgz: clean $(GOFILES) utils.c utils.h VERSION LICENSE
	@$(ECHO) "  TAR\t\t$@"
	tar --transform "s/^\./gopy-$(VERSION)/" --exclude '*.fixup' \
		--exclude '.gitignore' --exclude 'gopy-*.tgz' \
		--exclude 'bootstrap.mk' -czf $@ .

tarball: gopy-$(VERSION).tgz

deepclean: clean
	@$(ECHO) "  DEEPCLEAN"
	rm -f $(GOFILES) utils.[ch] VERSION LICENSE

.PHONY: gofiles tarball deepclean
