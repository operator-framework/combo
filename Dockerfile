# Build the binary
FROM golang:1.17-buster AS builder

WORKDIR /build

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY api api
COPY cmd cmd
COPY pkg pkg
COPY main.go main.go
COPY Makefile Makefile

# copy git-related information for binary version information
COPY .git/refs .git/refs
COPY .git/HEAD .git/HEAD
COPY .git/objects .git/objects

RUN make build-cli

# Copy the binary over to a distroless image and run it
FROM gcr.io/distroless/static:nonroot

WORKDIR /bin

COPY --from=builder /build/bin/combo .

EXPOSE 8080

CMD ["/bin/combo"]
