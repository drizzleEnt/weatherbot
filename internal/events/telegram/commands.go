package telegram

import (
	"fmt"
	"strings"
)

const (
	HelpCmd    = "/help"
	StartCmd   = "/start"
	WeatherCmd = "/weather"
)

func (p *processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	fmt.Printf("text: %v\n", text)
	switch text {
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID, username)
	case WeatherCmd:
		return p.sendWeather(chatID, text)
	default:
		p.sendUnknownCommand(chatID, username)
	}

	return nil
}

func (p *processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(msgHelp, chatID)
}

func (p *processor) sendHello(chatID int, username string) error {
	msg := fmt.Sprintf(msgHello, username)
	return p.tg.SendMessage(msg, chatID)
}

func (p *processor) sendWeather(chatID int, income string) error {
	incomeArr := strings.Split(income, " ")

	if len(incomeArr) < 2 {
		// Find city in repo

		//if non send err

		return p.tg.SendMessage(msgCityErr, chatID)
	}

	// add city in repo

	//send forecast

	return nil
}

func (p *processor) sendUnknownCommand(chatID int, username string) error {
	return nil
}
