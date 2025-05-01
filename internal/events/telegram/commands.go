package telegram

import (
	"fmt"
	"strings"
)

func (p *processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	fmt.Printf("get new message %s from %s in %v", text, username, chatID)

	return p.tg.SendMessage(text, chatID)
}
