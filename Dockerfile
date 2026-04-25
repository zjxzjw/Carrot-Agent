# 构建 API 服务
FROM golang:1.22-alpine AS api-builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o carrot-agent-api ./cmd/api

# 构建 CLI 服务
FROM golang:1.22-alpine AS cli-builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o carrot-agent ./cmd/cli

# 构建前端
FROM node:18-alpine AS frontend-builder

WORKDIR /app

COPY ui/package*.json ./ui/
RUN cd ui && npm install

COPY ui/ ./ui/
RUN cd ui && npm run build

# 最终镜像
FROM alpine:3.19

RUN apk add --no-cache ca-certificates sqlite-libs

WORKDIR /app

RUN addgroup -g 1000 carrot && \
    adduser -u 1000 -G carrot -s /bin/sh -D carrot

# 复制 API 服务
COPY --from=api-builder /app/carrot-agent-api /app/
# 复制 CLI 服务
COPY --from=cli-builder /app/carrot-agent /app/
# 复制配置
COPY --from=api-builder /app/config/ /app/config/
# 复制前端构建产物
COPY --from=frontend-builder /app/ui/dist /app/ui/dist

RUN mkdir -p /home/carrot/.carrot/data /home/carrot/.carrot/skills /home/carrot/.carrot/memories /home/carrot/.carrot/sessions && \
    chown -R carrot:carrot /home/carrot/.carrot

USER carrot

ENV CARROT_DATA_DIR=/home/carrot/.carrot
ENV CARROT_CONFIG=/app/config/config.yaml

# 暴露 API 服务端口
EXPOSE 8080

# 默认启动 API 服务
ENTRYPOINT ["/app/carrot-agent-api"]