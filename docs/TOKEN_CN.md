# 获取Token
---
# 参考Pandora项目的作者指导

https://github.com/pengzhile/pandora

获取Token的技术原理 https://zhile.io/2023/05/19/how-to-get-chatgpt-access-token-via-pkce.html

## 第三方接口获取Token
http://ai.fakeopen.com/auth 

你需要在这个新的网站的指导下安装浏览器插件，官方说明的有效期是14天。支持谷歌微软等第三方登录。（我谷歌注册的OpenAI就可以用这个）    

## 官网获取 Token
https://chat.openai.com/api/auth/session

打开后是个JSON，你需要先登录官方的ChatGPT网页版。里面有一个参数就是AccessToken。

# 参考go-chatgpt-api项目的作者指导
https://github.com/linweiyuan/go-chatgpt-api

ChatGPT 登录（返回 accessToken）（目前仅支持 ChatGPT 账号，谷歌或微软账号没有测试）

```POST /chatgpt/login```