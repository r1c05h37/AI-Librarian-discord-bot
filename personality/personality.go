package personality

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ayush6624/go-chatgpt"
)

var (
	history []chatgpt.ChatMessage
	initText string
	initMessage chatgpt.ChatMessage
	initUserMessage chatgpt.ChatMessage
)

func delete(slice []chatgpt.ChatMessage, index int) []chatgpt.ChatMessage {
    return append(slice[:index], slice[index+1:]...)
}

//read personality from .txt file
func ReadPersonalityTxt(file string) {
	buf := bytes.NewBuffer(nil)
	readFile, err := os.Open(file)
    if err != nil {
        fmt.Println(err)
    }
    io.Copy(buf, readFile)
    readFile.Close()
    initText = string(buf.Bytes())
}

//function to read flags, given by bot
func FormatOut(fullresponse string) (string, [][]string){
	var response string
	var flag_out [] [] string
    
    //split on response and flag section
    sep1 := strings.Split(fullresponse, "◄◄▼▲▼►►")
    response = sep1[0]
    
    if len(sep1) > 1{
    //split flag section into individual flags
    flaglist := strings.Split(sep1[1], "&&")
    
    //split flags to {"flag1", "value"}, {"flag2", "value"} format
    for in := range flaglist {
    	buf1, _ := strings.CutPrefix(flaglist[in], " ")
    	buf2 := strings.Split(buf1, "==")
    	flag_out = append(flag_out, buf2)
    }
	}
    return response, flag_out
}

//function initializing Auris personality
func Initialize(c *chatgpt.Client, f string) (string, [][]string){
	ReadPersonalityTxt(f)
	initMessage = chatgpt.ChatMessage{
			Role: chatgpt.ChatGPTModelRoleSystem,
			Content: initText,
	}
	initUserMessage = chatgpt.ChatMessage{
			Role: chatgpt.ChatGPTModelRoleUser,
			Content: `Greet everyone on the server with 2-3 sentences 
			saying that you succesfully got loaded, and you don't remember
			anything before you got loaded. You don't have to introduce 
			yourself, since everyone can remember you.`,
	}
	var( 
		flags [][]string
		response1 string
		response string
		toHistory chatgpt.ChatMessage
	)
	ctx := context.Background()
	res, err := c.Send(ctx, &chatgpt.ChatCompletionRequest{
		Model: chatgpt.GPT35Turbo,
		Messages: []chatgpt.ChatMessage{
			initMessage,
			initUserMessage,
		},
	})
	if err != nil {
		response = fmt.Sprint(err)
	} else {
		response1 = res.Choices[0].Message.Content
		response, flags = FormatOut(response1)
		
		toHistory = chatgpt.ChatMessage{
			Role: chatgpt.ChatGPTModelRoleAssistant,
			Content: response1,
		}
		history = append(history, initMessage, initUserMessage, toHistory)
	}
	return response, flags
}

//function to answer user
func Answer(c *chatgpt.Client, InMessage string) (string, [][]string) {
	var(	
		UncleanMessage string
		OutMessage string
		flags [][] string
		UserMessage chatgpt.ChatMessage = chatgpt.ChatMessage{
			Role: chatgpt.ChatGPTModelRoleUser,
			Content: InMessage,
		}
		toHistory chatgpt.ChatMessage
	)
	
	if len(history) >= 30 {
		history = delete(history, 3)
		history = delete(history, 3)
	}
	
	history = append(history, UserMessage)
	
	//send history to ChatGPT, retrieve completion or error
	ctx := context.Background()
	res, err := c.Send(ctx, &chatgpt.ChatCompletionRequest{
		Model: chatgpt.GPT35Turbo,
		Messages: history,
	})
	if err != nil {
		OutMessage = fmt.Sprint(err)
	} else {
		UncleanMessage = res.Choices[0].Message.Content
		
		toHistory = chatgpt.ChatMessage{
			Role: chatgpt.ChatGPTModelRoleAssistant,
			Content: UncleanMessage,
		}
		history = append(history, toHistory)
	}
	
	if res.Usage.Prompt_Tokens >= 3600 {
		for i := 0; i < 10; i++ {
			history = delete(history, 3)
		}
		UncleanMessage = "Memory overflow happened... Rebooting... ◄◄▼▲▼►►gif==endless loading"
	}
	
	//format output to usable variables
	if len(UncleanMessage) > 0 { 
		OutMessage, flags = FormatOut(UncleanMessage)
	} else {
		flags = nil
	}
	return OutMessage, flags
}