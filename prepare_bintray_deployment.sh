#!/bin/bash

cd /home/travis/gopath/bin

echo a
ls -R

if [ "${GIMME_OS}" != "linux" ] || [ "${GIMME_ARCH}" != "amd64" ] ; then
	cd "${GIMME_OS}_${GIMME_ARCH}" || exit
fi

echo b
ls -R
