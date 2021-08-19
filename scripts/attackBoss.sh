#!/bin/bash




echo "bossID:$1\n"
echo "token:$2\n"
echo "url:$3\n"


function rand(){   
    min=$1   
    max=$(($2-$min+1))   
    num=$(($RANDOM+1000000000)) #增加一个10位的数再求余   
    echo $(($num%$max+$min))   
}   

function getNow(){
    current=`date "+%Y-%m-%d %H:%M:%S"`  
    timeStamp=`date -j -f "%Y-%m-%d %H:%M:%S" "${current}" "+%s"`
    now=$((timeStamp*1000+`openssl rand -base64 8 |cksum |cut -c 1-3`)) 
    echo $(($now))
}
     




for ((i=1; i<=3; i++))
do

sleep 3

now=$(getNow)
damage=$(rand 95 100)   
echo "first now:$now,damage is : $damage"
result=$(curl --connect-timeout 15 -m 20 -s "$3/game?cmd=attackBoss&token=$2&bossID=$1&damage=$damage&isPerfect=0&isDouble=0&now=$now" | jq '.boss.leftHp')

if [ "$result" = "null" ]
then
  result=$(curl --connect-timeout 15 -m 20 -s "$3/game?cmd=attackBoss&token=$2&bossID=$1&damage=$damage&isPerfect=0&isDouble=0&now=$now" | jq '.boss.leftHp')
  echo "attackBoss fail, attack again , leftHp is $result \n"
else  
  echo -e " normal attackBoss leftHp is $result \n"
fi

done


for ((i=1; i<=2; i++))
do
sleep 4
now=$(getNow)
damage=$(rand 195 200)   
echo "second now:$now,damage is : $damage"
result=$(curl --connect-timeout 15 -m 20 -s "$3/game?cmd=attackBoss&token=$2&bossID=$1&damage=$damage&isPerfect=0&isDouble=1&now=$now" | jq '.boss.leftHp')

if [ "$result" = "null" ]
then
  result=$(curl --connect-timeout 15 -m 20 -s "$3/game?cmd=attackBoss&token=$2&bossID=$1&damage=$damage&isPerfect=0&isDouble=1&now=$now" | jq '.boss.leftHp')
  echo "double attackBoss fail, attack again , leftHp is $result \n"
else  
  echo -e " double attackBoss leftHp is $result \n"
fi


done



