package clients

import (
	"weatherbot/internal/domain/telegram"

	"github.com/jackc/pgx/v5"
)

type TelegramClient interface {
	Updates(offset int, limit int) ([]telegram.Update, error)
	SendMessage(text string, chatID int) error
}

type DBClient interface {
	DB() *pgx.Conn
	Close() error
}
