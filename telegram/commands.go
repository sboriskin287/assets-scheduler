package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sboriskin287/assets-scheduler/core"
)

type Command interface {
	HandleMsg(msg *tgbotapi.Message) *tgbotapi.MessageConfig
}

var startCmdMarkup = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/create_period")),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/get_period_list")))

type StartCommand struct {
}

func NewStartBotCommand() *StartCommand {
	return &StartCommand{}
}

func (cmd StartCommand) HandleMsg(msg *tgbotapi.Message) *tgbotapi.MessageConfig {
	outMsg := tgbotapi.NewMessage(msg.Chat.ID, "Select action")
	outMsg.ReplyMarkup = startCmdMarkup
	return &outMsg
}

type CreatePeriodCommand struct {
	periodService *core.PeriodService
}

func NewCreatePeriodCommand(periodService *core.PeriodService) *CreatePeriodCommand {
	return &CreatePeriodCommand{
		periodService: periodService,
	}
}

func (cmd CreatePeriodCommand) HandleMsg(msg *tgbotapi.Message) *tgbotapi.MessageConfig {
	var outMsg tgbotapi.MessageConfig
	return &outMsg
}

type GetPeriodListCommand struct {
	periodService *core.PeriodService
}

func NewGetPeriodListCommand(periodService *core.PeriodService) *GetPeriodListCommand {
	return &GetPeriodListCommand{
		periodService: periodService,
	}
}

func (cmd GetPeriodListCommand) HandleMsg(msg *tgbotapi.Message) *tgbotapi.MessageConfig {
	var outMsg tgbotapi.MessageConfig
	return &outMsg
}

type DefaultCommand struct {
}

func NewDefaultCommand() *DefaultCommand {
	return &DefaultCommand{}
}

func (cmd DefaultCommand) HandleMsg(msg *tgbotapi.Message) *tgbotapi.MessageConfig {
	outMsg := tgbotapi.NewMessage(msg.Chat.ID, "Unknown or not implemented yet command")
	return &outMsg
}
