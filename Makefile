INCLUDEFLAGS		:= -I$(JAVA_HOME)/include -I$(JAVA_HOME)/include/darwin
LINKERFLAGS     	:= -L$(JAVA_HOME)/lib/server -L$(JAVA_HOME)/lib -ljvm
CGO_CFLAGS       	:= $(INCLUDEFLAGS)
CGO_LDFLAGS      	:= $(LINKERFLAGS)
CLASSPATH 			:= .

JAVAC				:= javac

JNIConnector.class:
	$(JAVAC) internal/java/JNIConnector.java

.PHONY: build clean vendor test integration_test

build: JNIConnector.class
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" go build -o gmx ./cmd/main 

clean:
	go clean
	rm JNIConnector.class java/JNIConnector.class.go gmx  2>/dev/null || exit 0

vendor:
	go mod tidy
	go mod vendor

test: clean
	go test pkg/...

integration_test: clean
	docker run -d trixter1394/jniexample-snapshot
	go test internal/java/...
