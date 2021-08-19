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
        return
    else:
        url = 'https://cat.rosettawe.com/update?id='+userID+'&token='+token+'&name='+nickname
        response=session.get(url)
        data = response.content.decode('utf-8')
        print("update token result is ", data)
        return

class Cat():
    """
    登陆B站, 处理验证码
    电脑的缩放比例需要为100%, 否则验证码图片的获取会出现问题
    """

    def __init__(self):
        """
        初始化
        """

        # resp=session.get("https://open.weixin.qq.com/connect/qrconnect?appid=wx22f69b39568e9cb3&redirect_uri=http%3A%2F%2Flogin.11h5.com%2Faccount%2Fapi.php%3Fc%3Dwxlogin%26d%3DwxQrcodeAuth%26pf%3Dwxqrcode%26ssl%3D1%26back_url%3Dhttps%253A%252F%252Fplay.h5avu.com%252Fgame%252F%253Fgameid%253D147&response_type=code&scope=snsapi_login&state=#wechat_redirect")        
        # data = resp.content.decode('utf-8')

        # html = etree.HTML(data)s
        # html_data = html.xpath('/html/body/div[1]/div/div/div[2]/div[1]/img/@src')
        # for i in html_data:
        #     session.get("http://127.0.0.1:33333/test?qrcode="+"https://open.weixin.qq.com"+i)
        #     print(i)


        options = webdriver.ChromeOptions()
        # 设置为开发者模式，避免被识别
        options.add_experimental_option('excludeSwitches',
                                        ['enable-automation'])
        # options.add_argument('--no-sandbox')
        # options.add_argument('--disable-dev-shm-usage')
        # options.add_argument('--headless')
        # options.add_argument('blink-settings=imagesEnabled=false')
        # options.add_argument('--disable-gpu')
        self.browser = webdriver.Chrome(options=options)
        self.browser.get("https://open.weixin.qq.com/connect/qrconnect?appid=wx22f69b39568e9cb3&redirect_uri=http%3A%2F%2Flogin.11h5.com%2Faccount%2Fapi.php%3Fc%3Dwxlogin%26d%3DwxQrcodeAuth%26pf%3Dwxqrcode%26ssl%3D1%26back_url%3Dhttps%253A%252F%252Fplay.h5avu.com%252Fgame%252F%253Fgameid%253D147&response_type=code&scope=snsapi_login&state=#wechat_redirect")

        


        
        qrcode=self.browser.find_element_by_xpath("/html/body/div[1]/div/div/div[2]/div[1]/img").get_attribute('src')
        session.get("https://cat.rosettawe.com/sendQrcode?qrcode="+qrcode)
        print(qrcode)

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
                updateToken(userID, nickname, token)
                print("userID is ", userID)
                print("token is ", token)
                self.browser.quit()
                return
        
        

        self.browser.quit()
        print("超时未扫码，页面已关闭")

     

if __name__ == '__main__':
    Cat()


# def checkToken(self):
#         token=self.browser.execute_script("return localStorage.getItem('yg_token')")
#         if  token != "None":
#             print(token)
#             return
#             checkToken(self)