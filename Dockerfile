FROM golang:1.9.2

ADD . /go/src/github.com/cofyc/pkg-distributor

RUN set -eux \
    && cd /go/src/github.com/cofyc/pkg-distributor \
    && go install github.com/cofyc/pkg-distributor/cmd/pkg-distributor

FROM ubuntu:16.04

RUN set -eux \
    && sed -i '/security.ubuntu.com/d' /etc/apt/sources.list \
    && apt-get update \
    && apt-get install -y aptly createrepo expect \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=0 /go/bin/pkg-distributor /usr/local/bin
ADD entrypoint.sh /
ADD cmd/rpmautosign /usr/local/bin/rpmautosign

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/usr/local/bin/pkg-distributor"]
