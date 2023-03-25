package personality

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ayush6624/go-chatgpt"
)

var (
	history []chatgpt.ChatMessage
	initText string
)

func delete_at_index(slice []chatgpt.ChatMessage, index int) []chatgpt.ChatMessage {
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
    
    //split flag section into individual flags
    flaglist := strings.Split(sep1[1], "&&")
    
    //split flags to {"flag1", "value"}, {"flag2", "value"} format
    for in := range flaglist {
    	buf := strings.Split(flaglist[in], "==")
    	flag_out = append(flag_out, buf)
    }
    return response, flag_out
}


//function initializing Auris personality
func Initialize(c *chatgpt.Client, f string) string{
	ReadPersonalityTxt(f)
	var response string
	ctx := context.Background()
	res, err := c.Send(ctx, &chatgpt.ChatCompletionRequest{
		Model: chatgpt.GPT35Turbo,
		Messages: []chatgpt.ChatMessage{
			{
				Role: chatgpt.ChatGPTModelRoleSystem,
				Content: initText,
			},
		},
	})
	if err != nil {
		response = fmt.Sprint(err)
	} else {
		a, _ := json.MarshalIndent(res, "", "  ");
		response = string(a)	
	}
	return response
}

//function to answer user
func Answer(c *chatgpt.Client, InMessage string) (string, [][]string) {
	var(	
		UncleanMessage string
		OutMessage string
		UserMessage chatgpt.ChatMessage = chatgpt.ChatMessage{
			Role: chatgpt.ChatGPTModelRoleUser,
			Content: InMessage,
		}
		OutHistory chatgpt.ChatMessage
	)
		
	//add user message to history
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
		a, _ := json.MarshalIndent(res, "", "  ");
		UncleanMessage = string(a)	
	}
	
	//format response before saving to history
	OutHistory = chatgpt.ChatMessage{
		Role: chatgpt.ChatGPTModelRoleUser,
		Content: OutMessage,
	}
	
	//add response to history
	history = append(history, OutHistory)
	
	//delete 2 oldest non-system messages from history, if there's too much history.
	if len(history) >= 100 {
		delete_at_index(history, 2)
		delete_at_index(history, 2)
	}
	
	//format output to usable variables
	OutMessage, flags := FormatOut(UncleanMessage)
	
	return OutMessage, flags
}