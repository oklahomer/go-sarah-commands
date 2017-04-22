package runtimestats

import (
	"github.com/oklahomer/go-sarah"
	"github.com/oklahomer/go-sarah/slack"
	"github.com/oklahomer/golack/rtmapi"
	"github.com/oklahomer/golack/webapi"
	"golang.org/x/net/context"
	"regexp"
	"runtime"
	"strconv"
)

type ScheduleConfig struct {
	TaskSchedule string           `json:"schedule" yaml:"schedule"`
	ChannelID    rtmapi.ChannelID `json:"channel" yaml:"channel"`
}

func (c *ScheduleConfig) Schedule() string {
	return c.TaskSchedule
}

func (c *ScheduleConfig) DefaultDestination() sarah.OutputDestination {
	if c.ChannelID == "" {
		return nil
	}

	return c.ChannelID
}

func SlackScheduledTaskProps(config *ScheduleConfig) *sarah.ScheduledTaskProps {
	return sarah.NewScheduledTaskPropsBuilder().
		BotType(slack.SLACK).
		ConfigurableFunc(config, func(_ context.Context, conf sarah.TaskConfig) ([]*sarah.ScheduledTaskResult, error) {
			conf = conf.(*ScheduleConfig)
			return []*sarah.ScheduledTaskResult{{
				Content:     webapi.NewPostMessageWithAttachments(config.DefaultDestination().(rtmapi.ChannelID).String(), "", messageAttachments()),
				Destination: config.DefaultDestination(),
			}}, nil
		}).
		Identifier("runtime").
		MustBuild()
}

var SlackProps = sarah.NewCommandPropsBuilder().
	BotType(slack.SLACK).
	Identifier("runtime").
	InputExample(".runtime").
	MatchPattern(regexp.MustCompile(`^\.runtime`)).
	Func(func(_ context.Context, input sarah.Input) (*sarah.CommandResponse, error) {
		return slack.NewPostMessageResponse(input, "", messageAttachments()), nil
	}).
	MustBuild()

func messageAttachments() []*webapi.MessageAttachment {
	return []*webapi.MessageAttachment{
		{
			Fallback: "Current stats",
			Pretext:  "Stats:",
			Title:    "",
			Color:    "#32CD32",
			Fields: []*webapi.AttachmentField{
				{
					Title: "# of CPU",
					Value: strconv.Itoa(runtime.NumCPU()),
					Short: true,
				},
				{
					Title: "# of goroutines",
					Value: strconv.Itoa(runtime.NumGoroutine()),
					Short: true,
				},
			},
		},
	}
}