

# other makefile calls this PHONY to verify tool
# check, dispatch and install related tool
tools.verify.%:
	@if ! which $* &>/dev/null; \
	then \
		@echo -e "\033[1;31m-->[Missing $* tool, try to install]\033[0m" \
		$(MAKE) install.$*; \
	fi

#===================================  tool install CMD  ====================================#
# Because golangci-lint release too fast to adapt the code, only use the specific version(now is v1.41.1).
# More install detail in: https://golangci-lint.run/usage/install/#local-installation
# Note: such go install/go get installation aren't guaranteed to work. We recommend using binary installation.
# binary will be $(go env GOPATH)/bin/golangci-lint
.PHONY: install.golangci-lint
install.golangci-lint:
	@echo -e "\033[1;36m-->[install-tool:golangci-lint, to:$(go env GOPATH)/bin/golangci-lint]\033[0m"
	@echo -e "\033[1;36m--> This step may take a long time, please be patient and wait."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.41.1
