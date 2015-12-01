#!/bin/bash

cd /home/travis/gopath/bin

echo ${OS}
echo ${ARCH}

echo $GOOS
echo $GOARCH

if [ "${OS}" = "windows" ] ; then
	mv art art.exe
fi
