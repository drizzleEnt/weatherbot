package telegram

import (
	"fmt"
	"weatherbot/internal/clients/telegram"
	"weatherbot/internal/events"
)

var _ events.Processor = (*processor)(nil)
var _ events.Fetcher = (*processor)(nil)

type Meta struct {
	ChatID   int
	Username string
}

type processor struct {
	tg     *telegram.Client
	offset int
}

func New(tgClient *telegram.Client) *processor {
	return &processor{
		tg: tgClient,
	}
}

// Process implements events.Processor.
func (p *processor) Process(e events.Event) error {
	switch e.Type {
	case events.Message:
		return processMessage(e)
	case events.Unknown:
	default:
	}

	return nil
}

// Fetch implements events.Fetcher.
func (p *processor) Fetch() ([]events.Event, error) {
	panic("unimplemented")
}

func processMessage(e events.Event) error {
	meta, err := meta(e)
	if err != nil {
		return err
	}

	fmt.Printf("meta: %v\n", meta)
	return nil

}

func meta(e events.Event) (Meta, error) {
	res, ok := e.Meta.(Meta)
	if !ok {
		return Meta{}, fmt.Errorf("failed get meta data from message")
	}

	return res, nil
}
