FROM alpine

RUN apk -U --no-cache upgrade \
 && apk -U --no-cache add ca-certificates

COPY amqp-to-kafka /bin/amqp-to-kafka
RUN chmod 755 /bin/amqp-to-kafka \
 && adduser -s /bin/nologin -H -D amqp-to-kafka amqp-to-kafka

ENV METRICS_ADDRESS=:8080
EXPOSE 8080

USER amqp-to-kafka
CMD amqp-to-kafka
