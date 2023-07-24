# ChatGPT-to-API
ChatGPT のウェブサイトを使って偽 API を作る

> ## 重要
> このリポジトリに対する無償のサポートは受けられません。これは私個人の使用のために作られたもので、ドキュメントは本当に必要ないので、ドキュメントは制限され続けます。貢献者による中国語のドキュメントに、より詳細なドキュメントがあります。

**API エンドポイント: http://127.0.0.1:8080/v1/chat/completions.**

[英語ドキュメント（English Docs）](README.md)
[中国語ドキュメント（Chinese Docs）](https://github.com/xqdoo00o/ChatGPT-to-API/blob/master/README_ZH.md)
## セットアップ

### 認証

アクセストークンの取得は [OpenAIAuth](https://github.com/acheong08/OpenAIAuth/) により、アカウントのメールアドレスとパスワードで自動化されています。

`accounts.txt` - 改行で区切られたアカウントのリスト

フォーマット:
```
email:password
...
```

すべての認証されたアクセストークンは `access_tokens.json` に保存されます

アクセストークンは 14 日後に自動更新されます

注意！認証にはブロックされていない ip を使用してください。可能であれば、まず `https://chat.openai.com/` にログインして ip の可用性を確認してください。

### API認証（オプション）

OpenAI の API と同じような、この偽 API 用のカスタム API キー

`api_keys.txt` - 改行で区切られた API キーのリスト

フォーマット:
```
sk-123456
88888888
...
```

## 準備
```
git clone https://github.com/acheong08/ChatGPT-to-API
cd ChatGPT-to-API
go build
./freechatgpt
```

### 環境変数
  - `PUID` - chat.openai.com の Plus ユーザー向けのクッキーです。これは Cloudflare のレート制限を回避します
  - `SERVER_HOST` - デフォルトで 127.0.0.1 に設定
  - `SERVER_PORT` - デフォルトで 8080 に設定
  - `OPENAI_EMAIL` と `OPENAI_PASSWORD` - PUID が設定されている場合、自動的に更新されます
  - `ENABLE_HISTORY` - デフォルトで true に設定

### ファイル（オプション）
  - `proxies.txt` - 改行で区切られたプロキシのリスト

    ```
    http://127.0.0.1:8888
    ...
    ```
  - `access_tokens.json` - サイクリング用のアクセストークンの JSON 配列（あるいは、[正しいエンドポイント](https://github.com/acheong08/ChatGPT-to-API/blob/master/docs/admin.md)に PATCH リクエストを送る）
    ```
    ["access_token1", "access_token2"...]
    ```

## Admin API ドキュメント
https://github.com/acheong08/ChatGPT-to-API/blob/master/docs/admin.md

## API 使用方法ドキュメント
https://platform.openai.com/docs/api-reference/chat
