# About

Prox is a simple and easy-to-use forward proxy server on Go.

# Features
 - HTTP forwarding;
 - HTTPS tunneling through CONNECT;
 - JSON configuration;
 - Pure Go and no external dependencies;
 - Requests filtering rules.

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
  "port": 1080,
  "request": {
    "timeout": "2s",
    "forwardedHeader": false,
    "bufferSize": 2048,
    "rules": {
      "from": "(.*)",
      "to": "(.*)"      
    }
  },
  "log": {
    "level": "ERROR"
  }
 }
```

### All available configuration properties

| Property | Description |
| --- | --- |
| [port](#port) | The server port. It should be from 0 to 65535. The default value is 1080. |
| [request.timeout](#requesttimeout) | The request timeout. It has `time.Duration` syntax. The default value is 2 seconds. |
| [request.forwardedHeader](#requestforwardedheader) | If the value is true, Prox will add the `X-Forwarded-For` header. The default value is false. |
| [request.bufferSize](#requestbuffersize) | The buffer size for copying request data. The default value is 2048 bytes (2KB). |
| [request.rules.from](#requestrulesfrom) | The rule for checking the client's address by expressions and blocking the request if necessary. |
| [request.rules.to](#requestrulesto) | The rule for checking the target host by expressions and blocking the request if necessary. |
| [log.level](#loglevel) | The logger level. It should be INFO, ERROR or DEBUG. The default value is ERROR. |

#### port
The server port. It should be from 0 to 65535. The default value is 1080.
```json
{
  "port": 1080
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

#### request.rules.from
The rule for checking the client's address by expressions and blocking the request if necessary. By default, all requests will be allowed.

Examples:
 
 Only requests from `127.0.0.1` and `128.0.0.1` will be allowed.
```json
{
  "request": {
    "rules": {
      "from": {
        "any": ["127.0.0.1", "128.0.0.1"]
      }
    }
  }
}
```

Only requests from IP addresses that start with `127` will be allowed.
```json
{
  "request": {
    "rules": {
      "from": "(127.*)"
    }
  }
}
```

Requests from `127.0.0.1` and IP addresses that start with `128` will not be allowed.
```json
{
  "request": {
    "rules": {
      "from": {
        "all": ["!127.0.0.1", "!(128.*)"]
      }
    }
  }
}
```

#### request.rules.to
The rule for checking the target host by expressions and blocking the request if necessary. By default, all requests will be allowed.

Examples:

Only requests to `mysuper.site` or `mycool.site` will be allowed.
```json
{
  "request": {
    "rules": {
      "to": {
        "any": ["mysuper.site", "mycool.site"]
      }
    }
  }
}
```

Only requests to host that starts with `my` will be allowed.
```json
{
  "request": {
    "rules": {
      "to": "(my.*)"
    }
  }
}
```

Requests to `notsuper.site` and `notcool.site` will not be allowed.
```json
{
  "request": {
    "rules": {
      "to": {
        "all": ["!notsuper.site", "!notcool.site"]
      }
    }
  }
}
```

#### log.level
The logger level. It should be INFO, ERROR or DEBUG. The default value is ERROR.

```json
{
  "log": {
    "level": "ERROR"
  }
}
```