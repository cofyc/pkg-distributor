# Package Distributor

[![Build Status](https://travis-ci.org/cofyc/pkg-distributor.svg?branch=master)](https://travis-ci.org/cofyc/pkg-distributor)
[![Docker Repository on Quay](https://quay.io/repository/cofyc/pkg-distributor/status "Docker Repository on Quay")](https://quay.io/repository/cofyc/pkg-distributor)

## Table of Contents

* [Supported Package management Systems](#supported-package-management-systems)
* [Directory](#directory)
* [References](#references)

## Supported Package Management Systems

* apt
* yum

## Directory

* /apt - Debian packages.
  * /conf
  * /dists
  * /pool
* /yum - RPM packages.
  * /repos/<release>-<arch>/repodata - e.g. /repos/el7-x86_64/repodata

## How to use

### Ubuntu/Debian

```
wget -q https://repo.example.com/apt/conf/gpg.key -O - | apt-key add -
echo 'deb https://repo.example.com/apt xenial main' > /etc/apt/sources.list.d/<repo>.list
```

### CentOS 7

```
[<repo>]
name=<repo>
baseurl=https://repo.example.com/yum/repos/el7-x86_64/
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://repo.example.com/yum/doc/gpg.key
```

## References

- https://wiki.debian.org/DebianRepository/SetupWithReprepro
- https://debian-administration.org/article/286/Setting_up_your_own_APT_repository_with_upload_support
- https://wiki.debian.org/DebianRepository/Setup?action=show&redirect=HowToSetupADebianRepository#aptly
- http://yum.baseurl.org/wiki/RepoCreate
