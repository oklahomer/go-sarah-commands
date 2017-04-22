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
        options := sarah.NewRunnerOptions()

        // Setup Bot
        configBuf, _ := ioutil.ReadFile("/path/to/adapter/config.yaml")
        slackConfig := slack.NewConfig()
        yaml.Unmarshal(configBuf, slackConfig)
        storage := sarah.NewUserContextStorage(sarah.NewCacheConfig())
        slackBot, _ := sarah.NewBot(slack.NewAdapter(slackConfig), sarah.BotWithStorage(storage))
        options.Append(sarah.WithBot(slackBot))
        
        // Registering commands
        options.Append(sarah.WithCommandProps(giphy.SlackProps))
        options.Append(sarah.WithCommandProps(pick.SlackProps))
        options.Append(sarah.WithCommandProps(randomuser.SlackProps))
        options.Append(sarah.WithCommandProps(urlextractor.SlackProps))
        
        // Initialize Runner and start bot interaction.
        runner, _ := sarah.NewRunner(sarah.NewConfig(), options.Arg())

        // Start interaction
        rootCtx := context.Background()
        runnerCtx, _ := context.WithCancel(rootCtx)
        runner.Run(runnerCtx)
}
```