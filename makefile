PACKAGES ?= ./...
FLAGS ?=
DEVTOOLS ?= $(shell $(GOPATH)/bin/golangci-lint --version 2>/dev/null | grep 1.24)

test:
	go test -race -count=1 $(FLAGS) $(PACKAGES) -cover | tee coverage.out
	echo "\e[1m====================================="
	grep -Po "[0-9]+\.[0-9]+(?=%)" coverage.out | awk '{ SUM += $$1; PKGS += 1} END { print "  Total Coverage (" PKGS " pkg/s) : " SUM/PKGS "%"}'
	echo "=====================================\e[0m"
	rm -f coverage.out
.PHONY: test
.SILENT: test

cover:
	go test -race -count=1 $(PACKAGES) -coverprofile=coverage.out && go tool cover -html=coverage.out
	rm -f coverage.out
.SILENT: cover

lint: devtools
	$(GOPATH)/bin/golangci-lint run $(PACKAGES) -c ./.golangci.yml
.SILENT: lint

devtools:
ifeq ($(strip $(DEVTOOLS)),)
	echo "\e[1mDEVTOOLS not present, installing...\e[0m"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.24.0
endif
.SILENT: devtools

clean:
	rm -f *.exe *.out
