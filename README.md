# ChatGPT-to-API
Create a fake API using ChatGPT's website

**API endpoint: http://127.0.0.1:8080/v1/chat/completions.**

[中文文档（Chinese Docs）](README_CN.md)

## Help needed
- Documentation.

## Setup

<details>
  <summary>
    
### Authentication
  </summary>
  
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

</details>

## Getting set up
  
`git clone https://github.com/acheong08/ChatGPT-to-API`
`cd ChatGPT-to-API`
`go build`
`./freechatgpt`

### Environment variables
  - `PUID` - A cookie found on chat.openai.com for Plus users. This gets around Cloudflare rate limits
  - `http_proxy` - SOCKS5 or HTTP proxy. `socks5://HOST:PORT`
  - `SERVER_HOST` - Set to 127.0.0.1 by default
  - `SERVER_PORT` - Set to 8080 by default

### Files (Optional)
  - `access_tokens.json` - A JSON array of access tokens for cycling (Alternatively, send a PATCH request to the [correct endpoint](https://github.com/acheong08/ChatGPT-to-API/blob/master/docs/admin.md))
  - `proxies.txt` - A list of proxies separated by new line (Format: `USERNAME:PASSWORD:HOST:PORT`)
  


## Admin API docs
https://github.com/acheong08/ChatGPT-to-API/blob/master/docs/admin.md

## API usage docs
https://platform.openai.com/docs/api-reference/chat
