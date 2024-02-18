FROM alpine:3.18.4

ADD app-query /usr/local/da/query/app-query

RUN chmod +x /usr/local/da/query/app-query

WORKDIR /usr/local/da/query

CMD ["./app-query"]