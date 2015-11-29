SRC = $(wildcard *.go */*.go)
V := $(shell git describe --tags)
T := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GOFLAGS += -ldflags "-X main.version=$V -X main.buildtime=$T"

test: cli
	go test -v ./...
	./cli --version

get:
	go get -t -v ./...

cli: get
cli: ${SRC}
	go build ${GOFLAGS}

release: mktmpio-v$V-windows-x64.zip mktmpio-v$V-linux-x64.tgz mktmpio-v$V-darwin-x64.tgz
release: mktmpio-v$V-windows-x86.zip mktmpio-v$V-linux-x86.tgz mktmpio-v$V-darwin-x86.tgz

mktmpio-v$V-windows-x64 mktmpio-v$V-linux-x64 mktmpio-v$V-darwin-x64 mktmpio-v$V-windows-x86 mktmpio-v$V-linux-x86 mktmpio-v$V-darwin-x86: README.md
	mkdir -p $@
	cp $< $@

mktmpio-v$V-windows-x64/mktmpio.exe: get
mktmpio-v$V-windows-x86/mktmpio.exe: get
mktmpio-v$V-linux-x64/mktmpio: get
mktmpio-v$V-linux-x86/mktmpio: get
mktmpio-v$V-darwin-x64/mktmpio: get
mktmpio-v$V-darwin-x86/mktmpio: get

mktmpio-v$V-windows-x64.zip: mktmpio-v$V-windows-x64/mktmpio.exe
	zip -r $@ $(basename $@)

mktmpio-v$V-windows-x86.zip: mktmpio-v$V-windows-x86/mktmpio.exe
	zip -r $@ $(basename $@)

mktmpio-v$V-linux-x64.tgz: mktmpio-v$V-linux-x64/mktmpio
	tar -czf $@ $(basename $@)

mktmpio-v$V-linux-x86.tgz: mktmpio-v$V-linux-x86/mktmpio
	tar -czf $@ $(basename $@)

mktmpio-v$V-darwin-x64.tgz: mktmpio-v$V-darwin-x64/mktmpio
	tar -czf $@ $(basename $@)

mktmpio-v$V-darwin-x86.tgz: mktmpio-v$V-darwin-x86/mktmpio
	tar -czf $@ $(basename $@)

mktmpio-v$V-windows-x64/mktmpio.exe: ${SRC}
	GOOS=windows GOARCH=amd64 go build -o $@ ${GOFLAGS} $<

mktmpio-v$V-windows-x86/mktmpio.exe: ${SRC}
	GOOS=windows GOARCH=386 go build -o $@ ${GOFLAGS} $<

mktmpio-v$V-linux-x64/mktmpio: ${SRC}
	GOOS=linux GOARCH=amd64 go build -o $@ ${GOFLAGS} $<

mktmpio-v$V-linux-x86/mktmpio: ${SRC}
	GOOS=linux GOARCH=386 go build -o $@ ${GOFLAGS} $<

mktmpio-v$V-darwin-x64/mktmpio: ${SRC}
	GOOS=darwin GOARCH=amd64 go build -o $@ ${GOFLAGS} $<

mktmpio-v$V-darwin-x86/mktmpio: ${SRC}
	GOOS=darwin GOARCH=386 go build -o $@ ${GOFLAGS} $<
