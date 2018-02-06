package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/nlopes/slack"
)

// The GET request to the Kubernetes event watch API returns a JSON object
// which unmarshals into this Response type.
type Response struct {
	Type   string `json:"type"`
	Object Event  `json:"object"`
}

// The Event type and its child-types, contain only the values of the response
// that our alerts currently care about.
type Event struct {
	Source         EventSource         `json:"source"`
	InvolvedObject EventInvolvedObject `json:"involvedObject"`
	Metadata       EventMetadata       `json:"metadata"`
	Reason         string              `json:"reason"`
	Message        string              `json:"message"`
	FirstTimestamp time.Time           `json:"firstTimestamp"`
	LastTimestamp  time.Time           `json:"lastTimestamp"`
	Count          int                 `json:"count"`
}

type EventMetadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type EventSource struct {
	Component string `json:"component"`
	Host      string `json:"host"`
}

type EventInvolvedObject struct {
	Kind string `json:"kind"`
}

// Sends a message to the Slack channel about the Event.
func send_message(e Event, emoji string) error {
	api := slack.New(os.Getenv("SLACK_TOKEN"))

	msg := fmt.Sprintf("%s %s on %s", emoji, e.Message, os.Getenv("CLUSTER_NAME"))
	params := slack.PostMessageParameters{}

	channelID, timestamp, err := api.PostMessage(os.Getenv("SLACK_CHANNEL"), msg, params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}

	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	return nil
}

func main() {
	url := fmt.Sprintf("http://localhost:8001/api/v1/events?watch=true")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)

	if resp.StatusCode != 200 {
		log.Printf(string(resp.Status) + ": " + string(resp.StatusCode))
		log.Fatal("Non 200 status code returned from Kubernetes API.")
	}
	for {
		var r Response
		if err := dec.Decode(&r); err == io.EOF {
			log.Printf("EOF detected.")
			break
		} else if err != nil {
			// Debug output to help when we've failed to decode.
			htmlData, er := ioutil.ReadAll(resp.Body)
			if er != nil {
				log.Printf("Already failed to decode, but also failed to read response for log output.")
			}
			log.Printf(string(htmlData))
			log.Fatal("Decode: ", err)
		}
		e := r.Object
		log.Printf("Type: %s", r.Type)
		// Log all events for now.
		log.Printf("Reason: %s\nMessage: %s\nCount: %s\nFirstTimestamp: %s\nLastTimestamp: %s\n\n", e.Reason, e.Message, strconv.Itoa(e.Count), e.FirstTimestamp, e.LastTimestamp)

		send := false
		emoji := ""

		// @todo refactor the configuration of which things to post.
		if e.Reason == "SuccessfulCreate" {
			send = true
			emoji = ":white_check_mark:"
		} else if e.Reason == "NodeReady" {
			send = true
			emoji = ":white_check_mark:"
		} else if e.Reason == "Failed" {
			send = true
			emoji = ":warning:"			
		} else if e.Reason == "NodeNotReady" {
			send = true
			emoji = ":warning:"
		} else if e.Reason == "NodeOutOfDisk" {
			send = true
			emoji = ":skull:"
		}

		// For now, dont alert multiple times, except if it's a backoff
		if e.Count > 1 {
			send = false
		}
		if e.Reason == "BackOff" && e.Count == 3 {
			send = true
			emoji = ":warning:"
		}

		// Do not send any events that are more than 1 minute old.
		// This assumes events are processed quickly (very likely)
		// in exchange for not re-notifying of events after a crash
		// or fresh start.
		diff := time.Now().Sub(e.LastTimestamp)
		diffMinutes := int(diff.Minutes())
		if diffMinutes > 1 {
			log.Printf("Supressed %s minute old message: %s", strconv.Itoa(diffMinutes), e.Message)
			send = false
		}

		if send {
			err = send_message(e, emoji)
			if err != nil {
				log.Fatal("send_message: ", err)
			}
		}
	}
}
