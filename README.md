Some commands for [sarah](https://github.com/oklahomer/go-sarah).

```go
package main

import (
    "context"
    _ "github.com/oklahomer/go-sarah-commands/giphy"
    _ "github.com/oklahomer/go-sarah-commands/goproverbs"
    _ "github.com/oklahomer/go-sarah-commands/pick"
    _ "github.com/oklahomer/go-sarah-commands/randomuser"
    _ "github.com/oklahomer/go-sarah-commands/urlextractor"
    "github.com/oklahomer/go-sarah/v2"
    "github.com/oklahomer/go-sarah/v2/log"
    "github.com/oklahomer/go-sarah/v2/slack"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    // Setup Bot
    configBuf, _ := ioutil.ReadFile("/path/to/adapter/config.yaml")
    slackConfig := slack.NewConfig()
    yaml.Unmarshal(configBuf, slackConfig)
    adapter, _ := slack.NewAdapter(slackConfig)
    storage := sarah.NewUserContextStorage(sarah.NewCacheConfig())
    slackBot, _ := sarah.NewBot(adapter, sarah.BotWithStorage(storage))
    sarah.RegisterBot(slackBot)

    // Start interaction
    rootCtx := context.Background()
    runnerCtx, cancel := context.WithCancel(rootCtx)
    err := sarah.Run(runnerCtx, sarah.NewConfig())
    if err!= nil {
        panic(err)
    }

    // Wait till signal is sent
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	select {
	case <-c:
		log.Info("Stopping due to signal reception.")
		cancel()

	}
}
```