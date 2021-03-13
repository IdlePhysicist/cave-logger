# Environment Variables
SRC   =cmd
BUILD =build
PREFIX=$(GOPATH)/bin/

version?=`if [ -d ./.git ]; then git describe --tags; else echo release-build; fi`
date    =`date "+%Y-%m-%d"`
package =main
ldflags ="-X $(package).version=$(version) -X $(package).date=$(date)"

default: build

build: clean
	env go build -ldflags $(ldflags) -o $(BUILD)/cave-logger $(SRC)/main.go

clean:
	rm -f $(BUILD)/*
	touch $(BUILD)/.keep

install:
	mv $(BUILD)/cave-logger $(PREFIX)
