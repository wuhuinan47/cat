#!/bin/bash


if [[ $(screen -ls | grep "cat" | wc -l ) -gt 0 ]];then
       echo "CAT IS RUNNING"
       exit 1
fi


rm -rf cat_run
rm -rf screenlog.0
cp cat cat_run
start_file=./cat_run
screen -dmSL cat -s ${start_file}
    screen -r cat
