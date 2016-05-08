GB ?= gb
GO ?= go
UPX ?= upx
PACK ?= 0

DEPS = \
	go get github.com/constabulary/gb/...


all: bin/q3stats bin/q3simport

bin/q3stats: $(shell find src/q3stats -name '*.go')
	$(GB) build q3stats

packed: bin/q3stats.upx bin/q3simport.upx

%.upx: %
	$(UPX) $^ -o $@

bin/q3simport: $(shell find src/cmd/q3simport -name '*.go')
	$(GB) build cmd/q3simport

clean:
	rm -rf bin
	rm -f *.upx

get-deps:
	for d in $(DEPS); do \
		$(GO) get $$d; \
	done

test:
	$(GO) test -v ./...

.PHONY: clean all packed get-deps test
