package util

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type SenderOptions struct {
	// Sender display name
	Name string
	// Sender image as an emoji
	Emoji string
	// Sender image as an image URL
	Image string
}

type BodyElement struct {
	Text   string `json:"text"`
	Blocks []any  `json:"blocks"`
}

func PlainMessage(message string) BodyElement {
	return BodyElement{Text: message}
}

func SendMessage(channel string, body BodyElement, options *SenderOptions) error {
	if options == nil {
		options = &SenderOptions{}
	}

	requestBody := struct {
		Channel string `json:"channel"`
		Text    string `json:"text"`
		Blocks  []any  `json:"blocks,omitempty"`
		Name    string `json:"username,omitempty"`
		Emoji   string `json:"icon_emoji,omitempty"`
		Image   string `json:"icon_url,omitempty"`
	}{
		Channel: channel,
		Text:    body.Text,
		Blocks:  body.Blocks,
		Name:    options.Name,
		Emoji:   options.Emoji,
		Image:   options.Image,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	r, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	r.Header.Set("Content-Type", "application/json; charset=utf8")
	r.Header.Set("Authorization", "Bearer "+os.Getenv("SLACK_API_KEY"))

	client := &http.Client{}

	res, err := client.Do(r)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	jsonResponse := &struct {
		Ok    bool   `json:"ok"`
		Error string `json:"error"`
	}{}

	err = json.NewDecoder(res.Body).Decode(jsonResponse)
	if err != nil {
		return err
	}

	return nil
}
