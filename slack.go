package notify

// From https://dev.to/arunx2/simple-slack-notification-with-golang-55i2

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/projectdiscovery/retryablehttp-go"
)

// DefaultSlackTimeout to conclude operations
const DefaultSlackTimeout = 5 * time.Second

// SlackClient holding the slack communication logic
type SlackClient struct {
	client     *retryablehttp.Client
	WebHookURL string
	UserName   string
	Channel    string
	TimeOut    time.Duration
}

// SimpleSlackRequest basic request
type SimpleSlackRequest struct {
	Text      string
	IconEmoji string
}

// SlackJobNotification structure
type SlackJobNotification struct {
	Color     string
	IconEmoji string
	Details   string
	Text      string
}

// SlackMessage structure
type SlackMessage struct {
	Username    string       `json:"username,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	Text        string       `json:"text,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// Attachment of slack message
type Attachment struct {
	Color         string `json:"color,omitempty"`
	Fallback      string `json:"fallback,omitempty"`
	CallbackID    string `json:"callback_id,omitempty"`
	ID            int    `json:"id,omitempty"`
	AuthorID      string `json:"author_id,omitempty"`
	AuthorName    string `json:"author_name,omitempty"`
	AuthorSubname string `json:"author_subname,omitempty"`
	AuthorLink    string `json:"author_link,omitempty"`
	AuthorIcon    string `json:"author_icon,omitempty"`
	Title         string `json:"title,omitempty"`
	TitleLink     string `json:"title_link,omitempty"`
	Pretext       string `json:"pretext,omitempty"`
	Text          string `json:"text,omitempty"`
	ImageURL      string `json:"image_url,omitempty"`
	ThumbURL      string `json:"thumb_url,omitempty"`
	// Fields and actions are not defined.
	MarkdownIn []string    `json:"mrkdwn_in,omitempty"`
	TS         json.Number `json:"ts,omitempty"`
}

// SendSlackNotification will post to an 'Incoming Webook' url setup in Slack Apps. It accepts
// some text and the slack channel is saved within Slack.
func (sc *SlackClient) SendSlackNotification(sr SimpleSlackRequest) error {
	slackRequest := &SlackMessage{
		Text:      sr.Text,
		Username:  sc.UserName,
		IconEmoji: sr.IconEmoji,
		Channel:   sc.Channel,
	}
	return sc.sendHTTPRequest(slackRequest)
}

// SendJobNotification will post a job notification to slack
func (sc *SlackClient) SendJobNotification(job SlackJobNotification) error {
	attachment := Attachment{
		Color: job.Color,
		Text:  job.Details,
		TS:    json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	slackRequest := &SlackMessage{
		Text:        job.Text,
		Username:    sc.UserName,
		IconEmoji:   job.IconEmoji,
		Channel:     sc.Channel,
		Attachments: []Attachment{attachment},
	}
	return sc.sendHTTPRequest(slackRequest)
}

// SendError message
func (sc *SlackClient) SendError(message string, options ...string) (err error) {
	return sc.funcName("danger", message, options)
}

// SendInfo message
func (sc *SlackClient) SendInfo(message string, options ...string) (err error) {
	return sc.funcName("good", message, options)
}

// SendWarning message
func (sc *SlackClient) SendWarning(message string, options ...string) (err error) {
	return sc.funcName("warning", message, options)
}

func (sc *SlackClient) funcName(color, message string, options []string) error {
	emoji := ":hammer_and_wrench"
	if len(options) > 0 {
		emoji = options[0]
	}
	sjn := SlackJobNotification{
		Color:     color,
		IconEmoji: emoji,
		Details:   message,
	}
	return sc.SendJobNotification(sjn)
}

func (sc *SlackClient) sendHTTPRequest(slackRequest *SlackMessage) error {
	slackBody, err := json.Marshal(slackRequest)
	if err != nil {
		return err
	}
	req, err := retryablehttp.NewRequest(http.MethodPost, sc.WebHookURL, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	if sc.TimeOut == 0 {
		sc.TimeOut = DefaultSlackTimeout
	}

	resp, err := sc.client.Do(req)
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	//nolint:errcheck // silent fail
	defer resp.Body.Close()

	if string(buf) != ok {
		return err
	}
	return nil
}
