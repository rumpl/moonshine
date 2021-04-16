FROM golang:1.16-alpine AS base
WORKDIR /app
ENV GO111MODULE=on
COPY go.* .
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

FROM base AS builder
ENV CGO_ENABLED=0
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o /moonshine ./cmd/moonshine

FROM scratch
COPY --from=builder /moonshine /bin/moonshine
ENTRYPOINT ["/bin/moonshine"]
