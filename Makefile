INCLUDEFLAGS		:= -I$(JAVA_HOME)/include -I$(JAVA_HOME)/include/darwin
LINKERFLAGS     	:= -L$(JAVA_HOME)/lib/server -L$(JAVA_HOME)/lib -ljvm
CGO_CFLAGS       	:= $(INCLUDEFLAGS)
CGO_LDFLAGS      	:= $(LINKERFLAGS)
CLASSPATH 			:= .

JAVAC				:= javac

.PHONY: build clean vendor test integration_test mocks jniexample

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
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" TEST_ENV=$(TEST_ENV) go test ./...

integration_test: TEST_ENV=IT 
integration_test: jniexample test stop
	@echo "Test run complete"

mocks: _mock_gen vendor

_mock_gen:
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" mockery --all --inpackage

jniexample:
	docker compose -f ./test-docker-compose.yml up -d

stop:
	docker compose -f ./test-docker-compose.yml down