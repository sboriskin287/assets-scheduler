package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sboriskin287/assets-scheduler/core"
	"github.com/sboriskin287/assets-scheduler/mongo"
	"log"
	"strconv"
)

var periodService = core.NewPeriodService(mongo.CreateMongoClient())
var startCmd = NewStartBotCommand()
var createPeriodCmd = NewCreatePeriodCommand(periodService)
var getPeriodListCmd = NewGetPeriodListCommand(periodService)
var defaultCmd = NewDefaultCommand()

type Listener struct {
	bot    *tgbotapi.BotAPI
	token  string
	chatId string
}

func NewListener(token string, chatId string) (*Listener, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Listener{
		bot:    bot,
		chatId: chatId,
	}, nil
}

func (l *Listener) Listen() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := l.bot.GetUpdatesChan(u)
	for upd := range updates {
		if strconv.FormatInt(upd.Message.Chat.ID, 10) != l.chatId {
			continue
		}
		if !upd.Message.IsCommand() {
			continue
		}

		var cmd Command
		var outMsg *tgbotapi.MessageConfig
		switch upd.Message.Command() {
		case "start":
			cmd = startCmd
		case "create_period":
			cmd = createPeriodCmd
		case "get_period_list":
			cmd = getPeriodListCmd
		default:
			cmd = defaultCmd
		}
		outMsg = cmd.HandleMsg(upd.Message)
		_, err := l.bot.Send(outMsg)
		if err != nil {
			log.Panic(err)
		}
	}

}
