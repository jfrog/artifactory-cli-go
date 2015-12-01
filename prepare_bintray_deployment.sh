#!/bin/bash

# cd /home/travis/gopath/bin

# echo $GOOS
# echo $GOARCH

if [ "${GIMME_OS}" = "windows" ] ; then
	mv art art.exe
fi
