Some commands for [sarah](https://github.com/oklahomer/go-sarah).

```go
package main

import (
        "github.com/oklahomer/go-sarah"
        "github.com/oklahomer/go-sarah/slack"
        "github.com/oklahomer/go-sarah-commands/giphy"
        "github.com/oklahomer/go-sarah-commands/pick"
        "github.com/oklahomer/go-sarah-commands/randomuser"
        "github.com/oklahomer/go-sarah-commands/urlextractor"
        "golang.org/x/net/context"
        "gopkg.in/yaml.v2"
        "io/ioutil"
)

func main() {
        // Basic setup
        configBuf, _ := ioutil.ReadFile("/path/to/adapter/config.yaml")
        slackConfig := slack.NewConfig()
        yaml.Unmarshal(configBuf, slackConfig)
        slackBot := sarah.NewBot(slack.NewAdapter(slackConfig), sarah.NewCacheConfig())
        
        // Registering commands
        slackBot.AppendCommand(giphy.SlackCommand)
        slackBot.AppendCommand(pick.SlackCommand)
        slackBot.AppendCommand(randomuser.SlackCommand)
        slackBot.AppendCommand(urlextractor.SlackCommand)
        
        // Initialize Runner and start bot interaction.
        runner := sarah.NewRunner(sarah.NewConfig())
        runner.RegisterBot(slackBot)

        // Start interaction
        rootCtx := context.Background()
        runnerCtx, _ := context.WithCancel(rootCtx)
        runner.Run(runnerCtx)
}
```