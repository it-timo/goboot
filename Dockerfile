ARG GO_VERSION=1.26.5

FROM golang:${GO_VERSION} AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o /out/goboot \
    ./cmd/goboot

FROM golang:${GO_VERSION}

WORKDIR /workdir

COPY --from=build /out/goboot /usr/local/bin/goboot

ENTRYPOINT ["goboot"]
