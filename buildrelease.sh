#!/bin/bash
# Copyright 2022 Mike Bell, Albion Research Ltd.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This builds binaries for selected architectures
# For possible OS and ARCH names, run  
#     go tool dist list
# Each build description consists of OS / Architecture / (optional) file extension
# -ldflags="-s -w" strips out debugging/reflection info to reduce executable size
#
mkdir -p release
for build in windows/amd64/.exe windows/386/.exe linux/amd64/ ; do
  os=`echo $build  | cut -d '/' -f1`
  arch=`echo $build  | cut -d '/' -f2`
  ext=`echo $build  | cut -d '/' -f3`
  echo "Building for $os / $arch"
  target="release/syslogqd-$os-$arch$ext"
  GOOS=$os GOARCH=$arch go build -ldflags="-s -w" -o $target
done
ls -lh release/*

