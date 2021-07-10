#!/bin/sh

set -e

export TAG=latest

docker build -t ahojukka5/dwk-controller-debug:$TAG 01-debug && docker push ahojukka5/dwk-controller-debug:$TAG
docker build -t ahojukka5/dwk-controller:$TAG 02-controller && docker push ahojukka5/dwk-controller:$TAG
