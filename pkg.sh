#!/bin/bash

GOOS=linux GOARCH=amd64 go build . 
date
# scp 1369-36.json useMiningItem.json ctrl.html whn.html maolaile.html userInfo.html familyEnergy.html build/linuxWechat.py root@192.168.10.21:/data/cat/
# scp cat root@192.168.10.21:/data/cat/
# scp favicon.ico root@192.168.10.21:/data/cat/
# scp 1369-46.json root@192.168.10.21:/data/cat/
scp ctrl.html whn.html maolaile.html userInfo.html familyEnergy.html familyDayTask.html cat root@192.168.10.21:/data/cat/
# bs-scp 'master0->twbi->tcp://127.0.0.1:22' ctrl.html whn.html maolaile.html userInfo.html familyEnergy.html familyDayTask.html cat root@127.0.0.1:/data/cat/