#============================= detail for Golang ==============================#
# vet all file, include sub-package:   go vet ./...
# vet **only** this level package:     go vet ./
# vet the specified package:           go vet package-name-1 package-name-2 
# vet the default package:             go vet

#=================================== config ===================================#
# path to main source code
TARGET_PACKAGE = MarineSnow
BINARY_NAME = demoApp


#===================================  CMD  ====================================#
.PHONY: debug build run clean vet fmt help checkout

debug:
# go build -gcflags="-N -l" $(TARGET_PACKAGE)
	@echo -e "\033[1;36m====== make debug ======\033[0m"
	go build -v -gcflags="all=-N -l" -o ${BINARY_NAME} .

build:
# go build $(TARGET_PACKAGE)
	@echo -e "\033[1;36m====== make all ======\033[0m"
	go build -v -o ${BINARY_NAME} .

run: debug
	@echo -e "\033[1;36m====== start to run ======\033[0m"
	@./${BINARY_NAME} app start

clean:
# go clean -i $(TARGET_PACKAGE)
	@echo -e "\033[1;36m====== make clean ======\033[0m"
	go clean -i ./...
	@if [ -f ${BINARY_NAME} ] ; then rm ${BINARY_NAME} ; fi

vet:
	@echo -e "\033[1;36m====== make vet ======\033[0m"
	go vet ./framework/...
	go vet ./app/...
	go vet ./provider/...
	go vet main.go

fmt:
	@echo -e "\033[1;36m====== make fmt ======\033[0m"
	go fmt ./...

checkout: clean vet fmt
	rm -f ./storage/log/*
	rm -f ./storage/runtime/*
	@echo -e "\033[1;36m====== all format and check have finish, you can add and commit now ======\033[0m"
	git status

help:
	@echo 'Usage: make <OPTIONS>'
	@echo ''
	@echo 'Available OPTIONS are:'
	@echo '    help               Show this help screen'
	@echo '    clean              Remove binaries, artifacts and releases'
	@echo '    vet                Run go vet'
	@echo '    fmt                Run go fmt'
	@echo '    build              Build project for current platform'
	@echo '    debug              Build project for current platform in debug mode'
	@echo '    run                Build project for current platform in debug mode, and run it'
	@echo '    checkout           Do before checkout to git, perform clean, vet and fmt.'


# # memoty escapes analy
# memAnaly:
# 	go build -gcflags '-m -l' $(CODE_PATH)

# # you can do also like: go run -race main.go
# # Only when run(./main) can find some data race
# raceAnaly:
# 	go build -race $(CODE_PATH)

# # TODO: complate this
# # all DIR
# all_go_file = $(wildcard *.go) $(wildcard cmd/bookstore/*.go)
