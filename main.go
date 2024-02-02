package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	commands "github.com/nurpy2004/ser_bot/bot"
)

// bot's parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

// bot struct
type Bot struct {
	token  string
	apiKey string

	session       *discordgo.Session
	removeHandler func()
}

var dg *discordgo.Session

func init() {
	var err error
	dg, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	// Add a handler to handle interaction commands
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// Handle panics within command handlers
		defer func() {
			if r := recover(); r != nil {
				log.Println("Command handler panic:", r)
			}
		}()

		// Iterate through options and handle them accordingly
		for _, opt := range i.ApplicationCommandData().Options {
			switch opt.Type {
			case discordgo.ApplicationCommandOptionBoolean:
				// Handle boolean option
				boolValue := opt.BoolValue()
				// Your logic for boolean value here
				fmt.Println("Boolean value:", boolValue)
			case discordgo.ApplicationCommandOptionString:
				// Handle string option
				stringValue := opt.StringValue()
				// Your logic for string value here
				fmt.Println("String value:", stringValue)
			// Add cases for other option types as needed
			default:
				log.Printf("Unhandled option type: %s\n", opt.Type)
			}
		}

		if h, ok := commands.CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})

	err := dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	dg.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMembers |
		discordgo.IntentsGuildPresences

	for _, v := range commands.Commands {
		_, err := dg.ApplicationCommandCreate(dg.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	defer dg.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutdowning")
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
