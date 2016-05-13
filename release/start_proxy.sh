#/bin/bash

export GIN_MODE=release
./proxy > /var/log/proxy.log
