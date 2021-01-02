package main

import (
    "os"
    "os/signal"
    "syscall"

    "github.com/bwmarrin/discordgo"

    "flag"
    "fmt"
    "strings"
)

var Token string
const Command string = "!randomize"

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
    // if the command's format is correct, continue
    if (strings.HasPrefix(mess.Content, Command + " ") || mess.Content == Command) {
        words := strings.Split(mess.Content, " ")
        if len(words) > 1 {
            getPeople(sess, mess, words[1:]...)
        } else {
            getPeople(sess, mess)
        }
    } else if(strings.Contains(mess.Content, Command)) {
        sendHelpMessage(sess, mess)
    }
}

func sendHelpMessage(sess *discordgo.Session, mess *discordgo.MessageCreate) {
    help := "Bot use: \n!randomize <channel 1> <channel 2> <...>\n!randomize alone uses all channels"
    sess.ChannelMessageSend(mess.ChannelID, help)
}

// Get list of people in channel(s). No given channel names mean all of them.
func getPeople(sess *discordgo.Session, mess *discordgo.MessageCreate, channelNames ...string) {
    // get correct guild from message
    guild := mess.GuildID
    members := sess.RequestGuildMembers(guild, "", 0, true)
    for _, member := range members {
        fmt.Println(member.User)
    }
    // get all channels from guild
    // channels, _ := sess.GuildChannels(guild)
    // if len(channelNames) == 0 {
    // }
}

// attach handlers
func engageHandlers(bot *discordgo.Session) {
    bot.AddHandler(messageCreate)
}
