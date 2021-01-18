# syntax = docker/dockerfile:experimental
FROM golang as builder

WORKDIR /cave

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

WORKDIR /cave/cmd/cave-demos/

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -v -o app

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=builder /cave/cmd/cave-demos/app /opt/app
COPY --from=builder /cave/cmd/cave-demos/layout.html /opt/
WORKDIR /opt
CMD ["/opt/app"]
