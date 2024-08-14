package model

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sc-bot/internal/config"
	"sc-bot/internal/disk"
	"sc-bot/internal/messages"
	"strings"

	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
)

var (
	Token    string
	ModelURL string
	GuildId  string

	MessageHistory messages.MessageHistory
)

func init() {
	config := config.MustLoad()

	Token = config.Model.Token
	ModelURL = config.Model.ModelURL
	GuildId = config.Application.GuildID

	MessageHistory = *messages.New()
}

func Dialog(newMessage string) string {
	// todo проверить использует ли бот историю диалога
	// MessageHistory.AppendToHistory("user", newMessage)

	return Request(newMessage)
}

func Play(s *discordgo.Session, id string, channelId string) string {
	audioURL := "https://drive.google.com/uc?export=download&id=" + id

	voice, err := s.ChannelVoiceJoin(GuildId, channelId, false, false)
	if err != nil {
		log.Fatalf("Error joining voice channel: %v", err)
	}
	defer voice.Disconnect()

	ffmpeg := exec.Command("ffmpeg", "-i", audioURL, "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	ffmpegStdout, err := ffmpeg.StdoutPipe()
	if err != nil {
		log.Fatalf("Error creating StdoutPipe for FFmpeg: %v", err)
	}

	err = ffmpeg.Start()
	if err != nil {
		log.Fatalf("Error starting FFmpeg: %v", err)
	}

	voice.Speaking(true)
	defer voice.Speaking(false)

	opusEncoder, err := gopus.NewEncoder(48000, 2, gopus.Audio)
	if err != nil {
		log.Fatalf("Error creating Opus encoder: %v", err)
	}

	buffer := make([]byte, 960*2*2) // Buffer 20ms 48kHz stereo PCM
	for {
		n, err := ffmpegStdout.Read(buffer)
		if n > 0 {
			pcmData := make([]int16, n/2)
			for i := 0; i < len(pcmData); i++ {
				pcmData[i] = int16(binary.LittleEndian.Uint16(buffer[i*2 : (i+1)*2]))
			}

			opusData, err := opusEncoder.Encode(pcmData, 960, 4000)
			if err != nil {
				log.Fatalf("Error encoding PCM to Opus: %v", err)
			}
			voice.OpusSend <- opusData
		}
		if err == io.EOF {
			fmt.Println("End of stream")
			break
		}
		if err != nil {
			log.Fatalf("Error reading from FFmpeg stdout: %v", err)
		}
	}

	ffmpeg.Wait()
	return "Finished playing audio"
}

func GetChoices() []*discordgo.ApplicationCommandOptionChoice {
	service, err := disk.GetService()
	if err != nil {
		fmt.Printf("Error getting service: %v", err)
	}

	list := disk.ListFilesInFolder(service, "1KaLJMxkFQ8daK39Sl8Do6jgeFDTkDlD7") // Rainy Nights Of 1988
	choices := []*discordgo.ApplicationCommandOptionChoice{}

	for _, v := range list {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  formatName(v.Name),
			Value: v.Id,
		})
	}

	return choices
}

func formatName(name string) string {
	formName := strings.Split(name, "-")

	if len(formName) == 2 {
		return strings.TrimSpace(formName[1])
	}

	return "can't read a name"
}
