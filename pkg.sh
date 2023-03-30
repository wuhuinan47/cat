#!/bin/bash

GOOS=linux GOARCH=amd64 go build . 
date
scp useMiningItem.json ctrl.html whn.html maolaile.html build/linuxWechat.py root@110.238.110.71:/data/cat/
scpo cat root@110.238.110.71:/data/cat/
