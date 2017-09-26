#!/bin/bash

set -e

# Add default command if needed.
if [ "${1:0:1}" = '-' ]; then
    set -- /usr/local/bin/pkg-distributor "$@"
fi

DATA_DIR=${DATA_DIR:-/data}

mkdir -p ${DATA_DIR}/aptly/public/conf
mkdir -p ${DATA_DIR}/public
mkdir -p ${DATA_DIR}/files

# Import GPG key pair.
if [ -z "$GPG_PUBLIC_KEY" ]; then
    echo "error: GPG_PUBLIC_KEY environment not set"
    exit
fi
gpg --import - <<< "$GPG_PUBLIC_KEY"

if [ -z "$GPG_PRIVATE_KEY" ]; then
    echo "error: GPG_PRIVATE_KEY environment not set"
    exit
fi
gpg --import --allow-secret-key-import - <<< "$GPG_PRIVATE_KEY"

GPG_KEY_ID=$(gpg --with-colons <<< "$GPG_PUBLIC_KEY" | head -n 1 | cut -d ':' -f 5)

if [ -z "$APT_CONF_ORIGIN" ]; then
    echo "error: APT_CONF_ORIGIN environment not set"
    exit
fi

if [ -z "$APT_CONF_LABEL" ]; then
    echo "error: APT_CONF_LABEL environment not set"
    exit
fi

APT_CONF_CODENAMES=${APT_CONF_CODENAMES:-trusty xenial}
APT_CONF_VERSION=${APT_CONF_VERSION:-1.0}
APT_CONF_COMPONENTS=${APT_CONF_COMPONENTS:-main}
APT_CONF_DESCRIPTION=${APT_CONF_DESCRIPTION:-"APT default description, please set by 'APT_CONF_DESCRIPTION' environment."}

# Generate distributions file.
echo -n > ${DATA_DIR}/aptly/public/conf/distributions
for codename in $APT_CONF_CODENAMES; do

    cat <<EOF >> ${DATA_DIR}/aptly/public/conf/distributions
Origin: $APT_CONF_ORIGIN
Label: $APT_CONF_LABEL
Codename: $codename
Version: $APT_CONF_VERSION
Architectures: source x86_64 amd64 i386 i686
Components: $APT_CONF_COMPONENTS
SignWith: $GPG_KEY_ID
Description: $APT_CONF_DESCRIPTION

EOF

done

gpg --armor --export $GPG_KEY_ID > ${DATA_DIR}/aptly/public/conf/gpg.key

cat <<EOF >> /etc/aptly.conf
{
    "rootDir": "${DATA_DIR}/aptly",
    "downloadConcurrency": 4,
	"downloadSpeedLimit": 0,
	"architectures": [],
	"dependencyFollowSuggests": false,
	"dependencyFollowRecommends": false,
	"dependencyFollowAllVariants": false,
	"dependencyFollowSource": false,
	"gpgDisableSign": false,
	"gpgDisableVerify": false,
	"downloadSourcePackages": false,
	"ppaDistributorID": "ubuntu",
	"ppaCodename": "",
	"S3PublishEndpoints": {},
	"SwiftPublishEndpoints": {}
}
EOF

test -d ${DATA_DIR} || mkdir ${DATA_DIR}/public
ln -fs ${DATA_DIR}/aptly/public ${DATA_DIR}/public/apt

# Use sha256 as default digest algorithm.
# See https://github.com/smira/aptly/pull/366.
echo 'digest-algo sha256' > ~/.gnupg/gpg.conf

exec "$@"
