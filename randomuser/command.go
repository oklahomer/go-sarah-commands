package randomuser

import (
	"encoding/json"
	"fmt"
	"github.com/oklahomer/go-sarah"
	"github.com/oklahomer/go-sarah/slack"
	"github.com/oklahomer/go-sarah/slack/webapi"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// Command provides default setup of random user command.
// If different setup with another identifier, match pattern, etc. directly feed CommandFunc to preferred CommandBuilder
var Command = sarah.NewCommandBuilder().
	Identifier("random_user").
	InputExample(".randomuser | .random user").
	MatchPattern(regexp.MustCompile(`^\.random\s*user`)).
	Func(CommandFunc).
	MustBuild()

// CommandFunc provides the core function of random user.
func CommandFunc(ctx context.Context, input sarah.Input) (*sarah.CommandResponse, error) {
	resp, err := ctxhttp.Get(ctx, http.DefaultClient, "http://api.randomuser.me/")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status error. status: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := &APIResponse{}
	json.Unmarshal(body, data)
	if err := json.Unmarshal(body, data); err != nil {
		return nil, err
	}

	user := data.Results[0]
	limeGreen := "#32CD32"
	address := fmt.Sprintf("%s, %s, %s, %d", user.Location.Street, user.Location.City, user.Location.State, user.Location.Postcode)
	attachments := []*webapi.MessageAttachment{
		{
			Fallback: fmt.Sprintf("%s. %s %s", user.Name.Title, user.Name.First, user.Name.Last),
			Title:    "Name",
			Color:    limeGreen,
			ImageURL: user.Picture.Thumbnail,
			Fields: []*webapi.AttachmentField{
				{
					Title: "First Name",
					Value: strings.Title(user.Name.First),
					Short: true,
				},
				{
					Title: "Last Name",
					Value: strings.Title(user.Name.Last),
					Short: true,
				},
				{
					Title: "Title",
					Value: strings.Title(user.Name.Title),
					Short: true,
				},
			},
		},
		{
			Fallback: user.Gender,
			Title:    "Gender",
			Color: func(gender string) string {
				if gender == "male" {
					return "#0000ff"
				} else {
					return "#ff66cc"
				}
			}(user.Gender),
			Fields: []*webapi.AttachmentField{
				{
					Value: strings.Title(user.Gender),
					Short: false,
				},
			},
		},
		{
			Fallback: user.BirthDate,
			Title:    "Date of Birth",
			Fields: []*webapi.AttachmentField{
				{
					Value: user.BirthDate,
					Short: false,
				},
			},
		},
		{
			Fallback: address,
			Title:    "Address",
			Fields: []*webapi.AttachmentField{
				{
					Title: "Street",
					Value: strings.Title(user.Location.Street),
					Short: true,
				},
				{
					Title: "City",
					Value: strings.Title(user.Location.City),
					Short: true,
				},
				{
					Title: "State",
					Value: strings.Title(user.Location.State),
					Short: true,
				},
				{
					Title: "Postal Code",
					Value: strconv.Itoa(user.Location.Postcode),
					Short: true,
				},
			},
		},
	}

	return slack.NewPostMessageResponse(input, "", attachments), nil
}

type APIResponse struct {
	Results []*User `json:"results"`
	Info    *Info   `json:"info"`
}

type Name struct {
	Title string `json:"title"`
	First string `json:"first"`
	Last  string `json:"last"`
}

type Location struct {
	Street   string `json:"street"`
	City     string `json:"city"`
	State    string `json:"state"`
	Postcode int    `json:"postcode"`
}

type Login struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type Picture struct {
	Large     string `json:"large"`
	Medium    string `json:"medium"`
	Thumbnail string `json:"thumbnail"`
}

type Info struct {
	Seed    string `json:"seed"`
	Results uint   `json:"results"`
	Page    uint   `json:"page"`
	Version string `json:"version"`
}

type User struct {
	Gender    string    `json:"gender"`
	Name      *Name     `json:"name"`
	Location  *Location `json:"location"`
	Email     string    `json:"email"`
	Login     *Login    `json:"login"`
	BirthDate string    `json:"dob"`
	Phone     string    `json:"phone"`
	CellPhone string    `json:"cell"`
	Picture   *Picture  `json:"picture"`
}
