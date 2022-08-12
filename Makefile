INCLUDEFLAGS		:= -I$(JAVA_HOME)/include -I$(JAVA_HOME)/include/darwin
LINKERFLAGS 		:= -L$(JAVA_HOME)/lib/server -L$(JAVA_HOME)/lib -ljvm
CGO_CFLAGS       	:= $(INCLUDEFLAGS)
CGO_LDFLAGS 		:= $(LINKERFLAGS)
JAVAC				:= javac

.PHONY: build clean vendor test integration_test mocks jniexample

build:
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" go build -o gmx
	@chmod +x ./gmx

clean:
	go clean
	@rm -r dist
	@rm gmx  2>/dev/null || exit 0

snapshot:
	goreleaser release --snapshot --rm-dist -f .mac.goreleaser.yaml

package_linux:
	docker build --file build/linux.Dockerfile --tag gmx_linux .
	docker run -e GITHUB_TOKEN=$(GITHUB_TOKEN) --mount source=dist,target=/go/src/gmx/dist -it gmx_linux /bin/sh

vendor:
	go mod tidy
	go mod vendor

lint:
	go fmt ./...
	golint -set_exit_status cmd/... internal/... pkg/...

test: clean
	CGO_CFLAGS="$(CGO_CFLAGS)" go test --short ./...

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