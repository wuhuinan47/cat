# cat game 

# install google-chrome
yum -y install google-chrome-stable --nogpgcheck

# install chromedriver
downloads url is http://npm.taobao.org/mirrors/chromedriver/92.0.4515.107/chromedriver_linux64.zip
unzip ./build/chromedriver_linux64.zip
mv ./build/chromedriver /usr/bin/


# python env

pip3 install selenium
pip3 install Pillow