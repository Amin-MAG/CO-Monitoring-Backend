FROM golang:1.17.7

ENV GO11MODULE="on"

WORKDIR /go/src/comonitoring

COPY . .
RUN go install -v ./...

EXPOSE 8080

CMD ["comonitoring"]