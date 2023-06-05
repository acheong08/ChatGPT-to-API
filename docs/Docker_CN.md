# 使用阿里源实现Docker安装

移除旧的

```yum remove -y docker docker-common docker-selinux docker-engine```

安装依赖

```yum install -y yum-utils device-mapper-persistent-data lvm2```

配置Docker安装源（阿里）

```yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo```

检查可用的Docker-CE版本

```yum list docker-ce --showduplicates | sort -r```

安装Docker-CE

```yum -y install docker-ce```

运行Docker（默认不运行）

```systemctl start docker```

配置开机启动Docker

```systemctl enable docker```

# 使用官方二进制包安装Docker-Compase

下载 Docker-Compose 的二进制文件
```sudo curl -L "https://github.com/docker/compose/releases/download/v2.18.1/docker-compose-linux-x86_64" -o /usr/local/bin/docker-compose```

添加可执行权限

```sudo chmod +x /usr/local/bin/docker-compose```

验证 Docker-Compose 是否安装成功

```docker-compose --version```

启动容器

```docker-compose up -d```

关闭容器

```docker-compose down```

查看容器（如果启动了这里没有说明启动失败）

```docker ps```

# ChatGPT-TO-API的Docker-Compase文件

```
    ports:
      - '31480:31480'
    environment:
      SERVER_HOST: 0.0.0.0
      SERVER_PORT: 31480
      ADMIN_PASSWORD: TotallySecurePassword
      # Reverse Proxy - Available on accessToken
      #API_REVERSE_PROXY: https://bypass.churchless.tech/api/conversation
      #API_REVERSE_PROXY: https://ai.fakeopen.com/api/conversation
      PUID: user-7J4tdvHySlcilVgjFIrAtK1k

```

- 这里的ports，左边是外部端口，用于外部访问。右边的Docker端口，需要匹配下面程序设置的监听Port。
- 如果参数`API_REVERSE_PROXY`为空，则默认的请求URL为`https://chat.openai.com/backend-api/conversation`，并且需要提供PUID。PUID的获取参考 [README_CN.md](../README_CN.md)
- 这个密码需要自定义，我们构建请求的时候需要它来鉴权。默认是```TotallySecurePassword```

