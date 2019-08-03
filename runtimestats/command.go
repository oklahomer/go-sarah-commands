package runtimestats

import (
	"context"
	"github.com/oklahomer/go-sarah/v2"
	"github.com/oklahomer/go-sarah/v2/slack"
	"github.com/oklahomer/golack/slackobject"
	"github.com/oklahomer/golack/webapi"
	"regexp"
	"runtime"
	"strconv"
)

type ScheduleConfig struct {
	TaskSchedule string                `json:"schedule" yaml:"schedule"`
	ChannelID    slackobject.ChannelID `json:"channel" yaml:"channel"`
}

func (c *ScheduleConfig) Schedule() string {
	return c.TaskSchedule
}

func (c *ScheduleConfig) DefaultDestination() sarah.OutputDestination {
	return c.ChannelID
}

func SlackScheduledTaskProps(config *ScheduleConfig) *sarah.ScheduledTaskProps {
	return sarah.NewScheduledTaskPropsBuilder().
		BotType(slack.SLACK).
		ConfigurableFunc(config, func(_ context.Context, conf sarah.TaskConfig) ([]*sarah.ScheduledTaskResult, error) {
			typedConfig := conf.(*ScheduleConfig)
			return []*sarah.ScheduledTaskResult{{
				Content:     webapi.NewPostMessage(typedConfig.ChannelID, "").WithAttachments(messageAttachments()),
				Destination: typedConfig.ChannelID,
			}}, nil
		}).
		Identifier("runtime").
		MustBuild()
}

var SlackProps = sarah.NewCommandPropsBuilder().
	BotType(slack.SLACK).
	Identifier("runtime").
	Instruction(".runtime").
	MatchPattern(regexp.MustCompile(`^\.runtime`)).
	Func(func(_ context.Context, input sarah.Input) (*sarah.CommandResponse, error) {
		return slack.NewResponse(input, "", slack.RespWithAttachments(messageAttachments()));
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
