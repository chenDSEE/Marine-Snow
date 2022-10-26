
#===================================  include  ====================================#
include script/make-rule/common.mk
include script/make-rule/golang.mk
include script/make-rule/tool.mk


#===================================  target  ====================================#
# Build all by default, even if it's not first
.DEFAULT_GOAL := all

.PHONY: all
all: clean go.fmt go.lint build


#===================================  CMD  ====================================#
## build: Build and compile source code to program.
.PHONY: build
build:
	@echo -e "\033[1;36m====== make build ======\033[0m"
	@$(MAKE) go.build

## run: Build and compile source code, and then run the program.
.PHONY: run
run: build
	@echo -e "\033[1;36m====== start to run ======\033[0m"
	@$(MAKE) go.run

## clean: Remove all files that are created by building.
.PHONY: clean
clean:
# go clean -i $(TARGET_PACKAGE)
	@echo -e "\033[1;36m====== make clean ======\033[0m"
	@$(MAKE) go.clean

## checkout: Remove and clean unnecessary files before committing files.
.PHONY: checkout
checkout: clean
#	rm -f ./storage/log/*
#	rm -f ./storage/runtime/*
	@$(MAKE) go.fmt
	@$(MAKE) go.lint
	@echo -e "\033[1;36m====== all format and check have finish, you can add and commit now ======\033[0m"
	git status
	@echo -e "\033[1;36m--> 'git cz' is recommended for committing above files.\033[0m"

#===================================  help information  ====================================#
define USAGE_OPTIONS

Options:
  DEBUG        Whether to generate debug symbols. Default is 0.
               (example: make build DEBUG=1; make run DEBUG=1)
  RACE         Whether to enable data race detection. Default is 0.
               (example: make build RACE=1; make run RACE=1)
  LINT         Whether to enable golangci-lint. Default is 1.
               When golangci-lint disable, 'go vet' will be enabled.
               (example: make build LINT=0; make run LINT=0)
endef
export USAGE_OPTIONS

## help: Show this help information.
.PHONY: help
help: Makefile
	@echo -e "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"
