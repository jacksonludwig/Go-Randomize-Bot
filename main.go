package main

import (
    "os"

    "github.com/bwmarrin/discordgo"

    "flag"
    "fmt"
)

var Token string

func init() {
    flag.StringVar(&Token, "token", "", "Bot Token")
    flag.Parse()
}

func main() {
}

// message handler
func messageCreate(sess *discordgo.Session, mess *discordgo.MessageCreate) {
    if mess.Author.ID == sess.State.User.ID {
        return
    }

    fmt.Printf("Message received: %s\n", mess.Content)
}

func startBot() {
    bot, err := discordgo.New("Bot " + Token)
    if err != nil {
        fmt.Println("error creating Discord session,", err)
        os.Exit(1)
    }
}

func engageHandlers(bot &discordgo.Session) {
    
}
