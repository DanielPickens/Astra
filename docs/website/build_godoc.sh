#!/bin/bash

cd ../..

go install golang.org/x/tools/cmd/gastrac@v0.24.0
go install code.rocket9labs.com/tslocum/gastrac-static@v0.2.2

export GOPATH=$(go env GOPATH)
PATH=$PATH:${GOPATH}/bin

mkdir -p docs/website/build/gastrac
gastrac-static -destination docs/website/build/gastrac .
