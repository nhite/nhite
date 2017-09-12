LINUX=bin/nhite-linux-amd64
WINDOWS64=bin/nhite-amd64.exe
DARWIN=bin/nhite-darwin64

# These are the values we want to pass for VERSION and BUILD
# # git tag 1.0.1
# # git commit -am "One more change after the tags"
VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`

LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION} -X main.Build=${BUILD}"

all: $(LINUX) $(WINDOWS64) $(WINDOWS32) $(DARWIN)

$(LINUX): *.go
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o $(LINUX) *.go

$(WINDOWS64): *.go
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o $(WINDOWS64) *.go

$(DARWIN): *.go
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o $(DARWIN) *.go


