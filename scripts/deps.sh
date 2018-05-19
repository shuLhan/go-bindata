#!/bin/sh

DEPS=$1

VENDOR_DIR=${PWD}/vendor
VENDOR_SRC=${VENDOR_DIR}/src
VENDOR_BIN=${VENDOR_DIR}/bin

export GOPATH=${VENDOR_DIR}

while read -r SRC DST VER; do
	REPO_DIR="${VENDOR_SRC}/${DST}"

	echo ""
	echo ">>> Fetch ${SRC} ${VER}"
	echo "    into ${REPO_DIR}"

	if [ -d ${REPO_DIR} ]; then
		CUR_VER=`git -C ${REPO_DIR} describe --tags --exact-match 2>/dev/null || git -C ${REPO_DIR} rev-parse --short HEAD`

		if [ "${CUR_VER}" == "${VER}" ]; then
			echo ">>> Already up to date..."
		else
			git -C ${REPO_DIR} checkout master
			git -C ${REPO_DIR} reset --hard HEAD
			git -C ${REPO_DIR} fetch --tags
			git -C ${REPO_DIR} pull -q
		fi
	else
		mkdir -p ${REPO_DIR}
		git -C ${REPO_DIR} clone ${SRC} .
	fi

	git -C ${REPO_DIR} checkout -q ${VER}

	echo ">>> Rebuild ..."
	go install -a ${DST}/...

	if [[ "${DST}" == *gometalinter ]]; then
		echo ">>> Install gometalinter tools ..."
		${VENDOR_BIN}/gometalinter --install;
	fi
done < ${DEPS};
