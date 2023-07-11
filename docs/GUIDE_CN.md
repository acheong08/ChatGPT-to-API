# 中文指导手册

本中文手册由 [@BlueSkyXN](https://github.com/BlueSkyXN) 编写


[中文文档（Chinese Docs）](../README_CN.md)  
 [英文文档（English Docs）](../README.md)

# 基本配置

有关docker的指导请阅读 [DOCKER中文手册](Docker_CN.md)

有关Token的指导请阅读 [TOKEN中文手册](TOKEN_CN.md)

## Docker-Compose配置

```
version: '3'

services:
  app:
    image: acheong08/chatgpt-to-api
    container_name: chatgpttoapi
    restart: unless-stopped
    ports:
      - '10080:10080'
    environment:
      SERVER_HOST: 0.0.0.0
      SERVER_PORT: 10080
      ADMIN_PASSWORD: TotallySecurePassword
      API_REVERSE_PROXY: https://ai.fakeopen.com/api/conversation
      PUID: user-X

```
- ports 左边是外部端口，右边是内部端口，内部端口要和下面环境变量的Server port一致。
- Server host/port：监听配置，默认0000监听某一端口。
- ADMIN_PASSWORD：管理员密码，HTTP请求时候需要验证。
- API_REVERSE_PROXY:接口的反向代理，具体介绍请看下文的后端代理介绍部分。
- PUID: user-X，请看[中文文档（Chinese Docs）](../README_CN.md) 的介绍

其他可以不需要设置，包括预设的AccessToken和代理表、HTTP/S5代理。


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

我是注释掉 ##- GO_CHATGPT_API_PROXY= 的环境变量、换个外部端口后用Docker-Compose启动即可。然后不需要对这个代理接口做其他操作，包括登录。

搭建好之后最好测试下基本调用能不能用，下面是一个示例，你需要根据实际情况修改。


```
curl http://127.0.0.1:8080/chatgpt/backend-api/conversation \
  -H "Content-Type: application/json" \
  -d '{
     "model": "gpt-3.5-turbo",
     "messages": [{"role": "user", "content": "Say this is a test!"}],
     "temperature": 0.7
   }'

```

如果得到缺少认证的提示比如 ```{"errorMessage":"Missing accessToken."}``` 就说明已经正常跑了

# 用例

## 基本提问
```
curl http://127.0.0.1:10080/v1/chat/completions \
  -d '{
     "model": "text-davinci-002-render-sha",
     "messages": [{"role": "user", "content": "你是什么模型，是GPT3.5吗"}]
   }'
```

参考回复如下

```
{"id":"chatcmpl-QXlha2FBbmROaXhpZUFyZUF3XXXXXX","object":"chat.completion","created":0,"model":"gpt-3.5-turbo-0301","usage":{"prompt_tokens":0,"completion_tokens":0,"total_tokens":0},"choices":[{"index":0,"message":{"role":"assistant","content":"是的，我是一个基于GPT-3.5架构的语言模型，被称为ChatGPT。我可以回答各种问题，提供信息和进行对话。尽管我会尽力提供准确和有用的回答，但请记住，我并不是完美的，有时候可能会出现错误或者误导性的答案。"},"finish_reason":null}]}
```

请注意无论什么模型提问都只会显示为模型是GPT3.5T-0301。你在网页版看不到消息记录（可能是删除了），Chat不支持并发提问，你需要token轮询。

## 提交Token
通过文件提交

```
curl -X PATCH \
     -H "Content-Type: application/json" \
     -H "Authorization: TotallySecurePassword" \
     -d "@/root/access_tokens.json" \
     http://127.0.0.1:10080/admin/tokens

```

直接提交

```
curl -X PATCH \
  -H "Content-Type: application/json" \
  -H "Authorization: TotallySecurePassword" \
  -d '["eyJhbXXX"]' \
  http://127.0.0.1:10080/admin/tokens
```

要清理Token直接停用删除Docker容器后重新构建运行容器即可
