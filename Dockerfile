FROM golang:1.17-buster AS builder

WORKDIR /build

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY api api
COPY cmd cmd
COPY pkg pkg
COPY main.go main.go
COPY tools.go tools.go
COPY Makefile Makefile

RUN make build-cli

FROM golang:1.17-buster

WORKDIR /

COPY --from=builder /build/bin/combo /bin

EXPOSE 8080

CMD ["combo"]