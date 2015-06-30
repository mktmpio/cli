SRC = $(wildcard *.go)
PLATFORMS = linux/amd64 windows/amd64 darwin/amd64
V = $(shell git describe --tags)
GOFLAGS += -ldflags "-X main.version $V"

mktmpio: ${SRC}
	go get -t -v ./...
	go test -v ./...
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

mktmpio-v$V-windows-x64/mktmpio.exe: ${SRC} mktmpio-v$V-windows-x64
	gox ${GOFLAGS} -osarch="windows/amd64" -output $(basename $@)

mktmpio-v$V-linux-x64/mktmpio: ${SRC} mktmpio-v$V-linux-x64
	gox ${GOFLAGS} -osarch="linux/amd64" -output $@

mktmpio-v$V-darwin-x64/mktmpio: ${SRC} mktmpio-v$V-darwin-x64
	gox ${GOFLAGS} -osarch="darwin/amd64" -output $@
