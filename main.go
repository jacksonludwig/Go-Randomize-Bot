package main

import (
    "os"
    "os/signal"
    "syscall"

    "github.com/bwmarrin/discordgo"

    "flag"
    "fmt"
    "strings"
    "errors"
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
            fmt.Println("command was: " + mess.Content)
        } else {
            getPeople(sess, mess)
            fmt.Println("command was: " + mess.Content)
        }
    } else if(strings.Contains(mess.Content, Command)) {
        sendHelpMessage(sess, mess)
    }
}

func sendHelpMessage(sess *discordgo.Session, mess *discordgo.MessageCreate) {
    help := "``Bot use: \n!randomize <channel 1> <channel 2> <...>\n!randomize alone uses all channels``"
    sess.ChannelMessageSend(mess.ChannelID, help)
}

// TODO
// Get list of people in given channel(s). No given channel names mean all of them.
func getPeople(sess *discordgo.Session, mess *discordgo.MessageCreate, channelNames ...string) ([]*discordgo.User, error) {
    guild,_ := getGuild(sess, mess)
    members := guild.Members
    if len(channelNames) > 0 {
        channels := Map(guild.Channels, func(v *discordgo.Channel) string {
            return v.Name
        })
        for _,member := range members {
            voiceState,_ := findUserVoiceState(guild, member.User.ID)
            userChannel := voiceState.ChannelID
        }

    }
    return nil, nil
}

// gets guild using message
func getGuild(sess *discordgo.Session, mess *discordgo.MessageCreate) (*discordgo.Guild, error) {
    guild, err := sess.State.Guild(mess.GuildID)
    if err != nil {
        fmt.Println(err)
        return nil, errors.New("failed to retrieve guild data")
    }
    return guild, nil
}

// find voice channel user is currently in
func findUserVoiceState(guild *discordgo.Guild, userid string) (*discordgo.VoiceState, error) {
    for _, person := range guild.VoiceStates {
        if person.UserID == userid {
            return person, nil
        }
    }

    return nil, errors.New("Could not find user's voice state")
}

// attach handlers
func engageHandlers(bot *discordgo.Session) {
    bot.AddHandler(messageCreate)
}

// helper filter function
func Filter(vs []string, f func(string) bool) []string {
    vsf := make([]string, 0)
    for _, v := range vs {
        if f(v) {
            vsf = append(vsf, v)
        }
    }
    return vsf
}

// helper map function
func Map(vs []*discordgo.Channel, f func(*discordgo.Channel) string) []string {
    vsm := make([]string, len(vs))
    for i, v := range vs {
        vsm[i] = f(v)
    }
    return vsm
}
