FROM golang:alpine AS build

RUN mkdir -p /go/src/github.com/zdnscloud/cddemo
COPY . /go/src/github.com/zdnscloud/cddemo

WORKDIR /go/src/github.com/zdnscloud/cddemo
RUN CGO_ENABLED=0 GOOS=linux go build cddemo.go 

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build /go/src/github.com/zdnscloud/cddemo/cddemo /usr/local/bin/

ENTRYPOINT ["cddemo"]
