package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	flag.Parse()

	apiKey := os.Getenv("TG_BOT_SECRET")
	if apiKey == "" {
		log.Fatalln("Need to set the Telegram bot secret in envar TG_BOT_SECRET, got empty string")
	}

	pipeName := flag.Arg(0)
	if pipeName == "" {
		log.Fatalln("no pipe was provided to read from")
	}

	absPipeName, err := filepath.Abs(pipeName)
	if err != nil {
		log.Fatalf("could not canonicalize pipe path %s: %v\n", pipeName, err)
	}

	pipe, err := os.Open(absPipeName)
	if err != nil {
		log.Fatalf("could not open pipe %s: %v\n", absPipeName, err)
	}
	defer pipe.Close()

	// Hold the pipe open for writing but write nothing into it
	// this allows us to keep the FIFO alive for as long as this
	// program is running
	pipeWrite, err := os.Create(absPipeName)
	if err != nil {
		log.Fatalf("could not open pipe %s for writing: %v\n", absPipeName, err)
	}
	defer pipeWrite.Close()

	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		log.Fatalf("error connecting to bot api: %v\n", err)
	}

	scanner := bufio.NewScanner(pipe)

PipeScanLoop:
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineSplit := strings.SplitN(line, " ", 2)

		var chatID int64
		var text string
		for i, s := range lineSplit {
			if i == 0 {
				id, err := strconv.Atoi(s)
				if err != nil {
					log.Printf("error parsing chat id %s: %v\n", s, err)
					break
				}
				chatID = int64(id)
			} else if i == 1 {
				text = s
			}
		}

		if chatID == 0 {
			// Cannot send a message if there's no chat id
			continue PipeScanLoop
		}

		log.Printf("Sending to chat %d: %s", chatID, text)
		msg := tgbotapi.NewMessage(chatID, text)
		bot.Send(msg)
	}

	if err = scanner.Err(); err != nil {
		log.Println("pipe reading finished")
		log.Fatalf("error when reading from pipe: %v\n", err)
	}
}
