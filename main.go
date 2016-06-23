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

	// Setup a watcher for pods
	eventClient := kubeClient.Events(api.NamespaceAll)
	options := api.ListOptions{LabelSelector: labels.Everything()}
	w, err := eventClient.Watch(options)
	if err != nil {
		log.Fatalf("Failed to set up watch: %v", err)
	}
	select {
	case watchEvent, _ := <-w.ResultChan():
		send := false
		color := ""
		if watchEvent.Type == watch.Added {
			send = true
			color = "good"
		} else if watchEvent.Type == watch.Deleted {
			send = true
			color = "warning"
		}
		fmt.Printf("%+v\n", watchEvent)

		if send {
			event, _ := watchEvent.Object.(*api.Event)
			err = sendMessage(event, color)
			if err != nil {
				log.Fatal("sendMessage: ", err)
			}
		}
	}
}
