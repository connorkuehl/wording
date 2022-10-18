FROM golang:1.18.0-alpine AS build
RUN apk add git
RUN mkdir /builddir
ADD . /builddir
WORKDIR /builddir
RUN go build .

FROM alpine:latest
COPY --from=build /builddir/wording /usr/local/bin/wording

ENTRYPOINT ["/usr/local/bin/wording"]
