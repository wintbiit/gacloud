FROM golang:1.22-alpine as backend

RUN apk add --no-cache git

WORKDIR /go/src/gacloud
ADD go.mod go.sum ./
RUN go mod download

ADD . .

FROM backend as daemon
RUN version=$(git describe --tags --always --dirty) && \
    go build -ldflags "-X utils.version='$version' -s -w" -o ./out/gacloud-server ./cmd/daemon/main.go

FROM backend as cli
RUN version=$(git describe --tags --always --dirty) && \
    go build -ldflags "-X utils.version='$version' -s -w" -o ./out/gacloud ./cmd/cli/main.go

FROM node:20-slim as frontend
RUN npm install -g pnpm

WORKDIR /app
ADD ./web/package.json ./web/package-lock.json ./
RUN pnpm install --frozen-lockfile

COPY ./web .

RUN pnpm build

FROM alpine:3.14
ENV GACLOUD_DATA_DIR=/var/lib/gacloud
ENV GACLOUD_WEB_DIR=/usr/local/share/gacloud-web
ENV GACLOUD_LOG_DIR=/var/log/gacloud

RUN apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    groupadd -r gacloud && useradd -r -g gacloud gacloud && \
    mkdir -p /home/gacloud && chown -R gacloud:gacloud /home/gacloud && \
    mkdir -p /var/lib/gacloud && chown -R gacloud:gacloud /var/lib/gacloud && \
    mkdir -p /usr/local/share/gacloud-web && chown -R gacloud:gacloud /usr/local/share/gacloud-web \
    mkdir -p /var/log/gacloud && chown -R gacloud:gacloud /var/log/gacloud

COPY --from=daemon /go/src/gacloud/out/gacloud-server /usr/local/bin/gacloud-server
COPY --from=cli /go/src/gacloud/out/gacloud /usr/local/bin/gacloud
COPY --from=frontend /app/dist /usr/local/share/gacloud-web

USER gacloud
WORKDIR /home/gacloud

ENTRYPOINT ["/usr/local/bin/gacloud-server"]