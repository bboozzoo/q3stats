GO ?= go

all: q3stats q3simport

q3stats:
	$(GO) build -v

q3simport:
	$(GO) build -v ./cmd/q3simport

clean:
	$(GO) clean

.PHONY: clean q3stats q3simport all
