#!/bin/bash

# cd /home/travis/gopath/bin

if [ "${GIMME_OS}" = "windows" ] ; then
	mv /home/travis/gopath/bin/windows_amd64/art /home/travis/gopath/bin/windows_amd64/art.exe
fi
