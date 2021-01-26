SHELL = /bin/bash -o pipefail
# make -n
# make -np 2>&1 | less

go-deps-reset:
	@git checkout -- go.mod
	@go mod tidy

# -d flag ...download the source code needed to build ...
# -t flag ...consider modules needed to build tests ...
# -u flag ...use newer minor or patch releases when available 
go-deps-upgrade:
	@go get -d -t -u -v ./...
	@go mod tidy
go-mod-tidy:
	@go mod tidy
