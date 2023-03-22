# AI-Librarian-discord-bot
### Handling command identification in users messages using AI
The idea of this bot is that you can use AI to understand what user wants, and respond to that using flags at the end of AI response.

All you need is to set up AI to respond to user in a few sentences, then with separator, and with list of flags after the separator. 

In example we need to make that if user asks AI to add an image, it will answer with something and adds an image (for simplicity we just assume that user can ask only for 1 exact image of beaver)
As we know, language models do not have any possibility of generating and/or adding images to the answer, so how to do that?
Actually the answer is simple - make it so AI doesn't do that, **The code itself does!**
Initialise the AI with something like:

```
You are text interpreter. On any text sent tou you, you will generate an up to three sentences long response with flag section at the end. Flag section starts with "$" symbol. You cannot use this symbol anywhere else in the response, only on the start of the flag section. There can be only one instance of every type of flag. List of flags and in which cases they are used, is such as follows:
```

And there is fun part: setting up flags. For them to work properly you have to set rules for them. You will understand better after the example:

```
- img:"keywords" - when user in any way searches, asks to generate or just asks for a picture. After a colon you define keywords for an asked image. In example: "Generate me an image of dog" will trigger the "img" flag, with value "dog", and example response would be "Here goes your image of a dog $ img:'dog'"
```

And just because I can, there is example with two values and three cases:

```
-music:"link/search","type" - when user asks to play him some music in any way. In these types of texts user usually includes a link to a youtube (which has youtube.com in it), name of the song or its description. In case if there is a link, copy this link into the "link/search" area, and put number 1 into the "type" area. If there is no link, define search keywords in "link/search" area, and put number 2 in the "type" area. If there is link that is not from youtube, reply that you can't play music outside of youtube and do not use the flag.
```

Following similar rules, you can make AI add any functionality based on users input.


### Using this in my discord bot
Basically what I want, is to handle commands without the need of any - you can text something vague, like "@Bot Me want potato image" and bot will understand it anyway.
But after thinking a bit, the realisation came in - *I can make it roleplay in chat, while it being actually useful bot* 
Or even better - with some work, this is possible to make it *behave as an individual, that can do much more than just text. It can be as almost-alive helper bot.*
But this is for later. Right now I have to implement basics - interactions and flags. I sure do hope that I will not fuck up in process of coding, and do not lose interest as I usually do.