#!/bin/bash

cd /home/travis/gopath/bin

echo a
ls


if [ "${GIMME_OS}" != "linux" ] || ["${GIMME_ARCH}" != "amd64"] ; then
	cd "${GIMME_OS}_${GIMME_ARCH}"
fi


echo b
ls

if [ "${GIMME_OS}" = "windows" ] ; then
	mv art art.exe
fi

echo c
ls
