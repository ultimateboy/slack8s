package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	//	"strings"

	"github.com/gorilla/schema"
)

const (
	KubectlToken = "yuLR7s35eJEoS76k5C0xj6Ot"
)

type SlashCommand struct {
	Token       string `schema:"token"`
	TeamID      string `schema:"team_id"`
	TeamDomain  string `schema:"team_domain"`
	ChannelID   string `schema:"channel_id"`
	ChannelName string `schema:"channel_name"`
	UserID      string `schema:"user_id"`
	UserName    string `schema:"user_name"`
	Command     string `schema:"command"`
	Text        string `schema:"text"`
	ResponseURL string `schema:"response_url"`
}

type Response struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Text  string `json:"text"`
	Color string `json:"color"`
}

// Responds to /kubectl slack command POST.
func kubectl(w http.ResponseWriter, r *http.Request) {
	// Parse the posted data.
	err := r.ParseForm()
	if err != nil {
		log.Fatal("ParseForm: ", err)
	}

	// Decode into a SlashCommand type.
	cmd := new(SlashCommand)
	decoder := schema.NewDecoder()
	err = decoder.Decode(cmd, r.Form)
	if err != nil {
		log.Fatal("Decode: ", err)
	}

	// Validate the token.
	if cmd.Token != KubectlToken {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Invalid token.")
		return
	}

	// Respond.
	var res Response
	res.Text = "I'm alive!"
	w.Header().Set("content-type", "application/json")
	resJson, err := json.Marshal(res)
	if err != nil {
		log.Fatal("Marshal: ", err)
	}
	fmt.Fprintf(w, string(resJson))
}

func main() {
	http.HandleFunc("/kubectl", kubectl)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
