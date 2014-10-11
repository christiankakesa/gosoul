#!/bin/sh

PROGNAME=gosoul
PLATFORMS="darwin/386 darwin/amd64 freebsd/386 freebsd/amd64 freebsd/arm linux/386 linux/amd64 linux/arm windows/386 windows/amd64"
VERSION=$(git tag -l | sort | tail -n1)

for PLATFORM in $PLATFORMS; do
    OS=${PLATFORM%/*}
    ARCH=${PLATFORM#*/}
    GOOS=$OS CGO_ENABLED=0 GOARCH=$ARCH go build -o $PROGNAME
    ARCHIVE=$PROGNAME-$VERSION-$OS-$ARCH.tar.gz
    tar -czf $ARCHIVE $PROGNAME
    echo $ARCHIVE
    sleep 1
done
