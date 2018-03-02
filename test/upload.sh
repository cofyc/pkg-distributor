#!/bin/bash

ROOT=$(unset CDPATH && cd $(dirname "${BASH_SOURCE[0]}")/.. && pwd)
cd $ROOT

# debs
curl -u user:pass -F file=@test/testdata/wget_1.19.1-3ubuntu1_amd64.deb http://localhost:1973/v1/upload
curl -u user:pass -F file=@test/testdata/wget_1.16.1-1ubuntu1_amd64.deb http://localhost:1973/v1/upload

# rpms
curl -u user:pass -F file=@test/testdata/wget-1.14-15.el7.x86_64.rpm http://localhost:1973/v1/upload
