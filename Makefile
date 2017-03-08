.PHONY: clean functional-test get-deps unit-test test

test: get-deps unit-test functional-test

clean:
	rm -f farnsworth
	go clean

functional-test:
	go build
	cram tests/*.t

get-deps:
	go get ./...
	pip install -r tests/requirements.txt

unit-test:
	go test ./...

