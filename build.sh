#!/usr/bin/env bash

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o PPGo_ApiAdmin .
docker build -t wolcen/xopadmin:1.0 .
#docker push wolcen/xopadmin:1.0