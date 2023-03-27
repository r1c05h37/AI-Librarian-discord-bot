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
	MemoryFile string
	
	WelcomeChannel string
)

var (
	DiscordToken string
	ChatGPTToken string	
)

var GPTclient *chatgpt.Client

var BotIDmention string

//init flags
func init() {
	//flags 
	flag.StringVar(&KeyFile, "tok", "", "Tokens file: first line - Discord, second - chatGPT")
	flag.StringVar(&InitFile, "init", "", "Init file for chatGPT personality")
	flag.StringVar(&MemoryFile, "mem", "", "File for AI to use as long-time memory.")
	flag.StringVar(&WelcomeChannel, "wel", "", "Id of a channel to wake up to")
	flag.Parse()
}

//init keys
func init() {	
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
    
    //assign keys
    DiscordToken = Keys[0]
    ChatGPTToken = Keys[1]
}

//init memory

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
	
	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// To receive messages, userJoinsVc's and userExitsVC's
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	
	//Initialize ChatGPT's personality
	initresponse, _ := personality.Initialize(GPTclient, InitFile)
	if len(initresponse) > 0 { 
	fmt.Println(initresponse)
	dg.ChannelMessageSend(WelcomeChannel, initresponse)
	
	//print if initialized fully
	//dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
    //    fmt.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
    //})
    BotIDmention = fmt.Sprint("<@", dg.State.User.ID, ">")
	
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	
	// Cleanly close down the Discord session.
	dg.Close()
	}
}


//on new message
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    
    //do nothing if author is a bot
	if m.Author.Bot {
		return
	}
	
	GPTclient, err := chatgpt.NewClient(ChatGPTToken)
		if err != nil {
		fmt.Println(err)
	}
	
	//if contains mention of Auri
	if strings.Contains(m.Content, BotIDmention)  {
		var editedContent1 string
		var editedContent2 string
		var responce string
		var flags [][]string
		editedContent1 = strings.ReplaceAll(m.Content, BotIDmention, "Auri")
		editedContent2 = fmt.Sprint("( <@", m.Author.ID, "> says to you: )", editedContent1)
		
		//get response from AI and send it
		if len(editedContent2) > 0 {
			responce, flags = personality.Answer(GPTclient, editedContent2)
		}
		if len(responce) > 0 {
			s.ChannelMessageSend(m.ChannelID, responce)
		}
		
		//if responce has flags, react accordingly
		for id := range flags {
			switch flags[id][0]{
				case "img": {
					var imageURL string = "//*Images are not implemented yet*//"
					//TODO: imageURL := FindImageByKeywords(flags[id][1])
					s.ChannelMessageSend(m.ChannelID, imageURL)
				}
				case "gif": {
					var gifURL string = "//*GIFs are not implemented yet*//"
					//TODO: gifURL := FindGIFByKeywords(flags[id][1])
					s.ChannelMessageSend(m.ChannelID, gifURL)
				}
			}
		}
	}
}
