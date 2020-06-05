# Environment Variables
CGO   =1
SRC   =cmd
BUILD =build
PREFIX=$(GOPATH)/bin/

version?="0.0.0"
commit  =`if [ -d ./.git ]; then git rev-list -1 HEAD | head -c 8; else echo "release build"; fi`
date    =`date "+%Y-%m-%d"`
package =main
ldflags ="-X $(package).commit=$(commit) -X $(package).version=$(version) -X $(package).date=$(date)"

default: darwin

main: darwin

linux: clean
	env CGO_ENABLED=$(CGO) GOOS=$@ go build -ldflags $(ldflags) -o $(BUILD)/cave-logger $(SRC)/main.go

darwin: clean
	env CGO_ENABLED=$(CGO) GOOS=$@ go build -ldflags $(ldflags) -o $(BUILD)/cave-logger $(SRC)/main.go

clean:
	rm -f $(BUILD)/*
	touch $(BUILD)/.keep

install:
	mv $(BUILD)/cave-logger $(PREFIX)
