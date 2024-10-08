#!/bin/bash

shellExit()
{
if [ $1 -eq 1 ]; then
    printf "\nfailed!!!\n\n"
    exit 1
fi
}

printf "\nRegenerating file\n\n"
time go run -v ./cmd/gormmd/main.go -driver ${1:-mysql} -host $2 -port $3 -user $4 -pass $5 -db $6 -tables ${7:-"*"}
shellExit $?

printf "\ncreate curd code : \n"
time go build -o gormgen ./cmd/gormgen/main.go
shellExit $?

if [ ! -d $GOPATH/bin ];then
   mkdir -p $GOPATH/bin
fi

mv gormgen $GOPATH/bin
shellExit $?

go generate ./...
shellExit $?

printf "\nFormatting code\n\n"
time go run -v ./cmd/mfmt/main.go
shellExit $?

printf "\nDone.\n\n"
