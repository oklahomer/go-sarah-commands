package pick

import (
	"github.com/oklahomer/go-sarah"
	"github.com/oklahomer/go-sarah/slack"
	"golang.org/x/net/context"
	"math/rand"
	"regexp"
	"strings"
)

var matchPattern = regexp.MustCompile(`^\.pick\s+.*`)
var Command = sarah.NewCommandBuilder().
	Identifier("pick").
	InputExample(".pick Foo, Bar").
	MatchPattern(matchPattern).
	Func(CommandFunc).
	MustBuild()

func CommandFunc(_ context.Context, input sarah.Input) (*sarah.CommandResponse, error) {
	candidates := strings.Split(sarah.StripMessage(matchPattern, input.Message()), ",")
	if len(candidates) == 1 {
		msg := "Please input comma separated candidates. e.g. Foo, Bar, Buzz.\nOr input .abort to quit."
		return slack.NewStringResponseWithNext(msg, CommandFunc), nil
	}
	chosen := candidates[rand.Intn(len(candidates))]
	return slack.NewStringResponse(chosen), nil
}
