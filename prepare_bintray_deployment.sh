#!/bin/bash
echo a
ls

cd /home/travis/gopath/bin
echo "${GOOS}"
echo "${GOARCH}"
ls

if [ "${GOOS}" = "windows" ] ; then
	mv /home/travis/gopath/bin/art /home/travis/gopath/bin/art.exe
fi
echo c
ls
