FROM golang:alpine AS build
RUN apk update && apk add git

WORKDIR /go/src/github.com/darookee/adguard_exporter/

ADD . .
RUN go get -t -v ./...
RUN go build -o /adguard_exporter main.go && chmod +x /adguard_exporter

FROM alpine:latest

RUN apk add curl
COPY --from=build /adguard_exporter /adguard_exporter
WORKDIR /
ENTRYPOINT ["./adguard_exporter"]
CMD ["-h"]
EXPOSE 9311
