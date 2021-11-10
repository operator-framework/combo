FROM golang:1.17

WORKDIR /

COPY api api
COPY cmd cmd
COPY pkg pkg
COPY test test
COPY main.go main.go
COPY Makefile Makefile
COPY tools.go tools.go
COPY go.mod go.mod
COPY go.sum go.sum

RUN make build-cli

CMD ./bin/combo run
