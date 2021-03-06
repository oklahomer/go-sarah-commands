package goproverbs

import (
	"context"
	"github.com/oklahomer/go-sarah/v2"
	"github.com/oklahomer/go-sarah/v2/slack"
	"github.com/oklahomer/golack/slackobject"
	"github.com/oklahomer/golack/webapi"
	"math/rand"
)

func init() {
	props := sarah.NewScheduledTaskPropsBuilder().
		BotType(slack.SLACK).
		Identifier("goproverbs").
		ConfigurableFunc(newTaskConfig(), func(_ context.Context, config sarah.TaskConfig) ([]*sarah.ScheduledTaskResult, error) {
			typedConfig := config.(*taskConfig)
			return []*sarah.ScheduledTaskResult{
				{
					Content:     webapi.NewPostMessage(typedConfig.ChannelID, "").WithAttachments(messageAttachments()),
					Destination: typedConfig.ChannelID,
				},
			}, nil
		}).
		MustBuild()
	sarah.RegisterScheduledTaskProps(props)
}

var proverbs = []*struct {
	text string
	link string
}{
	{
		text: "Don't communicate by sharing memory, share memory by communicating.",
		link: "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=2m48s",
	},
	{
		text: "Concurrency is not parallelism.",
		link: "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=3m42s",
	},
	{
		text: "Channels orchestrate; mutexes serialize.",
		link: "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=4m20s",
	},
	{
		text: "The bigger the interface, the weaker the abstraction.",
		link: "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=5m17s",
	},
	{
		text: "Make the zero value useful.",
		link: "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=6m25s",
	},
	{
		text: "interface{} says nothing.",
		link: "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=7m36s",
	},
	{
		text: "Gofmt's style is no one's favorite, yet gofmt is everyone's favorite.",
		link: "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=8m43s",
	},
	{
		text: "A little copying is better than a little dependency.",
		link: "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=9m28s",
	},
	{
		text: "Syscall must always be guarded with build tags.",
		link: "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=11m10s",
	},
	{
		text: "Cgo must always be guarded with build tags.",
		link: "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=11m53s",
	},
	{
		text: "Cgo is not Go.",
		link: "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=12m37s",
	},
	{
		text: "With the unsafe package there are no guarantees.",
		link: "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=13m49s",
	},
	{
		text: "Clear is better than clever.",
		link: "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=14m35s",
	},
	{
		text: "Reflection is never clear.",
		link: "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=15m22s",
	},
	{
		text: "Errors are values.",
		link: "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=16m13s",
	},
	{
		text: "Don't just check errors, handle them gracefully.",
		link: "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=17m25s",
	},
	{
		text: "Design the architecture, name the components, document the details.",
		link: "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=18m09s",
	},
	{
		text: "Documentation is for users.",
		link: "https://www.youtube.com/watch?v=PAAkCSZUG1c&t=19m07s",
	},
	{
		text: "Don't panic.",
		link: "https://github.com/golang/go/wiki/CodeReviewComments#dont-panic",
	},
}

type taskConfig struct {
	TaskSchedule string                `json:"schedule" yaml:"schedule"`
	ChannelID    slackobject.ChannelID `json:"channel" yaml:"channel"`
}

func (c *taskConfig) Schedule() string {
	return c.TaskSchedule
}

func (c *taskConfig) DefaultDestination() sarah.OutputDestination {
	return c.ChannelID
}

func newTaskConfig() *taskConfig {
	return &taskConfig{
		TaskSchedule: "",
		ChannelID:    "",
	}
}

func messageAttachments() []*webapi.MessageAttachment {
	proverb := proverbs[rand.Intn(len(proverbs))]
	return []*webapi.MessageAttachment{
		{
			Pretext:    "Golang's proverb",
			Fallback:   proverb.text,
			Title:      proverb.text,
			TitleLink:  proverb.link,
			Color:      "#006400",
			AuthorName: "Rob Pike",
			AuthorLink: "https://twitter.com/rob_pike",
		},
	}
}
