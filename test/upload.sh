#!/bin/bash

ROOT=$(unset CDPATH && cd $(dirname "${BASH_SOURCE[0]}")/.. && pwd)
cd $ROOT

curl -u user:pass -F file=@test/testdata/wget_1.19.1-3ubuntu1_amd64.deb http://192.168.99.100:1973/v1/upload
