all: Main parallel.so simple.so rgo.so pyg

gopy="../gopy"

Main: main.go
	$(gopy) -exe $^

pyg: pyg.go
	$(gopy) -exe $^

parallel.so: parallel.go
	$(gopy) $^

simple.so: simple.go
	$(gopy) $^

rgo.so: rgo.go
	$(gopy) $^

clean:
	rm -f Main pyg *.so

.PHONY: clean
