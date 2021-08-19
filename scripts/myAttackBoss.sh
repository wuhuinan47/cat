#!/bin/bash

clear


function getNow(){
    current=`date "+%Y-%m-%d %H:%M:%S"`  
    timeStamp=`date -j -f "%Y-%m-%d %H:%M:%S" "${current}" "+%s"`
    now=$((timeStamp*1000+`openssl rand -base64 8 |cksum |cut -c 1-3`)) 
    echo $(($now))
}

function rand(){   
    min=$1   
    max=$(($2-$min+1))   
    num=$(($RANDOM+1000000000)) #增加一个10位的数再求余   
    echo $(($num%$max+$min))   
} 





for bossID in 382292124_1629361389 439943689_1629361374 374289806_1629361358;
do


now=$(getNow)

url=$(curl --connect-timeout 15 -m 20 -s  "http://159.138.158.52:33333/getServerURL")
token=$(curl --connect-timeout 15 -m 20 -s "http://159.138.158.52:33333/getZoneToken?id=302691822&serverURL=$url")


echo "url is :$url"
echo "token is :$token"
for ((i=1; i<=3; i++))
do

damage=$(rand 195 200)   
curl "$url//game?cmd=attackBoss&token=$token&bossID=$bossID&damage=$damage&isPerfect=0&isDouble=1&now=$now"
echo -e " \n"

sleep 3
done


for ((i=1; i<=1; i++))
do
damage=$(rand 390 400)   

curl "$url//game?cmd=attackBoss&token=$token&bossID=$bossID&damage=$damage&isPerfect=1&isDouble=1&now=$now"
echo -e " \n"

sleep 3
done


done



