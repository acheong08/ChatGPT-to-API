# ChatGPT-to-API
创建一个模拟API（通过ChatGPT网页版）。使用AccessToken把ChatGPT模拟成OpenAI API，从而在各类应用程序中使用OpenAI的API且不需要为API额外付费，因为模拟成网页版的使用了，和官方API基本互相兼容。

本中文手册由 [@BlueSkyXN](https://github.com/BlueSkyXN) 编写

[英文文档（English Docs）](README.md)

## 认证和各项准备工作

在使用之前，你需要完成一系列准备工作

1. 准备ChatGPT账号，最好的PLUS订阅的，有没有开API不重要
2. 完善的运行环境和网络环境（否则你总是要寻找方法绕过）
3. Access Token和PUID，下面会教你怎么获取
4. 选择一个代理后端或者自行搭建
5. 你可以在 https://github.com/BlueSkyXN/OpenAI-Quick-DEV 项目找到一些常用组件以及一些快速运行的教程或程序。

### 获取PUID

`_puid` cookie.

### 获取Access Token
目前有多种方法和原理，这部分内容可以参考 [TOKEN中文手册](docs/TOKEN_CN.md)

## 安装和运行
  
作者在[英文版介绍](README.md) 通过GO编译来构建二进制程序，但是我猜测这可能需要一个GO编译环境。所以我建议基于作者的Compose配置文件来Docker运行。 

有关docker的指导请阅读 [DOCKER中文手册](docs/Docker_CN.md)

安装好Docker和Docker-Compase后，通过Compase来启动

```docker-compose up -d```

注意，启动之前你需要配置 yml 配置文件，主要是端口和环境变量，各项参数、用法请参考 [中文指导手册](docs/GUIDE_CN.md)

最后的API端点（Endpoint）是

```http://127.0.0.1:8080/v1/chat/completions```

注意域名/IP和端口要改成你自己的

### 环境变量
  - `PUID` - 用户ID
  - `http_proxy` - SOCKS5 或 HTTP 代理 `socks5://HOST:PORT`
  - `SERVER_HOST` - (default)比如 127.0.0.1
  - `SERVER_PORT` - (default)比如 8080 by

### 文件选项
  - `access_tokens.json` - 附带AccessToken的Json文件
  - `proxies.txt` - 代理表 (格式: `USERNAME:PASSWORD:HOST:PORT`)
  
