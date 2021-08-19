#!/bin/bash


xiaohao1=301807377
xiaohao2=309433834
xiaohao3=326941142
xiaohao4=406378614
xiaohao5=690708340
xiaohao6=693419844


for bossID in 693419844_1629360416;
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


