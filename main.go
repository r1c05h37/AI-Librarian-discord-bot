package main

//Auri's main process

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"bufio"

	"github.com/bwmarrin/discordgo"
	"github.com/ayush6624/go-chatgpt"
	
	"github.com/r1c05h37/AI-Librarian-discord-bot/personality"
	//"github.com/r1c05h37/AI-Librarian-discord-bot/searchapis"
)

// Variables
var (
	KeyFile string
	InitFile string
	DiscordToken string
	ChatGPTToken string
	GPTclient *chatgpt.Client
)

func init() {
	//flags 
	flag.StringVar(&KeyFile, "k", "", "Tokens file")
	flag.StringVar(&InitFile, "i", "", "Init file for chatgpt")
	flag.Parse()
	
	//read file for tokens
	var Keys [] string
	readFile, err := os.Open(KeyFile)
    if err != nil {
        fmt.Println(err)
    }
    fileScanner := bufio.NewScanner(readFile)
    
    fileScanner.Split(bufio.ScanLines)
  
    for fileScanner.Scan() {
        Keys = append(Keys, fileScanner.Text())
    }
    readFile.Close()
    DiscordToken = Keys[0]
    ChatGPTToken = Keys[1]
}

func main() {
	// New Discord session
	dg, err := discordgo.New("Bot " + DiscordToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// New ChatGPT session
	GPTclient, err := chatgpt.NewClient(ChatGPTToken)
		if err != nil {
		fmt.Println(err)
	}
	//Initialize ChatGPT's personality
	initresponse := personality.Initialize(GPTclient, InitFile)
	dg.ChannelMessageSend("1087031653429420113", initresponse)
	
	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// To receive messages, userJoinsVc's and userExitsVC's
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

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
	
	// Cleanly close down the Discord session.
	dg.Close()
}


//on new message
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    
    //do nothing if author is a bot
	if m.Author.Bot {
		return
	}
	
	//if contains mention of Auri
	if strings.Contains(m.Content, "<@1088044714491654205>")  {
		responce, flags := personality.Answer(GPTclient, m.Content)
		s.ChannelMessageSend(m.ChannelID, responce)
		for id := range flags {
			switch flags[id][0]{
				case "img": {
					
				}
				case "gif": {
					
				}
				case "vid": {
					
				}
			}
		}
	}
}
