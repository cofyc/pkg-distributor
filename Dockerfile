FROM ubuntu:16.04

ENV GOLANG_VERSION 1.9

RUN set -eux \
    && sed -i '/security.ubuntu.com/d' /etc/apt/sources.list \
    && apt-get update \
    && apt-get install -y reprepro wget \
    && \
		dpkgArch="$(dpkg --print-architecture)"; \
		case "${dpkgArch##*-}" in \
			amd64) goRelArch='linux-amd64'; goRelSha256='d70eadefce8e160638a9a6db97f7192d8463069ab33138893ad3bf31b0650a79' ;; \
			armhf) goRelArch='linux-armv6l'; goRelSha256='f52ca5933f7a8de2daf7a3172b0406353622c6a39e67dd08bbbeb84c6496f487' ;; \
			arm64) goRelArch='linux-arm64'; goRelSha256='0958dcf454f7f26d7acc1a4ddc34220d499df845bc2051c14ff8efdf1e3c29a6' ;; \
			i386) goRelArch='linux-386'; goRelSha256='7cccff99dacf59162cd67f5b11070d667691397fd421b0a9ad287da019debc4f' ;; \
			ppc64el) goRelArch='linux-ppc64le'; goRelSha256='10b66dae326b32a56d4c295747df564616ec46ed0079553e88e39d4f1b2ae985' ;; \
			s390x) goRelArch='linux-s390x'; goRelSha256='e06231e4918528e2eba1d3cff9bc4310b777971e5d8985f9772c6018694a3af8' ;; \
			*) echo >&2 "warning: current architecture ($dpkgArch) does not have a corresponding Go binary release, exit"; exit -1 ;; \
	   	esac; \
      	wget -O go.tgz "https://golang.org/dl/go${GOLANG_VERSION}.${goRelArch}.tar.gz"; \
		echo "${goRelSha256} *go.tgz" | sha256sum -c -; \
		tar -C /usr/local -xzf go.tgz; \
		rm go.tgz \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

ADD . /go/src/github.com/cofyc/pkg-distributor

RUN set -e \
	&& export PATH="/usr/local/go/bin:$PATH" \
	&& export GOPATH="/go" \
    && go version \
	&& go build -o /usr/local/bin/pkg-distributor github.com/cofyc/pkg-distributor/cmd/pkg-distributor \
	&& rm -rf /go && rm -rf /usr/local/go

ADD entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/usr/local/bin/pkg-distributor", "--dir=/data/repo"]
