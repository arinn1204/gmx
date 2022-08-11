INCLUDEFLAGS		:= -I$(JAVA_HOME)/include -I$(JAVA_HOME)/include/darwin
CGO_CFLAGS       	:= $(INCLUDEFLAGS)

JAVAC				:= javac

.PHONY: build clean vendor test integration_test mocks jniexample

build:
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" go build -o gmx
	chmod +x ./gmx

clean:
	go clean
	@rm gmx  2>/dev/null || exit 0

vendor:
	go mod tidy
	go mod vendor

lint:
	go fmt ./...
	golint -set_exit_status cmd/... internal/... pkg/...

test: clean
	CGO_CFLAGS="$(CGO_CFLAGS)" TEST_ENV=$(TEST_ENV) go test --short ./...

integration_test: jniexample clean
	go clean -testcache
	CGO_CFLAGS="$(CGO_CFLAGS)" TEST_ENV=$(TEST_ENV) go test  ./...
	docker compose -f ./test-docker-compose.yml down

mocks: _mock_gen vendor

_mock_gen:
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" mockery --all --inpackage

jniexample:
	docker compose -f ./test-docker-compose.yml up -d

stop:
	docker compose -f ./test-docker-compose.yml down