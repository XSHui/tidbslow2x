GOVER := $(shell go version)

GOOS    := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
GOARCH  := $(if $(GOARCH),$(GOARCH),amd64)
#GOENV   := GO111MODULE=on CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH)
GOENV   := GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH)
GO      := $(GOENV) go
GOBUILD := $(GO) build $(BUILD_FLAG)

tidbslow2x:
	$(GOBUILD) -o bin/tidbslow2x 

clean:
	rm ./bin/*
