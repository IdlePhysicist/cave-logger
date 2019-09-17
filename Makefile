# Environment Variables
CGO=1
OS=darwin
SRCDIR=cmd
BUILDDIR=build

default: main

main: clean api 

api:
	env CGO_ENABLED=$(CGO) GOOS=$(OS) go build $(SRCDIR)/main.go

clean:
	rm -f $(BUILDDIR)/*
	touch $(BUILDDIR)/.keep
