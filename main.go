package main

import (
	"os"
	"os/signal"
	"syscall"

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
    bot := startBot()
    openConnection(bot)
    engageHandlers(bot)
    awaitTermination()
    bot.Close()
}

func startBot() *discordgo.Session {
    bot, err := discordgo.New("Bot " + Token)
    if err != nil {
        fmt.Println("error creating Discord session,", err)
        os.Exit(1)
    }
    bot.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
    return bot
}

// enable connect to server
func openConnection(bot *discordgo.Session) {
    err := bot.Open()
    if err != nil {
        fmt.Println("error opening connection,", err)
        os.Exit(2)
    }
}

func awaitTermination() {
    fmt.Println("Bot is now running.  Press CTRL-C to exit.")
    closeChannel := make(chan os.Signal, 1)
    signal.Notify(closeChannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-closeChannel
}

// message handler
func messageCreate(sess *discordgo.Session, mess *discordgo.MessageCreate) {
    if mess.Author.ID == sess.State.User.ID {
        return
    }
    fmt.Printf("Message received: %s\n", mess.Content)
}

// attach handlers
func engageHandlers(bot *discordgo.Session) {
    bot.AddHandler(messageCreate)
}
