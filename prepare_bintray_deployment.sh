#!/bin/bash

cd /home/travis/gopath/bin

if [ "${GIMME_OS}" = "windows" ] ; then
	mv art art.exe
fi
