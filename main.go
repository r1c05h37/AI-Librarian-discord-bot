package main

//Auri's main process

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ayush6624/go-chatgpt"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/ini.v1"
	
	"github.com/r1c05h37/AI-Librarian-discord-bot/personality"
	"github.com/r1c05h37/AI-Librarian-discord-bot/searchapis"
)

var (
	DiscordToken string
	ChatGPTToken string	
	PixabayToken string
	TenorToken string
	InitFile string
	WelcomeChannel string
)

var GPTclient *chatgpt.Client

var BotIDmention string

//init keys
func init() {
	cfg, err := ini.Load("config.ini")
    if err != nil {
        fmt.Printf("Failed to load config file: %v", fmt.Sprint(err))
    }

    // Assign variables from the config file
    DiscordToken = cfg.Section("keys").Key("discord_token").String()
    ChatGPTToken = cfg.Section("keys").Key("chatgpt_token").String()
    PixabayToken = cfg.Section("keys").Key("pixabay_token").String()
    TenorToken = cfg.Section("keys").Key("tenor_token").String()

    InitFile = cfg.Section("personality").Key("init_file").String()

    WelcomeChannel = cfg.Section("channels").Key("welcome_channel").String()
}

func main() {	
	
	// New Discord session
	dg, err := discordgo.New("Bot " + DiscordToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	fmt.Println("Inited discord")

	// New ChatGPT session
	GPTclient, err := chatgpt.NewClient(ChatGPTToken)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Inited chatgpt")
	
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
	fmt.Println("Inited discord websocket")
	
	
	//Initialize ChatGPT's personality
	initresponse, _ := personality.Initialize(GPTclient, InitFile)
	fmt.Println("Got response")
	if len(initresponse) > 0 { 
	fmt.Println(initresponse)
	dg.ChannelMessageSend(WelcomeChannel, initresponse)
	
	//print if initialized fully
	//dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
    //    fmt.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
    //})
    BotIDmention = fmt.Sprint("<@", dg.State.User.ID, ">")
	
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.", "\n", "Press CTRL-C to exit.")
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

	

	//if contains mention of Auri
	if strings.Contains(m.Content, BotIDmention)  {
		var editedContent1 string
		var editedContent2 string
		var responce string
		var respIMG searchapis.PixabayResult
		var respGIF searchapis.TenorResult
		var flags [][]string
		editedContent1 = strings.ReplaceAll(m.Content, BotIDmention, "Auri")
		editedContent2 = fmt.Sprint("( <@", m.Author.ID, "> says to you: )", editedContent1)
		
		GPTclient, err := chatgpt.NewClient(ChatGPTToken)
			if err != nil {
			fmt.Println(err)
		}
		
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
					var cid int
					
					PixabayClient, err := searchapis.NewPixabayClient(PixabayToken)
					if err != nil {
						fmt.Println(err)
					}
					ctx := context.Background()
					respIMG, err = PixabayClient.PixabayImageById(ctx, flags[id][1])
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, fmt.Sprint(err))
					}
					
					id = rand.Intn(len(respIMG.Hits))
					
					s.ChannelMessageSend(m.ChannelID, fmt.Sprint(respIMG.Hits[cid].LargeImageURL))
				}
				case "gif": {
					var gid int
					
					TenorClient, err := searchapis.NewTenorClient(TenorToken)
					if err != nil {
						fmt.Println(err)
					}
					ctx := context.Background()
					respGIF, err = TenorClient.TenorGifById(ctx, flags[id][1])
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, fmt.Sprint(err))
					}
					
					id = rand.Intn(len(respGIF.Results))
					s.ChannelMessageSend(m.ChannelID, fmt.Sprint(string(respGIF.Results[gid].MediaFormats.Gif.URL)))
				}
			}
		}
	}
}
