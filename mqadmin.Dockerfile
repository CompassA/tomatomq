# 构建镜像
# podman build -t compassa/tomatomq-admin:1.0.0-dev  -f ./mqadmin.Dockerfile .

# 查看本地etcd、mysql容器的ip信息
# podman inspect --format='{{.NetworkSettings.IPAddress}}' [容器id]

# 启动容器
# podman run -e APP_ENV=dev -d compassa/tomatomq-admin:1.0.0-dev
# podman logs -f [容器id] 查看日志
FROM docker.io/library/golang:1.26.1 AS builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o mqadmin cmd/admin/main.go


FROM scratch
ARG CONFIG_FILE_NAME=admin-dev.yaml
WORKDIR /
COPY --from=builder /workspace/mqadmin .
COPY --from=builder /workspace/config/$CONFIG_FILE_NAME .
EXPOSE 8080 8090
ENTRYPOINT ["/mqadmin"]