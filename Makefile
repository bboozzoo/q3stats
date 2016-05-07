GO ?= go
UPX ?= upx
PACK ?= 0

all: q3stats q3simport

q3stats:
	$(GO) build -v

packed: q3stats.upx q3simport.upx

%.upx: %
	$(UPX) $^ -o $@

q3simport:
	$(GO) build -v ./cmd/q3simport

clean:
	$(GO) clean
	rm -f *.upx q3simport q3stats

.PHONY: clean q3stats q3simport all packed
