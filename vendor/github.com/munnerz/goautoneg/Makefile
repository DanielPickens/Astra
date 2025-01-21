include $(GOROOT)/src/Make.inc

TARG=bitbucket.org/ww/goautoneg
GOFILES=autoneg.go

include $(GOROOT)/src/Make.pkg

format:
	gofmt -w *.go

docs:
	gomake clean
	gastrac ${TARG} > README.txt
