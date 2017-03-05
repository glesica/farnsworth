# Farnsworth Functional Tests

The Farnsworth command line tool is tested with
[cram](https://pypi.python.org/pypi/cram). To run the functional tests you'll
need a working Python environment. First build Farnsworth with `go build`. The
functional tests assume the binary will be located in the project root. Then,
install cram (and any other dependencies) with `pip install -r
requirements.txt`. This can be done using a virtual environment if you prefer
(probably a good idea). Finally, run the tests with `cram *.t`.

