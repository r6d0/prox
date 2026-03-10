# About

Prox is a simple and easy-to-use forward proxy server on Go.

# Features
 - HTTP forwarding;
 - HTTPS tunneling through CONNECT;
 - JSON configuration;
 - Pure Go and no external dependencies;

 # How to use
 ## Build from source code
 ```
 git clone https://github.com/r6d0/prox.git
 cd ./prox
 go run ./cmd/main.go
 ```

## How to configure

### All available configuration properties: 

| Property | Description |
| --- | --- |
| port | The server port. It should be from 0 to 65535. The default value is 9999. |
| protocol | The server protocol. It should be HTTP or HTTPS. The default value is HTTP. |
| request.timeout | The request timeout. It has `time.Duration` syntax. The default value is 2 seconds. |
| request.forwardedHeader | If the value is true, Prox will add the `X-Forwarded-For` header. The default value is false. |
| tls.certFile | The certificate file location. The default value is empty. |
| tls.keyFile | The private key file location. The default value is empty. |
| log.level | The logger level. It should be INFO, ERROR or DEBUG. The default value is INFO. |

### JSON configuration example:

 ```json
 {
  "port": 9999,
  "protocol": "HTTP",
  "request": {
    "timeout": "2s",
    "forwardedHeader": false
  },
  "tls": {
    "certFile": "/usr/cert/fullchain.pem",
    "keyFile": "/usr/cert/privatekey.pem"
  },
  "log": {
    "level": "INFO"
  }
 }
 ```