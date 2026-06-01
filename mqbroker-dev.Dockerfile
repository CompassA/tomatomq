# 构建镜像
# podman build --build-arg CONFIG_FILE_NAME=broker-config-podman-a.yaml -t compassa/tomatomq-broker-a:1.0.0-dev  -f ./mqbroker-dev.Dockerfile .
# podman build --build-arg CONFIG_FILE_NAME=broker-config-podman-b.yaml -t compassa/tomatomq-broker-b:1.0.0-dev  -f ./mqbroker-dev.Dockerfile .

# 查看本地etcd、mysql容器的ip信息
# podman inspect --format='{{.NetworkSettings.IPAddress}}' [容器id]

# 启动容器
# podman run -e APP_ENV=dev -p 6778:6778 -e POD_IP=192.168.112.130 -d compassa/tomatomq-broker-a:1.0.0-dev
# podman run -e APP_ENV=dev -p 6779:6779 -e POD_IP=192.168.112.130 -d compassa/tomatomq-broker-b:1.0.0-dev
# podman logs -f [容器id] 查看日志
# podman exec my_etcd etcdctl get --prefix /broker   查看注册信息
FROM docker.io/library/golang:1.26.1 AS builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o broker cmd/broker/main.go


FROM scratch
ARG CONFIG_FILE_NAME
WORKDIR /
COPY --from=builder /workspace/broker .
COPY --from=builder /workspace/config/$CONFIG_FILE_NAME ./broker-config.yaml
EXPOSE 6778 6779
ENTRYPOINT ["/broker"]