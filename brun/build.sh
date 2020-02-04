#!/usr/bin/env bash

#cpath=`pwd`
#PROJECT_PATH=${cpath%src*}
#echo $PROJECT_PATH
#export GOPATH=$GOPATH:${PROJECT_PATH}

SOURCE_FILE_NAME=main
TARGET_FILE_NAME=logagent

rm -fr ${TARGET_FILE_NAME}*

build(){
    echo $GOOS $GOARCH
    tname=${TARGET_FILE_NAME}_${GOOS}_${GOARCH}${EXT}
    env GOOS=$GOOS GOARCH=$GOARCH \
    go build -o ${tname} \
    -v ${SOURCE_FILE_NAME}.go
    chmod +x ${tname}
    mv ${tname} ${TARGET_FILE_NAME}${EXT}
    if [ ${GOOS} == "windows" ];then
        zip ${tname}.zip ${TARGET_FILE_NAME}${EXT} conf.ini
    else
        tar  --exclude=*.log --exclude=*.gz  --exclude=*.zip  --exclude=*.git -czvf ${tname}.tar.gz ${TARGET_FILE_NAME}${EXT} conf.ini *.sh  -C ./ .
    fi
    mv ${TARGET_FILE_NAME}${EXT} ${tname}

}
CGO_ENABLED=0
#mac os 64
GOOS=darwin
GOARCH=amd64
build

#linux 64
GOOS=linux
GOARCH=amd64
build

#windows
#64
GOOS=windows
GOARCH=amd64
build

GOARCH=386
build