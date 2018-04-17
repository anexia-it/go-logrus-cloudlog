# go-logrus-cloudlog

go-logrus-cloud is a hook for [logrus](https://github.com/sirupsen/logrus) logging framework which logs to Anexia CloudLog using the [go cloudlog client](https://github.com/anexia-it/go-cloudlog).

## Install
```
go get -u github.com/anexia-it/go-logrus-cloudlog
```

## Quickstart
Please also see the [cloudlog client README](https://github.com/anexia-it/go-cloudlog/blob/master/README.md)
```go
package main

import "github.com/anexia-it/go-logrus-cloudlog"

func main() {

    // new cloudlog hook
    hook := cloudlogrus.Must("index", "ca.pem", "cert.pem", "cert.key")
    
    // add the new hook to logrus
    log.AddHook(hook)

    // use logrus
    log.Info("my first cloudlog log message")
}

```

### Custom Log level
A custom log level can be set on the hook.
This way you can for example log everything to the console and only Errors to cloudlog
```go
hook.SetLevel(logrus.WarnLevel)
```

### Custom client and mapper
Custom cloudlog client and a custom map function can be provided.
The map function controls the format of the logs within cloudlog
```go
cloudlogrus.NewCustomHook(customClient, customMapFunction, logrus.ErrorLevel)
```