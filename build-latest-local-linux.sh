#!/bin/bash

# ensures local build that prevents error message
# standard_init_linux.go:178: exec user process caused "no such file or directory"

# when you are only changing tests, it is enough to call
# test/test.sh local-test
# subsequently
rm amqp-to-kafka
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o amqp-to-kafka
docker build -t "meteogroup/amqp-to-kafka" .
