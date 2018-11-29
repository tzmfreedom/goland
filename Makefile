ANTLR := java -Xmx500M -cp "/usr/local/lib/antlr-4.7.1-complete.jar:$(CLASSPATH)" org.antlr.v4.Tool

.PHONY: run
run: format
	@go run . run -a "Foo#action" -f example.cls

.PHONY: run/format
run/format:
	@go run . format -f example.cls

.PHONY: test
test: format
	@go test ./...

.PHONY: build
build: format
	@go build

.PHONY: format
format: import
	@gofmt -w .

.PHONY: import
import:
ifneq ($(shell command -v goimports 2> /dev/null),)
	@goimports -w compiler/ ast/ visitor/ interpreter/ builtin/ server/ ./land.go
endif

.PHONY: generate
generate:
	cd ./parser; \
	$(ANTLR) -Dlanguage=Go -visitor apex.g4

.PHONY: dep
dep:
ifeq ($(shell command -v dep 2> /dev/null),)
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif

.PHONY: deps
deps: dep
	dep ensure

