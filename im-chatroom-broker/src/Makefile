#GO_ENV := "dev"
#GOOS := 'linux'
#VERSION := "v0.0.1"

# 二进制文件名
MAIN_PATH :="broker.go"
VERSION := "imchatroom-broker v0.0.1"
DATE := `date +%FT%T%z`

#
#OS := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
#ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))

version:
	@echo $(VERSION)

# .PHONY 有 build 文件，不影响 build 命令执行
.PHONY: build

resource:
	@echo Current env haha
	@echo Copy resources

build-mac:
	GOOS := 'Mac'
	echo version:$(VERSION) date:$(DATE) os: $(GOOS)
	OS := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
	ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))
	GOOS=$(OS) GOARCH=$(ARCH) go build broker.go

install:
	go mod download
	go mod tidy

# 交叉编译运行在linux系统环境
build-linux:
	GOOS := 'Linux'
	echo version:$(VERSION) date:$(DATE) os: $(GOOS)
	OS := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
	ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))
	GOOS=$(OS) GOARCH=$(ARCH) go build broker.go

clean:
	rm broker
