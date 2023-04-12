API Documentation:

puidHandler:

This API endpoint receives a POST request with a JSON body that contains a "puid" field. The API updates the value of the PUID variable with the value provided in the request body.

HTTP method: POST

Endpoint: /puid

Request body:

```json
{
    "puid": string
}
```

Response status codes:
- 200 OK: The PUID variable was successfully updated.
- 400 Bad Request: The "puid" field is missing or not provided in the request body.

passwordHandler:

This API endpoint receives a POST request with a JSON body that contains a "password" field. The API updates the value of the ADMIN_PASSWORD variable with the value provided in the request body.

HTTP method: POST

Endpoint: /password

Request body:

```json
{
    "password": string
}
```

Response status codes:
- 200 OK: The ADMIN_PASSWORD variable was successfully updated.
- 400 Bad Request: The "password" field is missing or not provided in the request body.

tokensHandler:

This API endpoint receives a POST request with a JSON body that contains an array of request tokens. The API updates the value of the ACCESS_TOKENS variable with a new access token generated from the request tokens provided in the request body.

HTTP method: POST

Endpoint: /tokens

Request body:

```json
[
    string
]
```

Response status codes:
- 200 OK: The ACCESS_TOKENS variable was successfully updated.
- 400 Bad Request: The request tokens are missing or not provided in the request body.
