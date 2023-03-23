package commands

//Slash commands loader

import (

    "github.com/bwmarrin/discordgo"
)

var s *discordgo.Session

var (
    commands = []*discordgo.ApplicationCommand{
		{
		    Name:        "image",
			Description: "Find images (without AI)",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "keywords",
					Description: "keywords to search for image",
					Required:    false,
				},
            },
        },
        {
            Name:        "Count",
			Description: "Counts number of notes in archive",
        },
    }
    commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
        "image": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Images are unavailable at the time",
				},
			})
		},   
        "count": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Your mother",
				},
			})
		},   
    }
)

func Initialize() {
	// Add handler for commands
    s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	
	// Register commands
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
}

func Clean() {
	for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}