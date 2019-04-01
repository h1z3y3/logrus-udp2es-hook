# Logrus-udp2es-hook

用于写入udp地址走后续自动入ES流程

## Usage

```go
import (
    "github.com/sirupsen/logrus"
    "github.com/h1z3y3/logrus-udp2es-hook"
)

func main() {

    // logrus 必须设置为json输出
    logrus.SetFormatter(&logrus.JSONFormatter{})

    // 初始化hook
    hook, err := logrus_udp2es.NewUdp2EsHook(&logrus_udp2es.Hook{
        Host: "127.0.0.1", // your udp server host
        Port: 12345, // your udp server port
        ESIndex: "test-index", // your es index
    })

    // 只记录指定Level的日志
    hook.SetLevels([]logrus.Level{logrus.ErrorLevel, logrus.WarnLevel})

    if err == nil {
    	logrus.AddHook(hook)
    } else {
    	logrus.Error("add hook error:", err)
    }

    logrus.Warning("Here is your message")
}
```

## ChangeLog

