# API Documentation:
## openaiHandler:

This API endpoint receives a POST request with a JSON body that contains OpenAI credentials. The API updates the environment variables "OPENAI_EMAIL" and "OPENAI_PASSWORD" with the provided credentials.

HTTP method: PATCH

Endpoint: /openai

Request body:

```json
{
    "OpenAI_Email": "string",
    "OpenAI_Password": "string"
}
```

Response status codes:
- 200 OK: The OpenAI credentials were successfully updated.
- 400 Bad Request: The JSON in the request body is invalid.

## passwordHandler:

This API endpoint receives a POST request with a JSON body that contains a new password. The API updates the global variable "ADMIN_PASSWORD" and the environment variable "ADMIN_PASSWORD" with the new password.

HTTP method: PATCH

Endpoint: /password

Request body:

```json
{
    "password": "string"
}
```

Response status codes:
- 200 OK: The password was successfully updated.
- 400 Bad Request: The password is missing or not provided in the request body.

## puidHandler:

This API endpoint receives a POST request with a JSON body that contains a new PUID (Personal User ID). The API updates the environment variable "PUID" with the new PUID.

HTTP method: PATCH

Endpoint: /puid

Request body:

```json
{
    "puid": "string"
}
```

Response status codes:
- 200 OK: The PUID was successfully updated.
- 400 Bad Request: The PUID is missing or not provided in the request body.

## tokensHandler:

This API endpoint receives a POST request with a JSON body that contains an array of request tokens. The API updates the value of the global variable "ACCESS_TOKENS" with a new access token generated from the request tokens provided in the request body.

HTTP method: PATCH

Endpoint: /tokens

Request body:

```json
[
    "string", "..."
]
```

Response status codes:
- 200 OK: The ACCESS_TOKENS variable was successfully updated.
- 400 Bad Request: The request tokens are missing or not provided in the request body.


