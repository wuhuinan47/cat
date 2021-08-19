#!/bin/bash

clear


function getNow(){
    current=`date "+%Y-%m-%d %H:%M:%S"`  
    timeStamp=`date -j -f "%Y-%m-%d %H:%M:%S" "${current}" "+%s"`
    now=$((timeStamp*1000+`openssl rand -base64 8 |cksum |cut -c 1-3`)) 
    echo $(($now))
}


url=$(curl --connect-timeout 15 -m 20 -s  "http://159.138.158.52:33333/getServerURL")
token=$(curl --connect-timeout 15 -m 20 -s "http://159.138.158.52:33333/getZoneToken?id=406378614&serverURL=$url")

now=$(getNow)
helpList=$(curl -s "$url/game?cmd=getGoldMineHelpList&token=$token&now=$now")

for i in `seq 1 25`;
do
quality=$(echo $helpList|jq ".helpList[$i].quality")
fuid=$(echo $helpList|jq ".helpList[$i].uid")



if [ $quality != 1 ]
then
echo "q:$quality, u:$fuid"

for ((i=1; i<=1; i++))
do
now=$(getNow)
result=$(curl -s "$url/game?cmd=enterGoldMine&token=$token&fuid=$fuid&type=0&now=$now" |jq '.goldMine.treasureList."21"')
echo "result is $result"
if [ "$result" = "null" ]
then
  echo "result is not set!"
else  
    now=$(getNow)
    curl "$url/game?cmd=goldMineFish&token=$token&fuid=$fuid&id=21&now=$now"
    echo -e " \n"
    sleep 3
  echo "dmin is set !"
fi
done


for ((i=1; i<=1; i++))
do
now=$(getNow)
result=$(curl -s "$url/game?cmd=enterGoldMine&token=$token&fuid=$fuid&type=0&now=$now" |jq '.goldMine.treasureList."22"')
echo "result is $result"
if [ "$result" = "null" ]
then
  echo "result is not set!"
else  
    now=$(getNow)
    curl "$url/game?cmd=goldMineFish&token=$token&fuid=$fuid&id=22&now=$now"
    echo -e " \n"
  echo "dmin is set !"
  sleep 3
fi
done

else
echo "jaha"
fi


done
