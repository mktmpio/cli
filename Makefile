SRC = $(wildcard *.go */*.go)
V := $(shell git describe --tags)
T := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GOFLAGS += -ldflags "-X main.version=$V -X main.buildtime=$T"

OSES ?= darwin linux windows
darwin_ARCH := 386 amd64
linux_ARCH := 386 amd64
windows_ARCH := 386 amd64
windows_EXT := .exe
windows_TGZ := zip
default_TGZ := tgz
NAME_386 := x86
NAME_amd64 := x64

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

test: cli get $(BINARIES)
	go test -v ./...
	./cli help
	./cli legal
	./cli --version

get:
	go get -t -v ./...

cli: | get
cli: ${SRC}
	go build ${GOFLAGS}

release: get $(TARBALLS)

# All binaries are built using the same recipe
$(BINARIES): ${SRC} | get
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $@ $(GOFLAGS) mktmpio.go

$(DIRS): README.md LICENSE
	mkdir -p $@
	cp $^ $@

# How to build an archive from a directory
%.zip : %
	zip -r $@ $<
%.tgz : %
	tar -czf $@ $<
