# List of implements server
SERVERS = health/v1/health.proto sample/v1/sample.proto
# Enable grpc gateway
GATEWAY = true

TARGET = bin
TARGET_BIN = rpc-runtime
GO_CMD_MAIN = cmd/main.go

RPC_FOLDER = shared/rpc

SERVER_PACKAGE_NAME = server
SERVER_OUT_FOLDER = rpcimpl

####################  DOES NOT EDIT BELLOW  ############################
.PHONY = build generate all clean

GO_TOOLS = git.teko.vn/shared/rpc-framework/cmd/protoc-gen-rpc-server

$(GO_TOOLS):
	GOSUMDB=off go get -u $@

# support fresh install on osx, not sure if it can't run properly
install-osx: $(GO_TOOLS)
	brew install protobuf

# support fresh install on linux, not sure if it can't run properly
PROTOC_LINUX_VERSION = 3.11.4
PROTOC_LINUX_ZIP = protoc-$(PROTOC_LINUX_VERSION)-linux-x86_64.zip

install-linux: $(GO_TOOLS)
	curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_LINUX_VERSION)/$(PROTOC_LINUX_ZIP)
	sudo unzip -o $(PROTOC_LINUX_ZIP) -d /usr/local bin/protoc
	sudo unzip -o $(PROTOC_LINUX_ZIP) -d /usr/local 'include/*'
	rm -f $(PROTOC_LINUX_ZIP)

$(RPC_FOLDER):
	@echo '# update $@'
	@[[ -d $@ ]] || git clone git@git.teko.vn:$@.git $@
	@cd $@ && git checkout master && git pull origin master

.ONESHELL:
common-update: $(RPC_FOLDER)
	GOSUMDB=off go get -u go.tekoapis.com/kitchen/...
	GOSUMDB=off go get -u rpc.tekoapis.com/...

prepare:
	mkdir -p $(SERVER_OUT_FOLDER)/$(SERVER_PACKAGE_NAME)

photon-server:
	@echo \# generating photon-server....
	protoc -I $(RPC_FOLDER)/proto \
		-I $(RPC_FOLDER)/.third_party/googleapis \
		-I $(RPC_FOLDER)/.third_party/envoyproxy \
		--rpc-server_out=gateway=$(GATEWAY):$(SERVER_OUT_FOLDER)/$(SERVER_PACKAGE_NAME) \
		$(SERVERS)

generate: prepare photon-server
	@echo \# source code is generated

build: generate
	go build -o $(TARGET)/$(TARGET_BIN) $(GO_CMD_MAIN)

run: generate
	go run $(GO_CMD_MAIN) server

migrate:
	echo \# make migrate name="$(name)"
	go run $(GO_CMD_MAIN) migrate create $(name)

migrate-up:
	go run $(GO_CMD_MAIN) migrate up

migrate-down-1:
	go run $(GO_CMD_MAIN) migrate down 1

all: common-update build
	@echo what is done is done!

clean:
	rm -rf $(SERVER_OUT_FOLDER)/$(SERVER_PACKAGE_NAME)
	rm -rf $(TARGET)
