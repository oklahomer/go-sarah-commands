package pick

import (
	"context"
	"github.com/oklahomer/go-sarah/v2"
	"github.com/oklahomer/go-sarah/v2/slack"
	"math/rand"
	"regexp"
	"strings"
)

func init() {
	props := sarah.NewCommandPropsBuilder().
		BotType(slack.SLACK).
		Identifier("pick").
		Instruction(`Input ".pick Foo, Bar" to pick one option.`).
		MatchPattern(matchPattern).
		Func(slackCommandFunc).
		MustBuild()
	sarah.RegisterCommandProps(props)
}

var matchPattern = regexp.MustCompile(`^\.pick\s+.*`)

func slackCommandFunc(_ context.Context, input sarah.Input) (*sarah.CommandResponse, error) {
	candidates := strings.Split(sarah.StripMessage(matchPattern, input.Message()), ",")
	if len(candidates) == 1 {
		msg := "Please input comma separated candidates. e.g. Foo, Bar, Buzz.\nOr input .abort to quit."
		return slack.NewResponse(input, msg, slack.RespWithNext(slackCommandFunc))
	}
	chosen := candidates[rand.Intn(len(candidates))]
	return slack.NewResponse(input, chosen)
}
