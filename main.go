package main

//Auri's main process

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/ayush6624/chatgpt-go/chatgpt"
	
	"github.com/r1c05h37/AI-Librarian-discord-bot/commands"
	"github.com/r1c05h37/AI-Librarian-discord-bot/actions"
	"github.com/r1c05h37/AI-Librarian-discord-bot/chatgpt"
	"github.com/r1c05h37/AI-Librarian-discord-bot/chatsonic"
)

// Variables used for command line parameters
var (
	DiscordToken string
	ChatgptToken string
	ChatsonicToken string

	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
	
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer

	commandlist = commands.list();
)

func init() {

	flag.StringVar(&DiscordToken, "discord", "", "Bot Token")
	flag.StringVar(&ChatgptToken, "chatgpt", "", "Bot Token")
	flag.StringVar(&ChatsonicToken, "chatsonic", "", "Bot Token")
	flag.Parse()
}

func main() {
	// New Discord session
	dg, err := discordgo.New("Bot " + DiscordToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// New ChatGPT session
	client, err := chatgpt.NewClient(ChatgptToken)
		if err != nil {
		log.Fatal(err)
	}
	
	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// To receive messages
	dg.Identify.Intents = discordgo.IntentsGuildMessages
	
	// Initialize commands
	commands.Initialize()

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	
	
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	commands.Clean()
	
	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    
	if m.Author.Bot {
		return
	}
	
	
}