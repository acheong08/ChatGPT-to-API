# ChatGPT-to-API
Create a fake API using ChatGPT's website

## A fork. For personal use. Check upstream for details.

**API endpoint: http://127.0.0.1:8080/v1/chat/completions.**

**When calling the API, you must include the authorization parameter in the request header: `'Authorization':'Bearer ' + accessToken`.**

**You can get your accessToken from the following link: [ChatGPT](https://chat.openai.com/api/auth/session)**

**This API can be used with the project [BetterChatGPT](https://github.com/huangzt/BetterChatGPT)**

## Docker build & Run

```bash
docker build -t chatgpt-to-api .

# 后台运行
docker run --name chatgpttoapi -d -p 127.0.0.1:8080:8080 chatgpt-to-api

# API地址
http://127.0.0.1:8080/v1/chat/completions

```

## Docker compose

[Hub 地址](https://hub.docker.com/repository/docker/huangzhenting/chatgpt-to-api/general)

```yml
version: '3'

services:
  app:
    image: huangzhenting/chatgpt-to-api # 总是使用latest,更新时重新pull该tag镜像即可
    container_name: chatgpttoapi
    restart: unless-stopped
    ports:
      - '8080:8080'
    environment:
      SERVER_HOST: 0.0.0.0
      SERVER_PORT: 8080
      ADMIN_PASSWORD: TotallySecurePassword
      # Reverse Proxy - Available on accessToken
      API_REVERSE_PROXY: https://bypass.churchless.tech/api/conversation
      # If the parameter API_REVERSE_PROXY is empty, the default request URL is https://chat.openai.com/backend-api/conversation, and the PUID is required.
      # You can get your PUID for Plus account from the following link: https://chat.openai.com/api/auth/session.
      PUID: xxx

```