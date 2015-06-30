SRC = $(wildcard *.go)
PLATFORMS = linux/amd64 windows/amd64 darwin/amd64
V = $(shell git describe --tags)
GOFLAGS += -ldflags "-X main.version $V"

test: cli
	go test -v ./...
	./cli --version

cli: ${SRC}
	go get -t -v ./...
	go build ${GOFLAGS}

release: mktmpio-v$V-windows-x64.zip mktmpio-v$V-linux-x64.tgz mktmpio-v$V-darwin-x64.tgz

gox:
	go get github.com/mitchellh/gox
	gox -osarch="${PLATFORMS}" -build-toolchain

mktmpio-v$V-windows-x64 mktmpio-v$V-linux-x64 mktmpio-v$V-darwin-x64: README.md
	mkdir -p $@
	cp $< $@

mktmpio-v$V-windows-x64.zip: mktmpio-v$V-windows-x64/mktmpio.exe
	zip -r $@ $(basename $@)

mktmpio-v$V-linux-x64.tgz: mktmpio-v$V-linux-x64/mktmpio
	tar -czf $@ $(basename $@)

mktmpio-v$V-darwin-x64.tgz: mktmpio-v$V-darwin-x64/mktmpio
	tar -czf $@ $(basename $@)

mktmpio-v$V-windows-x64/mktmpio.exe: gox-build
mktmpio-v$V-linux-x64/mktmpio: gox-build
mktmpio-v$V-darwin-x64/mktmpio: gox-build

gox-build: ${SRC} mktmpio-v$V-windows-x64 mktmpio-v$V-linux-x64 mktmpio-v$V-darwin-x64
	gox ${GOFLAGS} -osarch="windows/amd64 linux/amd64 darwin/amd64" -output "mktmpio-v$V-{{ .OS }}-x64/mktmpio"
