# ChatGPT-to-API
Create a fake API using ChatGPT's website

**API endpoint: http://127.0.0.1:8080/v1/chat/completions.**

## Help needed
- Documentation.

## Setup

### Authentication
Access token retrieval has been automated:
https://github.com/acheong08/ChatGPT-to-API/tree/master/tools/authenticator

Converting from a newline delimited list of access tokens to `access_tokens.json`
```bash
#!/bin/bash     

START="["
END="]"

TOKENS=""

while read -r line; do
  if [ -z "$TOKENS" ]; then
    TOKENS="\"$line\""
  else
    TOKENS+=",\"$line\""
  fi
done < access_tokens.txt

echo "$START$TOKENS$END" > access_tokens.json
```

### Cloudflare annoyances
`export PUID="user-..."`
or
`export API_REVERSE_PROXY="https://bypass.churchless.tech/api/conversation"`

## Docker build & Run

```bash
docker build -t chatgpt-to-api .

# Running the API
docker run --name chatgpttoapi -d -p 127.0.0.1:8080:8080 chatgpt-to-api

# API path
http://127.0.0.1:8080/v1/chat/completions

```

## Admin API docs
https://github.com/acheong08/ChatGPT-to-API/blob/master/docs/admin.md

## API usage docs
https://platform.openai.com/docs/api-reference/chat


## Docker compose

[Hub address](https://hub.docker.com/repository/docker/acheong08/chatgpt-to-api/general)

```yml
version: '3'

services:
  app:
    image: acheong08/chatgpt-to-api # Use latest tag
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
