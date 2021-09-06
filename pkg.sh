#!/bin/bash

GOOS=linux GOARCH=amd64 go build . 

scp ctrl.html cat root@159.138.158.52:/mnt/app/cat/
