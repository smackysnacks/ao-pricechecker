package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"pricechecker/marketdata"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const maxMessageLength = 2000

type commandHandler func(*discordgo.Session, *discordgo.MessageCreate)

type Bot struct {
	token       string
	bot         *discordgo.Session
	cmdHandlers map[string]commandHandler
}

// New creates and initializes a new Bot instance
func New(token string) (*Bot, error) {
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		err = fmt.Errorf("error creating Discord session: %v", err)
		return nil, err
	}
	b := Bot{token, bot, make(map[string]commandHandler)}

	b.setupHandlers()
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

func (b *Bot) setupHandlers() {
	b.cmdHandlers["/help"] = b.commandHelp
	b.cmdHandlers["/pc"] = b.commandPriceCheck
}

func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	tokens := strings.Split(m.Content, " ")
	if f, ok := b.cmdHandlers[tokens[0]]; ok {
		f(s, m)
	}
}

func (b *Bot) commandHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	if _, err := s.ChannelMessageSend(m.ChannelID, "help message"); err != nil {
		fmt.Printf("error sending to channel: %v\n", err)
	}
}

func (b *Bot) commandPriceCheck(s *discordgo.Session, m *discordgo.MessageCreate) {
	// get the rest of the inbound message after the "/<command>"
	rest := m.Content[4:] //m.Content[len(strings.Split(m.Content, " ")[0]):]
	log.Printf("%v\n", []byte(rest))

	q := marketdata.Closest(rest)
	log.Printf("%s: %s\n", rest, q.FriendlyName)
	if q.UniqueName == "" {
		s.ChannelMessageSend(m.ChannelID, "Sorry, I'm not sure what you meant")
		return
	}
	marketData, err := marketdata.Query(q.UniqueName)
	if err != nil {
		fmt.Printf("failed to query market data: %v\n", err)
		return
	}

	var itemOverview string
	if q.FriendlyName != "" && q.Description != "" {
		itemOverview = "```" + q.FriendlyName + " - " + q.Description + "```"
	}
	if q.FriendlyName != "" && q.Description == "" {
		itemOverview = "```" + q.FriendlyName + "```"
	}
	if len(itemOverview) > 0 {
		s.ChannelMessageSend(m.ChannelID, itemOverview)
	}

	if len(marketData) > 0 {
		itemTable := marketData.Table()
		if len(itemTable) > maxMessageLength-len(itemOverview) {
			itemTable = itemTable[:maxMessageLength-len(itemOverview)]
		}
		itemTable = "```" + itemTable + "```"
		if _, err := s.ChannelMessageSend(m.ChannelID, itemTable); err != nil {
			fmt.Println("err:", err)
		}
	}
}
