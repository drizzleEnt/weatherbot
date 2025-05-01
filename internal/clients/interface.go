package clients

import (
	"database/sql"
	"weatherbot/internal/domain/telegram"
)

type TelegramClient interface {
	Updates(offset int, limit int) ([]telegram.Update, error)
	SendMessage(text string, chatID int) error
}

type DBClient interface {
	DB() *sql.DB
	Close() error
}
