TEST?=./...
COMMIT = $$(git describe --always)
NAME = "$(shell awk -F\" '/^const Name/ { print $$2; exit }' version.go)"
VERSION = "$(shell awk -F\" '/^const Version/ { print $$2; exit }' version.go)"

default: test

deps:
	go get -d -t ./...
	go get golang.org/x/tools/cmd/cover

depsdev:
	go get -u github.com/linyows/mflag
	go get -u github.com/mattn/go-shellwords
	go get -u github.com/mitchellh/gox
	go get -u github.com/tcnksm/ghr

test: deps
	go test -v $(TEST) $(TESTARGS) -timeout=30s -parallel=4
	go test -race $(TEST) $(TESTARGS)
	go vet .

cover: deps
	go test $(TEST) -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

bin: depsdev
	@sh -c "'$(CURDIR)/scripts/build.sh' $(NAME)"

dist: bin
	ghr v$(VERSION) pkg

.PHONY: default bin dist test testrace deps
