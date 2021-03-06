# Copyright 2016 The Rook Authors. All rights reserved.
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

# See https://git.fedorahosted.org/cgit/libaio.git/
set(Aio_VERSION 0.3.110)
set(Aio_URL http://ftp.us.debian.org/debian/pool/main/liba/libaio/libaio_${Aio_VERSION}.orig.tar.gz)

message(STATUS "External: Building libaio ${Aio_VERSION}")

#
# Build
#


ExternalProject_Add(aio
  PREFIX ${EXTERNAL_ROOT}

  URL ${Aio_URL}

  DOWNLOAD_DIR ${EXTERNAL_DOWNLOAD_DIR}
  BUILD_IN_SOURCE 1

  PATCH_COMMAND patch -p0 < ${CMAKE_CURRENT_SOURCE_DIR}/00-makefile.patch
  CONFIGURE_COMMAND true
  BUILD_COMMAND $(MAKE) CFLAGS=${EXTERNAL_C_FLAGS} ${EXTERNAL_TOOLCHAIN}
  INSTALL_COMMAND $(MAKE) prefix=<INSTALL_DIR> install

  LOG_DOWNLOAD ${EXTERNAL_LOGGING}
  LOG_PATCH ${EXTERNAL_LOGGING}
  LOG_CONFIGURE ${EXTERNAL_LOGGING}
  LOG_BUILD ${EXTERNAL_LOGGING}
  LOG_INSTALL ${EXTERNAL_LOGGING})
