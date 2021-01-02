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
        } else {
            getPeople(sess, mess)
        }
    } else if(strings.Contains(mess.Content, Command)) {
        sendHelpMessage(sess, mess)
    }
}

func sendHelpMessage(sess *discordgo.Session, mess *discordgo.MessageCreate) {
    help := "``Bot use: \n!randomize <channel 1> <channel 2> <...>\n!randomize alone uses all channels``"
    sess.ChannelMessageSend(mess.ChannelID, help)
}

// Get list of people in given channel(s). No given channel names mean all of them.
func getPeople(sess *discordgo.Session, mess *discordgo.MessageCreate, channelNames ...string) ([]discordgo.User, error) {
    guild,_ := getGuild(sess, mess)
    members := guild.Members
    channels := guild.Channels

    var channelIDs []string

    if len(channelNames) > 0 {
        for _,channel := range channels {
            if contains(channelNames, channel.Name) {
                channelIDs = append(channelIDs, channel.ID)
            }
        }
    } else {
        channelIDs = Map(channels, func(v *discordgo.Channel) string {
            return v.ID
        })
    }

    return createUserList(guild, members, channelIDs)
}

// Helper for getPeople.
// Checks users in all voice channels against the given list of channels of interest.
func createUserList(guild *discordgo.Guild, members []*discordgo.Member, channelIDs []string) ([]discordgo.User, error) {
    var users []discordgo.User
    for _,member := range members {
        voiceState,err := findUserVoiceState(guild, member.User.ID)
        if err != nil {
            return nil, err
        }
        userChannel := voiceState.ChannelID

        if contains(channelIDs, userChannel) {
            users = append(users, *member.User)
        }
    }

    return users,nil
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

// helper contains function
func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

// helper map function
func Map(vs []*discordgo.Channel, f func(*discordgo.Channel) string) []string {
    vsm := make([]string, len(vs))
    for i, v := range vs {
        vsm[i] = f(v)
    }
    return vsm
}
