FROM golang:1.22-alpine AS builder

WORKDIR /usr/local/src

COPY ["go.mod", "./"]
RUN --mount=type=cache,target=/go/pkg/mod/ \
    go mod download -x
COPY . .

RUN go build -o ./bin/app ./*.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /app

ENTRYPOINT [ "/app" ]