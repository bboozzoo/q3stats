GO ?= go
UPX ?= upx
PACK ?= 0

DEPS = \
	github.com/codegangsta/cli \
	github.com/gorilla/mux \
	github.com/gorilla/handlers \
	github.com/jinzhu/gorm \
	github.com/jinzhu/gorm/dialects/sqlite \
	gopkg.in/gcfg.v1


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

get-deps:
	for d in $(DEPS); do \
		$(GO) get $$d; \
	done

.PHONY: clean q3stats q3simport all packed get-deps
