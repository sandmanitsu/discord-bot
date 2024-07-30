package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const (
	BotToken = "MTI2Nzc2NjM0Mzg1NjQyNzExOQ.GYCBy3.FgHYq9hvTQX26UaJoA6vNeKBeD_XdjTuMkvaV8"
	AppID    = "1267766343856427119"
	GuildID  = "1216051459053977621" // tested server Id
)

var sess *discordgo.Session

var (
	commands = []discordgo.ApplicationCommand{
		{
			Name: "hello",
			Type: discordgo.ChatApplicationCommand,
		},
	}

	commandsHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"hello": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.ChannelMessageSend(i.ChannelID, "world!")
		},
	}
)

func main() {
	sess, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatal(err)
	}

	// sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
	// 	if m.Author.ID == s.State.User.ID {
	// 		return
	// 	}

	// 	if m.Content == "hello" {
	// 		s.ChannelMessageSend(m.ChannelID, "world!")
	// 	}
	// })

	cmdIDs := make(map[string]string, len(commands))

	for _, cmd := range commands {
		rcmd, err := sess.ApplicationCommandCreate(AppID, GuildID, &cmd)
		if err != nil {
			log.Fatal(err)
		}

		cmdIDs[rcmd.ID] = rcmd.Name
	}

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	fmt.Println("the bot is online!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
