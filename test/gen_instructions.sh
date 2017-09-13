#!/bin/bash

wget -q http://192.168.99.100:1973/apt/conf/gpg.key -O - | apt-key add -
echo 'deb http://192.168.99.100:1973/apt xenial main' > /etc/apt/sources.list.d/test.list
