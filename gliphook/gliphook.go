// Package gliphook provides simple notification feature
package gliphook

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// GlipMessageSimple structure for simple webhook
const (
	contentTypeJson = "application/json"
)

type GlipMessageSimple struct {
	Activity string `json:"activity"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	IconUrl  string `json:"icon,omitempty"`
}

// GlipMessageCard refers to card view in glip
type GlipMessageCard struct {
	// https://developers.ringcentral.com/guide/team-messaging/manual/formatting
	Title       string                      `json:"title,omitempty"`
	Body        string                      `json:"body,omitempty"`
	IconUrl     string                      `json:"icon,omitempty"`
	Attachments []GlipMessageCardAttachment `json:"attachments"`
}

// GlipMessageCardAttachment is used to define multiple cards
type GlipMessageCardAttachment struct {
	AuthorName string                 `json:"author_name,omitempty"`
	AuthorLink string                 `json:"author_link,omitempty"`
	AuthorIcon string                 `json:"author_icon,omitempty"`
	Color      string                 `json:"color,omitempty"`
	ImageUrl   string                 `json:"image_url,omitempty"`
	Fallback   string                 `json:"fallback"`
	Fields     []GlipMessageCardField `json:"fields,omitempty"`
	Footer     string                 `json:"footer,omitempty"`
	FooterIcon string                 `json:"footer_icon,omitempty"`
	Pretext    string                 `json:"pretext,omitempty"`
	Text       string                 `json:"text"`
	Type       string                 `json:"type"`
	ThumbUrl   string                 `json:"thumb_url,omitempty"`
	Title      string                 `json:"title"`
	TitleLink  string                 `json:"title_link,omitempty"`
}

// GlipMessageCardField is used to define multiple fields
type GlipMessageCardField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// GlipMessage interface connects multiple glip message types
type GlipMessage interface {
	glipMessageMarshal() (data []byte, err error)
}

// GlipSendHook posts webhook message to passed url
func GlipSendHook(webhookUrl string, message GlipMessage) error {
	messageType := contentTypeJson

	data, err := message.glipMessageMarshal()
	if err != nil {
		return err
	}
	_, err = http.Post(webhookUrl, messageType, bytes.NewBuffer(data))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Simple functions to use single interface
func (message *GlipMessageSimple) glipMessageMarshal() (data []byte, err error) {
	return (json.Marshal(message))
}

func (message *GlipMessageCard) glipMessageMarshal() (data []byte, err error) {
	return (json.Marshal(message))
}
