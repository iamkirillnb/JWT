FROM golang:alpine as builder

RUN apk --no-cache  add tzdata git gcc g++ libwebp

ARG SERVICE_NAME

WORKDIR /src
COPY . .

RUN GIT_COMMIT=$(git rev-list -1 HEAD) && BUILD_AT=$(date +"%d-%m-%Y,%T") && VERSION="latest" && \
    CGO_ENABLE=0 go build -ldflags "-X main.BuildAt=$BUILD_AT -X main.Commit=$GIT_COMMIT) -X main.Version=$VERSION" -v -o /bundle ./cmd/${SERVICE_NAME}/...

FROM alpine:latest

RUN GRPC_HEALTH_PROBE_VERSION=v0.3.6 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /bundle /bundle
ENV TZ=Europe/Moscow

CMD ["/bundle"]