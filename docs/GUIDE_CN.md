# 后端代理
目前使用PUID+官网URL的方式不是很可靠，建议使用第三方程序或者网站绕过这个WAF限制。


## 公共代理
温馨提醒，由于OpenAI用的强力CloudFlareWAF，所以7层转发是无效的（不过4层在浏览器还是可以的）

目前根据几个大项目的介绍，我找到了这个介绍页 https://github.com/transitive-bullshit/chatgpt-api#reverse-proxy
最后得知主要是这两个

| Reverse Proxy URL                              | Author       | Rate Limits                    | Last Checked |
|-----------------------------------------------|--------------|--------------------------------|--------------|
| https://ai.fakeopen.com/api/conversation       | @pengzhile   | 5 req / 10 seconds by IP       | 4/18/2023    |
| https://api.pawan.krd/backend-api/conversation | @PawanOsman  | 50 req / 15 seconds (~3 r/s)   | 3/23/2023    |


## 自建方案

我经过测试，发现Pandora的API不行，原因可能是发起对话后的返回值会一次性返回一堆信息导致提取失败。不过我亲测GO-ChatGPT-API是可以的。

GO-ChatGPT-API项目 https://github.com/linweiyuan/go-chatgpt-api

我是注释掉 ##- GO_CHATGPT_API_PROXY= 的环境变量、换个外部端口后用Docker-Compase启动即可。然后不需要对这个代理接口做其他操作，包括登录。

搭建好之后最好测试下基本调用能不能用，下面是一个示例，你需要根据实际情况修改。


```
curl http://127.0.0.1:8080/chatgpt/conversation \
  -H "Content-Type: application/json" \
  -d '{
     "model": "gpt-3.5-turbo",
     "messages": [{"role": "user", "content": "Say this is a test!"}],
     "temperature": 0.7
   }'

```

如果得到缺少认证的提示比如 ```{"errorMessage":"Missing accessToken."}``` 就说明已经正常跑了
