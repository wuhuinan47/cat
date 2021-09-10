from selenium import webdriver
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By
from selenium.webdriver.common.action_chains import ActionChains
from selenium.common.exceptions import TimeoutException
from PIL import Image
from io import BytesIO
from time import sleep
from getpass import getpass
import random
import requests
import json
import time

from scrapy import Selector

"""
info:
author:CriseLYJ
github:https://github.com/CriseLYJ/
update_time:2019-3-7
"""

session = requests.session()
headers = {
    'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36'
}

def getGameInfo():
    url = 'https://api.11h5.com/conf?cmd=getGameInfo&gameid=147&'+str(time.time())+'274'
    response = session.get(url)
    data = response.content.decode('utf-8')
    dic = json.loads(data)
    serverURL = dic['ext']['serverURL']
    return serverURL


def getZoneToken(serverURL, token):
    url = serverURL+'/zone?cmd=enter&token='+token+'&yyb=0&inviteId=null&share_from=null&cp_shareId=null&now='+str(time.time())+'435'
    response = session.get(url)
    data = response.content.decode('utf-8')
    dic = json.loads(data)
    zoneToken = dic['zoneToken']
    nickname = dic['nickname']
    return nickname, zoneToken

def getGoldMineHelpList(serverURL, zoneToken):
    url = serverURL+'/game?cmd=getGoldMineHelpList&token='+zoneToken+'&now='+str(time.time())+'425'
    response = session.get(url)
    data = response.content.decode('utf-8')
    dic = json.loads(data)
    helpList = dic['helpList']
    return helpList


def updateToken(userID, nickname, token):
    if userID == None:
        print("empty userID ", userID)
        return 0
    else:
        url = 'https://mcps.51yizhuan.com:13010/update?id='+userID+'&token='+token+'&name='+nickname
        response=session.get(url)
        data = response.content.decode('utf-8')
        print("update token result is ", data)
        return 1

class Cat():
    """
    登陆B站, 处理验证码
    电脑的缩放比例需要为100%, 否则验证码图片的获取会出现问题
    """

    def __init__(self):
        """
        初始化
        """
        options = webdriver.ChromeOptions()
        # 设置为开发者模式，避免被识别
        options.add_experimental_option('excludeSwitches',
                                        ['enable-automation'])
        self.browser = webdriver.Chrome(options=options)
        self.browser.get("https://graph.qq.com/oauth2.0/show?which=Login&display=pc&response_type=token&client_id=101206450&state=&redirect_uri=http%3A%2F%2Flogin.vutimes.com%2Faccount%2Fpage%2FqqAuthCallback.html%3FswitchVersion%3D1%26pf%3Dqq%26ssl%3D1%26back_url%3Dhttps%253A%252F%252Fplay.h5avu.com%252Fgame%252F%253Fgameid%253D147")
        



        iframe=self.browser.find_element_by_css_selector("#ptlogin_iframe")
        self.browser.switch_to_frame(iframe)
        img=self.browser.find_element_by_id("qrlogin_img")


        location=img.location
        size=img.size

        print(location)
        print(size)

        location['x']=400
        location['y']=300
        # top,bottom,left,right=location['y'],location['y']+size['height'],location['x'],location['x']+size['width']
        top,bottom,left,right=400,532,300,432


        scrennshot=self.browser.get_screenshot_as_png()
        scrennshot=Image.open(BytesIO(scrennshot))
        scrennshot=scrennshot.crop((left,top,right,bottom))
        file_name='wechatQrcode.png'
        scrennshot.save(file_name)




        # qrcode=self.browser.find_element_by_id("qrlogin_img").get_attribute('src')


        # qrcode=self.browser.find_element_by_xpath("/html/body/div[2]/div[1]/div")

        # session.get("https://mcps.51yizhuan.com:13010/sendQrcode?qrcode="+qrcode)
       
        # print(qrcode)

        for x in range(100):
            token=self.browser.execute_script("return localStorage.getItem('yg_token')")
            if token == None:
                sleep(2)
            else:
                userID=self.browser.execute_script("return localStorage.getItem('__TD_userID')")
                # 获取serverURL
                serverURL = getGameInfo()
                print("serverURL is ", serverURL)

                # 获取nickname,zoneToken
                nickname,zoneToken=getZoneToken(serverURL, token)
                print("zoneToken is ", zoneToken)

                # 更新token到数据库
                result=updateToken(userID, nickname, token)
                if result==1:
                    print("userID is ", userID)
                    print("token is ", token)
                    return
                
                sleep(2)

        
        print("超时未扫码，页面已关闭")

     

if __name__ == '__main__':
    Cat()


# def checkToken(self):
#         token=self.browser.execute_script("return localStorage.getItem('yg_token')")
#         if  token != "None":
#             print(token)
#             return
#             checkToken(self)