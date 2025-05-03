package telegram

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"weatherbot/internal/domain"
)

const (
	HelpCmd    = "/help"
	StartCmd   = "/start"
	WeatherCmd = "/w"
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

	weather, err := p.s.GetWeather(context.Background(), userInfo)
	if err != nil {
		if errors.Is(err, domain.MainCityNotSetErr) {
			return p.tg.SendMessage(msgCityErr, chatID)
		}

		return p.tg.SendMessage("error "+err.Error(), chatID)
	}
	msg, err := convertWeatherToMessage(weather)
	if err != nil {
		return p.tg.SendMessage("internal error", chatID)
	}

	return p.tg.SendMessage(msg, chatID)
}

func (p *processor) sendUnknownCommand(chatID int, username string) error {
	return nil
}

func convertWeatherToMessage(wd *domain.WeatherDataResponse) (string, error) {
	if len(wd.Hourly.Time) == 0 {
		return "", fmt.Errorf("failed get weather")
	}

	if len(wd.Hourly.Temperature2M) == 0 {
		return "", fmt.Errorf("failed get weather")
	}
	const inputFormat = "2006-01-02T15:04"
	var message bytes.Buffer

	message.WriteString("*Прогноз погоды*\n\n")

	message.WriteString("| Время     | Температура |\n")
	message.WriteString("|-----------|-------------|\n")

	for i := range wd.Hourly.Temperature2M {
		parsedTime, err := time.Parse(inputFormat, wd.Hourly.Time[i])
		if err != nil {
			return "", fmt.Errorf("Failed parse time: %v", err)
		}
		message.WriteString(fmt.Sprintf("| `%s` | %.1f°C |\n", parsedTime.Format("15:04"), wd.Hourly.Temperature2M[i]))
	}

	message.WriteString(`\n☀️ Хорошего дня! 🌡️`)

	return message.String(), nil
}
