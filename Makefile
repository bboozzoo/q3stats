GO ?= go
ESC ?= esc
UPX ?= upx
PACK ?= 0

DEPS = \
	github.com/codegangsta/cli \
	github.com/gorilla/mux \
	github.com/gorilla/handlers \
	github.com/jinzhu/gorm \
	github.com/jinzhu/gorm/dialects/sqlite \
	gopkg.in/gcfg.v1 \
	github.com/pkg/errors \
	github.com/stretchr/testify/assert \
	github.com/mjibson/esc


all: q3stats q3simport

q3stats: assets/assets.go
	$(GO) build -v

assets/assets.go:
	cd webroot && \
		$(ESC) -pkg assets -o ../$@ .

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

test:
	$(GO) test -v ./...

.PHONY: clean q3stats q3simport \
	all packed get-deps test
