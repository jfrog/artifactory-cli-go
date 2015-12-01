#!/bin/bash
echo a
ls

# cd /home/travis/gopath/bin
echo b
ls

if [ "${GIMME_OS}" = "windows" ] ; then
	mv /home/travis/gopath/bin/art /home/travis/gopath/bin/art.exe
fi
echo c
ls
