FROM golang:alpine

RUN apk -U --no-cache add git ca-certificates

COPY integration_test.go /test/
COPY run-test.sh /test/

RUN adduser -s /bin/nologin -H -D go-test go-test \
 && chown -R go-test:go-test /test

USER go-test
ENV GOPATH=/test/lib

CMD /test/run-test.sh
