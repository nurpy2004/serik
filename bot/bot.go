package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/bregydoc/gtranslate"
)

const (
	CmdPoll      = "poll"
	CmdPolList   = "pollist"
	CmdPollHelp  = "pollhelp"
	CmdClosePoll = "closepoll"
	CmdTranslate = "translate"
)

// list of all available commands
var commandList = []string{CmdPoll, CmdPolList, CmdPollHelp, CmdClosePoll, CmdTranslate}

// Commands is a slice of Discord Application Commands
var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        CmdPoll,
			Description: "basic command route for starting a poll",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "question",
					Description: "question for the poll",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "multiple-options",
					Description: "able to cast multiple votes",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "answer1",
					Description: "first answer",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "answer2",
					Description: "second answer",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "answer-3",
					Description: "third answer",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "answer-4",
					Description: "fourth answer",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "answer-5",
					Description: "fifth answer",
					Required:    false,
				},
			},
		},
		{
			Name:        CmdPolList,
			Description: "List all open polls",
		},
		{
			Name:        CmdPollHelp,
			Description: "get help on all commands",
		},
		{
			Name:        CmdClosePoll,
			Description: "Close a poll by id",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "poll-id",
					Description: "Poll id",
					Required:    true,
				},
			},
		},

		{
			Name:        CmdTranslate,
			Description: "Translate text to a specified language",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "text",
					Description: "Text to translate",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "target-language",
					Description: "Target language",
					Required:    true,
				},
			},
		},
	}
	// CommandHandlers is a map of command names to their corresponding handler functions
	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		CmdPoll: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			margs := []interface{}{}
			msgformat := "New poll: \n"
			if len(i.ApplicationCommandData().Options) >= 3 {
				for j, opt := range i.ApplicationCommandData().Options {
					if opt.Name == "question" {
						msgformat += "question: %s \n"
						margs = append(margs, opt.StringValue())
					} else if opt.Name == "multipleOptions" {
						msgformat += "> multipleOptions: %v\n"
						margs = append(margs, opt.BoolValue())
					} else {
						msgformat += fmt.Sprintf("answer %d", j)
						msgformat += ": %v\n"
						margs = append(margs, opt.StringValue())
					}
				}
				margs = append(margs, i.ApplicationCommandData().Options[0].StringValue())
				msgformat += "> poll-id: <#%s>\n"
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						msgformat,
						margs...,
					),
				},
			})
		},
		CmdPolList: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "List of all open polls",
				},
			})
		},
		CmdPollHelp: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			msgFormat := "All available commands: \n"
			cmdFormat := "/%s \n"
			for _, c := range commandList {
				msgFormat += fmt.Sprintf(cmdFormat, c)
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msgFormat,
				},
			})
		},
		CmdClosePoll: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			margs := []interface{}{
				// Here we need to convert raw interface{} value to wanted type.
				// Also, as you can see, here is used utility functions to convert the value
				// to particular type. Yeah, you can use just switch type,
				// but this is much simpler
				i.ApplicationCommandData().Options[0].StringValue(),
			}
			msgformat :=
				` Attempting to close:
				> poll-id: %s
`
			if len(i.ApplicationCommandData().Options) >= 2 {
				margs = append(margs, i.ApplicationCommandData().Options[0].StringValue())
				msgformat += "> poll-id: <#%s>\n"
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Ignore type for now, we'll discuss them in "responses" part
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						msgformat,
						margs...,
					),
				},
			})
		},
		CmdTranslate: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Your logic for handling the translation command
			margs := []interface{}{}
			msgformat := "Translation result: \n"

			// Extract parameters from interaction data
			textToTranslate := i.ApplicationCommandData().Options[0].StringValue()
			targetLanguage := i.ApplicationCommandData().Options[1].StringValue()

			// Call the translation function
			translated, err := gtranslate.TranslateWithParams(
				textToTranslate,
				gtranslate.TranslationParams{
					From: "en", // Assuming input text is in English
					To:   targetLanguage,
				},
			)
			if err != nil {
				panic(err)
			}

			// Include translated text in the response
			msgformat += "Source text: %s\nTranslated text: %s \n"
			margs = append(margs, textToTranslate, translated)

			// Add the translated text to the response
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(msgformat, margs...),
				},
			})
		},
	}
)
