#!/bin/bash

GOOS=linux GOARCH=amd64 go build . 
date
scp ctrl.html maolaile.html build/linuxWechat.py cat root@110.238.110.71:/data/cat/
