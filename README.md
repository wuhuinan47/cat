# cat game 

# install google-chrome
yum -y install google-chrome-stable --nogpgcheck

# install chromedriver
downloads url is https://chromedriver.storage.googleapis.com/109.0.5414.25/chromedriver_linux64.zip
unzip ./build/chromedriver_linux64.zip
mv ./build/chromedriver /usr/bin/


# python env

pip3 install selenium
pip3 install Pillow



# 8=水雷 7=鱼雷 6=鱼叉 5=20分 4=10分 3=5分 2=1次碎石 1=2次碎石 0=可挖但是啥也没 -1=不能挖


现在有一个挖宝的游戏，路径是一个7乘8的矩形，横向是7个格子 ，纵向有8格子，每个格子是否可以挖是看格子的状态，-1表示不能挖，0或1表示可以挖但是没有收益，2表示需要挖两次才会变成-1，3表示挖了可以得5分并且变为-1,4表示挖了可以得10分并且变为-1,5表示挖了可以得20分并且变为-1,6表示挖了可以得到一个鱼叉，鱼叉用6表示，7表示挖了可以得到一个鱼雷，鱼雷用7表示，8表示挖了可以获得一个水雷，水雷用8表示。鱼叉可以挖一个格子，鱼雷可以横向炸掉除了2以外，如果有2需要炸2次才行，水雷表示可以炸范围是3乘3的区域。在进行挖掘时，最后一行被挖时会自动生成新的一行，并且第一行会被顶掉，永远都是一个7乘8的区域。写一个算法，在输入一行格子后，判断在使用鱼叉还是鱼雷还是水雷时，那个收益最高，并且标明在哪个位置使用。用go写一个算法