# Copyright 2016 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


GP := .gopath
PARENT := github.com/ktsakalozos
PKG := lxdlet
PKGPATH := ${PWD}/${GP}/src/${PARENT}/${PKG}
export GOPATH=${PWD}/${GP}

all: build

build: path-setup
	cd "${PKGPATH}" && \
	go build -o bin/lxdlet ./lxdlet/cmd/server/main.go

path-setup:
	@if [ ! -d "${GP}" ]; then mkdir -p "${GP}/src/${PARENT}" "${GP}/pkg" "${GP}/bin"; fi && \
	if [ ! -e "${PKGPATH}" ]; then ln -s "${PWD}" "${PKGPATH}"; fi && \
	echo "Local GOPATH set up at ${GOPATH}"

test: path-setup
	cd "${PKGPATH}" && \
	go test ./lxdlet/...

integration:
	nohup ./bin/lxdlet &
	sleep 3
	(cd hack; ./integration-tests.sh)
	killall lxdlet

clean:
	rm -rf ./bin/lxdlet ./${GP} ./nohup.out

.PHONY: all test clean
