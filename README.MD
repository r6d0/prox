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

Prox can be run with the default configuration:
```
go run ./cmd/main.go
```

Also, it can be run with user's configuration:
```
go run ./cmd/main.go ./prox_config.json
```

### JSON configuration example

JSON configuration is given below:

```json
{
  "port": 9999,
  "request": {
    "timeout": "2s",
    "bufferSize": 2048,
    "forwardedHeader": false,
    "rules": {
      "fromAddr": ["regexp(.*)"]
    }
  },
  "log": {
    "level": "DEBUG"
  }
}
```

### All available configuration properties

| Property | Description |
| --- | --- |
| [port](#port) | The server port. It should be from 0 to 65535. The default value is 9999. |
| [request.timeout](#requesttimeout) | The request timeout. It has `time.Duration` syntax. The default value is 2 seconds. |
| [request.forwardedHeader](#requestforwardedheader) | If the value is true, Prox will add the `X-Forwarded-For` header. The default value is false. |
| [request.bufferSize](#requestbuffersize) | The buffer size for copying request data. The default value is 2048 bytes (2KB). |
| [request.rules.fromAddr](#requestrulesfromaddr) | The rule for checking the client's address by expressions and blocking the request if necessary. |
| [log.level](#loglevel) | The logger level. It should be INFO, ERROR or DEBUG. The default value is INFO. |

#### port
The server port. It should be from 0 to 65535. The default value is 9999.
```json
{
  "port": 9999
}
```

#### request.timeout
The request timeout. It has `time.Duration` syntax. The default value is 2 seconds.
```json
{
  "request": {
    "timeout": "2s"
  }
}
```

 #### request.forwardedHeader
 If the value is true, Prox will add the `X-Forwarded-For` header. The default value is false.
```json
{  
  "request": {
    "forwardedHeader": false
  }
}
```

 #### request.bufferSize
 The buffer size for copying request data. The default value is 2048 bytes (2KB).
```json
{
  "request": {
    "bufferSize": 2048
  }
}
```

 #### request.rules.fromAddr
 The rule for checking the client's address by expressions and blocking the request if necessary. By default, all requests will be allowed.

 Examples:
 
 Only requests from `127.0.0.1` and `127.0.0.1` will be allowed.
```json
{
  "request": {
    "rules": {
      "fromAddr": ["127.0.0.1", "127.0.0.2"]
    }
  }
}
```

Only requests from IP addresses that start with `127` will be allowed.
```json
{
  "request": {
    "rules": {
      "fromAddr": ["regexp(127.*)"]
    }
  }
}
```

Requests from `127.0.0.1` and IP addresses that start with `128` will not be allowed.
```json
{
  "request": {
    "rules": {
      "fromAddr": ["not(127.0.0.1)", "not(regexp(128.*))"]
    }
  }
}
```

#### log.level
The logger level. It should be INFO, ERROR or DEBUG. The default value is INFO.

```json
{
  "log": {
    "level": "DEBUG"
  }
}
```