# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" "$GOPATH/src/generator" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH/src/generator

ADD . ./
RUN go get ./...

# Run the outyet command by default when the container starts.
ENTRYPOINT ["/go/bin/generator"]
