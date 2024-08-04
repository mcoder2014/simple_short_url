#!/bin/bash
RUN_NAME=simple_short_url
mkdir -p output/bin
mkdir -p output/conf
cp conf/* output/conf/

cp script/* output 2>/dev/null
chmod +x output/bootstrap.sh
go build -o output/bin/${RUN_NAME}