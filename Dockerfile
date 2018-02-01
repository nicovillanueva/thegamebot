FROM golang:1.9.3-alpine3.7 AS build
RUN apk --no-cache add git && \
    go get github.com/golang/dep/cmd/dep

COPY Gopkg.toml Gopkg.lock /go/src/project/
WORKDIR /go/src/project/
RUN dep ensure -vendor-only

COPY . /go/src/project/
RUN go build -o /bin/compiled

# TODO: migrate to scratch
FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=build /bin/compiled /bin/gamebot
STOPSIGNAL 9
ENTRYPOINT ["/bin/gamebot"]
