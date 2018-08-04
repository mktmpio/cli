SRC = $(wildcard *.go */*.go)
Q ?= @
vV := $(shell git describe --tags)
V := $(vV:v%=%)
T := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GOFLAGS += -ldflags "-X main.version=$V -X main.buildtime=$T"
GO ?= go

OSES ?= darwin linux windows
darwin_ARCH := 386 amd64
linux_ARCH := 386 amd64
windows_ARCH := 386 amd64
windows_EXT := .exe
windows_TGZ := zip
default_TGZ := tgz
NAME_386 := x86
NAME_amd64 := x64

default: quicktest

name = mktmpio-v$(V)-$(1)-$(firstword $(NAME_$(2)) $(2))
define build_t
$(3)/mktmpio$($(1)_EXT): $(3)
$(3)/mktmpio$($(1)_EXT): GOOS=$(1)
$(3)/mktmpio$($(1)_EXT): GOARCH=$(2)
$(3).$(firstword $($(1)_TGZ) $(default_TGZ)): $(3)/mktmpio$($(1)_EXT)
$(eval TARBALLS += $(3).$(firstword $($(1)_TGZ) $(default_TGZ)))
$(eval BINARIES += $(3)/mktmpio$($(1)_EXT))
$(eval DIRS += $(3))
endef

# Generate targets and variables for all the supported OS/ARCH combinations
$(foreach os,$(OSES), \
	$(foreach arch,$($(os)_ARCH), \
		$(eval \
			$(call build_t,$(os),$(arch),$(call name,$(os),$(arch))) \
		) \
	) \
)

test: Q =
test: quicktest $(BINARIES)

quicktest: cli
	$Q $(GO) test ./commands && echo "ok - test"
	$Q ./cli help | grep -q "mktmpio" && echo "ok - help"
	$Q ./cli legal | grep -q "Artistic" && echo "ok - legal"
	$Q ./cli --version | grep -q "mktmpio" && echo "ok - version"
	$Q echo "ok"

get:
	$(GO) get -t -v ./...

cli: ${SRC}
	$Q $(GO) build -o $@ ${GOFLAGS} ./cmd/mktmpio

release: $(TARBALLS)

# All binaries are built using the same recipe
$(BINARIES): ${SRC}
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -o $@ $(GOFLAGS) ./cmd/mktmpio

# When doing a big parallel build ensure that cli finishes first so that we
# don't have multiple 'go build' processes trying to fetch the same
# dependencies at the same time and tripping over each other's lock files.
$(BINARIES): | cli

$(DIRS): README.md LICENSE
	mkdir -p $@
	cp $^ $@

# How to build an archive from a directory
%.zip : %
	zip -r $@ $<
%.tgz : %
	tar -czf $@ $<
