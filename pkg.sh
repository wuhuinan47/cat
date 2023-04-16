#!/bin/bash

GOOS=linux GOARCH=amd64 go build . 
date
scp useMiningItem.json ctrl.html whn.html maolaile.html build/linuxWechat.py root@192.168.10.21:/data/cat/
scp cat root@192.168.10.21:/data/cat/
