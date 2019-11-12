# Environment Variables
CGO=1
OS=darwin
SRC=cmd
BUILD=build

default: main

main: clean api 

api:
	env CGO_ENABLED=$(CGO) GOOS=$(OS) go build -o $(BUILD)/cave-logger $(SRC)/main.go

clean:
	rm -f $(BUILD)/*
	touch $(BUILD)/.keep

install: api
	mv $(BUILD)/cave-logger $(GOPATH)/bin/.