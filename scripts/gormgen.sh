#!/bin/bash

shellExit()
{
if [ $1 -eq 1 ]; then
    printf "\nfailed!!!\n\n"
    exit 1
fi
}

printf "\nRegenerating file\n\n"
time go run -v ./cmd/database/main.go  -host $1 -port $2 -user $3 -pass $4 -name $5 -tables $6
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
