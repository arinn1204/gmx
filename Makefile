INCLUDEFLAGS		:= -I$(JAVA_HOME)/include -I$(JAVA_HOME)/include/darwin
LINKERFLAGS     	:= -L$(JAVA_HOME)/lib/server -L$(JAVA_HOME)/lib -ljvm
CGO_CFLAGS       	:= $(INCLUDEFLAGS)
CGO_LDFLAGS      	:= $(LINKERFLAGS)
CLASSPATH 			:= .

JAVAC				:= javac

.PHONY: build clean vendor test integration_test mocks

build:
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" go build -o gmx ./cmd/main 

clean:
	go clean
	rm gmx  2>/dev/null || exit 0

vendor:
	go mod tidy
	go mod vendor

lint:
	go fmt ./...
	golint cmd/... internal/... pkg/...

test: clean
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" TEST_ENV=UT go test ./...

integration_test: clean jniexample secondary_example
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" TEST_ENV=IT go test ./...

mocks: _mock_gen vendor

_mock_gen:
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" mockery --all --inpackage

secondary_example: name=jniexample2
secondary_example: port=9002
secondary_example: example

jniexample: name=jniexample
jniexample: port=9001
jniexample: example

example:
	docker rm --force $(name) 2>/dev/null
	docker run -d -p $(port):9001 --name $(name) trixter1394/jniexample