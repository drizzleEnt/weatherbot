package telegram

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"weatherbot/internal/domain"
)

const (
	HelpCmd    = "/help"
	StartCmd   = "/start"
	WeatherCmd = "/weather"
)

func (p *processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	incomeArr := strings.Split(text, " ")
	cmd := incomeArr[0]

	switch cmd {
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID, username)
	case WeatherCmd:
		if len(incomeArr) < 2 {
			incomeArr = append(incomeArr, "")
		}
		return p.sendWeather(chatID, incomeArr[1], username)
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

func (p *processor) sendWeather(chatID int, income string, username string) error {
	userInfo := domain.UserInfo{
		Username: username,
		ChatID:   chatID,
		City:     income,
	}

	err := p.s.GetWeather(context.Background(), userInfo)
	if err != nil {
		if errors.Is(err, domain.MainCityNotSetErr) {
			return p.tg.SendMessage(msgCityErr, chatID)
		}

		return p.tg.SendMessage("error "+err.Error(), chatID)
	}
	//send forecast

	return p.tg.SendMessage("weather forecast", chatID)
}

func (p *processor) sendUnknownCommand(chatID int, username string) error {
	return nil
}
