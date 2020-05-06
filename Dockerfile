FROM golang:1.14.2 as builder
ENV PROTOC_VERSION 3.11.4
WORKDIR /rpc
# Install recommendation
RUN apt-get update && apt-get install -y --no-install-recommends unzip \
    &&export LOG_LEVEL=debug \
    && export LOG_MODE=development \
    && export GOSUMDB=off \
    && git config --global user.name "go-builder" \
    && git config --global user.email "builder@teko.vn" \
    && git config --global credential.helper store \
    && echo 'https://golang-builder:YRMt66DQZ4nXzCztJNy6@git.teko.vn' > ~/.git-credentials \
    && curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip \
    && unzip -o protoc-${PROTOC_VERSION}-linux-x86_64.zip -d /usr/local bin/protoc \
    && unzip -o protoc-${PROTOC_VERSION}-linux-x86_64.zip -d /usr/local 'include/*' \
    && go get -u git.teko.vn/shared/rpc-framework/cmd/protoc-gen-rpc-server
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN git clone https://git.teko.vn/shared/rpc.git shared/rpc
RUN make all

## Today ubuntu is using minimalized image by default, using ubuntu for better compatible than alpine
FROM ubuntu:20.04
WORKDIR /rpc/bin/
COPY --from=builder /rpc/bin/rpc-runtime /rpc/bin/
COPY migrations ./
EXPOSE 10080 10433
CMD ["/rpc/bin/rpc-runtime", "server"]
