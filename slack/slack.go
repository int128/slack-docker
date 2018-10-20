// Package slack provides Slack Incoming Webhooks.
// See also https://api.slack.com/docs/message-formatting
package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Message is a message sent to an incoming webhook.
type Message struct {
	Username    string      `json:"username,omitempty"`
	IconEmoji   string      `json:"icon_emoji,omitempty"`
	IconURL     string      `json:"icon_url,omitempty"`
	Text        string      `json:"text,omitempty"`
	Attachments Attachments `json:"attachments,omitempty"`
}

// Attachment is an attachment of a message.
type Attachment struct {
	Fallback   string           `json:"fallback,omitempty"`
	Color      string           `json:"color,omitempty"`
	Pretext    string           `json:"pretext,omitempty"`
	AuthorName string           `json:"author_name,omitempty"`
	AuthorLink string           `json:"author_link,omitempty"`
	AuthorIcon string           `json:"author_icon,omitempty"`
	Title      string           `json:"title,omitempty"`
	TitleLink  string           `json:"title_link,omitempty"`
	Text       string           `json:"text,omitempty"`
	Fields     AttachmentFields `json:"fields,omitempty"`
	ImageURL   string           `json:"image_url,omitempty"`
	ThumbURL   string           `json:"thumb_url,omitempty"`
	Footer     string           `json:"footer,omitempty"`
	FooterIcon string           `json:"footer_icon,omitempty"`
	Timestamp  int64            `json:"ts,omitempty"`
	MrkdwnIn   []string         `json:"mrkdwn_in,omitempty"` // pretext, text, fields
}

// Attachments is an array of Attachment.
type Attachments []Attachment

// AttachmentField is a field in the attachment.
type AttachmentField struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}

// AttachmentFields is an array of AttachmentField.
type AttachmentFields []AttachmentField

// Send sends the message to the incoming webhook.
func Send(webhookURL string, message *Message) error {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(message)
	resp, err := http.Post(webhookURL, "application/json", b)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API error: %d %s: %s", resp.StatusCode, resp.Status, string(b))
	}
	return nil
}
