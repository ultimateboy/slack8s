package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"

	"k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/watch"
)

// Sends a message to the Slack channel about the Event.
func sendMessage(e *api.Event, color string) error {
	api := slack.New(os.Getenv("SLACK_TOKEN"))
	params := slack.PostMessageParameters{}
	metadata := e.GetObjectMeta()
	attachment := slack.Attachment{
		// The fallback message shows in clients such as IRC or OS X notifications.
		Fallback: e.Message,
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Message",
				Value: e.Message,
			},
			slack.AttachmentField{
				Title: "Object",
				Value: e.InvolvedObject.Kind,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Name",
				Value: metadata.GetName(),
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

	// Use a color if provided, otherwise try to guess.
	if color != "" {
		attachment.Color = color
	} else if strings.HasPrefix(e.Reason, "Success") {
		attachment.Color = "good"
	} else if strings.HasPrefix(e.Reason, "Fail") {
		attachment.Color = "danger"
	}
	params.Attachments = []slack.Attachment{attachment}

	channelID, timestamp, err := api.PostMessage(os.Getenv("SLACK_CHANNEL"), "", params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}

	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	return nil
}

func main() {

	kubeClient, err := client.NewInCluster()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Setup a watcher for events.
	eventClient := kubeClient.Events(api.NamespaceAll)
	options := api.ListOptions{LabelSelector: labels.Everything()}
	w, err := eventClient.Watch(options)
	if err != nil {
		log.Fatalf("Failed to set up watch: %v", err)
	}
	select {
	case watchEvent, _ := <-w.ResultChan():

		e, _ := watchEvent.Object.(*api.Event)
		// Log all events for now.
		log.Printf("Reason: %s\nMessage: %s\nCount: %s\nFirstTimestamp: %s\nLastTimestamp: %s\n\n", e.Reason, e.Message, strconv.Itoa(e.Count), e.FirstTimestamp, e.LastTimestamp)

		send := false
		color := ""
		if watchEvent.Type == watch.Added {
			send = true
			color = "good"
		} else if watchEvent.Type == watch.Deleted {
			send = true
			color = "warning"
		} else if e.Reason == "SuccessfulCreate" {
			send = true
			color = "good"
		} else if e.Reason == "NodeReady" {
			send = true
			color = "good"
		} else if e.Reason == "NodeNotReady" {
			send = true
			color = "warning"
		} else if e.Reason == "NodeOutOfDisk" {
			send = true
			color = "danger"
		}

		// For now, dont alert multiple times, except if it's a backoff
		if e.Count > 1 {
			send = false
		}
		if e.Reason == "BackOff" && e.Count == 3 {
			send = true
			color = "danger"
		}

		if send {
			err = sendMessage(e, color)
			if err != nil {
				log.Fatalf("sendMessage: %v", err)
			}
		}
	}
}
