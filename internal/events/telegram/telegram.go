package telegram

import (
	"fmt"
	"weatherbot/internal/clients"
	tgdomain "weatherbot/internal/domain/telegram"
	"weatherbot/internal/events"
	"weatherbot/internal/service"
)

var _ events.Processor = (*processor)(nil)
var _ events.Fetcher = (*processor)(nil)

type Meta struct {
	ChatID   int
	Username string
}

type processor struct {
	tg     clients.TelegramClient
	s      service.WeatherService
	offset int
}

func New(tgClient clients.TelegramClient, s service.WeatherService) *processor {
	return &processor{
		tg: tgClient,
		s:  s,
	}
}

// Process implements events.Processor.
func (p *processor) Process(e events.Event) error {
	switch e.Type {
	case events.Message:
		return p.processMessage(e)
	case events.Unknown:
	default:
	}

	return nil
}

// Fetch implements events.Fetcher.
func (p *processor) Fetch(limit int) ([]events.Event, error) {
	upds, err := p.tg.Updates(p.offset, limit)

	if err != nil {
		return nil, err
	}

	if len(upds) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(upds))

	for _, u := range upds {
		res = append(res, event(u))
	}

	p.offset = upds[len(upds)-1].ID + 1

	return res, nil
}

func (p *processor) processMessage(e events.Event) error {
	meta, err := meta(e)
	if err != nil {
		return err
	}

	if err := p.doCmd(e.Text, meta.ChatID, meta.Username); err != nil {
		return fmt.Errorf("failed do cmd %w", err)
	}

	return nil

}

func meta(e events.Event) (Meta, error) {
	res, ok := e.Meta.(Meta)
	if !ok {
		return Meta{}, fmt.Errorf("failed get meta data from message")
	}

	return res, nil
}

func event(upd tgdomain.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

func fetchText(upd tgdomain.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd tgdomain.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}
