
#===================================  include  ====================================#
include script/make-rule/common.mk
include script/make-rule/golang.mk


#===================================  target  ====================================#
# Build all by default, even if it's not first
.DEFAULT_GOAL := all

.PHONY: all
all: clean go.fmt go.vet build





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

## checkout: Remove and clean unnecessary files before promoting code.
.PHONY: checkout
checkout: clean
	rm -f ./storage/log/*
	rm -f ./storage/runtime/*
	@$(MAKE) go.fmt
	@$(MAKE) go.vet
	@echo -e "\033[1;36m====== all format and check have finish, you can add and commit now ======\033[0m"
	git status

#===================================  help information  ====================================#
define USAGE_OPTIONS

Options:
  DEBUG        Whether to generate debug symbols. Default is 0.
               (example: make build DEBUG=1; make run DEBUG=1)
endef
export USAGE_OPTIONS

## help: Show this help information.
.PHONY: help
help: Makefile
	@echo -e "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"
