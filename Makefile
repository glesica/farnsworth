.PHONY: clean functional-test unit-test test

clean:
	rm -f farnsworth
	go clean

functional-test:
	go build
	cram tests/*.t

unit-test:
	go test ./...

test: unit-test functional-test

