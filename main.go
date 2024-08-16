package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sc-bot/internal/config"
	"sc-bot/internal/model"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	BotToken string
	AppID    string
	GuildID  string
)

var (
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "message",
			Description: "ask your question",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "text",
					Description: "type your question",
					Required:    true,
				},
			},
		},
		{
			Name:        "play",
			Description: "choose a song",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "track",
					Description: "take one",
					Required:    true,
					Choices:     model.GetChoices(),
				},
			},
		},
		{
			Name:        "stop",
			Description: "stop playing",
		},
		{
			Name:        "random",
			Description: "play random tracks",
		},
	}

	commandsHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"message": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Thinking....",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}

			data := i.ApplicationCommandData()
			newMessage := data.Options[0].StringValue()

			s.ChannelMessageSend(i.ChannelID, model.Dialog(newMessage))
		},
		"play": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			userID := i.Member.User.ID
			member, err := s.GuildMember(GuildID, userID)
			if err != nil {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "User not found....",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				if err != nil {
					panic(err)
				}

				return
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Start play to " + member.User.Username,
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}

			guild, err := s.State.Guild(GuildID)
			if err != nil {
				return
			}

			audioID := i.ApplicationCommandData().Options[0].StringValue()

			for _, vs := range guild.VoiceStates {
				if vs.UserID == member.User.ID {
					s.ChannelMessageSend(i.ChannelID, model.Play(s, audioID, vs.ChannelID))

					return
				}
			}

			s.ChannelMessageSend(i.ChannelID, "Join any channel!")
		},
		"stop": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Stop play audio...",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}

			s.ChannelMessageSend(i.ChannelID, model.Stop())
		},
		"random": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			userID := i.Member.User.ID
			member, err := s.GuildMember(GuildID, userID)
			if err != nil {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "User not found....",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				if err != nil {
					panic(err)
				}

				return
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Start play to " + member.User.Username,
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}

			guild, err := s.State.Guild(GuildID)
			if err != nil {
				return
			}

			audioID := model.GetRandomTrack()

			for _, vs := range guild.VoiceStates {
				if vs.UserID == member.User.ID {
					s.ChannelMessageSend(i.ChannelID, model.Play(s, audioID, vs.ChannelID))

					return
				}
			}

			s.ChannelMessageSend(i.ChannelID, "Join any channel!")
		},
	}
)

func init() {
	config := config.MustLoad()

	BotToken = config.Application.BotToken
	AppID = config.Application.AppID
	GuildID = config.Application.GuildID
}

func main() {
	flag.Parse()

	sess, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatal(err)
	}

	sess.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandsHandler[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	sess.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := sess.ApplicationCommandCreate(sess.State.User.ID, GuildID, v)
		if err != nil {
			log.Fatal(err)
		}

		registeredCommands[i] = cmd
	}

	fmt.Println("the bot is online!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	if *RemoveCommands {
		log.Println("Removing commands...")

		for _, v := range registeredCommands {
			err := sess.ApplicationCommandDelete(sess.State.User.ID, GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Shutting down.")
}
