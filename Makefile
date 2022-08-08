INCLUDEFLAGS		:= -I$(JAVA_HOME)/include -I$(JAVA_HOME)/include/darwin
LINKERFLAGS     	:= -L$(JAVA_HOME)/lib/server -L$(JAVA_HOME)/lib -ljvm
CGO_CFLAGS       	:= $(INCLUDEFLAGS)
CGO_LDFLAGS      	:= $(LINKERFLAGS)
CLASSPATH 			:= .

JAVAC				:= javac

.PHONY: build clean vendor test integration_test

build:
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" go build -o gmx ./cmd/main 

clean:
	go clean
	rm gmx  2>/dev/null || exit 0

vendor:
	go mod tidy
	go mod vendor

test: clean
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" TEST_ENV=UT go test ./...

integration_test: clean
#	docker rm --force jniexample
#	docker run -d -p 9001:9001 --name jniexample trixter1394/jniexample-snapshot 
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" TEST_ENV=IT go test ./...

mock_gen:
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" mockery --all --inpackage
	go mod tidy
	go mod vendor