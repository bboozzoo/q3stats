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

TOPDIR = $(shell pwd)

GIT_VERSION = $(shell git describe --always --tags --dirty)
BUILD_DATE = $(shell date +%Y-%m-%d\ %H:%M)
VERSION = $(GIT_VERSION) $(BUILD_DATE)

LDFLAGS = -ldflags \
	"-X 'github.com/bboozzoo/q3stats/version.versionInfo=$(VERSION)'"

all: q3stats q3simport

q3stats: assets/static/assets.go assets/templates/assets.go
	$(GO) build $(LDFLAGS) -v

assets/static/assets.go: $(shell find webroot/static)
	cd webroot/static && \
		$(ESC) -pkg static -o $(TOPDIR)/$@ .

assets/templates/assets.go: $(shell find webroot/templates)
	cd webroot/templates && \
		$(ESC) -pkg templates -o $(TOPDIR)/$@ .

packed: q3stats.upx q3simport.upx

%.upx: %
	$(UPX) $^ -o $@

q3simport:
	$(GO) build $(LDFLAGS) -v ./cmd/q3simport

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
