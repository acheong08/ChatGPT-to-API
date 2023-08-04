# ChatGPT-to-API
从ChatGPT网站模拟使用API

**模拟API地址: http://127.0.0.1:8080/v1/chat/completions.**

## 使用
    
### 设置

配置账户邮箱和密码，自动生成和更新Access tokens 和 PUID（仅PLUS账户）（使用[OpenAIAuth](https://github.com/acheong08/OpenAIAuth/)）

`accounts.txt` - 存放OpenAI账号邮箱和密码的文件

格式:
```
邮箱:密码
邮箱:密码
...
```

所有登录后的Access tokens和PUID会存放在`access_tokens.json`

每7天自动更新Access tokens和PUID

注意！ 请使用未封锁的ip登录账号，请先打开浏览器登录`https://chat.openai.com/`以检查ip是否可用

### GPT-4 设置（可选）

如果配置PLUS账户并使用GPT-4模型，则需要HAR文件（`chat.openai.com.har`）以完成captcha验证

1. 使用基于chromium的浏览器（Chrome，Edge）或Safari浏览器 登录`https://chat.openai.com/`，然后打开浏览器开发者工具（F12），并切换到网络标签页。

2. 新建聊天并选择GPT-4模型，随意问一个问题，点击网络标签页下的导出HAR按钮，导出文件`chat.openai.com.har`

### API 密钥（可选）

如OpenAI的官方API一样，可给模拟的API添加API密钥认证

`api_keys.txt` - 存放API密钥的文件

格式:
```
sk-123456
88888888
...
```

## 开始
```  
git clone https://github.com/acheong08/ChatGPT-to-API
cd ChatGPT-to-API
go build
./freechatgpt
```

### 环境变量
  - `PUID` - Plus账户可在`chat.openai.com`的cookies里找到，用于绕过cf的频率限制
  - `SERVER_HOST` - 默认127.0.0.1
  - `SERVER_PORT` - 默认8080
  - `ENABLE_HISTORY` - 默认true，允许网页端历史记录

### 可选文件配置
  - `proxies.txt` - 存放代理地址的文件

    ```
    http://127.0.0.1:8888
    socks5://127.0.0.1:9999
    ...
    ```
  - `access_tokens.json` - 一个存放Access tokens 和PUID JSON数组的文件 (可使用 PATCH请求更新Access tokens [correct endpoint](https://github.com/acheong08/ChatGPT-to-API/blob/master/docs/admin.md))
    ```
    [{token:"access_token1", puid:"puid1"}, {token:"access_token2", puid:"puid2"}...]
    ```

## 用户管理文档
https://github.com/acheong08/ChatGPT-to-API/blob/master/docs/admin.md

## API使用说明
https://platform.openai.com/docs/api-reference/chat
