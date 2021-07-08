#!/bin/sh

set -e

export TAG=4.06

docker build -t ahojukka5/dwk-mainapp-gen-timestamp:$TAG mainapp/gen-timestamp && docker push ahojukka5/dwk-mainapp-gen-timestamp:$TAG
docker build -t ahojukka5/dwk-mainapp-read-timestamp:$TAG mainapp/read-timestamp && docker push ahojukka5/dwk-mainapp-read-timestamp:$TAG
docker build -t ahojukka5/dwk-pingpong:$TAG pingpong && docker push ahojukka5/dwk-pingpong:$TAG
docker build -t ahojukka5/dwk-todo-backend:$TAG todo-backend && docker push ahojukka5/dwk-todo-backend:$TAG
docker build -t ahojukka5/dwk-todo-frontend:$TAG todo-frontend && docker push ahojukka5/dwk-todo-frontend:$TAG
docker build -t ahojukka5/dwk-todo-cronjob:$TAG todo-cronjob && docker push ahojukka5/dwk-todo-cronjob:$TAG
docker build -t ahojukka5/dwk-todo-broadcaster:$TAG todo-broadcaster && docker push ahojukka5/dwk-todo-broadcaster:$TAG
