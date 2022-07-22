GO_ENV := "dev"
GOOS := 'linux'

OS := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))

build:
	GOOS=$(OS) GOARCH=$(ARCH) go build imchatroombroker.go




# 二进制文件名
PROJECT="main"
# 启动文件PATH
MAIN_PATH="main.go"
VERSION="v0.0.1"
DATE= `date +%FT%T%z`

# 可以执行函数，但是函数不能存在 command
ifeq (${VERSION}, "v0.0.1")
    VERSION=VERSION = "v0.0.1"
endif

version:
    @echo ${VERSION}

# .PHONY 有 build 文件，不影响 build 命令执行
.PHONY: build
build:
    @echo version: ${VERSION} date: ${DATE} os: Mac OS
    @go  build -o ${PROJECT} ${MAIN_PATH}

install:
    @echo download package
    @go mod download

# 交叉编译运行在linux系统环境
build-linux:
    @echo version: ${VERSION} date: ${DATE} os: linux-centOS
    @GOOS=linux go build -o ${PROJECT} ${MAIN_PATH}

run:   build
    @./${PROJECT}

clean:
    rm -rf ./log