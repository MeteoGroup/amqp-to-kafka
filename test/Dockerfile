FROM golang:alpine

RUN apk -U --no-cache add git ca-certificates

COPY ./ /test/

RUN adduser -s /bin/nologin -H -D go-test go-test \
 && chown -R go-test:go-test /test
 && pip install virtualenv \

RUN pip install requests[security]

USER go-test
ENV GOPATH=/test/lib

CMD /test/run-test.sh
