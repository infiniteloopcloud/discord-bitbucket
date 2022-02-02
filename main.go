package main

import (
	"fmt"
	"net/http"
	"os"

	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
)

const (
	Token = "TOKEN"
)

var session *discordgo.Session

func hello(w http.ResponseWriter, req *http.Request) {
	m, err := getSession().ChannelMessageSendEmbed("938346303165980692", embed.NewGenericEmbed("Example", "This is an example embed!"))
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%+v", m)
	fmt.Fprintf(w, "thanks\n")
}

func main() {
	http.HandleFunc("/hello", hello)

	http.ListenAndServe(":8090", nil)
}

func getSession() *discordgo.Session {
	if session == nil {
		var err error
		session, err = discordgo.New("Bot " + os.Getenv(Token))
		if err != nil {
			panic(err)
		}
	}
	return session
}
