#!/bin/bash

GOOS=linux GOARCH=amd64 go build . 
date
scp ctrl.html maolaile.html build/linuxWechat.py cat root@47.115.55.167:/mnt/app/cat/
