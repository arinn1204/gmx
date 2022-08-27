INCLUDEFLAGS		:= -I$(JAVA_HOME)/include
LINKERFLAGS 		:= -L$(JAVA_HOME)/lib -ljvm
CGO_CFLAGS       	:= $(INCLUDEFLAGS)
CGO_LDFLAGS 		:= $(LINKERFLAGS)
JAVAC				:= javac

.PHONY: build clean vendor test integration_test mocks jniexample

all: clean mocks vendor lint integration_test stop release_patch

install:
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" go install github.com/arinn1204/gmx

build:
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" go build -o gmx
	@chmod +x ./gmx

clean:
	go clean
	@rm -r dist || exit 0
	@rm gmx  2>/dev/null || exit 0

snapshot:
	goreleaser release --snapshot --rm-dist -f .mac.goreleaser.yaml

release_patch: clean mocks vendor lint integration_test
	bash scripts/dirty.sh
	bash scripts/increment_tag.sh
	goreleaser release --rm-dist -f .mac.goreleaser.yaml

release: clean mocks vendor lint integration_test
	bash scripts/dirty.sh
	goreleaser release --rm-dist -f .mac.goreleaser.yaml

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