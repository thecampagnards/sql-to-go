FROM golang:alpine

COPY sql-to-go /usr/bin/sql-to-go

ENTRYPOINT ["/usr/bin/sql-to-go"]