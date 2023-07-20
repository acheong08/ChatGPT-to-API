# ChatGPT-to-API
从ChatGPT网站模拟使用API

**模拟API地址: http://127.0.0.1:8080/v1/chat/completions.**

## 使用
    
### 设置

配置账户邮箱和密码，自动生成和更新Access tokens（使用[OpenAIAuth](https://github.com/acheong08/OpenAIAuth/)）

`accounts.txt` - 存放OpenAI账号邮箱和密码的文件

格式:
```
邮箱:密码
邮箱:密码
...
```

所有登录后的Access tokens会存放在`access_tokens.json`

每14天自动更新Access tokens

注意！ 请使用未封锁的ip登录账号，请先打开浏览器登录`https://chat.openai.com/`以检查ip是否可用

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
git clone https://github.com/xqdoo00o/ChatGPT-to-API
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
  - `access_tokens.json` - 一个存放Access tokens JSON数组的文件 (可使用 PATCH请求更新Access tokens [correct endpoint](https://github.com/acheong08/ChatGPT-to-API/blob/master/docs/admin.md))
    ```
    ["access_token1", "access_token2"...]
    ```

## 用户管理文档
https://github.com/acheong08/ChatGPT-to-API/blob/master/docs/admin.md

## API使用说明
https://platform.openai.com/docs/api-reference/chat
