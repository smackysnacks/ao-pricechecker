package bot

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	token string
	bot   *discordgo.Session
}

// New creates and initializes a new Bot instance
func New(token string) (*Bot, error) {
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		err = fmt.Errorf("error creating Discord session: %v", err)
		return nil, err
	}
	b := Bot{token, bot}

	bot.AddHandler(b.messageCreate)

	return &b, nil
}

// Run starts the bot
func (b *Bot) Run() {
	err := b.bot.Open()
	if err != nil {
		fmt.Printf("error opening connection: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Bot is running... Press Ctrl-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	b.bot.Close()
}

func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("messageCreate")
	// ignore messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	fmt.Printf("content of message: %s\n", m.Content)

	if m.Content == "ping" {
		msg := `**PONG**
this is one another line
something else
## try a heading
# bigger
***PONG***`
		fmt.Println("My token is:", b.token)
		_, err := s.ChannelMessageSend(m.ChannelID, msg)
		if err != nil {
			fmt.Println("There was an error")
			return
		}
	}
}
