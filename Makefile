all:
	go install -v ./...

install: all
	@echo

test:
	go test ./...

clean:
	@echo

