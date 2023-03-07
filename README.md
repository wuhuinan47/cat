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