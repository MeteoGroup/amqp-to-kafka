#!/bin/sh -ex

cd "`dirname \"$0\"`"

go get -d -t
go test
