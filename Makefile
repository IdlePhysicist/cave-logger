# Environment Variables
CGO=1
SRC=cmd
BUILD=build

default: darwin

main: darwin

linux: clean
	env CGO_ENABLED=$(CGO) GOOS=$@ go build -o $(BUILD)/cave-logger $(SRC)/main.go

darwin: clean
	env CGO_ENABLED=$(CGO) GOOS=$@ go build -o $(BUILD)/cave-logger $(SRC)/main.go

clean:
	rm -f $(BUILD)/*
	touch $(BUILD)/.keep

install:
	mv $(BUILD)/cave-logger $(GOPATH)/bin/.
