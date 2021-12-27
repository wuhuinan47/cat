#!/bin/bash

GOOS=linux GOARCH=amd64 go build . 
date
scp ctrl.html maolaile.html build/linuxWechat.py cat root@159.138.158.52:/mnt/app/cat/
