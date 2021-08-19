#!/bin/bash


xiaohao1=694981971
xiaohao2=374289806
xiaohao3=439943689
xiaohao4=382292124
xiaohao5=385498006
xiaohao6=381909995


for bossID in 309392050_1629361911 381909995_1629361447 385498006_1629361435 382292124_1629361389 439943689_1629361374 374289806_1629361358;
do

clear

echo "1 start.\n"


serverURL=$(curl --connect-timeout 15 -m 20 -s "http://159.138.158.52:33333/getServerURL")
zoneToken=$(curl --connect-timeout 15 -m 20 -s "http://159.138.158.52:33333/getZoneToken?id=$xiaohao1&serverURL=$serverURL")
echo "serverURLserverURL is :$serverURL"
echo "zoneTokenzoneToken is :$zoneToken"
./attackBoss.sh $bossID $zoneToken $serverURL
echo "1 end.\n"

clear

echo "2 start.\n"

serverURL=$(curl --connect-timeout 15 -m 20 -s "http://159.138.158.52:33333/getServerURL")
zoneToken=$(curl --connect-timeout 15 -m 20 -s "http://159.138.158.52:33333/getZoneToken?id=$xiaohao2&serverURL=$serverURL")
./attackBoss.sh $bossID $zoneToken $serverURL
echo "2 end.\n"

clear

echo "3 start.\n"

serverURL=$(curl --connect-timeout 15 -m 20 -s "http://159.138.158.52:33333/getServerURL")
zoneToken=$(curl --connect-timeout 15 -m 20 -s "http://159.138.158.52:33333/getZoneToken?id=$xiaohao3&serverURL=$serverURL")
./attackBoss.sh $bossID $zoneToken $serverURL

echo "3 end.\n"

clear

echo "4 start.\n"

serverURL=$(curl --connect-timeout 15 -m 20 -s "http://159.138.158.52:33333/getServerURL")
zoneToken=$(curl --connect-timeout 15 -m 20 -s "http://159.138.158.52:33333/getZoneToken?id=$xiaohao4&serverURL=$serverURL")
./attackBoss.sh $bossID $zoneToken $serverURL

echo "4 end.\n"

clear

echo "5 start.\n"

serverURL=$(curl --connect-timeout 15 -m 20 -s "http://159.138.158.52:33333/getServerURL")
zoneToken=$(curl --connect-timeout 15 -m 20 -s "http://159.138.158.52:33333/getZoneToken?id=$xiaohao5&serverURL=$serverURL")
./attackBoss.sh $bossID $zoneToken $serverURL

echo "5 end.\n"

clear

echo "6 start.\n"

serverURL=$(curl --connect-timeout 15 -m 20 -s "http://159.138.158.52:33333/getServerURL")
zoneToken=$(curl --connect-timeout 15 -m 20 -s "http://159.138.158.52:33333/getZoneToken?id=$xiaohao6&serverURL=$serverURL")
./attackBoss.sh $bossID $zoneToken $serverURL

echo "6 end.\n"


done


