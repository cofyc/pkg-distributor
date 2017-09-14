#!/bin/bash

ROOT=$(unset CDPATH && cd $(dirname "${BASH_SOURCE[0]}")/.. && pwd)
cd $ROOT

set -e

docker build -t pkg-distributor .
docker run --name pkg-distributor --rm -it -p 1973:1973  \
    -e GPG_PUBLIC_KEY="$(cat test/testdata/testonly.key)" \
    -e GPG_PRIVATE_KEY="$(cat test/testdata/testonly_private.key)" \
    -e APT_CONF_ORIGIN='Yecheng Fu' \
    -e APT_CONF_LABEL='yechengfu' \
    pkg-distributor --dir=/data/repo --basic-auth user:pass
