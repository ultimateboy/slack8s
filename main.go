package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/nlopes/slack"
)

const (
	KubeHost       = ""
	SlackToken     = ""
	SlackChannelID = ""
)

type Response struct {
	Type   string `json:"type"`
	Object Event  `json:"object"`
}

type Event struct {
	Source         EventSource         `json:"source"`
	InvolvedObject EventInvolvedObject `json:"involvedObject"`
	Metadata       EventMetadata       `json:"metadata"`
	Reason         string              `json:"reason"`
	Message        string              `json:"message"`
}

type EventMetadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type EventSource struct {
	Component string `json:"component"`
}

type EventInvolvedObject struct {
	Kind string `json:"kind"`
}

func send_message(e Event) error {
	api := slack.New(SlackToken)
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		AuthorName: e.Message,
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Object",
				Value: e.InvolvedObject.Kind,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Name",
				Value: e.Metadata.Name,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Reason",
				Value: e.Reason,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Component",
				Value: e.Source.Component,
				Short: true,
			},
		},
	}

	// Colors!
	if strings.HasPrefix(e.Reason, "Successful") {
		attachment.Color = "good"
	} else if strings.HasPrefix(e.Reason, "Failed") {
		attachment.Color = "danger"
	}
	params.Attachments = []slack.Attachment{attachment}

	channelID, timestamp, err := api.PostMessage(SlackChannelID, "", params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}

	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	return nil
}

func main() {
	url := fmt.Sprintf("%s/api/v1/events?watch=true", KubeHost)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	for {
		var r Response
		if err := dec.Decode(&r); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		e := r.Object
		// Log all events for now.
		fmt.Printf("%s: %s\n", e.Reason, e.Message)

		send := false

		// @todo refactor the configuration of which things to post.
		if e.Reason == "SuccessfulCreate" {
			send = true
		}
		if send {
			err = send_message(e)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
