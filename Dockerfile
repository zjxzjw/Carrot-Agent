FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o carrot-agent ./cmd/cli

FROM alpine:3.19

RUN apk add --no-cache ca-certificates sqlite-libs

WORKDIR /app

RUN addgroup -g 1000 carrot && \
    adduser -u 1000 -G carrot -s /bin/sh -D carrot

COPY --from=builder /app/carrot-agent /app/
COPY --from=builder /app/config/ /app/config/

RUN mkdir -p /root/.carrot/data /root/.carrot/skills /root/.carrot/memories /root/.carrot/sessions && \
    chown -R carrot:carrot /root/.carrot

USER carrot

ENV CARROT_DATA_DIR=/root/.carrot
ENV CARROT_CONFIG=/app/config/config.yaml

ENTRYPOINT ["/app/carrot-agent"]