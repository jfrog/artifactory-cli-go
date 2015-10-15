#!/bin/bash

cd /home/travis/gopath/bin

if [ "${GIMME_OS}" != "linux" ] || [ "${GIMME_ARCH}" != "amd64" ] ; then
	cd "${GIMME_OS}_${GIMME_ARCH}" || exit
fi

if [ "${GIMME_OS}" = "windows" ] ; then
	mv artifactory-cli-go.exe art.exe
else 
	mv artifactory-cli-go art
fi