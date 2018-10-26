// Package slack provides Slack Incoming Webhooks.
// See also https://api.slack.com/docs/messages
package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Message represents a message sent via Incoming WebHook API.
// See https://api.slack.com/docs/message-formatting for details.
type Message struct {
	Username    string       `json:"username,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	IconURL     string       `json:"icon_url,omitempty"`
	Text        string       `json:"text,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// Attachment represents an attachment of a message.
// See https://api.slack.com/docs/message-attachments for details.
type Attachment struct {
	Fallback   string            `json:"fallback,omitempty"`
	Color      string            `json:"color,omitempty"`
	Pretext    string            `json:"pretext,omitempty"`
	AuthorName string            `json:"author_name,omitempty"`
	AuthorLink string            `json:"author_link,omitempty"`
	AuthorIcon string            `json:"author_icon,omitempty"`
	Title      string            `json:"title,omitempty"`
	TitleLink  string            `json:"title_link,omitempty"`
	Text       string            `json:"text,omitempty"`
	Fields     []AttachmentField `json:"fields,omitempty"`
	ImageURL   string            `json:"image_url,omitempty"`
	ThumbURL   string            `json:"thumb_url,omitempty"`
	Footer     string            `json:"footer,omitempty"`
	FooterIcon string            `json:"footer_icon,omitempty"`
	Timestamp  int64             `json:"ts,omitempty"`
	MrkdwnIn   []string          `json:"mrkdwn_in,omitempty"` // pretext, text, fields
}

// AttachmentField represents a field in an attachment.
// See https://api.slack.com/docs/message-attachments for details.
type AttachmentField struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}

// Send sends the message to Slack via Incomming WebHook API.
// It returns an error if the API did not return 2xx status.
func Send(webHookURL string, message *Message) error {
	if message == nil {
		return fmt.Errorf("message is nil")
	}
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(message); err != nil {
		return fmt.Errorf("Could not encode JSON: %s", err)
	}
	resp, err := http.Post(webHookURL, "application/json", &b)
	if err != nil {
		return fmt.Errorf("Could not send the request: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Slack API returned %s (could not read body: %s)", resp.Status, err)
		}
		return fmt.Errorf("Slack API returned %s: %s", resp.Status, string(b))
	}
	return nil
}
