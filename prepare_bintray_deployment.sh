#!/bin/bash

cd /home/travis/gopath/bin

if [ "${OS}" = "windows" ] ; then
	mv art art.exe
fi
