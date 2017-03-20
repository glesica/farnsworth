.PHONY: clean functional-test get-deps get-deps-user unit-test test

test: unit-test functional-test

clean:
	rm -f farnsworth
	go clean

functional-test:
	echo "Running functional tests..."
	go build
	cram tests/*.t

get-deps:
	go get -t ./...
	pip install -r tests/requirements.txt

get-deps-user:
	go get -t ./...
	pip install --user -r tests/requirements.txt

unit-test:
	echo "Running unit tests..."
	go test ./...

