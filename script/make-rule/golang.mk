#===================================  target  ====================================#
BINARY_NAME = demoApp

# path to main source code
ROOT_PACKAGE = MarineSnow

GO := go

#===================================  flag and option  ====================================#
GO_BUILD_FLAGS := -v

# debug flag for dlv
ifeq ($(origin DEBUG),command line)
	GO_BUILD_FLAGS += -gcflags="all=-N -l"
endif

# data race detection for go
ifeq ($(origin RACE),command line)
	GO_BUILD_FLAGS += -race
endif

# golangci-lint
LINT_ENABLE := 1
ifeq ($(origin LINT),command line)
	LINT_ENABLE := 0
endif

#===================================  CMD target  ====================================#
# Usage:
# 1. make build
# 2. make build DEBUG=1
.PHONY: go.build
go.build: go.lint
	@echo -e "\033[1;36m-->[CMD:go.build, go-flag:$(GO_BUILD_FLAGS), bins:$(BINARY_NAME), root-dir:$(ROOT_DIR)]\033[0m"
	$(GO) build $(GO_BUILD_FLAGS) -o $(BINARY_NAME) .

# Usage:
# 1. make run
# 2. make run DEBUG=1
RUN_CMD := ./$(BINARY_NAME) app start
.PHONY: go.run
go.run:
	@echo -e "\033[1;36m-->[CMD:$(RUN_CMD)]\033[0m"
	$(RUN_CMD)

# Usage:
# 1. make checkout
# 2. make go.lint
# 2. make run LINT=0
.PHONY: go.lint
go.lint: tools.verify.golangci-lint
ifeq ($(LINT_ENABLE), 1)
	@echo -e "\033[1;36m-->[CMD:golangci-lint run, root-dir:$(ROOT_DIR)]\033[0m"
	golangci-lint run ./...
else
	@echo -e "\033[1;36m-->[CMD:go.vet, root-dir:$(ROOT_DIR)]\033[0m"
#	    $(GO) vet ./framework/...
#	    $(GO) vet ./app/...
#	    $(GO) vet ./provider/...
	$(GO) vet main.go
endif

# Usage:
# 1. make checkout
# 2. make go.fmt
.PHONY: go.fmt
go.fmt:
	@echo -e "\033[1;36m-->[CMD:go.fmt, root-dir:$(ROOT_DIR)]\033[0m"
	$(GO) fmt ./...

# Usage:
# 1. make checkout
# 2. make clean
# 2. make go.clean
.PHONY: go.clean
go.clean:
	@echo -e "\033[1;36m-->[CMD:go.clean, root-dir:$(ROOT_DIR)]\033[0m"
	go clean -i ./...
	@if [ -f $(BINARY_NAME) ] ; then rm $(BINARY_NAME) ; fi
