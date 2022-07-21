GO_ENV := "dev"
GOOS := 'linux'

OS := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build imchatroombroker.go